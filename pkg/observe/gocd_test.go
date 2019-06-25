package observe

import (
	"testing"

	"github.com/ashwanthkumar/go-gocd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdatePipelineCounterFirstRun(t *testing.T) {
	metrics := NewGocdMetrics()
	pipelineInstance := gocd.PipelineInstance{
		Name:    "test-pipeline",
		Counter: 12,
	}
	updatePipelineCounter(&metrics, pipelineInstance)

	pipeline, ok := metrics.pipelines["test-pipeline"]
	require.True(t, ok)
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
	require.True(t, ok)
	counter := pipeline.counter
	assert.Equal(t, counter, 9)
}

func TestUpdateStageStatusFirstRun(t *testing.T) {
	metrics := NewGocdMetrics()
	stages := []gocd.StageRun{
		gocd.StageRun{
			Name:   "test-stage",
			Result: "Passed",
		},
	}
	pipelineInstance := gocd.PipelineInstance{
		Name:    "test-pipeline",
		Counter: 12,
		Stages:  stages,
	}

	updateStageStatus(&metrics, pipelineInstance)

	pipeline, ok := metrics.pipelines["test-pipeline"]
	require.True(t, ok)
	assert.Equal(t, len(pipeline.stages), 1)
	assert.Equal(t, pipeline.stages[0].name, "test-stage")
	assert.Equal(t, pipeline.stages[0].result, "Passed")
}

func TestUpdateStageStatusRemovesOldStages(t *testing.T) {
	metrics := NewGocdMetrics()
	stages := []gocd.StageRun{
		gocd.StageRun{
			Name:   "test-stage",
			Result: "Passed",
		},
	}
	pipelineInstance := gocd.PipelineInstance{
		Name:    "test-pipeline",
		Counter: 12,
		Stages:  stages,
	}

	updateStageStatus(&metrics, pipelineInstance)
	pipelineInstance.Stages = make([]gocd.StageRun, 0)
	updateStageStatus(&metrics, pipelineInstance)

	pipeline, ok := metrics.pipelines["test-pipeline"]
	require.True(t, ok)
	assert.Equal(t, len(pipeline.stages), 0)
}
