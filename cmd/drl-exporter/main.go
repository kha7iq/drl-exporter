package main

import (
	prom "github.com/m47ik/drl-exporter/internal/prom"
	"github.com/m47ik/drl-exporter/internal/vars"
	"github.com/nicholasjackson/env"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	log.Printf("Starting server on port %v\n", *vars.BindAddress)

	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
