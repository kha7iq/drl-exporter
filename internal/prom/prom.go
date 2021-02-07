package prom

import (
	"github.com/m47ik/drl-exporter/internal/collector"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var (
	maxRequestTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "dockerhub_limit_max_requests_total", Help: "Dockerhub rate limit maximum requests in given time"})
	maxRequestTotalTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "dockerhub_limit_max_requests_time", Help: "Dockerhub rate limit maximum requests total time seconds"})
	remainingRequestTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "dockerhub_limit_remaining_requests_total", Help: "Dockerhub rate limit remaining requests in given time"})
	remainingRequestTotalTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "dockerhub_limit_remaining_requests_time", Help: "Dockerhub rate limit remaining requests time seconds"})
)

func RegisterCollectors()  {
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
			maxRequestTotal.Set(collector.DockerMetrics["maxRequestTotal"])
			maxRequestTotalTime.Set(collector.DockerMetrics["maxRequestTotalTime"])
			remainingRequestTotal.Set(collector.DockerMetrics["remainingRequestTotal"])
			remainingRequestTotalTime.Set(collector.DockerMetrics["remainingRequestTotalTime"])
			time.Sleep(10 * time.Second)
		}
	}()
}