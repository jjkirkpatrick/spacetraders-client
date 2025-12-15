package client

import "context"

// contextKey is a private type for context keys to avoid collisions
type contextKey string

const (
	// MetricLabelsKey is the context key for custom metric labels
	MetricLabelsKey contextKey = "st_metric_labels"
)

// WithMetricLabels adds custom labels to a context for metric labeling.
// Labels are merged with any existing labels in the context.
// This allows consumers to propagate arbitrary metadata (tree_name, action_name, etc.)
// through API calls for metric attribution.
func WithMetricLabels(ctx context.Context, labels map[string]string) context.Context {
	existing := GetMetricLabels(ctx)

	// Merge new labels with existing (new values override existing)
	merged := make(map[string]string, len(existing)+len(labels))
	for k, v := range existing {
		merged[k] = v
	}
	for k, v := range labels {
		merged[k] = v
	}

	return context.WithValue(ctx, MetricLabelsKey, merged)
}

// WithMetricLabel adds a single label to a context for metric labeling.
// This is a convenience function for adding one label at a time.
func WithMetricLabel(ctx context.Context, key, value string) context.Context {
	return WithMetricLabels(ctx, map[string]string{key: value})
}

// GetMetricLabels extracts metric labels from context.
// Returns an empty map if no labels are set.
func GetMetricLabels(ctx context.Context) map[string]string {
	if v := ctx.Value(MetricLabelsKey); v != nil {
		if labels, ok := v.(map[string]string); ok {
			return labels
		}
	}
	return make(map[string]string)
}
