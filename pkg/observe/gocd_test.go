package observe

import (
	"testing"

	"github.com/ashwanthkumar/go-gocd"
	"github.com/stretchr/testify/assert"
)

func TestUpdatePipelineCounterFirstRun(t *testing.T) {
	metrics := NewGocdMetrics()
	pipelineInstance := gocd.PipelineInstance{
		Name:    "test-pipeline",
		Counter: 12,
	}
	updatePipelineCounter(&metrics, pipelineInstance)

	pipeline, ok := metrics.pipelines["test-pipeline"]
	assert.True(t, ok)
	counter := pipeline.counter
	assert.Equal(t, counter, 12)
}

func TestUpdatePipelineCounterLaterRun(t *testing.T) {
	metrics := NewGocdMetrics()
	firstPipelineInstance := gocd.PipelineInstance{
		Name:    "test-pipeline",
		Counter: 8,
	}
	secondPipelineInstance := gocd.PipelineInstance{
		Name:    "test-pipeline",
		Counter: 9,
	}
	updatePipelineCounter(&metrics, firstPipelineInstance)
	updatePipelineCounter(&metrics, secondPipelineInstance)

	pipeline, ok := metrics.pipelines["test-pipeline"]
	assert.True(t, ok)
	counter := pipeline.counter
	assert.Equal(t, counter, 9)
}
