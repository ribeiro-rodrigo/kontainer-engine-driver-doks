package state

import (
	"testing"

	"github.com/rancher/kontainer-engine/types"
	"github.com/stretchr/testify/assert"
)

var stateBuilder = NewBuilder()

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

func TestGetStateFromOptsKeysSnakeCase(t *testing.T) {

	const token = "tkalal1234761"
	const displayName = "cluster-test"
	const name = "digitalOceanCluster"
	const regionSlug = "nyc3"
	const versionSlug = "1.17.5-do.0"
	const vpcID = "18iahdddoaerr"
	const nodePoolName = "node-pool-1"
	const nodePoolSize = "s-2vcpu-2gb"
	const autoUpgraded = true
	const nodePoolAutoScale = true
	var tags = []string{"tag1", "tag2"}
	var nodePoolLabels = []string{"key1=label1", "key2=label2"}

	var nodePoolCount int = 3
	var nodePoolMin int = 2
	var nodePoolMax int = 4

	driverOptions := types.DriverOptions{
		StringOptions: map[string]string{
			"token": token,
			"display-name":   displayName,
			"name":           name,
			"region-slug":    regionSlug,
			"version-slug":   versionSlug,
			"vpc-id":         vpcID,
			"node-pool-name": nodePoolName,
			"node-pool-size": nodePoolSize,
		},
		BoolOptions: map[string]bool{
			"auto-upgraded":       autoUpgraded,
			"node-pool-autoscale": nodePoolAutoScale,
		},
		StringSliceOptions: map[string]*types.StringSlice{
			"tags":             {Value: tags},
			"node-pool-labels": {Value: nodePoolLabels},
		},
		IntOptions: map[string]int64{
			"node-pool-min":   int64(nodePoolMin),
			"node-pool-max":   int64(nodePoolMax),
			"node-pool-count": int64(nodePoolCount),
		},
	}

	state, err := stateBuilder.BuildStateFromOpts(&driverOptions)

	assert.Nil(t, err, "Not error in getStateFromOpts")
	assert.Equal(t, token, state.Token, "Token equals")
	assert.Equal(t, displayName, state.DisplayName, "DisplayName equals")
	assert.Equal(t, name, state.Name, "Name equals")
	assert.Equal(t, regionSlug, state.RegionSlug, "RegionSlug equals")
	assert.Equal(t, versionSlug, state.VersionSlug, "VersionSlug equals")
	assert.Equal(t, vpcID, state.VPCID, "VPCID equals")
	assert.Equal(t, nodePoolName, state.NodePool.Name, "nodePoolName equals")
	assert.Equal(t, nodePoolSize, state.NodePool.Size, "nodePoolSize equals")
	assert.Equal(t, autoUpgraded, state.AutoUpgrade, "autoUpgraded equals")
	assert.Equal(t, nodePoolAutoScale, state.NodePool.AutoScale, "nodePoolAutoScale equals")
	assert.Equal(t, tags, state.Tags, "tags equals")
	assert.Equal(t, map[string]string{"key1": "label1", "key2": "label2"}, state.NodePool.Labels, "nodePoolLabels equals")
	assert.Equal(t, nodePoolCount, state.NodePool.Count, "nodePoolCount equals")
	assert.Equal(t, nodePoolMin, state.NodePool.MinNodes, "nodePoolMin equals")
	assert.Equal(t, nodePoolMax, state.NodePool.MaxNodes, "nodePoolMax equals")

}

func TestGetStateFromOptsKeysCamelCase(t *testing.T) {

	const token = "tkalal1234761"
	const displayName = "cluster-test"
	const name = "digitalOceanCluster"
	const regionSlug = "nyc3"
	const versionSlug = "1.17.5-do.0"
	const vpcID = "18iahdddoaerr"
	const nodePoolName = "node-pool-1"
	const nodePoolSize = "s-2vcpu-2gb"
	const autoUpgraded = true
	const nodePoolAutoScale = true
	var tags = []string{"tag1", "tag2"}
	var nodePoolLabels = []string{"key1=label1", "key2=label2"}

	var nodePoolCount int = 3
	var nodePoolMin int = 2
	var nodePoolMax int = 4

	driverOptions := types.DriverOptions{
		StringOptions: map[string]string{
			"token": token,
			"name": name,
			"displayName":  displayName,
			"regionSlug":   regionSlug,
			"versionSlug":  versionSlug,
			"vpcID":        vpcID,
			"nodePoolName": nodePoolName,
			"nodePoolSize": nodePoolSize,
		},
		BoolOptions: map[string]bool{
			"autoUpgraded":      autoUpgraded,
			"nodePoolAutoscale": nodePoolAutoScale,
		},
		StringSliceOptions: map[string]*types.StringSlice{
			"tags":           {Value: tags},
			"nodePoolLabels": {Value: nodePoolLabels},
		},
		IntOptions: map[string]int64{
			"nodePoolMin":   int64(nodePoolMin),
			"nodePoolMax":   int64(nodePoolMax),
			"nodePoolCount": int64(nodePoolCount),
		},
	}

	state, err := stateBuilder.BuildStateFromOpts(&driverOptions)

	assert.Nil(t, err, "Not error in getStateFromOpts")
	assert.Equal(t, token, state.Token, "Token equals")
	assert.Equal(t, displayName, state.DisplayName, "DisplayName equals")
	assert.Equal(t, name, state.Name, "Name equals")
	assert.Equal(t, regionSlug, state.RegionSlug, "RegionSlug equals")
	assert.Equal(t, versionSlug, state.VersionSlug, "VersionSlug equals")
	assert.Equal(t, vpcID, state.VPCID, "VPCID equals")
	assert.Equal(t, nodePoolName, state.NodePool.Name, "nodePoolName equals")
	assert.Equal(t, nodePoolSize, state.NodePool.Size, "nodePoolSize equals")
	assert.Equal(t, autoUpgraded, state.AutoUpgrade, "autoUpgraded equals")
	assert.Equal(t, nodePoolAutoScale, state.NodePool.AutoScale, "nodePoolAutoScale equals")
	assert.Equal(t, tags, state.Tags, "tags equals")
	assert.Equal(t, map[string]string{"key1": "label1", "key2": "label2"}, state.NodePool.Labels, "nodePoolLabels equals")
	assert.Equal(t, nodePoolCount, state.NodePool.Count, "nodePoolCount equals")
	assert.Equal(t, nodePoolMin, state.NodePool.MinNodes, "nodePoolMin equals")
	assert.Equal(t, nodePoolMax, state.NodePool.MaxNodes, "nodePoolMax equals")

}

func TestGetStateFromOptsKeysNotAutoScale(t *testing.T) {

	const nodePoolCount = 8

	driverOptions := types.DriverOptions{
		BoolOptions: map[string]bool{
			"nodePoolAutoscale": false,
		},
		IntOptions: map[string]int64{
			"nodePoolMin":   int64(5),
			"nodePoolMax":   int64(10),
			"nodePoolCount": int64(nodePoolCount),
		},
	}

	state, err := stateBuilder.BuildStateFromOpts(&driverOptions)

	assert.Nil(t, err, "Not error in getStateFromOpts")
	assert.Equal(t, 0, state.NodePool.MinNodes)
	assert.Equal(t, 0, state.NodePool.MaxNodes)
	assert.Equal(t, nodePoolCount, state.NodePool.Count)

}
