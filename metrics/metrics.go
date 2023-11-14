package metrics

import (
	metricsLib "github.com/tempmee/go-metrics-lib"
	"github.com/tempmee/go-metrics-lib/clients/prometheus"
)

var metricsInstance metricsLib.MetricsImpl

var prometheusInstance *prometheus.PrometheusClient

func NewMetricsInstance() metricsLib.MetricsImpl {
	if metricsInstance == nil {
		prometheusInstance = NewPrometheusInstance()
		initMetrics(prometheusInstance)
		metricsInstance = metricsLib.NewMetrics(prometheusInstance, 1)
	}
	return metricsInstance
}

func NewPrometheusInstance() *prometheus.PrometheusClient {
	if prometheusInstance == nil {
		prometheusInstance = prometheus.NewPrometheusClient()
		initMetrics(prometheusInstance)
	}
	return prometheusInstance
}

func initMetrics(prometheusInstance *prometheus.PrometheusClient) {
	prometheusInstance.CreateHistogramVec("", "graphql resolver millisecond", []string{"resolver", "service", "result"}, []float64{
		// create buckets 10000 split into 10 buckets
		100,
		200,
		300,
		400,
		500,
		600,
		700,
		800,
		900,
		1000,
	})
}
