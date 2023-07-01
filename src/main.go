package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/showwin/speedtest-go/speedtest"
)

func speedTest() {

	var speedtestClient = speedtest.New()

	serverList, _ := speedtestClient.FetchServers()
	targets, _ := serverList.FindServer([]int{})

	for _, s := range targets {
		s.PingTest(nil)
		s.DownloadTest()
		s.UploadTest()
		fmt.Printf("Latency: %s, Download: %f, Upload: %f\n", s.Latency, s.DLSpeed, s.ULSpeed)

		prometheus.MustRegister(speedTestResultLatency)
		speedTestResultLatency.Set(float64(s.Latency.Microseconds()) / 1000)
	}

}

var (
	// TODO: expose rest of the metrics (Download, Upload)
	speedTestResultLatency = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "speedtest_latency",
		Help: "Latency of the speedtest in ms",
	})
)

func main() {

	go speedTest()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
	fmt.Printf("Serving /metrics at :8080")
}
