package zlog

import (
	"context"
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
	_          log.Logger       = (*Zlog)(nil)
	_          plugin.Factory   = (*Factory)(nil)
	_          plugin.PluginIns = (*Zlog)(nil)
	_zaplogger *Zlog
	_factory   = &Factory{}
)

const (
	_name     = "log_zap"
	_logtmfmt = "2006-01-02 15:04:05.001"
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

type Factory struct {
}

func (f *Factory) Type() plugin.Type {
	return plugin.Log
}
func (f *Factory) Name() string {
	return _name
}

func (f *Factory) Setup(ctx context.Context, c map[string]interface{}) (plugin.PluginIns, error) {
	if _zaplogger != nil {
		return _zaplogger, nil
	}
	var cfg ZlogCfg
	if err := mapstructure.Decode(c, &cfg); err != nil {
		return nil, errors.Wrap(err, "Zlog Decode Cfg")
	}
	loglv := log.ConvertLogLv(cfg.LogLevel)

	var zw zapcore.WriteSyncer
	if cfg.Filename != "" {
		zw = zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxAge:     cfg.MaxAge,
			MaxBackups: cfg.MaxBackups,
			LocalTime:  cfg.LocalTime,
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
	if cfg.Encoding == "json" {
		newCore = zapcore.NewCore(zapcore.NewJSONEncoder(encoderConf), zw, zap.NewAtomicLevelAt(_loglvconvert[loglv]))
	} else {
		newCore = zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConf), zw, zap.NewAtomicLevelAt(_loglvconvert[loglv]))
	}
	_zaplogger = &Zlog{
		loglv:  loglv,
		logger: zap.New(newCore),
	}
	if _, err := log.New(_zaplogger); err != nil {
		return nil, err
	}
	return _zaplogger, nil
}

func (f *Factory) Destroy(plugin.PluginIns) error {
	return nil
}

func (f *Factory) Reload(plugin.PluginIns, map[string]interface{}) error {
	return nil
}

func (f *Factory) CanUnload(plugin.PluginIns) bool {
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
