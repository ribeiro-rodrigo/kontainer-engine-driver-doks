package state

import (
	"strings"

	"github.com/digitalocean/godo"
	"github.com/rancher/kontainer-engine/drivers/options"
	"github.com/rancher/kontainer-engine/types"
)

type State struct {
	Token string
	DisplayName string
	Name        string
	Tags        []string
	AutoUpgrade bool
	RegionSlug  string
	VPCID       string
	VersionSlug string
	NodePool    *godo.KubernetesNodePoolCreateRequest
	ClusterInfo types.ClusterInfo
}

type StateBuilder struct{}

func NewStateBuilder() StateBuilder {
	return StateBuilder{}
}

func (*StateBuilder) BuildStateFromOpts(driverOptions *types.DriverOptions) (State, error) {

	state := State{
		ClusterInfo: types.ClusterInfo{
			Metadata: map[string]string{},
		},
		Tags:     []string{},
		NodePool: &godo.KubernetesNodePoolCreateRequest{},
	}

	getValue := func(typ string, keys ...string) interface{} {
		return options.GetValueFromDriverOptions(driverOptions, typ, keys...)
	}

	state.Token = getValue(types.StringType, "token").(string)
	state.DisplayName = getValue(types.StringType, "display-name", "displayName").(string)
	state.Name = getValue(types.StringType, "name").(string)
	state.Tags = getValue(types.StringSliceType, "tags").(*types.StringSlice).Value
	state.AutoUpgrade = getValue(types.BoolType, "auto-upgraded", "autoUpgraded").(bool)
	state.RegionSlug = getValue(types.StringType, "region-slug", "regionSlug").(string)
	state.VPCID = getValue(types.StringType, "vpc-id", "vpcID").(string)
	state.VersionSlug = getValue(types.StringType, "version-slug", "versionSlug").(string)

	state.NodePool.Name = getValue(types.StringType, "node-pool-name", "nodePoolName").(string)
	state.NodePool.AutoScale = getValue(types.BoolType, "node-pool-autoscale", "nodePoolAutoscale").(bool)

	if state.NodePool.AutoScale {
		state.NodePool.MaxNodes = int(getValue(types.IntType, "node-pool-max", "nodePoolMax").(int64))
		state.NodePool.MinNodes = int(getValue(types.IntType, "node-pool-min", "nodePoolMin").(int64))
	}

	state.NodePool.Count = int(getValue(types.IntType, "node-pool-count", "nodePoolCount").(int64))

	nodePoolLabels := getLabelsFromStringSlice(
		getValue(types.StringSliceType, "node-pool-labels", "nodePoolLabels").(*types.StringSlice),
	)

	state.NodePool.Labels = nodePoolLabels
	state.NodePool.Size = getValue(types.StringType, "node-pool-size", "nodePoolSize").(string)

	return state, nil
}

func getLabelsFromStringSlice(labelsString *types.StringSlice) map[string]string {

	labels := map[string]string{}

	if labelsString == nil {
		return labels
	}

	for _, part := range labelsString.Value {
		kv := strings.Split(part, "=")

		if len(kv) == 2 {
			labels[kv[0]] = kv[1]
		}
	}

	return labels
}
