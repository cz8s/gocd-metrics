package observe

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	pipelineCountMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gocd_pipeline_count",
			Help: "The total number of times a pipeline has ran",
		},
		[]string{
			"pipeline",
		},
	)
)

// RegisterPrometheus adds the prometheus handler to the mux router
// Note you must register every metric with prometheus for it show up
// when the /metrics route is hit.
func RegisterPrometheus(m *mux.Router) *mux.Router {
	prometheus.MustRegister(pipelineCountMetric)

	m.Handle("/metrics", promhttp.Handler())
	return m
}

func UpdatePrometheus(metrics GocdMetrics) {
	for _, pipeline := range metrics.pipelines {
		pipelineCountMetric.WithLabelValues(pipeline.name).Set(float64(pipeline.counter))
	}
}
