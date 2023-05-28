package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	gatewayStatusDesc = prometheus.NewDesc(
		"opnsense_gateway_status",
		"OPNsense gateway status",
		[]string{"gateway", "address", "status_translated"},
		nil,
	)
	lossDesc = prometheus.NewDesc(
		"opnsense_gateway_loss_pct",
		"OPNsense gateway packet loss percentage",
		[]string{"gateway", "address"},
		nil,
	)
	stddevDesc = prometheus.NewDesc(
		"opnsense_gateway_stddev_ms",
		"OPNsense gateway standard deviation in milliseconds",
		[]string{"gateway", "stddev"},
		nil,
	)
	delayDesc = prometheus.NewDesc(
		"opnsense_gateway_delay_ms",
		"OPNsense gateway delay in milliseconds",
		[]string{"gateway", "address"},
		nil,
	)
)

func (e *opnSenseExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- gatewayStatusDesc
	ch <- lossDesc
	ch <- delayDesc
}

func (e *opnSenseExporter) Collect(ch chan<- prometheus.Metric) {
	gatewayStatus, err := e.getGatewayStatus()
	if err != nil {
		log.Println("Error getting gateway status:", err)
		return
	}

	for _, gateway := range gatewayStatus.Items {
		status := 0.0
		if gateway.StatusTranslated == "Online" {
			status = 1.0
		}

		labels := []string{gateway.Name, gateway.Address, gateway.StatusTranslated}
		ch <- prometheus.MustNewConstMetric(
			gatewayStatusDesc,
			prometheus.GaugeValue,
			status,
			labels...,
		)

		lossStr := strings.TrimSuffix(gateway.Loss, " %")
		loss, err := strconv.ParseFloat(lossStr, 64)
		if err != nil {
			log.Println("Error parsing loss:", err)
			continue
		}
		gateway.LossValue = loss

		delayStr := strings.TrimSuffix(gateway.Delay, " ms")
		delay, err := strconv.ParseFloat(delayStr, 64)
		if err != nil {
			log.Println("Error parsing delay:", err)
			continue
		}
		gateway.DelayValue = delay

		ch <- prometheus.MustNewConstMetric(
			lossDesc,
			prometheus.GaugeValue,
			gateway.LossValue,
			gateway.Name, gateway.Address,
		)

		ch <- prometheus.MustNewConstMetric(
			delayDesc,
			prometheus.GaugeValue,
			gateway.DelayValue,
			gateway.Name, gateway.Address,
		)

		ch <- prometheus.MustNewConstMetric(
			stddevDesc,
			prometheus.GaugeValue,
			gateway.DelayValue,
			gateway.Name, gateway.Address,
		)
	}
}

func start_prometheus(apiURL, apiKey, apiSecret string) {
	exporter := &opnSenseExporter{
		apiURL:    apiURL,
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}

	prometheus.MustRegister(exporter)

	// TODO: run asynchronously
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
