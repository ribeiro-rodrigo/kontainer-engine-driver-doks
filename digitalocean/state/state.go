package state

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/rancher/kontainer-engine/drivers/options"
	"github.com/rancher/kontainer-engine/types"
)

type State struct {
	ClusterID string `json:"cluster_id,omitempty"`
	Token string `json:"token,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Name        string `json:"name,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	AutoUpgrade bool `json:"auto_upgrade,omitempty"`
	RegionSlug  string `json:"region_slug,omitempty"`
	VPCID       string `json:"vpc_id,omitempty"`
	VersionSlug string `json:"version_slug,omitempty"`
	NodePool    *godo.KubernetesNodePoolCreateRequest `json:"node_pool,omitempty"`
}

func (state *State) Save(clusterInfo *types.ClusterInfo) error{
	bytes, err := json.Marshal(state)

	if err != nil {
		return errors.Wrap(err, "could not marshal state")
	}

	if clusterInfo.Metadata == nil {
		clusterInfo.Metadata = make(map[string]string)
	}

	clusterInfo.Metadata["state"] = string(bytes)

	return nil
}

type Builder interface {
	BuildStateFromOpts(driverOptions *types.DriverOptions) (State, error)
	BuildStateFromClusterInfo(clusterInfo *types.ClusterInfo)(State,error)
}

type builderImpl struct{}

func NewBuilder() Builder {
	return builderImpl{}
}

func (builderImpl) BuildStateFromOpts(driverOptions *types.DriverOptions) (State, error) {

	state := State{
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

func (builderImpl) BuildStateFromClusterInfo(clusterInfo *types.ClusterInfo)(State, error){
	stateJson, ok := clusterInfo.Metadata["state"]
	state := State{}

	if !ok{
		return state, errors.New("there is no state in the clusterInfo")
	}

	err := json.Unmarshal([]byte(stateJson),&state)

	return state, err

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
