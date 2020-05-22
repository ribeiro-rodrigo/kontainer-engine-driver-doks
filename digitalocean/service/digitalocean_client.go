package service

import (
	"context"
	"github.com/digitalocean/godo"
	"github.com/ghodss/yaml"
	"github.com/rancher/kontainer-engine/store"
	"github.com/ribeiro-rodrigo/kontainer-engine-driver-digitalocean/digitalocean/state"
)

type DigitalOcean struct {
	client *godo.Client
}

func NewDigitalOcean(token string) DigitalOcean {
	return DigitalOcean{
		client: godo.NewFromToken(token),
	}
}

func (do *DigitalOcean) CreateCluster(state state.State) (string, error){
	createClusterRequest := &godo.KubernetesClusterCreateRequest{
		Name: state.Name,
		Tags: state.Tags,
		AutoUpgrade: state.AutoUpgrade,
		RegionSlug: state.RegionSlug,
		VersionSlug: state.VersionSlug,
		NodePools: []*godo.KubernetesNodePoolCreateRequest{state.NodePool},
	}

	cluster, _, err := do.client.Kubernetes.Create(context.Background(),createClusterRequest)

	if err != nil {
		return "",err
	}

	return cluster.ID, nil
}

func (do *DigitalOcean) GetKubeConfig(clusterID string)(*store.KubeConfig,error){

	clusterKubeConfig, _, err := do.client.Kubernetes.GetKubeConfig(context.TODO(),clusterID)

	if err != nil {
		return nil,err
	}

	kubeConfig := &store.KubeConfig{}

	err = yaml.Unmarshal(clusterKubeConfig.KubeconfigYAML, kubeConfig)

	if err != nil {
		return nil, err
	}

	return kubeConfig, nil
}


