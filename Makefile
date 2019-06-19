build: vendor gocd-metrics

test: vendor
	go test ./...

vendor:
	dep ensure

gocd-metrics:
	go build ./cmd/...
