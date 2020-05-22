package options

import (
	"testing"

	"github.com/rancher/kontainer-engine/types"
	"github.com/stretchr/testify/assert"
)

var optionsBuilder = NewOptionsBuilder()

func TestGetCreateOptions(t *testing.T) {

	options := optionsBuilder.BuildCreateOptions()

	tokenFlag, ok := options.Options["token"]

	assert.True(t, ok, "Token is present")
	assert.Equal(t, types.StringType, tokenFlag.GetType(), "Token type is string")

	displayNameFlag, ok := options.Options["display-name"]

	assert.True(t, ok, "DisplayName flag is present")
	assert.Equal(t, types.StringType, displayNameFlag.GetType(), "DisplayName type is string")

	nameFlag, ok := options.Options["name"]

	assert.True(t, ok, "Name flag is present")
	assert.Equal(t, types.StringType, nameFlag.GetType(), "Name type is string")

	autoUpgradeFlag, ok := options.Options["auto-upgraded"]

	assert.True(t, ok, "AutoUpgrade flag is present")
	assert.Equal(t, types.BoolType, autoUpgradeFlag.GetType(), "AutoUpgrade type is bool")

	regionSlugFlag, ok := options.Options["region-slug"]

	assert.True(t, ok, "RegionSlug flag is present")
	assert.Equal(t, types.StringType, regionSlugFlag.GetType(), "RegionSlug type is string")

	versionSlugFlag, ok := options.Options["version-slug"]

	assert.True(t, ok, "VersionSlug flag is present")
	assert.Equal(t, types.StringType, versionSlugFlag.GetType(), "VersionSlug type is string")

	nodePoolNameFlag, ok := options.Options["node-pool-name"]

	assert.True(t, ok, "NodePoolName flag is present")
	assert.Equal(t, types.StringType, nodePoolNameFlag.GetType(), "NodePoolName type is string")

	nodePoolAutoScaleFlag, ok := options.Options["node-pool-autoscale"]

	assert.True(t, ok, "NodePoolAutoScale flag is present")
	assert.Equal(t, types.BoolType, nodePoolAutoScaleFlag.GetType(), "NodePoolAutoScale type is bool")

	nodePoolCountFlag, ok := options.Options["node-pool-count"]

	assert.True(t, ok, "NodePoolCount flag is present")
	assert.Equal(t, types.IntType, nodePoolCountFlag.GetType(), "NodePoolCount type is int")

	nodePoolMinFlag, ok := options.Options["node-pool-min"]

	assert.True(t, ok, "NodePoolMin flag is present")
	assert.Equal(t, types.IntType, nodePoolMinFlag.GetType(), "NodePoolMin type is int")

	nodePoolMaxFlag, ok := options.Options["node-pool-max"]

	assert.True(t, ok, "NodePoolMax flag is present")
	assert.Equal(t, types.IntType, nodePoolMaxFlag.GetType(), "NodePoolMax type is int")

	nodePoolSizeFlag, ok := options.Options["node-pool-size"]

	assert.True(t, ok, "NodePoolSize flag is present")
	assert.Equal(t, types.StringType, nodePoolSizeFlag.GetType(), "NodePoolSize type is string")

	tagsFlag, ok := options.Options["tags"]

	assert.True(t, ok, "Tags flag is present")
	assert.Equal(t, types.StringSliceType, tagsFlag.GetType(), "Tags type is []string")

	nodePoolLabelsFlag, ok := options.Options["node-pool-labels"]

	assert.True(t, ok, "NodePoolLabels flag is present")
	assert.Equal(t, types.StringSliceType, nodePoolLabelsFlag.GetType(), "NodePoolLabels type is []string")

	VPCIDFlag, ok := options.Options["vpc-id"]

	assert.True(t, ok, "VPCID flag is present")
	assert.Equal(t, types.StringType, VPCIDFlag.GetType(), "VPCID type is string")
}
