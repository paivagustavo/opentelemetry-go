package otlpmetric

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel/exporters/otlp/internal/transform"
	"go.opentelemetry.io/otel/metric"
	metricsdk "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"sync"
)

var (
	errAlreadyStarted = errors.New("already started")
)

type Exporter struct {
	client             Client
	exportKindSelector metricsdk.ExportKindSelector

	mu      sync.RWMutex
	started bool

	startOnce sync.Once
	stopOnce  sync.Once
}

var _ metricsdk.Exporter = (*Exporter)(nil)

// NewExporter constructs a new Exporter and starts it.
func NewExporter(ctx context.Context, client Client, opts ...ExporterOption) (*Exporter, error) {
	exp := NewUnstartedExporter(client, opts...)
	if err := exp.Start(ctx); err != nil {
		return nil, err
	}
	return exp, nil
}

func NewUnstartedExporter(client Client, opts ...ExporterOption) *Exporter {
	cfg := config{
		// Note: the default ExportKindSelector is specified
		// as Cumulative:
		// https://github.com/open-telemetry/opentelemetry-specification/issues/731
		exportKindSelector: metricsdk.CumulativeExportKindSelector(),
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	e := &Exporter{
		client:             client,
		exportKindSelector: cfg.exportKindSelector,
	}

	return e
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

func (e *Exporter) Export(ctx context.Context, checkpointSet metricsdk.CheckpointSet) error {
	rms, err := transform.CheckpointSet(ctx, e, checkpointSet, 1)
	if err != nil {
		return err
	}
	if len(rms) == 0 {
		return nil
	}

	return e.client.UploadMetrics(ctx, rms)
}

func (e *Exporter) ExportKindFor(descriptor *metric.Descriptor, aggregatorKind aggregation.Kind) metricsdk.ExportKind {
	return e.exportKindSelector.ExportKindFor(descriptor, aggregatorKind)
}

// NewExportPipeline sets up a complete export pipeline
// with the recommended TracerProvider setup.
func NewExportPipeline(ctx context.Context, client Client, exporterOpts ...ExporterOption) (*Exporter, *basic.Controller, error) {
	exp, err := NewExporter(ctx, client, exporterOpts...)
	if err != nil {
		return nil, nil, err
	}

	cntr := basic.New(
		processor.New(
			simple.NewWithInexpensiveDistribution(),
			exp,
		),
	)

	return exp, cntr, nil
}

// InstallNewPipeline instantiates a NewExportPipeline with the
// recommended configuration and registers it globally.
func InstallNewPipeline(ctx context.Context, client Client, exporterOpts ...ExporterOption) (*Exporter, *basic.Controller, error) {
	exp, cntr, err := NewExportPipeline(ctx, client, exporterOpts...)
	if err != nil {
		return nil, nil, err
	}

	err = cntr.Start(ctx)
	if err != nil {
		return nil, nil, err
	}

	return exp, cntr, err
}
