build: vendor gocd-metrics

test: vendor
	go test ./...

vendor:
	dep ensure

gocd-metrics: vendor
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/...

docker: gocd-metrics
	docker build -t gocd-metrics .

docker-run: docker
	docker run -d -p 9090:9090 gocd-metrics

clean:
	rm -r vendor
	rm gocd-metrics

lint:
	golint
	gosec ./...
