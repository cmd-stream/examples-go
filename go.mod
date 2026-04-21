module github.com/cmd-stream/examples-go

go 1.24.1

require (
	github.com/brianvoe/gofakeit v3.18.0+incompatible
	github.com/cmd-stream/cmd-stream-go v0.6.1
	github.com/cmd-stream/codec-json-go v0.0.0-20260421202202-cd29f9c4adfc
	github.com/cmd-stream/codec-mus-stream-go v0.0.0-20260415035208-ba6c3f57289a
	github.com/cmd-stream/codec-protobuf-go v0.0.0-20260421210102-d8dfaf03511a
	github.com/mus-format/common-go v0.0.0-20260324174526-3d8f1741b5a2
	github.com/mus-format/mus-gen-go v0.5.1
	github.com/mus-format/mus-stream-go v0.10.1
	github.com/ymz-ncnk/assert v0.0.0-20260108210721-155bc9aa4282
	github.com/ymz-ncnk/circbrk-go v0.0.0-20260117185435-2da55ad9387e
	go.opentelemetry.io/otel v1.36.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.36.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.36.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutlog v0.12.2
	go.opentelemetry.io/otel/log v0.12.2
	go.opentelemetry.io/otel/sdk v1.36.0
	go.opentelemetry.io/otel/sdk/log v0.12.2
	go.opentelemetry.io/otel/sdk/metric v1.36.0
	go.opentelemetry.io/otel/trace v1.36.0
	google.golang.org/protobuf v1.36.10
)

require github.com/cmd-stream/codec-go v0.0.0-20260421172244-cee5b400964e // indirect

require (
	github.com/cenkalti/backoff/v5 v5.0.2 // indirect
	github.com/cmd-stream/otelcmd-stream-go v0.2.3
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3 // indirect
	github.com/ymz-ncnk/jointwork-go v0.0.0-20240428103805-1ee224bde88a // indirect
	github.com/ymz-ncnk/multierr-go v0.0.0-20230813140901-5e9302c2e02a // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.36.0 // indirect
	go.opentelemetry.io/otel/metric v1.36.0 // indirect
	go.opentelemetry.io/proto/otlp v1.6.0 // indirect
	golang.org/x/exp v0.0.0-20250408133849-7e4ce0ab07d0 // indirect
	golang.org/x/mod v0.33.0 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	golang.org/x/tools v0.41.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250519155744-55703ea1f237 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250519155744-55703ea1f237 // indirect
	google.golang.org/grpc v1.72.1 // indirect
)
