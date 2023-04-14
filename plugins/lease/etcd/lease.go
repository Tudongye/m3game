package etcd

import (
	"context"
	"m3game/config"
	"m3game/meta/errs"
	"m3game/plugins/lease"
	"m3game/plugins/log"
	"m3game/runtime/plugin"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	_etcdlease *Lease
	_factory   = &Factory{}
)

func init() {
	plugin.RegisterFactory(_factory)
}

const (
	_name = "lease_etcd"
)

type etcdLeaseCfg struct {
	Endpoints         string `mapstructure:"Endpoints" validate:"required"`
	DialTimeout       int    `mapstructure:"DialTimeout" validate:"gt=0"`
	LeaseKeepLiveTime int    `mapstructure:"LeaseKeepLiveTime" validate:"gt=0"`
	PreExitTime       int    `mapstructure:"PreExitTime" validate:"gt=0"`
	EndpointsList     []string
}

type Factory struct {
}

func (f *Factory) Type() plugin.Type {
	return plugin.Lease
}
func (f *Factory) Name() string {
	return _name
}

func (f *Factory) Setup(ctx context.Context, c map[string]interface{}) (plugin.PluginIns, error) {
	var cfg etcdLeaseCfg
	if err := mapstructure.Decode(c, &cfg); err != nil {
		return nil, errs.EtcdSetupFail.Wrap(err, "Lease Decode Cfg")
	}
	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, errs.EtcdSetupFail.Wrap(err, "")
	}
	cfg.EndpointsList = strings.Split(cfg.Endpoints, ",")
	config := clientv3.Config{
		Endpoints:   cfg.EndpointsList,
		DialTimeout: time.Duration(cfg.DialTimeout) * time.Second,
	}
	_etcdlease = &Lease{}
	var err error
	if _etcdlease.client, err = clientv3.New(config); err != nil {
		return nil, errs.EtcdSetupFail.Wrap(err, "")
	}
	_etcdlease.lease = clientv3.NewLease(_etcdlease.client)
	if leaseGrantResp, err := _etcdlease.lease.Grant(ctx, int64(cfg.LeaseKeepLiveTime)); err != nil {
		return nil, errs.EtcdSetupFail.Wrap(err, "")
	} else {
		_etcdlease.leaseId = leaseGrantResp.ID
	}
	var keepRespChan <-chan *clientv3.LeaseKeepAliveResponse
	ctx, cancel := context.WithCancel(ctx)
	_etcdlease.cancel = cancel
	if keepRespChan, err = _etcdlease.lease.KeepAlive(ctx, _etcdlease.leaseId); err != nil {
		return nil, errs.EtcdSetupFail.Wrap(err, "")
	}
	go func() {
		t := time.NewTicker(1 * time.Second)
		_etcdlease.timeout = time.Now().Unix()
		for {
			select {
			case keepResp, ok := <-keepRespChan:
				if !ok || keepResp == nil {
					_etcdlease.safecancel(context.Background())
					return
				}
				_etcdlease.timeout = time.Now().Unix() + keepResp.TTL
			case <-t.C:
				if _etcdlease.timeout < time.Now().Unix()+int64(cfg.PreExitTime) {
					_etcdlease.safecancel(context.Background())
					return
				}
			}
		}
	}()
	_etcdlease.kv = clientv3.NewKV(_etcdlease.client)
	if _, err := lease.New(_etcdlease); err != nil {
		return nil, err
	}
	log.Info("EtcdLease...........")
	return _etcdlease, nil
}

func (f *Factory) Destroy(p plugin.PluginIns) error {
	l := p.(*Lease)
	l.safecancel(context.Background())
	return nil
}

func (f *Factory) Reload(plugin.PluginIns, map[string]interface{}) error {
	return nil
}

func (f *Factory) CanUnload(p plugin.PluginIns) bool {
	l := p.(*Lease)
	return l.isstoped
}

type Lease struct {
	client   *clientv3.Client
	lease    clientv3.Lease
	leaseId  clientv3.LeaseID
	kv       clientv3.KV
	timeout  int64
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
	r.isstoped = true
	r.cancel()
}

func (r *Lease) AllocLease(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.isstoped {
		return errs.EtcdIsClosed.New("Lease Closed")
	}
	txn := r.kv.Txn(ctx)
	txn.If(clientv3.Compare(clientv3.CreateRevision(id), "=", 0)).
		Then(clientv3.OpPut(id, config.GetAppID().String(), clientv3.WithLease(r.leaseId))).
		Else(clientv3.OpGet(id))
	if txnResp, err := txn.Commit(); err != nil {
		return err
	} else if !txnResp.Succeeded {
		return errs.EtcdAllocLeaseFail.New("Lease %s Value %s", id, string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
	}
	return nil
}

func (r *Lease) FreeLease(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.isstoped {
		return errs.EtcdIsClosed.New("Lease Closed")
	}
	if _, err := r.kv.Delete(context.Background(), id); err != nil {
		return err
	}
	return nil
}

func (r *Lease) WhereIsLease(ctx context.Context, id string) (string, error) {
	if getrsp, err := r.kv.Get(context.Background(), id); err != nil {
		return "", err
	} else if getrsp.Count == 0 {
		return "", nil
	} else {
		return string(getrsp.Kvs[0].Value), nil
	}
}
