package metrics

import (
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type MetricsClient struct {
	influxClient influxdb2.Client
	org          string
	bucket       string
}

func NewMetricsClient(url, token, org, bucket string) *MetricsClient {
	return &MetricsClient{
		influxClient: influxdb2.NewClient(url, token),
		org:          org,
		bucket:       bucket,
	}
}

type MetricsReporter interface {
	WritePoint(Metric)
	Increment(namespace string, tags map[string]string, amount int)
	Decrement(namespace string, tags map[string]string, amount int)
	Observe(namespace string, tags map[string]string, value float64)
}

// NoOpMetricsReporter is a no-operation metrics reporter.
type NoOpMetricsReporter struct{}

// WritePoint for NoOpMetricsReporter does nothing.
func (n *NoOpMetricsReporter) WritePoint(m Metric) {
}

func (n *NoOpMetricsReporter) Increment(namespace string, tags map[string]string, amount int) {
}

func (n *NoOpMetricsReporter) Observe(namespace string, tags map[string]string, value float64) {
}

func (n *NoOpMetricsReporter) Decrement(namespace string, tags map[string]string, amount int) {
}

func (mc *MetricsClient) WritePoint(m Metric) {
	// Create a point and add to batch
	p := influxdb2.NewPoint(m.Namespace, m.Tags, m.Fields, m.Timestamp)
	// Get non-blocking write client
	writeAPI := mc.influxClient.WriteAPI(mc.org, mc.bucket)
	// Write point asynchronously
	writeAPI.WritePoint(p)
	// Ensure all writes are done
	writeAPI.Flush()
}

func (mc *MetricsClient) Increment(namespace string, tags map[string]string, amount int) {
	metric, _ := NewMetricBuilder().
		Namespace(namespace).
		Tags(tags).
		Field("count", amount).
		Timestamp(time.Now()).
		Build()
	mc.WritePoint(metric)
}

func (mc *MetricsClient) Decrement(namespace string, tags map[string]string, amount int) {
	metric, _ := NewMetricBuilder().
		Namespace(namespace).
		Tags(tags).
		Field("count", -amount).
		Timestamp(time.Now()).
		Build()
	mc.WritePoint(metric)
}

func (mc *MetricsClient) Observe(namespace string, tags map[string]string, value float64) {
	metric, _ := NewMetricBuilder().
		Namespace(namespace).
		Tags(tags).
		Field("value", value).
		Timestamp(time.Now()).
		Build()
	mc.WritePoint(metric)
}
