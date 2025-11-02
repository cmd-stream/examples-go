module github.com/cmd-stream/examples-go

go 1.23.0

require (
	github.com/brianvoe/gofakeit v3.18.0+incompatible
	github.com/cmd-stream/cmd-stream-go v0.4.4
	github.com/cmd-stream/codec-json-go v0.0.0-20251102042300-5f13e0e191a7
	github.com/cmd-stream/codec-mus-stream-go v0.0.0-20251102044610-7e5e8a11b682
	github.com/cmd-stream/codec-protobuf-go v0.0.0-20251102042926-a83b3227fcc0
	github.com/cmd-stream/core-go v0.0.0-20251102020427-f23e62426486
	github.com/cmd-stream/delegate-go v0.0.0-20251102020741-164e6005aadf
	github.com/cmd-stream/handler-go v0.0.0-20251102020950-33189f2d8d28
	github.com/cmd-stream/transport-go v0.0.0-20251102021115-2f2d348f4122
	github.com/mus-format/common-go v0.0.0-20251026152644-9f5ac6728d8a
	github.com/mus-format/dts-stream-go v0.9.1
	github.com/mus-format/mus-stream-go v0.7.2
	github.com/mus-format/musgen-go v0.2.3
	github.com/ymz-ncnk/assert v0.0.0-20250528151733-c41b2fca7933
	github.com/ymz-ncnk/circbrk-go v0.0.0-20250912145433-3ecf61f801af
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

require github.com/cmd-stream/codec-generic-go v0.0.0-20251102041526-c9db158fec65 // indirect

require (
	github.com/cenkalti/backoff/v5 v5.0.2 // indirect
	github.com/cmd-stream/otelcmd-stream-go v0.1.1-0.20251102051623-25f8e8efd6e5
	github.com/cmd-stream/sender-go v0.0.0-20251102050957-6655c15cfddf
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
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	golang.org/x/tools v0.32.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250519155744-55703ea1f237 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250519155744-55703ea1f237 // indirect
	google.golang.org/grpc v1.72.1 // indirect
)
