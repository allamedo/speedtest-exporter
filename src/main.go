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

		prometheus.MustRegister(speedTestResultDLSpeed)
		speedTestResultDLSpeed.Set(s.DLSpeed)

		prometheus.MustRegister(speedTestResultULSpeed)
		speedTestResultULSpeed.Set(s.ULSpeed)
	}

}

var (
	// TODO: expose rest of the metrics (Download, Upload)
	speedTestResultLatency = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "speedtest_latency",
		Help: "Latency of the speedtest in ms",
	})
	speedTestResultDLSpeed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "speedtest_download_speed",
		Help: "Download speed of the speedtest in Mbps",
	})
	speedTestResultULSpeed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "speedtest_upnload_speed",
		Help: "Upnload speed of the speedtest in Mbps",
	})
)

func main() {

	go speedTest()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
	fmt.Printf("Serving /metrics at :2112")
}
