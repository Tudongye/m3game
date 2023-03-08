package metric

import (
	"fmt"
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

func Counter(key string, v float64) StatCounter {
	t, ok := _statCounters.Load(key)
	if ok {
		if t == nil {
			panic(fmt.Sprintf("StatCounter Load nil %s", key))
		}
		return t.(StatCounter)
	}
	c := Get().NewCounter(key, "")
	r, _ := _statCounters.LoadOrStore(key, c)
	if r == nil {
		panic(fmt.Sprintf("StatCounter Load nil %s", key))
	}
	return r.(StatCounter)
}

func Gauge(key string, v float64) StatGauge {
	t, ok := _statGauges.Load(key)
	if ok {
		if t == nil {
			panic(fmt.Sprintf("StatGauge Load nil %s", key))
		}
		return t.(StatGauge)
	}
	c := Get().NewGauge(key, "")
	r, _ := _statGauges.LoadOrStore(key, c)
	if r == nil {
		panic(fmt.Sprintf("StatGauge Load nil %s", key))
	}
	return r.(StatGauge)
}

func Histogram(key string, v float64) StatHistogram {
	t, ok := _statHistograms.Load(key)
	if ok {
		if t == nil {
			panic(fmt.Sprintf("StatHistogram Load nil %s", key))
		}
		return t.(StatHistogram)
	}
	c := Get().NewHistogram(key, "")
	r, _ := _statHistograms.LoadOrStore(key, c)
	if r == nil {
		panic(fmt.Sprintf("StatHistogram Load nil %s", key))
	}
	return r.(StatHistogram)
}

func Summary(key string, v float64) StatSummary {
	t, ok := _statSummarys.Load(key)
	if ok {
		if t == nil {
			panic(fmt.Sprintf("StatHistogram Load nil %s", key))
		}
		return t.(StatSummary)
	}
	c := Get().NewSummary(key, "")
	r, _ := _statSummarys.LoadOrStore(key, c)
	if r == nil {
		panic(fmt.Sprintf("StatHistogram Load nil %s", key))
	}
	return r.(StatSummary)
}
