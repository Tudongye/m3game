package metric

import (
	"fmt"
	"m3game/plugins/log"
)

var (
	_metric        Metric
	_defaultmetric = &defaultMetric{}
)

func New(me Metric) (Metric, error) {
	if _metric != nil {
		log.Fatal("Metric Only One")
		return nil, fmt.Errorf("Metric is newed ")
	}
	_metric = me
	return _metric, nil
}

func Instance() Metric {
	if _metric == nil {
		return _defaultmetric
	}
	return _metric
}

type Metric interface {
	NewCounter(key string, group string) StatCounter
	NewGauge(key string, group string) StatGauge
	NewHistogram(key string, group string) StatHistogram
	NewSummary(key string, group string) StatSummary
}

type defaultStat struct {
}

func (*defaultStat) Add(float64) {

}
func (*defaultStat) Set(float64) {

}
func (*defaultStat) Sub(float64) {

}
func (*defaultStat) Inc() {

}
func (*defaultStat) Dec() {

}

func (*defaultStat) Observe(float64) {

}

type defaultMetric struct {
}

func (*defaultMetric) NewCounter(key string, group string) StatCounter {
	return &defaultStat{}
}
func (*defaultMetric) NewGauge(key string, group string) StatGauge {
	return &defaultStat{}
}
func (*defaultMetric) NewHistogram(key string, group string) StatHistogram {
	return &defaultStat{}
}

func (*defaultMetric) NewSummary(key string, group string) StatSummary {
	return &defaultStat{}
}
