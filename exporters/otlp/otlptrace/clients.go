package otlptrace

import (
	"context"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
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
	// UploadTraces should transform the passed traces to the wire
	// format and send it to the collector. May be called
	// concurrently.
	UploadTraces(ctx context.Context, protoSpans []*tracepb.ResourceSpans) error
}
