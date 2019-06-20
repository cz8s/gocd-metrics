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
}

func NewGocdMetrics() GocdMetrics {
	metrics := GocdMetrics{
		pipelines: map[string]*Pipeline{},
	}
	return metrics
}

func NewGocdClient(host string, username string, password string, skipTlsVerify bool) gocd.Client {
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
	}
	return nil
}

func updatePipelineCounter(metrics *GocdMetrics, pipeline gocd.PipelineInstance) {
	cachedPipeline, ok := metrics.pipelines[pipeline.Name]
	if !ok {
		cachedPipeline = &Pipeline{
			name:    pipeline.Name,
			counter: 0,
		}
		metrics.pipelines[pipeline.Name] = cachedPipeline
	}
	cachedPipeline.counter = pipeline.Counter
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
