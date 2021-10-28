package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Scrin/prom-pinger/ping"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	sent     *prometheus.CounterVec
	received *prometheus.CounterVec
	latency  *prometheus.HistogramVec
)

func setup() {
	metricPrefix := "ping_"
	labels := []string{"target"}

	sent = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metricPrefix + "sent_count",
		Help: "Total number of sent ICMP Ping packets",
	}, labels)
	received = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metricPrefix + "received_count",
		Help: "Total number of received ICMP Ping packets",
	}, labels)
	latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    metricPrefix + "latency",
		Help:    "ICMP Ping latency to target",
		Buckets: prometheus.ExponentialBuckets(0.0001, 2, 20),
	}, labels)

	prometheus.MustRegister(sent)
	prometheus.MustRegister(received)
	prometheus.MustRegister(latency)
}

func runPinger(target string) {
	labels := prometheus.Labels{"target": target}
	for {
		pinger, err := ping.NewPinger(target)
		if err != nil { // This can happen ie. when a domain fails to resolve, unreachable dns etc
			log.Print(err)
			time.Sleep(10 * time.Second)
			continue
		}
		pinger.SetPrivileged(true)
		pinger.Count = 1
		pinger.Timeout = 55 * time.Second
		pinger.Run()
		stats := pinger.Statistics()

		sent.With(labels).Add(float64(stats.PacketsSent))
		received.With(labels).Add(float64(stats.PacketsRecv))
		if stats.PacketsRecv > 0 {
			latency.With(labels).Observe(float64(stats.AvgRtt) / float64(time.Second))
		}

		time.Sleep(time.Second)
	}
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("No targets specified. Usage: " + os.Args[0] + " <target> [<target>...]")
	}

	setup()

	for _, target := range os.Args[1:] {
		log.Print("Starting pinger for target: " + target)
		go runPinger(target)
		time.Sleep(100 * time.Millisecond)
	}

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
