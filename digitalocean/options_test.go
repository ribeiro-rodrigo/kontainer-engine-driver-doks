package digitalocean

import (
	"testing"

	"github.com/rancher/kontainer-engine/types"
	"github.com/stretchr/testify/assert"
)

func TestGetCreateOptions(t *testing.T) {
	options := getCreateOptions()

	DisplayNameFlag, ok := options.Options["display-name"]

	assert.True(t, ok, "DisplayName flag is present")
	assert.Equal(t, types.StringType, DisplayNameFlag.GetType(), "DisplayName type is string")

	NameFlag, ok := options.Options["name"]

	assert.True(t, ok, "Name flag is present")
	assert.Equal(t, types.StringType, NameFlag.GetType(), "Name type is string")

	AutoUpgradeFlag, ok := options.Options["auto-upgraded"]

	assert.True(t, ok, "AutoUpgrade flag is present")
	assert.Equal(t, types.BoolType, AutoUpgradeFlag.GetType(), "AutoUpgrade type is bool")

	RegionSlugFlag, ok := options.Options["region-slug"]

	assert.True(t, ok, "RegionSlug flag is present")
	assert.Equal(t, types.StringType, RegionSlugFlag.GetType(), "RegionSlug type is string")

	VersionSlugFlag, ok := options.Options["version-slug"]

	assert.True(t, ok, "VersionSlug flag is present")
	assert.Equal(t, types.StringType, VersionSlugFlag.GetType(), "VersionSlug type is string")

	NodePoolNameFlag, ok := options.Options["node-pool-name"]

	assert.True(t, ok, "NodePoolName flag is present")
	assert.Equal(t, types.StringType, NodePoolNameFlag.GetType(), "NodePoolName type is string")

	NodePoolAutoScaleFlag, ok := options.Options["node-pool-autoscale"]

	assert.True(t, ok, "NodePoolAutoScale flag is present")
	assert.Equal(t, types.BoolType, NodePoolAutoScaleFlag.GetType(), "NodePoolAutoScale type is bool")

	NodePoolCountFlag, ok := options.Options["node-pool-count"]

	assert.True(t, ok, "NodePoolCount flag is present")
	assert.Equal(t, types.IntType, NodePoolCountFlag.GetType(), "NodePoolCount type is int")

	NodePoolMinFlag, ok := options.Options["node-pool-min"]

	assert.True(t, ok, "NodePoolMin flag is present")
	assert.Equal(t, types.IntType, NodePoolMinFlag.GetType(), "NodePoolMin type is int")

	NodePoolMaxFlag, ok := options.Options["node-pool-max"]

	assert.True(t, ok, "NodePoolMax flag is present")
	assert.Equal(t, types.IntType, NodePoolMaxFlag.GetType(), "NodePoolMax type is int")

	NodePoolSizeFlag, ok := options.Options["node-pool-size"]

	assert.True(t, ok, "NodePoolSize flag is present")
	assert.Equal(t, types.StringType, NodePoolSizeFlag.GetType(), "NodePoolSize type is string")

	TagsFlag, ok := options.Options["tags"]

	assert.True(t, ok, "Tags flag is present")
	assert.Equal(t, types.StringSliceType, TagsFlag.GetType(), "Tags type is []string")

	NodePoolLabelsFlag, ok := options.Options["node-pool-labels"]

	assert.True(t, ok, "NodePoolLabels flag is present")
	assert.Equal(t, types.StringSliceType, NodePoolLabelsFlag.GetType(), "NodePoolLabels type is []string")

	VPCIDFlag, ok := options.Options["vpc-id"]

	assert.True(t, ok, "VPCID flag is present")
	assert.Equal(t, types.StringType, VPCIDFlag.GetType(), "VPCID type is string")
}
