package zlog

import (
	"fmt"
	"m3game/plugins/log"
	"m3game/runtime/plugin"
	"os"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	_         log.Logger       = (*Zlog)(nil)
	_         plugin.Factory   = (*Factory)(nil)
	_         plugin.PluginIns = (*Zlog)(nil)
	_instance *Zlog
	_factory  = &Factory{}
	_cfg      ZlogCfg
)

const (
	_factoryname = "log_zap"
	_logtmfmt    = "2006-01-02 15:04:05.001"
)

func init() {
	plugin.RegisterFactory(_factory)
}

type ZlogCfg struct {
	LogLevel   string `mapstructure:"LogLevel"`   //
	Encoding   string `mapstructure:"Encoding"`   //json or console
	Filename   string `mapstructure:"Filename"`   //log file name
	MaxSize    int    `mapstructure:"MaxSize"`    //max size of log.(MB)
	MaxAge     int    `mapstructure:"MaxAge"`     //time to keep, (day)
	MaxBackups int    `mapstructure:"MaxBackups"` //max file numbers
	LocalTime  bool   `mapstructure:"LocalTime"`  //(default UTC)
	Compress   bool   `mapstructure:"Compress"`   //default false
}

func (c *ZlogCfg) CheckVaild() error {

	return nil
}

type Factory struct {
}

func (f *Factory) Type() plugin.Type {
	return plugin.Log
}
func (f *Factory) Name() string {
	return _factoryname
}

func (f *Factory) Setup(c map[string]interface{}) (plugin.PluginIns, error) {
	if _instance != nil {
		return _instance, nil
	}
	if err := mapstructure.Decode(c, &_cfg); err != nil {
		return nil, errors.Wrap(err, "Zlog Decode Cfg")
	}
	if err := _cfg.CheckVaild(); err != nil {
		return nil, err
	}
	loglv := log.ConvertLogLv(_cfg.LogLevel)

	var zw zapcore.WriteSyncer
	if _cfg.Filename != "" {
		zw = zapcore.AddSync(&lumberjack.Logger{
			Filename:   _cfg.Filename,
			MaxSize:    _cfg.MaxSize,
			MaxAge:     _cfg.MaxAge,
			MaxBackups: _cfg.MaxBackups,
			LocalTime:  _cfg.LocalTime,
		})
	} else {
		zw = os.Stdout
	}

	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(_logtmfmt))
	}

	encoderConf := zapcore.EncoderConfig{
		TimeKey:        "ts",
		MessageKey:     "msg",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     customTimeEncoder, // 自定义时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	var newCore zapcore.Core
	if _cfg.Encoding == "json" {
		newCore = zapcore.NewCore(zapcore.NewJSONEncoder(encoderConf), zw, zap.NewAtomicLevelAt(_loglvconvert[loglv]))
	} else {
		newCore = zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConf), zw, zap.NewAtomicLevelAt(_loglvconvert[loglv]))
	}
	_instance = &Zlog{
		loglv:  loglv,
		logger: zap.New(newCore),
	}
	log.Set(_instance)
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

type Zlog struct {
	loglv  log.LogLv
	logger *zap.Logger
}

func (c *Zlog) Factory() plugin.Factory {
	return _factory
}

func (c *Zlog) Output(depth log.Depth, lv log.LogLv, plus log.LogPlus, format string, v ...interface{}) {
	file, funcName, line := depth.String()
	var logstr string
	if len(plus) == 0 {
		logstr = fmt.Sprintf("[%s] [%s:%d %s] %s", lv.String(), file, line, funcName, fmt.Sprintf(format, v...))
	} else {
		logstr = fmt.Sprintf("[%s] [%s:%d %s][%s] %s", lv.String(), file, line, funcName, plus.String(), fmt.Sprintf(format, v...))
	}
	switch lv {
	case log.LogLvDebug:
		c.logger.Sugar().Debug(logstr)
	case log.LogLvInfo:
		c.logger.Sugar().Info(logstr)
	case log.LogLvWarn:
		c.logger.Sugar().Warn(logstr)
	case log.LogLvError:
		c.logger.Sugar().Error(logstr)
	case log.LogLvFatal:
		c.logger.Sugar().Fatal(logstr)
	}
}

func (c *Zlog) SetLevel(level log.LogLv) {
	c.loglv = level
}

func (c *Zlog) GetLevel() log.LogLv {
	return c.loglv
}

var (
	_loglvconvert = map[log.LogLv]zapcore.Level{
		log.LogLvDebug: zap.DebugLevel,
		log.LogLvInfo:  zap.InfoLevel,
		log.LogLvWarn:  zap.WarnLevel,
		log.LogLvError: zap.ErrorLevel,
		log.LogLvFatal: zap.FatalLevel,
	}
)
