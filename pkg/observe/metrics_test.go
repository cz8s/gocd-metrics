package observe

import (
	"testing"

	testutil "github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestUpdatePrometheus(t *testing.T) {
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

	counter1, err := pipelineCountMetric.GetMetricWithLabelValues("pipeline1")
	assert.Nil(t, err)
	assert.Equal(t, testutil.ToFloat64(counter1), 5.0)
	counter2, err := pipelineCountMetric.GetMetricWithLabelValues("pipeline2")
	assert.Nil(t, err)
	assert.Equal(t, testutil.ToFloat64(counter2), 10.0)
}
