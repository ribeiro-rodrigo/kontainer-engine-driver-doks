package service

import (
	"context"
	"github.com/digitalocean/godo"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/rancher/kontainer-engine/store"
	"github.com/ribeiro-rodrigo/kontainer-engine-driver-digitalocean/digitalocean/state"
)

type DigitalOceanFactory func(token string)DigitalOcean

func NewDigitalOceanFactory()DigitalOceanFactory{
	return func(token string)DigitalOcean{
		return newDigitalOcean(token)
	}
}

type DigitalOcean interface {
	CreateCluster(ctx context.Context, state state.State) (string, error)
	GetKubeConfig(clusterID string)(*store.KubeConfig,error)
}

type digitalOceanImpl struct {
	client *godo.Client
}

func newDigitalOcean(token string) DigitalOcean {
	return &digitalOceanImpl{
		client: godo.NewFromToken(token),
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

	clusterKubeConfig, _, err := do.client.Kubernetes.GetKubeConfig(context.TODO(),clusterID)

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





