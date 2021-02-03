package main

import (
	"github.com/m47ik/drl-exporter/internal/collector"
	"github.com/m47ik/drl-exporter/internal/vars"
	"github.com/nicholasjackson/env"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

func main() {
	if err := env.Parse(); err != nil {
		log.Println(err)
	}
	mux := http.NewServeMux()
	s := &http.Server{
		Addr:         ":" + *vars.BindAddress,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	prometheus.Unregister(prometheus.NewGoCollector())
	prometheus.MustRegister(maxRequestTotal)
	prometheus.MustRegister(maxRequestTotalTime)
	prometheus.MustRegister(remainingRequestTotal)
	prometheus.MustRegister(remainingRequestTotalTime)

	collectMetrics()

	mux.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting server on port %v\n", *vars.BindAddress)

	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func collectMetrics() {
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
