package prom

import (
	"fmt"
	"log"
	"time"

	"github.com/kha7iq/drl-exporter/internal/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type Metrics struct {
	MaxRequestTotal           *prometheus.GaugeVec
	MaxRequestTotalTime       *prometheus.GaugeVec
	RemainingRequestTotal     *prometheus.GaugeVec
	RemainingRequestTotalTime *prometheus.GaugeVec
}

var metrics *Metrics

func NewMetrics() *Metrics {
	return &Metrics{
		MaxRequestTotal: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "dockerhub_limit_max_requests_total",
			Help: "Dockerhub rate limit maximum requests in given time"},
			[]string{"reqsource"},
		),
		MaxRequestTotalTime: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "dockerhub_limit_max_requests_time",
			Help: "Dockerhub rate limit maximum requests total time seconds"},
			[]string{"reqsource"},
		),
		RemainingRequestTotal: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "dockerhub_limit_remaining_requests_total",
			Help: "Dockerhub rate limit remaining requests in given time"},
			[]string{"reqsource"},
		),
		RemainingRequestTotalTime: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "dockerhub_limit_remaining_requests_time",
			Help: "Dockerhub rate limit remaining requests time seconds"},
			[]string{"reqsource"},
		),
	}
}

func (m *Metrics) Register() {
	prometheus.Unregister(collectors.NewGoCollector())
	prometheus.MustRegister(m.MaxRequestTotal)
	prometheus.MustRegister(m.MaxRequestTotalTime)
	prometheus.MustRegister(m.RemainingRequestTotal)
	prometheus.MustRegister(m.RemainingRequestTotalTime)
}

func (m *Metrics) CollectMetrics(interval time.Duration) {
	go func() {
		for {
			if err := m.collectAndSetMetrics(); err != nil {
				log.Printf("Error collecting metrics: %v", err)
			}
			time.Sleep(interval)
		}
	}()
}

func (m *Metrics) collectAndSetMetrics() error {
	collector.GetMetrics()
	if source, ok := collector.DockerLabels["reqsource"]; ok {
		m.setGaugeMetric(m.MaxRequestTotal, source, collector.DockerMetrics["maxRequestTotal"])
		m.setGaugeMetric(m.MaxRequestTotalTime, source, collector.DockerMetrics["maxRequestTotalTime"])
		m.setGaugeMetric(m.RemainingRequestTotal, source, collector.DockerMetrics["remainingRequestTotal"])
		m.setGaugeMetric(m.RemainingRequestTotalTime, source, collector.DockerMetrics["remainingRequestTotalTime"])
		return nil
	}
	return fmt.Errorf("no reqsource label found")
}

func (m *Metrics) setGaugeMetric(gauge *prometheus.GaugeVec, label string, value float64) {
	gauge.WithLabelValues(label).Set(value)
}

func InitMetrics() {
	metrics = NewMetrics()
	metrics.Register()
}

func StartMetricsCollection(interval time.Duration) {
	if metrics == nil {
		InitMetrics()
	}
	metrics.CollectMetrics(interval)
}
