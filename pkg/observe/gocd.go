package observe

import "github.com/ashwanthkumar/go-gocd"

type GocdMetrics struct {
	pipelines map[string]Pipeline
}

type Pipeline struct {
	name    string
	counter int
}

func NewGocdMetrics() GocdMetrics {
	metrics := GocdMetrics{
		pipelines: map[string]Pipeline{},
	}
	return metrics
}

func (metrics *GocdMetrics) update(gocd gocd.Client) error {
	pipelineNames, err := getPipelineNames(gocd)
	if err != nil {
		return err
	}
	for _, pipelineName := range pipelineNames {
		history, err := gocd.GetPipelineHistoryPage(pipelineName, 0)
		if err != nil {
			return err
		}
		if len(history.Pipelines) < 1 {
			continue
		}
		lastRun := history.Pipelines[0]
		updatePipelineCounter(metrics, lastRun)
	}
	return nil
}

func updatePipelineCounter(metrics *GocdMetrics, pipeline gocd.PipelineInstance) {
	cachedPipeline, ok := metrics.pipelines[pipeline.Name]
	if !ok {
		cachedPipeline = Pipeline{
			name:    pipeline.Name,
			counter: 0,
		}
	}
	cachedPipeline.counter = pipeline.Counter
	metrics.pipelines[pipeline.Name] = cachedPipeline
}

func getPipelineNames(gocd gocd.Client) ([]string, error) {
	pipelineNames := make([]string, 0)
	pipelineGroups, err := gocd.GetPipelineGroups()
	if err != nil {
		return nil, err
	}
	for _, pipelineGroup := range pipelineGroups {
		for _, pipeline := range pipelineGroup.Pipelines {
			pipelineNames = append(pipelineNames, pipeline.Name)
		}
	}
	return pipelineNames, nil
}
