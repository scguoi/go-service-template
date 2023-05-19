package impl

import "github.com/prometheus/client_golang/prometheus"

var ResponseTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "request_duration_seconds",
	Help:    "Request duration in seconds.",
	Buckets: prometheus.ExponentialBucketsRange(1, 60000, 15), // 桶的配置
}, []string{"method"})

func init() {
	prometheus.MustRegister(ResponseTime)
}
