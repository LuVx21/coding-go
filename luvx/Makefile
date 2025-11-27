SERVICE_NAME=luvx/go_app
SERVICE_VERSION=latest

gotool:
	go fmt ./
	go vet ./

build:
	CGO_ENABLED=0 go build -ldflags="-w -s" -v -o luvx main.go

docker_build:
	docker build -t $(SERVICE_NAME):$(SERVICE_VERSION) .

cgo_check:
	go list -f '{{if .CgoFiles}}{{.ImportPath}}{{end}}' all
	# go install github.com/visualfc/cgo-check@latest
	# cgo-check ./...

deploy:
	cp -r ../luvx_service_sdk .
	sed -i '' 's|=> \.\./|=> ./|g' go.mod
	docker compose up --build -d --scale app=1
	rm -fr luvx_service_sdk
	sed -i '' 's|=> ./|=> \.\./|g' go.mod
