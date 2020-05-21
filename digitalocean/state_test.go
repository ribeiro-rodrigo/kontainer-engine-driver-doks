package digitalocean

import (
	"testing"

	"github.com/rancher/kontainer-engine/types"
	"github.com/stretchr/testify/assert"
)

func TestGetLabelsFromStringSlice(t *testing.T) {
	labelsStringSlice := types.StringSlice{
		Value: []string{
			"key1=value1",
			"key2=value2",
			"key3=value3",
		},
	}

	labels := getLabelsFromStringSlice(&labelsStringSlice)

	expectedLabels := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	assert.Equal(t, expectedLabels, labels, "Labels equals expectedLabels")
}

func TestGetLabelsFromStringSlicePassNilSlice(t *testing.T) {
	labels := getLabelsFromStringSlice(nil)

	expectedEmptyLabels := map[string]string{}

	assert.Equal(t, expectedEmptyLabels, labels, "Slice nil expected empty labels")
}
