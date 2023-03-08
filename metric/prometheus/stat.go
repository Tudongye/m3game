package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

type pstatCounter struct {
	c prometheus.Counter
}

func (p *pstatCounter) Add(f float64) {
	p.c.Add(f)
}

func (p *pstatCounter) Inc() {
	p.c.Inc()
}

type pstatGauge struct {
	c prometheus.Gauge
}

func (p *pstatGauge) Set(f float64) {
	p.c.Set(f)
}

func (p *pstatGauge) Inc() {
	p.c.Inc()
}

func (p *pstatGauge) Add(f float64) {
	p.c.Add(f)
}

func (p *pstatGauge) Dec() {
	p.c.Dec()
}

func (p *pstatGauge) Sub(f float64) {
	p.c.Sub(f)
}

type pstatHistogram struct {
	c prometheus.Histogram
}

func (p *pstatHistogram) Observe(f float64) {
	p.c.Observe(f)
}

type pstatSummary struct {
	c prometheus.Summary
}

func (p *pstatSummary) Observe(f float64) {
	p.c.Observe(f)
}
