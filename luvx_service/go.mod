module luvx_service

go 1.24.2

replace luvx_service_sdk => ../luvx_service_sdk

require (
	github.com/dgraph-io/badger/v4 v4.7.0
	github.com/klauspost/compress v1.18.0
	github.com/luvx21/coding-go/coding-common v0.0.0-20250412082033-7367dd329551
	github.com/luvx21/coding-go/infra/infra_kv v0.0.0-20250412082033-7367dd329551
	go.etcd.io/etcd/client/v3 v3.5.21
	google.golang.org/grpc v1.72.0
	google.golang.org/protobuf v1.36.6
	luvx_service_sdk v0.0.0-00010101000000-000000000000
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/dgraph-io/ristretto/v2 v2.2.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/flatbuffers v25.2.10+incompatible // indirect
	go.etcd.io/etcd/api/v3 v3.5.21 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.21 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/exp v0.0.0-20250408133849-7e4ce0ab07d0 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250422160041-2d3770c4ea7f // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250422160041-2d3770c4ea7f // indirect
)
