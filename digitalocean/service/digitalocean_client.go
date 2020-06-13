package service

import (
	"context"
	"fmt"
	"github.com/digitalocean/godo"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/rancher/kontainer-engine/store"
	"github.com/ribeiro-rodrigo/kontainer-engine-driver-digitalocean/digitalocean/helper"
	"github.com/ribeiro-rodrigo/kontainer-engine-driver-digitalocean/digitalocean/state"
	"time"
)

type DigitalOceanFactory func(token string)DigitalOcean

func NewDigitalOceanFactory()DigitalOceanFactory{
	return func(token string)DigitalOcean{
		return newDigitalOcean(token, helper.NewTimerSleeper())
	}
}

type DigitalOcean interface {
	CreateCluster(ctx context.Context, state state.State) (string, string, error)
	GetKubernetesClusterVersion(ctx context.Context, clusterID string)(string,error)
	UpgradeKubernetesVersion(ctx context.Context, clusterID, version string)error
	DeleteCluster(ctx context.Context, clusterID string)error
	GetNodeCount(ctx context.Context, clusterID string) (int,error)
	UpdateNodePool(ctx context.Context, clusterID string, nodePool state.NodePool ) error
	GetKubeConfig(clusterID string)(*store.KubeConfig,error)
	WaitClusterCreated(ctx context.Context, clusterID string)error
	WaitClusterDeleted(ctx context.Context, clusterID string)error
}

type digitalOceanImpl struct {
	client *godo.Client
	sleeper helper.Sleeper
}

func newDigitalOcean(token string, sleeper helper.Sleeper) DigitalOcean {
	return &digitalOceanImpl{
		client: godo.NewFromToken(token),
		sleeper: sleeper,
	}
}

func (do *digitalOceanImpl) CreateCluster(ctx context.Context, state state.State) (string, string, error){
	createClusterRequest := &godo.KubernetesClusterCreateRequest{
		Name: state.Name,
		Tags: state.Tags,
		AutoUpgrade: *state.AutoUpgrade,
		RegionSlug: state.RegionSlug,
		VersionSlug: state.VersionSlug,
		NodePools: do.buildNodePoolCreateRequest(state.NodePool),
	}

	cluster, _, err := do.client.Kubernetes.Create(ctx,createClusterRequest)

	if err != nil {
		return "","",errors.Wrap(err,"error creating the cluster")
	}

	return cluster.ID, cluster.NodePools[0].ID, nil
}

func (do *digitalOceanImpl) DeleteCluster(ctx context.Context, clusterID string)error{
	_, err := do.client.Kubernetes.Delete(ctx, clusterID)

	if err != nil {
		return errors.Wrap(err,"error in delete cluster")
	}

	return nil
}

func (do *digitalOceanImpl) GetKubeConfig(clusterID string)(*store.KubeConfig,error){

	clusterKubeConfig, _, err := do.client.Kubernetes.GetKubeConfig(context.TODO(), clusterID)

	if err != nil {
		return nil, errors.Wrapf(err,"error get kubeConfig for cluster %s",clusterID)
	}

	kubeConfig := &store.KubeConfig{}

	err = yaml.Unmarshal(clusterKubeConfig.KubeconfigYAML, kubeConfig)

	if err != nil {
		return nil, errors.Wrapf(err,"error marshal kubeConfig from clusterID %s",clusterID)
	}

	return kubeConfig, nil
}

func (do digitalOceanImpl) WaitClusterCreated(ctx context.Context, clusterID string)error{
	_, err := do.waitCluster(ctx, clusterID, godo.KubernetesClusterStatusRunning)
	return err
}

func (do digitalOceanImpl) WaitClusterDeleted(ctx context.Context, clusterID string)error{
	response, err := do.waitCluster(ctx, clusterID, godo.KubernetesClusterStatusDeleted)

	if response != nil && response.StatusCode == 404 {
		return nil
	}

	return err
}

func (do digitalOceanImpl) GetNodeCount(ctx context.Context,
	clusterID string) (int, error){

	cluster, _, err := do.client.Kubernetes.Get(ctx, clusterID)

	if err != nil {
		return 0, errors.Wrap(err,fmt.Sprintf("error in get cluster %v",err))
	}

	count := 0

	for _, pool := range cluster.NodePools{
		count += pool.Count
	}

	return count, err
}

func (do digitalOceanImpl) UpdateNodePool(ctx context.Context, clusterID string, nodePool state.NodePool) error{

	updateRequest := &godo.KubernetesNodePoolUpdateRequest{
		Name: nodePool.Name,
		Labels: nodePool.Labels,
		AutoScale: nodePool.AutoScale,
		MaxNodes: &nodePool.MaxNodes,
		MinNodes: &nodePool.MinNodes,
		Count: &nodePool.Count,
		Tags: nodePool.Tags,
	}

	_, _, err :=  do.client.Kubernetes.UpdateNodePool(ctx,clusterID,nodePool.ID,updateRequest)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error in getNodePool: %v",err))
	}

	return nil
}

func (do digitalOceanImpl) GetKubernetesClusterVersion(ctx context.Context, clusterID string)(string,error){
	cluster, _, err := do.client.Kubernetes.Get(ctx,clusterID)

	if err != nil {
		return "", errors.Wrap(err,fmt.Sprintf("error in get cluster %v",err))
	}

	return cluster.VersionSlug, nil
}

func (do digitalOceanImpl) UpgradeKubernetesVersion(ctx context.Context, clusterID, version string)error{
	upgradeRequest := &godo.KubernetesClusterUpgradeRequest{VersionSlug: version}

	_, err := do.client.Kubernetes.Upgrade(ctx, clusterID, upgradeRequest)

	return err
}

func (do digitalOceanImpl) waitCluster(ctx context.Context, clusterID string,
	statusState godo.KubernetesClusterStatusState)(*godo.Response, error){

	for {
		cluster, response, err := do.client.Kubernetes.Get(ctx, clusterID)

		if err != nil {
			err = errors.Wrap(err, "error get cluster in waitCluster")
			return response, err
		}

		if cluster.Status.State == godo.KubernetesClusterStatusError {
			err = errors.New("cluster status error")
			return response, err
		}

		if cluster.Status.State != statusState{
			do.sleeper.Sleep(5 * time.Second)
			continue
		}

		return response, err
	}
}

func (do digitalOceanImpl) buildNodePoolCreateRequest(nodePool state.NodePool) []*godo.KubernetesNodePoolCreateRequest{
	return []*godo.KubernetesNodePoolCreateRequest{
		{
			Name: nodePool.Name,
			Size: nodePool.Size,
			Count: nodePool.Count,
			Tags: nodePool.Tags,
			Labels: nodePool.Labels,
			AutoScale: *nodePool.AutoScale,
			MinNodes: nodePool.MinNodes,
			MaxNodes: nodePool.MaxNodes,
		},
	}
}





