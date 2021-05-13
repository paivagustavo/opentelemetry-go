// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otlpmetrictest

import (
	"sort"

	collectormetricpb "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
)

// Collector is an interface that mock collectors should implements,
// so they can be used for the end-to-end testing.
type Collector interface {
	Stop() error
	GetMetrics() []*metricpb.Metric
}

// MetricsStorage stores the metrics. Mock collectors could use it to
// store metrics they have received.
type MetricsStorage struct {
	metrics []*metricpb.Metric
}

// NewMetricsStorage creates a new metrics storage.
func NewMetricsStorage() MetricsStorage {
	return MetricsStorage{}
}

// AddMetrics adds metrics to the metrics storage.
func (s *MetricsStorage) AddMetrics(request *collectormetricpb.ExportMetricsServiceRequest) {
	for _, rm := range request.GetResourceMetrics() {
		// TODO (rghetia) handle multiple resource and library info.
		if len(rm.InstrumentationLibraryMetrics) > 0 {
			s.metrics = append(s.metrics, rm.InstrumentationLibraryMetrics[0].Metrics...)
		}
	}
}

// GetMetrics returns the stored metrics.
func (s *MetricsStorage) GetMetrics() []*metricpb.Metric {
	// copy in order to not change.
	m := make([]*metricpb.Metric, 0, len(s.metrics))
	return append(m, s.metrics...)
}

func resourceString(res *resourcepb.Resource) string {
	sAttrs := sortedAttributes(res.GetAttributes())
	rstr := ""
	for _, attr := range sAttrs {
		rstr = rstr + attr.String()
	}
	return rstr
}

func sortedAttributes(attrs []*commonpb.KeyValue) []*commonpb.KeyValue {
	sort.Slice(attrs[:], func(i, j int) bool {
		return attrs[i].Key < attrs[j].Key
	})
	return attrs
}
