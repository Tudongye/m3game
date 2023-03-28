package shape

import (
	"m3game/meta/errs"
	"m3game/plugins/log"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type ShapeRuleFile struct {
	Rules []Rule `toml:"Rules"`
}

type Rule struct {
	Method     string      `toml:"Method"`
	FlowRules  []FlowRule  `toml:"FlowRules"`
	BreakRules []BreakRule `toml:"BreakRules"`
}

// 限流管理
type FlowRule struct {
	Threshold      int `toml:"Threshold"`      // 限流阈值
	StatIntervalMs int `toml:"StatIntervalMs"` // 统计周期
	MaxQueueWaitMs int `toml:"MaxQueueWaitMs"` // 限流等待
}

const (
	SlowRequst int = 1 // 慢调用, Threshold 千分比
	ErrorRatio int = 2 // 错误比例, Threshold 千分比
	ErrorCount int = 3 // 错误数
)

func IsSlowRequst(s string) bool {
	return s == "SlowRequst"
}

func IsErrorRatio(s string) bool {
	return s == "ErrorRatio"
}

func IsErrorCount(s string) bool {
	return s == "ErrorCount"
}

// 熔断管理
type BreakRule struct {
	Threshold      int    `toml:"Threshold"`      // 熔断阈值
	Strategy       string `toml:"Strategy"`       // 熔断策略
	SlowRequestMs  int    `toml:"SlowRequestMs"`  // 慢查询时长
	StatIntervalMs int    `toml:"StatIntervalMs"` // 统计周期
	RetryTimeOutMs int    `toml:"RetryTimeOutMs"` // 恢复时长
	MinRequestNum  int    `toml:"MinRequestNum"`  // 最小统计次数
}

type ShapeCfg struct {
	RuleConfigFile string `toml:"RuleConfigFile"`
}

func Setup(c map[string]interface{}) error {
	if Instance() == nil {
		log.Info("Shape.Setup.NoShape...")
		return nil
	}
	var scfg ShapeCfg
	if err := mapstructure.Decode(c, &scfg); err != nil {
		return errs.ShapeRuleInitFail.Wrap(err, "Shape Decode Cfg")
	}
	var cfg ShapeRuleFile
	v := viper.New()
	v.SetConfigFile(scfg.RuleConfigFile)
	v.SetConfigType("toml")
	if err := v.ReadInConfig(); err != nil {
		return errs.ShapeRuleInitFail.Wrap(err, "")
	}
	if err := v.Unmarshal(&cfg); err != nil {
		log.Error("UnMarshal ShapeRuleFile %s", err.Error())
		return errs.ShapeRuleInitFail.Wrap(err, "")
	}
	return Instance().RegisterRule(cfg.Rules)
}
