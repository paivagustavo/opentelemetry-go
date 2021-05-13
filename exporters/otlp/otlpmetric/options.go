package otlpmetric

import metricsdk "go.opentelemetry.io/otel/sdk/export/metric"

// ExporterOption are setting options passed to an Exporter on creation.
type ExporterOption func(*config)

type config struct {
	exportKindSelector metricsdk.ExportKindSelector
}

// WithMetricExportKindSelector defines the ExportKindSelector used
// for selecting AggregationTemporality (i.e., Cumulative vs. Delta
// aggregation). If not specified otherwise, exporter will use a
// cumulative export kind selector.
func WithMetricExportKindSelector(selector metricsdk.ExportKindSelector) ExporterOption {
	return func(cfg *config) {
		cfg.exportKindSelector = selector
	}
}
