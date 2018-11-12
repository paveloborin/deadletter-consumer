GOOS?=linux
APP?=consumer
FILES = 

.PHONY: *

#RUN LOCAL
run:
	export `cat .env` && go run cmd/consumer/*.go --pretty-logging

#BUILD
build_linux:
	CGO_ENABLED=0 GOOS=${GOOS} go build -a -installsuffix cgo \
		-o ./bin/${APP} ./cmd/consumer/*.go

#DOCKER
build_docker: build_linux
	docker build -t ${APP}:latest -f Dockerfile .

run_docker: build_docker
	docker stop consumer || true
	docker stop rabbit || true
	docker-compose up


#LINTERS
fmt:
	go fmt ./...

errcheck:
	errcheck ./...

lint:
	go fmt $$(go list ./... | grep -v ./vendor/)
	goimports -d -w $$(find . -type f -name '*.go' -not -path './vendor/*')
	golangci-lint run --skip-dirs tmp