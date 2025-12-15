package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithMetricLabels(t *testing.T) {
	ctx := context.Background()
	ctx = WithMetricLabels(ctx, map[string]string{
		"tree_name":   "mining",
		"action_name": "extract",
	})

	labels := GetMetricLabels(ctx)
	assert.Equal(t, "mining", labels["tree_name"])
	assert.Equal(t, "extract", labels["action_name"])
}

func TestWithMetricLabels_Merge(t *testing.T) {
	ctx := context.Background()
	ctx = WithMetricLabels(ctx, map[string]string{
		"tree_name": "mining",
	})
	ctx = WithMetricLabels(ctx, map[string]string{
		"action_name": "extract",
	})

	labels := GetMetricLabels(ctx)
	assert.Equal(t, "mining", labels["tree_name"])
	assert.Equal(t, "extract", labels["action_name"])
}

func TestWithMetricLabels_Override(t *testing.T) {
	ctx := context.Background()
	ctx = WithMetricLabels(ctx, map[string]string{
		"action_name": "orbit",
	})
	ctx = WithMetricLabels(ctx, map[string]string{
		"action_name": "extract",
	})

	labels := GetMetricLabels(ctx)
	assert.Equal(t, "extract", labels["action_name"])
}

func TestGetMetricLabels_EmptyContext(t *testing.T) {
	ctx := context.Background()
	labels := GetMetricLabels(ctx)
	assert.NotNil(t, labels)
	assert.Empty(t, labels)
}

func TestWithMetricLabel_Single(t *testing.T) {
	ctx := context.Background()
	ctx = WithMetricLabel(ctx, "ship_role", "hauler")

	labels := GetMetricLabels(ctx)
	assert.Equal(t, "hauler", labels["ship_role"])
}
