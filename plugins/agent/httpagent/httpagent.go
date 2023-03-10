package httpagent

import (
	"encoding/json"
	"fmt"
	"io"
	"m3game/config"
	"m3game/plugins/agent"
	"m3game/plugins/log"
	"m3game/plugins/store"
	"m3game/runtime/plugin"
	"m3game/util"
	"net/http"

	"github.com/hashicorp/consul/api"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

var (
	_         agent.Agent    = (*Agent)(nil)
	_         plugin.Factory = (*Factory)(nil)
	_cfg      httpAgentCfg
	_instance *Agent
	_factory  = &Factory{}
)

const (
	_factoryname = "agent_http"
)

func init() {
	plugin.RegisterFactory(_factory)
}

type httpAgentCfg struct {
	Port         int    `mapstructure:"Port"`
	LogicUrl     string `mapstructure:"LogicUrl"`
	AuthUrl      string `mapstructure:"AuthUrl"`
	SessionLiveS int64  `mapstructure:"SessionLiveS"`
}

func (c *httpAgentCfg) CheckVaild() error {
	if c.Port == 0 {
		return errors.New("Port cant be 0")
	}
	if c.LogicUrl == "" {
		return errors.New("LogicUrl cant be space")
	}
	if c.AuthUrl == "" {
		return errors.New("AuthUrl cant be space")
	}
	return nil
}

type Factory struct {
}

func (f *Factory) Type() plugin.Type {
	return plugin.Agent
}

func (f *Factory) Name() string {
	return _factoryname
}

func (f *Factory) Setup(c map[string]interface{}) (plugin.PluginIns, error) {
	if _instance != nil {
		return _instance, nil
	}
	if err := mapstructure.Decode(c, &_cfg); err != nil {
		return nil, errors.Wrap(err, "Agent Decode Cfg")
	}
	if err := _cfg.CheckVaild(); err != nil {
		return nil, err
	}
	_instance = &Agent{}
	_instance.mux = http.NewServeMux()
	_instance.mux.Handle(_cfg.AuthUrl, authHandle{})
	_instance.mux.Handle(_cfg.LogicUrl, logicHandle{})
	listenaddr := fmt.Sprintf(":%d", _cfg.Port)
	go func() {
		log.Info("HttpAgent Listen %s", listenaddr)
		if err := http.ListenAndServe(listenaddr, _instance.mux); err != nil {
			if err == http.ErrServerClosed {
				log.Info("GateWay Server Close")
			} else {
				panic(fmt.Sprintf("GateWay Server Err %s", err.Error()))
			}
		}
	}()

	return _instance, nil
}

func (f *Factory) Destroy(plugin.PluginIns) error {
	return nil
}

func (f *Factory) Reload(plugin.PluginIns, map[string]interface{}) error {
	return nil
}

func (f *Factory) CanDelete(plugin.PluginIns) bool {
	return false
}

type Agent struct {
	client *api.Client
	mux    *http.ServeMux
}

func (r *Agent) Factory() plugin.Factory {
	return _factory
}

type authHandle struct {
}

func (h authHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var ap agent.AuthPara
	if err := json.NewDecoder(r.Body).Decode(&ap); err != nil {
		http.Error(w, "Can't decode body", http.StatusBadRequest)
		return
	}
	if uid, err := agent.Auther()(ap); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		sessionid := util.GenSessionID(uid)
		session := store.NewSession(sessionid, _cfg.SessionLiveS)
		session.Set(agent.SessionKey_Uid, uid)
		rsp := &agent.AuthRsp{
			Uid:     uid,
			Session: sessionid,
			EnvID:   config.GetEnvID(),
			WorldID: config.GetWorldID(),
		}
		json.NewEncoder(w).Encode(rsp)
		return
	}
}

type logicHandle struct {
}

func (h logicHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	method := r.Header.Get("method")
	if method == "" {

		http.Error(w, "Not find Method", 512)
		return
	}
	sessionid := r.Header.Get("session")
	if sessionid == "" {
		http.Error(w, "Not find SessionID", 513)
		return
	}
	uid := r.Header.Get("uid")
	if uid == "" {
		http.Error(w, "Not find Uid", 514)
		return
	}
	if session := store.GetSession(sessionid); session == nil {
		http.Error(w, "SessionID invaild", 515)
		return
	} else if session.Get(agent.SessionKey_Uid) != uid {
		http.Error(w, "SessionID Uid not match", 516)
		return
	}
	if rsp, err := agent.Caller()(method, uid, r.Body); err != nil {
		type FailMsg struct {
			Tips string
		}
		s, _ := json.Marshal(FailMsg{Tips: fmt.Sprintf("Caller Fial:%s", err.Error())})
		w.Write(s)
		return
	} else {
		io.Copy(w, rsp)
	}
}
