SERVICE_NAME=luvx/go_app
SERVICE_VERSION=latest

gotool:
	go fmt ./
	go vet ./

build:
	CGO_ENABLED=0 go build -ldflags="-w -s" -v -o luvx main.go

docker_build:
	docker build -t $(SERVICE_NAME):$(SERVICE_VERSION) .
