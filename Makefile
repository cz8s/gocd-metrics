build: gocd-metrics

test: 
	go test ./...

gocd-metrics: 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/...

docker: gocd-metrics
	docker build -t gocd-metrics .

docker-push: 
	./docker-push.sh

docker-run: docker
	docker run -d -p 9090:9090 gocd-metrics

clean:
	rm gocd-metrics

lint:
	golint
	gosec -exclude=G104 ./...

integration-test: docker
	cd integration-test ; ./test-integration.sh

