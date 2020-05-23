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
	CreateCluster(ctx context.Context, state state.State) (string, error)
	GetNodeCount(ctx context.Context, clusterID string) (int,error)
	GetKubeConfig(clusterID string)(*store.KubeConfig,error)
	WaitCluster(ctx context.Context, clusterID string)error
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

func (do *digitalOceanImpl) CreateCluster(ctx context.Context, state state.State) (string, error){
	createClusterRequest := &godo.KubernetesClusterCreateRequest{
		Name: state.Name,
		Tags: state.Tags,
		AutoUpgrade: state.AutoUpgrade,
		RegionSlug: state.RegionSlug,
		VersionSlug: state.VersionSlug,
		NodePools: []*godo.KubernetesNodePoolCreateRequest{state.NodePool},
	}

	cluster, _, err := do.client.Kubernetes.Create(ctx,createClusterRequest)

	if err != nil {
		return "",errors.Wrap(err,"error creating the cluster")
	}

	return cluster.ID, nil
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

func (do digitalOceanImpl) WaitCluster(ctx context.Context, clusterID string)error{

	for {
		cluster, _, err := do.client.Kubernetes.Get(ctx, clusterID)

		if err != nil {
			return errors.Wrap(err, "error get cluster in waitCluster")
		}

		if cluster.Status.State == godo.KubernetesClusterStatusError {
			return errors.New("cluster is not being created")
		}

		if cluster.Status.State == godo.KubernetesClusterStatusProvisioning{
			do.sleeper.Sleep(5 * time.Second)
			continue
		}

		return nil
	}
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





