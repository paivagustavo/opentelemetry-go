package otlpmetric

import (
	"context"
	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"
)

//TODO: better specify these interfaces.

// Client is an interface used by OTLP exporter. It's
// responsible for connecting to and disconnecting from the collector,
// and for transforming traces and metrics into wire format and
// transmitting them to the collector.
type Client interface {
	// Start should establish connection(s) to endpoint(s). It is
	// called just once by the exporter, so the implementation
	// does not need to worry about idempotence and locking.
	Start(ctx context.Context) error
	// Stop should close the connections. The function is called
	// only once by the exporter, so the implementation does not
	// need to worry about idempotence, but it may be called
	// concurrently with ExportMetrics or ExportTraces, so proper
	// locking is required. The function serves as a
	// synchronization point - after the function returns, the
	// process of closing connections is assumed to be finished.
	Stop(ctx context.Context) error
	// UploadMetrics should transform the passed metrics to the
	// wire format and send it to the collector. May be called
	// concurrently with ExportTraces, so the manager needs to
	// take this into account by doing proper locking.
	UploadMetrics(ctx context.Context, protoMetrics []*metricpb.ResourceMetrics) error
}
