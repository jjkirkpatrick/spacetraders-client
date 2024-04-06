package metrics

import (
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
}

// NoOpMetricsReporter is a no-operation metrics reporter.
type NoOpMetricsReporter struct{}

// WritePoint for NoOpMetricsReporter does nothing.
func (n *NoOpMetricsReporter) WritePoint(m Metric) {
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
