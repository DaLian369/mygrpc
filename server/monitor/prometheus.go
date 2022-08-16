package monitor

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Reqs = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "rpc_server",
		Name:      "reqs",
		Help:      "reqs",
		Buckets:   []float64{3, 8, 15, 30, 50, 100, 300, 600, 1000},
	}, []string{"uri"})
)

func InitPrometheus() (err error) {
	prometheus.MustRegister(Reqs)
	http.Handle("/metrics", promhttp.Handler())
	return
}
