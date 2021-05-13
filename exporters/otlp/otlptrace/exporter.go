package otlptrace

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/internal/transform"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"sync"
)

var (
	errAlreadyStarted = errors.New("already started")
)

type Exporter struct {
	client Client

	mu      sync.RWMutex
	started bool

	startOnce sync.Once
	stopOnce  sync.Once
}

func (e *Exporter) ExportSpans(ctx context.Context, ss []tracesdk.ReadOnlySpan) error {
	protoSpans := transform.Spans(ss)
	if len(protoSpans) == 0 {
		return nil
	}

	return e.client.UploadTraces(ctx, protoSpans)
}

func (e *Exporter) Start(ctx context.Context) error {
	var err = errAlreadyStarted
	e.startOnce.Do(func() {
		e.mu.Lock()
		e.started = true
		e.mu.Unlock()
		err = e.client.Start(ctx)
	})

	return err
}

func (e *Exporter) Shutdown(ctx context.Context) error {

	e.mu.RLock()
	started := e.started
	e.mu.RUnlock()

	if !started {
		return nil
	}

	var err error

	e.stopOnce.Do(func() {
		err = e.client.Stop(ctx)
		e.mu.Lock()
		e.started = false
		e.mu.Unlock()
	})

	return err
}

var _ tracesdk.SpanExporter = (*Exporter)(nil)

// NewExporter constructs a new Exporter and starts it.
func NewExporter(ctx context.Context, client Client) (*Exporter, error) {
	exp := NewUnstartedExporter(client)
	if err := exp.Start(ctx); err != nil {
		return nil, err
	}
	return exp, nil
}

// NewUnstartedExporter constructs a new Exporter and does not start it.
func NewUnstartedExporter(client Client) *Exporter {
	e := &Exporter{
		client: client,
	}

	return e
}

// NewExportPipeline sets up a complete export pipeline
// with the recommended TracerProvider setup.
func NewExportPipeline(ctx context.Context, client Client) (*Exporter, *tracesdk.TracerProvider, error) {
	exp, err := NewExporter(ctx, client)
	if err != nil {
		return nil, nil, err
	}

	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
	)

	return exp, tracerProvider, nil
}

// InstallNewPipeline instantiates a NewExportPipeline with the
// recommended configuration and registers it globally.
func InstallNewPipeline(ctx context.Context, client Client) (*Exporter, *tracesdk.TracerProvider, error) {
	exp, tp, err := NewExportPipeline(ctx, client)
	if err != nil {
		return nil, nil, err
	}

	otel.SetTracerProvider(tp)
	return exp, tp, err
}
