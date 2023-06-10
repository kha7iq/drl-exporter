package prom

import (
	"time"

	"github.com/kha7iq/drl-exporter/internal/collector"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	maxRequestTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "dockerhub_limit_max_requests_total", Help: "Dockerhub rate limit maximum requests in given time"},
		[]string{
			"reqsource",
		},
	)
	maxRequestTotalTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "dockerhub_limit_max_requests_time", Help: "Dockerhub rate limit maximum requests total time seconds"},
		[]string{
			"reqsource",
		},
	)
	remainingRequestTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "dockerhub_limit_remaining_requests_total", Help: "Dockerhub rate limit remaining requests in given time"},
		[]string{
			"reqsource",
		},
	)
	remainingRequestTotalTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "dockerhub_limit_remaining_requests_time", Help: "Dockerhub rate limit remaining requests time seconds"},
		[]string{
			"reqsource",
		},
	)
)

func RegisterCollectors() {
	prometheus.Unregister(prometheus.NewGoCollector())
	prometheus.MustRegister(maxRequestTotal)
	prometheus.MustRegister(maxRequestTotalTime)
	prometheus.MustRegister(remainingRequestTotal)
	prometheus.MustRegister(remainingRequestTotalTime)
}

func CollectMetrics() {
	go func() {
		for {
			collector.GetMetrics()
			maxRequestTotal.WithLabelValues(collector.DockerLabels["reqsource"]).Set(collector.DockerMetrics["maxRequestTotal"])
			maxRequestTotalTime.WithLabelValues(collector.DockerLabels["reqsource"]).Set(collector.DockerMetrics["maxRequestTotalTime"])
			remainingRequestTotal.WithLabelValues(collector.DockerLabels["reqsource"]).Set(collector.DockerMetrics["remainingRequestTotal"])
			remainingRequestTotalTime.WithLabelValues(collector.DockerLabels["reqsource"]).Set(collector.DockerMetrics["remainingRequestTotalTime"])
			time.Sleep(10 * time.Second)
		}
	}()
}
