package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const namespace = "monitoring"

var (
	listenAddress = flag.String("web.listen-address", ":8080",
		"Address to listen on for telemetry")
	metricsPath = flag.String("web.telemetry-path", "/metrics",
		"Path under which to expose metrics")

	// Metrics
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"was the last network checker query successful",
		nil, nil,
	)
	networkStatus = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "status"),
		"current network status",
		[]string{"ip", "port", "status"}, nil,
	)
)

type Exporter struct{}
type NetworkChecker struct {
	status bool
	port   int
	ip     string
}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- networkStatus
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	results := tcpGather([]string{"8.8.8.8", "127.0.0.1"}, []int{53, 80, 443})

	ch <- prometheus.MustNewConstMetric(
		up, prometheus.GaugeValue, 1,
	)

	for i, fields := range results {
		fmt.Println(i, fields)
		ch <- prometheus.MustNewConstMetric(
			networkStatus, prometheus.GaugeValue, 1, fields.ip, strconv.Itoa(fields.port), strconv.FormatBool(fields.status),
		)
	}

	fmt.Println("Metrics collected")
}

func tcpGather(ips []string, ports []int) []NetworkChecker {
	results := make([]NetworkChecker, len(ips)*len(ports))
	k := 0
	for _, ip := range ips {
		for _, port := range ports {
			address := net.JoinHostPort(ip, strconv.Itoa(port))
			conn, err := net.DialTimeout("tcp", address, 1*time.Second)
			status := false
			if err == nil && conn != nil {
				status = true
				conn.Close()
			}
			results[k] = NetworkChecker{status: status, ip: ip, port: port}
			k++
		}
	}
	return results
}

func main() {
	log.SetOutput(os.Stdout)
	flag.Parse()

	exporter := NewExporter()
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Network port checker</title></head>
             <body>
             <h1>Network port checker</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
	fmt.Println("Exporter started")
}
