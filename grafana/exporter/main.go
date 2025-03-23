package main

import (
    "log"
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    sampleGauge = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "sample_gauge",
        Help: "A sample gauge metric",
    })
    sampleCounter = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "sample_counter_total",
        Help: "A sample counter metric",
    })
)

func init() {
    prometheus.MustRegister(sampleGauge)
    prometheus.MustRegister(sampleCounter)
}

func generateMetrics() {
    for {
        sampleGauge.Set(float64(time.Now().Unix() % 100)) // Random value 0-99
        sampleCounter.Inc()
        time.Sleep(5 * time.Second)
    }
}

func main() {
    go generateMetrics()
    http.Handle("/metrics", promhttp.Handler())
    log.Println("Exporter running on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}