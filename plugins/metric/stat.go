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

/*
监控指标类型参考Prometheus
*/
type StatCounter interface {
	Add(float64)
	Inc()
}

type StatGauge interface {
	Set(float64)
	Sub(float64)
	Dec()
	Add(float64)
	Inc()
}

type StatHistogram interface {
	Observe(float64)
}

type StatSummary interface {
	Observe(float64)
}

func Counter(key string) StatCounter {
	t, ok := _statCounters.Load(key)
	if ok {
		return t.(StatCounter)
	}
	c := Instance().NewCounter(key, "")
	r, _ := _statCounters.LoadOrStore(key, c)
	return r.(StatCounter)
}

func Gauge(key string) StatGauge {
	t, ok := _statGauges.Load(key)
	if ok {
		return t.(StatGauge)
	}
	c := Instance().NewGauge(key, "")
	r, _ := _statGauges.LoadOrStore(key, c)
	return r.(StatGauge)
}

func Histogram(key string) StatHistogram {
	t, ok := _statHistograms.Load(key)
	if ok {
		return t.(StatHistogram)
	}
	c := Instance().NewHistogram(key, "")
	r, _ := _statHistograms.LoadOrStore(key, c)
	return r.(StatHistogram)
}

func Summary(key string) StatSummary {
	t, ok := _statSummarys.Load(key)
	if ok {
		return t.(StatSummary)
	}
	c := Instance().NewSummary(key, "")
	r, _ := _statSummarys.LoadOrStore(key, c)
	return r.(StatSummary)
}
