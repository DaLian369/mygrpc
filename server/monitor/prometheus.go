package monitor

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Reqs = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "xxx",
		Name:      "yyy",
		Help:      "hhh",
	})
)

func InitPrometheus() (err error) {
	prometheus.Register(Reqs)
	http.Handle("/metrics", promhttp.Handler())
	return
}
