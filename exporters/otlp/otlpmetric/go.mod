module go.opentelemetry.io/otel/exporters/otlp/otlpmetric

go 1.14

replace (
	go.opentelemetry.io/otel => ../../../

	go.opentelemetry.io/otel/sdk => ./../../../sdk
)

replace go.opentelemetry.io/otel/exporters/otlp => ../

replace go.opentelemetry.io/otel/metric => ./../../../metric

replace go.opentelemetry.io/otel/oteltest => ./../../../oteltest

replace go.opentelemetry.io/otel/trace => ./../../../trace

replace go.opentelemetry.io/otel/sdk/export/metric => ../../../sdk/export/metric

replace go.opentelemetry.io/otel/sdk/metric => ../../../sdk/metric

require (
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/otel v0.20.0
	go.opentelemetry.io/otel/exporters/otlp v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel/metric v0.20.0
	go.opentelemetry.io/otel/sdk v0.20.0
	go.opentelemetry.io/otel/sdk/export/metric v0.20.0
	go.opentelemetry.io/otel/sdk/metric v0.20.0
	go.opentelemetry.io/proto/otlp v0.7.0

)
