package etcd

import (
	"context"
	"fmt"
	"m3game/config"
	"m3game/plugins/lease"
	"m3game/plugins/log"
	"m3game/runtime/plugin"
	"strings"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	_instance *Lease
	_cfg      etcdLeaseCfg
	_factory  = &Factory{}
)

func init() {
	plugin.RegisterFactory(_factory)
}

const (
	_factoryname = "lease_etcd"
)

type etcdLeaseCfg struct {
	Endpoints         string `mapstructure:"Endpoints"`
	DialTimeout       int    `mapstructure:"DialTimeout"`
	LeaseKeepLiveTime int    `mapstructure:"LeaseKeepLiveTime"`
	PreExitTime       int    `mapstructure:"PreExitTime"`
	EndpointsList     []string
}

func (c *etcdLeaseCfg) CheckVaild() error {
	if c.Endpoints == "" {
		return errors.New("Endpoints cant be space")
	}
	c.EndpointsList = strings.Split(c.Endpoints, ",")
	if c.DialTimeout == 0 {
		return errors.New("DialTimeout cant be 0")
	}
	if c.LeaseKeepLiveTime == 0 {
		return errors.New("LeaseKeepLiveTime cant be 0")
	}
	if c.PreExitTime == 0 {
		return errors.New("PreExitTime cant be 0")
	}
	return nil
}

type Factory struct {
}

func (f *Factory) Type() plugin.Type {
	return plugin.Lease
}
func (f *Factory) Name() string {
	return _factoryname
}

func (f *Factory) Setup(c map[string]interface{}) (plugin.PluginIns, error) {
	if err := mapstructure.Decode(c, &_cfg); err != nil {
		return nil, errors.Wrap(err, "Lease Decode Cfg")
	}
	if err := _cfg.CheckVaild(); err != nil {
		return nil, err
	}
	config := clientv3.Config{
		Endpoints:   _cfg.EndpointsList,
		DialTimeout: time.Duration(_cfg.DialTimeout) * time.Second,
	}
	_instance = &Lease{
		leasemap: make(map[string]lease.LeaseMoveOutFunc),
	}
	var err error
	if _instance.client, err = clientv3.New(config); err != nil {
		return nil, err
	}
	_instance.lease = clientv3.NewLease(_instance.client)
	if leaseGrantResp, err := _instance.lease.Grant(context.Background(), int64(_cfg.LeaseKeepLiveTime)); err != nil {
		return nil, err
	} else {
		_instance.leaseId = leaseGrantResp.ID
	}
	var keepRespChan <-chan *clientv3.LeaseKeepAliveResponse
	ctx, cancel := context.WithCancel(context.Background())
	_instance.cancel = cancel
	if keepRespChan, err = _instance.lease.KeepAlive(ctx, _instance.leaseId); err != nil {
		return nil, err
	}
	go func() {
		t := time.NewTicker(1 * time.Second)
		_instance.timeout = time.Now().Unix()
		for {
			select {
			case keepResp, ok := <-keepRespChan:
				if !ok || keepResp == nil {
					_instance.safecancel(context.Background())
					return
				}
				_instance.timeout = time.Now().Unix() + keepResp.TTL
			case <-t.C:
				if _instance.timeout < time.Now().Unix()+int64(_cfg.PreExitTime) {
					_instance.safecancel(context.Background())
					return
				}
			}
		}
	}()
	_instance.kv = clientv3.NewKV(_instance.client)
	lease.Set(_instance)
	log.Info("EtcdLease...........")
	return _instance, nil
}

func (f *Factory) Destroy(p plugin.PluginIns) error {
	l := p.(*Lease)
	l.safecancel(context.Background())
	return nil
}

func (f *Factory) Reload(plugin.PluginIns, map[string]interface{}) error {
	return nil
}

func (f *Factory) CanDelete(p plugin.PluginIns) bool {
	l := p.(*Lease)
	return l.isstoped
}

type Lease struct {
	client   *clientv3.Client
	lease    clientv3.Lease
	leaseId  clientv3.LeaseID
	kv       clientv3.KV
	timeout  int64
	leasemap map[string]lease.LeaseMoveOutFunc
	isstoped bool
	mutex    sync.Mutex
	cancel   context.CancelFunc
}

func (r *Lease) Factory() plugin.Factory {
	return _factory
}

func (r *Lease) safecancel(ctx context.Context) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.isstoped {
		return
	}
	for _, f := range r.leasemap {
		f(ctx)
	}
	r.isstoped = true
	r.cancel()
}

func (r *Lease) AllocLease(ctx context.Context, id string, f lease.LeaseMoveOutFunc) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.isstoped {
		return errors.New("Lease Closed")
	}
	if _, ok := r.leasemap[id]; ok {
		return fmt.Errorf("Lease %s alloced", id)
	}
	txn := r.kv.Txn(ctx)
	txn.If(clientv3.Compare(clientv3.CreateRevision(id), "=", 0)).
		Then(clientv3.OpPut(id, config.GetAppID().String(), clientv3.WithLease(r.leaseId))).
		Else(clientv3.OpGet(id))
	if txnResp, err := txn.Commit(); err != nil {
		return err
	} else if !txnResp.Succeeded {
		return fmt.Errorf("Lease %s Value %s", id, string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
	}
	r.leasemap[id] = f
	return nil
}

func (r *Lease) FreeLease(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.isstoped {
		return errors.New("Lease Closed")
	}
	if _, err := r.kv.Delete(context.Background(), id); err != nil {
		return err
	}
	if _, ok := r.leasemap[id]; ok {
		delete(r.leasemap, id)
	}
	return nil
}

func (r *Lease) RecvKickLease(ctx context.Context, id string) ([]byte, error) {
	f := func() lease.LeaseMoveOutFunc {
		r.mutex.Lock()
		defer r.mutex.Unlock()
		if r.isstoped {
			return nil
		}
		if f, ok := r.leasemap[id]; !ok {
			return func(ctx context.Context) ([]byte, error) {
				return nil, r.FreeLease(ctx, id)
			}
		} else {
			return f
		}
	}()
	if f == nil {
		return nil, nil
	}
	return f(ctx)
}

func (r *Lease) KickLease(ctx context.Context, id string) ([]byte, error) {
	if getrsp, err := r.kv.Get(context.Background(), id); err != nil {
		return nil, err
	} else if getrsp.Count == 0 {
		return nil, nil
	} else {
		return lease.GetReciver().SendKickLease(ctx, id, string(getrsp.Kvs[0].Value))
	}
}

func (r *Lease) GetLease(ctx context.Context, id string) ([]byte, error) {
	if getrsp, err := r.kv.Get(context.Background(), id); err != nil {
		return nil, err
	} else if getrsp.Count == 0 {
		return nil, nil
	} else {
		return getrsp.Kvs[0].Value, nil
	}
}
