package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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

	rqs := *vars.RequestInterval
	interval, err := strconv.Atoi(rqs)
	if err != nil {
		log.Printf("Invalid interval value '%s', defaulting to 10 seconds: %v\n", rqs, err)
		interval = 10
	}
	// Convert interval to a time.Duration
	duration := time.Duration(interval) * time.Second

	// Initialize and start metrics collection with a configurable interval
	prom.InitMetrics()
	prom.StartMetricsCollection(duration)

	mux.Handle("/metrics", promhttp.Handler())
	appInfo(version, buildDate, commitSha)

	err = s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func appInfo(version, buildDate, commit string) {
	cyan := "\x1b[36m"
	reset := "\x1b[0m"

	banner := `
	𝐃 𝐑 𝐋   𝐄 𝐗 𝐏 𝐎 𝐑 𝐓 𝐄 𝐑																				   
	`
	fmt.Println(cyan + banner + reset)

	fmt.Printf("Version: %v\n\nBuild Date: %v\n\nCommit: %v\n\n", cyan+version+reset, cyan+buildDate+reset, cyan+commit+reset)

	fmt.Printf("Starting server on Port: "+cyan+"%v\n\n"+reset, *vars.BindAddress)

}
