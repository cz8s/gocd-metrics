package observe

import (
	"crypto/tls"
	"time"

	"github.com/ashwanthkumar/go-gocd"
	"github.com/parnurzeal/gorequest"
)

type GocdMetrics struct {
	pipelines map[string]*Pipeline
}

type Pipeline struct {
	name    string
	counter int
	stages  []*Stage
}

type Stage struct {
	name   string
	result string
}

func NewGocdMetrics() GocdMetrics {
	metrics := GocdMetrics{
		pipelines: map[string]*Pipeline{},
	}
	return metrics
}

func NewGocdClient(host string, username string, password string, skipTlsVerify bool) gocd.Client {
	/* #nosec G402 */
	return &gocd.DefaultClient{
		Host:    host,
		Request: gorequest.New().Timeout(60*time.Second).SetBasicAuth(username, password).TLSClientConfig(&tls.Config{InsecureSkipVerify: skipTlsVerify}),
	}
}

func UpdateGocdMetrics(metrics *GocdMetrics, gocd gocd.Client) error {
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
		updateStageStatus(metrics, lastRun)
	}
	return nil
}

func updatePipelineCounter(metrics *GocdMetrics, pipeline gocd.PipelineInstance) {
	cachedPipeline := ensurePipelineCache(metrics, pipeline.Name)
	cachedPipeline.counter = pipeline.Counter
}

func updateStageStatus(metrics *GocdMetrics, pipeline gocd.PipelineInstance) {
	cachedPipeline := ensurePipelineCache(metrics, pipeline.Name)
	cachedStages := make([]*Stage, len(pipeline.Stages))
	for i, stage := range pipeline.Stages {
		cachedStages[i] = &Stage{
			name:   stage.Name,
			result: stage.Result,
		}
	}
	cachedPipeline.stages = cachedStages
}

func ensurePipelineCache(metrics *GocdMetrics, pipelineName string) *Pipeline {
	cachedPipeline, ok := metrics.pipelines[pipelineName]
	if !ok {
		cachedPipeline = &Pipeline{
			name:    pipelineName,
			counter: 0,
			stages:  make([]*Stage, 0),
		}
		metrics.pipelines[pipelineName] = cachedPipeline
	}
	return cachedPipeline
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
