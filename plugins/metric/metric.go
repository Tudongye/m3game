package metric

import (
	"sync"
)

var (
	_statCounters   sync.Map
	_statGauges     sync.Map
	_statHistograms sync.Map
	_statSummarys   sync.Map
)

var (
	_metric        Metric
	_defaultmetric = &defaultMetric{}
)

func Set(m Metric) {
	if _metric != nil {
		panic("Metric only one")
	}
	_metric = m
}

func Get() Metric {
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

type defaultMetric struct {
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
