package observe

import (
	"testing"

	testutil "github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdatePrometheusForPipelineCount(t *testing.T) {
	gocdMetrics := NewGocdMetrics()
	gocdMetrics.pipelines["pipeline1"] = &Pipeline{
		name:    "pipeline1",
		counter: 5,
	}
	gocdMetrics.pipelines["pipeline2"] = &Pipeline{
		name:    "pipeline2",
		counter: 10,
	}

	UpdatePrometheus(gocdMetrics)

	gauge1, err := pipelineCountMetric.GetMetricWithLabelValues("pipeline1")
	require.Nil(t, err)
	assert.Equal(t, testutil.ToFloat64(gauge1), 5.0)
	gauge2, err := pipelineCountMetric.GetMetricWithLabelValues("pipeline2")
	require.Nil(t, err)
	assert.Equal(t, testutil.ToFloat64(gauge2), 10.0)
}

func TestUpdatePrometheusForStageResult(t *testing.T) {
	gocdMetrics := NewGocdMetrics()
	gocdMetrics.pipelines["pipeline1"] = &Pipeline{
		name:    "pipeline1",
		counter: 5,
		stages: []*Stage{
			&Stage{name: "stage-1",
				result: "Passed",
			},
			&Stage{name: "stage-2",
				result: "Failed",
			},
		},
	}

	UpdatePrometheus(gocdMetrics)

	stage1PassedGauge, err := stageResultMetric.GetMetricWithLabelValues("pipeline1", "stage-1", "Passed")
	require.Nil(t, err)
	assert.Equal(t, testutil.ToFloat64(stage1PassedGauge), 1.0)
	stage1FailedGauge, err := stageResultMetric.GetMetricWithLabelValues("pipeline1", "stage-1", "Failed")
	require.Nil(t, err)
	assert.Equal(t, testutil.ToFloat64(stage1FailedGauge), 0.0)

	stage2PassedGauge, err := stageResultMetric.GetMetricWithLabelValues("pipeline1", "stage-2", "Passed")
	require.Nil(t, err)
	assert.Equal(t, testutil.ToFloat64(stage2PassedGauge), 0.0)
	stage2FailedGauge, err := stageResultMetric.GetMetricWithLabelValues("pipeline1", "stage-2", "Failed")
	require.Nil(t, err)
	assert.Equal(t, testutil.ToFloat64(stage2FailedGauge), 1.0)
}

func TestUpdatePrometheusForResettingOldStageResult(t *testing.T) {
	gocdMetrics := NewGocdMetrics()
	gocdMetrics.pipelines["pipeline1"] = &Pipeline{
		name:    "pipeline1",
		counter: 5,
		stages: []*Stage{
			&Stage{name: "stage-1",
				result: "Passed",
			},
		},
	}

	UpdatePrometheus(gocdMetrics)
	gocdMetrics.pipelines["pipeline1"].stages[0].result = "Failed"
	UpdatePrometheus(gocdMetrics)

	stage1PassedGauge, err := stageResultMetric.GetMetricWithLabelValues("pipeline1", "stage-1", "Passed")
	require.Nil(t, err)
	assert.Equal(t, testutil.ToFloat64(stage1PassedGauge), 0.0)
	stage1FailedGauge, err := stageResultMetric.GetMetricWithLabelValues("pipeline1", "stage-1", "Failed")
	require.Nil(t, err)
	assert.Equal(t, testutil.ToFloat64(stage1FailedGauge), 1.0)
}

func TestUpdatePrometheusForIgnoringUnknownResult(t *testing.T) {
	gocdMetrics := NewGocdMetrics()
	gocdMetrics.pipelines["pipeline1"] = &Pipeline{
		name:    "pipeline1",
		counter: 5,
		stages: []*Stage{
			&Stage{name: "stage-1",
				result: "Passed",
			},
		},
	}

	UpdatePrometheus(gocdMetrics)
	gocdMetrics.pipelines["pipeline1"].stages[0].result = "Building"
	UpdatePrometheus(gocdMetrics)

	stage1PassedGauge, err := stageResultMetric.GetMetricWithLabelValues("pipeline1", "stage-1", "Passed")
	require.Nil(t, err)
	assert.Equal(t, testutil.ToFloat64(stage1PassedGauge), 1.0)
	stage1UnknownGauge, err := stageResultMetric.GetMetricWithLabelValues("pipeline1", "stage-1", "Failed")
	require.Nil(t, err)
	assert.Equal(t, testutil.ToFloat64(stage1UnknownGauge), 0.0)
}
