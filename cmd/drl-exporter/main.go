package main

import (
	"log"
	"net/http"
	"time"

	prom "github.com/kha7iq/drl-exporter/internal/prom"
	"github.com/kha7iq/drl-exporter/internal/vars"
	"github.com/nicholasjackson/env"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	version   string
	buildDate string
	commitSha string
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

	prom.RegisterCollectors()
	prom.CollectMetrics()

	mux.Handle("/metrics", promhttp.Handler())
	log.Printf("DRL exporter Version %s, Commit %s, BuildDate %s\n", version, commitSha, buildDate)
	log.Printf("Starting server on port %v\n", *vars.BindAddress)

	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
