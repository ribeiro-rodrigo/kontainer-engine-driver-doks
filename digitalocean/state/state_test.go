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

	clusterState, nodePoolState, err := stateBuilder.BuildStatesFromOpts(&driverOptions)

	assert.Nil(t, err, "Not error in getStateFromOpts")
	assert.Equal(t, token, clusterState.Token, "Token equals")
	assert.Equal(t, displayName, clusterState.DisplayName, "DisplayName equals")
	assert.Equal(t, name, clusterState.Name, "Name equals")
	assert.Equal(t, regionSlug, clusterState.RegionSlug, "RegionSlug equals")
	assert.Equal(t, versionSlug, clusterState.VersionSlug, "VersionSlug equals")
	assert.Equal(t, vpcID, clusterState.VPCID, "VPCID equals")
	assert.Equal(t, nodePoolName, nodePoolState.Name, "nodePoolName equals")
	assert.Equal(t, nodePoolSize, nodePoolState.Size, "nodePoolSize equals")
	assert.Equal(t, autoUpgraded, *clusterState.AutoUpgrade, "autoUpgraded equals")
	assert.Equal(t, nodePoolAutoScale, *nodePoolState.AutoScale, "nodePoolAutoScale equals")
	assert.Equal(t, tags, clusterState.Tags, "tags equals")
	assert.Equal(t, map[string]string{"key1": "label1", "key2": "label2"}, nodePoolState.Labels, "nodePoolLabels equals")
	assert.Equal(t, nodePoolCount, nodePoolState.Count, "nodePoolCount equals")
	assert.Equal(t, nodePoolMin, nodePoolState.MinNodes, "nodePoolMin equals")
	assert.Equal(t, nodePoolMax, nodePoolState.MaxNodes, "nodePoolMax equals")

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

	clusterState, nodePoolState, err := stateBuilder.BuildStatesFromOpts(&driverOptions)

	assert.Nil(t, err, "Not error in getStateFromOpts")
	assert.Equal(t, token, clusterState.Token, "Token equals")
	assert.Equal(t, displayName, clusterState.DisplayName, "DisplayName equals")
	assert.Equal(t, name, clusterState.Name, "Name equals")
	assert.Equal(t, regionSlug, clusterState.RegionSlug, "RegionSlug equals")
	assert.Equal(t, versionSlug, clusterState.VersionSlug, "VersionSlug equals")
	assert.Equal(t, vpcID, clusterState.VPCID, "VPCID equals")
	assert.Equal(t, nodePoolName, nodePoolState.Name, "nodePoolName equals")
	assert.Equal(t, nodePoolSize, nodePoolState.Size, "nodePoolSize equals")
	assert.Equal(t, autoUpgraded, *clusterState.AutoUpgrade, "autoUpgraded equals")
	assert.Equal(t, nodePoolAutoScale, *nodePoolState.AutoScale, "nodePoolAutoScale equals")
	assert.Equal(t, tags, clusterState.Tags, "tags equals")
	assert.Equal(t, map[string]string{"key1": "label1", "key2": "label2"}, nodePoolState.Labels, "nodePoolLabels equals")
	assert.Equal(t, nodePoolCount, nodePoolState.Count, "nodePoolCount equals")
	assert.Equal(t, nodePoolMin, nodePoolState.MinNodes, "nodePoolMin equals")
	assert.Equal(t, nodePoolMax, nodePoolState.MaxNodes, "nodePoolMax equals")

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

	_, nodePoolState, err := stateBuilder.BuildStatesFromOpts(&driverOptions)

	assert.Nil(t, err, "Not error in getStateFromOpts")
	assert.Equal(t, 0, nodePoolState.MinNodes)
	assert.Equal(t, 0, nodePoolState.MaxNodes)
	assert.Equal(t, nodePoolCount, nodePoolState.Count)

}

func TestGetStateWithNilValues(t *testing.T){
	driverOptions := types.DriverOptions{}

	clusterState, nodePoolState, err := stateBuilder.BuildStatesFromOpts(&driverOptions)

	assert.Nil(t, err, "Not error in getStateFromOpts")
	assert.Equal(t, "", clusterState.Name, "Name empty value")
	assert.Equal(t, "", clusterState.Token, "Token equals")
	assert.Equal(t, "", clusterState.DisplayName, "DisplayName equals")
	assert.Equal(t, "", clusterState.Name, "Name equals")
	assert.Equal(t, "", clusterState.RegionSlug, "RegionSlug equals")
	assert.Equal(t, "", clusterState.VersionSlug, "VersionSlug equals")
	assert.Equal(t, "", clusterState.VPCID, "VPCID equals")
	assert.Equal(t, "", nodePoolState.Name, "nodePoolName equals")
	assert.Equal(t, "", nodePoolState.Size, "nodePoolSize equals")
	assert.Nil(t, clusterState.AutoUpgrade, "autoUpgraded equals")
	assert.Nil(t, nodePoolState.AutoScale, "nodePoolAutoScale equals")
	assert.Equal(t, []string{}, clusterState.Tags, "tags equals")
	assert.Equal(t, map[string]string{}, nodePoolState.Labels, "nodePoolLabels equals")
	assert.Equal(t, 0, nodePoolState.Count, "nodePoolCount equals")
	assert.Equal(t, 0, nodePoolState.MinNodes, "nodePoolMin equals")
	assert.Equal(t, 0, nodePoolState.MaxNodes, "nodePoolMax equals")
}


