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

	stageResultMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gocd_stage_result",
			Help: "The current result (Passed, Failed) for each stage",
		},
		[]string{
			"pipeline",
			"stage",
			"result",
		},
	)
)

// RegisterPrometheus adds the prometheus handler to the mux router
// Note you must register every metric with prometheus for it show up
// when the /metrics route is hit.
func RegisterPrometheus(m *mux.Router) *mux.Router {
	prometheus.MustRegister(pipelineCountMetric)
	prometheus.MustRegister(stageResultMetric)

	m.Handle("/metrics", promhttp.Handler())
	return m
}

func UpdatePrometheus(metrics GocdMetrics) {
	for _, pipeline := range metrics.pipelines {
		pipelineCountMetric.WithLabelValues(pipeline.name).Set(float64(pipeline.counter))

		for _, stage := range pipeline.stages {
			if stage.result == "Passed" {
				stageResultMetric.WithLabelValues(pipeline.name, stage.name, "Passed").Set(1.0)
				stageResultMetric.WithLabelValues(pipeline.name, stage.name, "Failed").Set(0.0)
			} else if stage.result == "Failed" {
				stageResultMetric.WithLabelValues(pipeline.name, stage.name, "Failed").Set(1.0)
				stageResultMetric.WithLabelValues(pipeline.name, stage.name, "Passed").Set(0.0)
			}
		}
	}
}
