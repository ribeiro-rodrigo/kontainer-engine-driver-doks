package service

import (
	"context"
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
	GetKubeConfig(clusterID string)(*store.KubeConfig,error)
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

	clusterKubeConfig, err := do.waitAndRetryGetKubeConfig(clusterID)

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

func (do digitalOceanImpl) waitAndRetryGetKubeConfig(clusterID string)(*godo.KubernetesClusterConfig, error){

	retry := 3
	var err error
	var clusterKubeConfig *godo.KubernetesClusterConfig

	for i:=0; i<retry; i++ {
		clusterKubeConfig, _, err = do.client.Kubernetes.GetKubeConfig(context.TODO(),clusterID)
		if err != nil {
			retry--
			do.sleeper.Sleep(2 * time.Second)
			continue
		}
		break
	}

	return clusterKubeConfig, nil
}





