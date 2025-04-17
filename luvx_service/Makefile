app := luvx_service

.PHONY: build run clean

build:
	go build -o $(app)

run:
	go run main.go

clean:
	rm -f $(app)

deploy:
	cp -r ../luvx_service_sdk . && cd luvx_service_sdk && make clean gen && cd ..
	cp ../luvx/Dockerfile .
	sed -i '' 's|=> \.\./|=> ./|g' go.mod
	docker compose -p luvx up --build -d
	sed -i '' 's|=> ./|=> \.\./|g' go.mod
	rm -fr luvx_service_sdk
	rm -fr Dockerfile
