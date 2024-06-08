package metrics

import "time"

type Metric struct {
	Namespace string
	Tags      map[string]string
	Fields    map[string]interface{}
	Timestamp time.Time
}

type MetricBuilder struct {
	namespace string
	tags      map[string]string
	fields    map[string]interface{}
	timestamp time.Time
}

func NewMetricBuilder() *MetricBuilder {
	return &MetricBuilder{
		tags:   make(map[string]string),
		fields: make(map[string]interface{}),
	}
}

func (b *MetricBuilder) Namespace(namespace string) *MetricBuilder {
	b.namespace = namespace
	return b
}

func (b *MetricBuilder) Tag(key, value string) *MetricBuilder {
	b.tags[key] = value
	return b
}

func (b *MetricBuilder) Tags(newTags map[string]string) *MetricBuilder {
	for key, value := range newTags {
		b.tags[key] = value
	}
	return b
}

func (b *MetricBuilder) Field(key string, value interface{}) *MetricBuilder {
	b.fields[key] = value
	return b
}

func (b *MetricBuilder) Timestamp(timestamp time.Time) *MetricBuilder {
	b.timestamp = timestamp
	return b
}

func (b *MetricBuilder) Build() (Metric, error) {
	// Validate the constructed metric and return it
	// For simplicity, let's assume Metric is a struct you define to hold these values
	return Metric{
		Namespace: b.namespace,
		Tags:      b.tags,
		Fields:    b.fields,
		Timestamp: b.timestamp,
	}, nil
}
