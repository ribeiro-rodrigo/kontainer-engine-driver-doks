package state

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strings"

	"github.com/rancher/kontainer-engine/drivers/options"
	"github.com/rancher/kontainer-engine/types"
)

type Cluster struct {
	ClusterID 	string `json:"cluster_id,omitempty"`
	Token 		string `json:"token,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Name        string `json:"name,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	AutoUpgrade *bool `json:"auto_upgrade,omitempty"`
	RegionSlug  string `json:"region_slug,omitempty"`
	VPCID       string `json:"vpc_id,omitempty"`
	VersionSlug string `json:"version_slug,omitempty"`
	NodePoolID  string `json:"node_pool_id,omitempty"`
}

type NodePool struct {
	Name      string
	Size      string
	Count     int
	Tags      []string
	Labels    map[string]string
	AutoScale *bool
	MinNodes  int
	MaxNodes  int
}

func (state *Cluster) Save(clusterInfo *types.ClusterInfo) error{
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
	BuildStatesFromOpts(driverOptions *types.DriverOptions) (Cluster, NodePool ,error)
	BuildClusterStateFromClusterInfo(clusterInfo *types.ClusterInfo)(Cluster,error)
}

type builderImpl struct{}

func NewBuilder() Builder {
	return builderImpl{}
}

func (builderImpl) BuildStatesFromOpts(driverOptions *types.DriverOptions) (Cluster, NodePool ,error) {

	clusterState := Cluster{Tags:[]string{}}
	nodePoolState := NodePool{}

	getValue := func(typ string, keys ...string) interface{} {
		return options.GetValueFromDriverOptions(driverOptions, typ, keys...)
	}

	clusterState.Token = getValue(types.StringType, "token").(string)
	clusterState.DisplayName = getValue(types.StringType, "display-name", "displayName").(string)
	clusterState.Name = getValue(types.StringType, "name").(string)
	clusterState.Tags = getTagsFromStringSlice(getValue(types.StringSliceType, "tags").(*types.StringSlice))
	clusterState.AutoUpgrade = getBoolPointer(getValue(types.BoolPointerType, "auto-upgraded", "autoUpgraded"))
	clusterState.RegionSlug = getValue(types.StringType, "region-slug", "regionSlug").(string)
	clusterState.VPCID = getValue(types.StringType, "vpc-id", "vpcID").(string)
	clusterState.VersionSlug = getValue(types.StringType, "version-slug", "versionSlug").(string)
	nodePoolState.Name = getValue(types.StringType, "node-pool-name", "nodePoolName").(string)
	nodePoolState.AutoScale = getBoolPointer(
		getValue(types.BoolPointerType, "node-pool-autoscale", "nodePoolAutoscale"),
	)

	if nodePoolState.AutoScale != nil && *nodePoolState.AutoScale {
		nodePoolState.MaxNodes = int(getValue(types.IntType, "node-pool-max", "nodePoolMax").(int64))
		nodePoolState.MinNodes = int(getValue(types.IntType, "node-pool-min", "nodePoolMin").(int64))
	}

	nodePoolState.Count = int(getValue(types.IntType, "node-pool-count", "nodePoolCount").(int64))

	nodePoolLabels := getLabelsFromStringSlice(
		getValue(types.StringSliceType, "node-pool-labels", "nodePoolLabels").(*types.StringSlice),
	)

	nodePoolState.Labels = nodePoolLabels
	nodePoolState.Size = getValue(types.StringType, "node-pool-size", "nodePoolSize").(string)

	return clusterState, nodePoolState, nil
}

func (builderImpl) BuildClusterStateFromClusterInfo(clusterInfo *types.ClusterInfo)(Cluster, error){
	stateJson, ok := clusterInfo.Metadata["state"]
	state := Cluster{}

	if !ok{
		return state, errors.New("there is no state in the clusterInfo")
	}

	err := json.Unmarshal([]byte(stateJson),&state)

	return state, err

}

func getTagsFromStringSlice(tagsString *types.StringSlice)[]string{
	if tagsString.Value == nil {
		return []string{}
	}

	return tagsString.Value
}

func getBoolPointer(boolPointer interface{})*bool{

	if boolPointer == nil {
		return nil
	}
	return boolPointer.(*bool)
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
