test: vendor
	go test ./...

build: gocd-metrics

vendor:
	dep ensure

gocd-metrics:
	go build ./cmd/..
