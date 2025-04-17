PROTO_EXISTS := $(shell command -v protoc protoc-gen-go protoc-gen-go-grpc 2> /dev/null)

.PHONY: gen clean

gen: clean
ifndef PROTO_EXISTS
	@echo "protoc not found, installing..."
	brew install protobuf
	go install -ldflags="-w -s" google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install -ldflags="-w -s" google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
endif
	protoc --go_out=. --go-grpc_out=. proto/*.proto

clean:
	rm -f $(app)
	# rm -f proto_gen/*.pb.go
	rm -rf proto_gen