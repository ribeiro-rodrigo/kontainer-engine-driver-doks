package digitalocean

import (
	"context"
	"errors"
	"github.com/ribeiro-rodrigo/kontainer-engine-driver-digitalocean/digitalocean/service"
	"github.com/sirupsen/logrus"

	"github.com/rancher/kontainer-engine/types"
	"github.com/ribeiro-rodrigo/kontainer-engine-driver-digitalocean/digitalocean/options"
	"github.com/ribeiro-rodrigo/kontainer-engine-driver-digitalocean/digitalocean/state"
)

type Driver struct {
	stateBuilder       state.Builder
	optionsBuilder     options.Builder
	digitalOceanFactory service.DigitalOceanFactory
	driverCapabilities types.Capabilities
}

func NewDriver() Driver {
	driver := Driver{
		stateBuilder:   state.NewBuilder(),
		optionsBuilder: options.NewBuilder(),
		digitalOceanFactory: service.NewDigitalOceanFactory(),
	}

	return driver
}

func (driver *Driver) GetDriverCreateOptions(_ context.Context) (*types.DriverFlags, error) {
	logrus.Debug("DigitalOcean.Driver.GetDriverCreateOptions(...) called")
	return driver.optionsBuilder.BuildCreateOptions(), nil
}

func (driver *Driver) GetDriverUpdateOptions(_ context.Context) (*types.DriverFlags, error) {
	logrus.Debug("DigitalOcean.Driver.GetDriverUpdateOptions(...) called")
	return driver.optionsBuilder.BuildUpdateOptions(), nil
}

func (driver *Driver) Create(ctx context.Context, opts *types.DriverOptions, info *types.ClusterInfo) (*types.ClusterInfo, error) {
	logrus.Debug("DigitalOcean.Driver.Create(...) called")
	clusterState, err := driver.stateBuilder.BuildStateFromOpts(opts)

	if err != nil{
		logrus.Debugf("Error building clusterState: %v",err)
		return nil, err
	}

	if clusterState.Token == ""{
		logrus.Debugf("Error token not found: %v",err)
		err = errors.New("token was not reported")
		return nil, err
	}

	digitalOceanService := driver.digitalOceanFactory(clusterState.Token)

	clusterID, nodePoolID, err := digitalOceanService.CreateCluster(ctx, clusterState)

	if err != nil {
		logrus.Debugf("Error crate cluster: %v",err)
		return nil, err
	}

	clusterState.ClusterID = clusterID
	clusterState.NodePool.ID = nodePoolID

	err = clusterState.Save(info)

	if err != nil {
		logrus.Debugf("Error save clusterState: %v",err)
		return nil, err
	}

	err = digitalOceanService.WaitClusterCreated(ctx,clusterID)

	if err != nil {
		logrus.Debugf("Error wait cluster: %v",err)
		return nil, err
	}

	return info, nil
}

func (driver *Driver) PostCheck(ctx context.Context, clusterInfo *types.ClusterInfo) (*types.ClusterInfo, error) {
	logrus.Debug("DigitalOcean.Driver.PostCheck(...) called")

	clusterState, err := driver.stateBuilder.BuildStateFromClusterInfo(clusterInfo)

	if err != nil {
		logrus.Debugf("Error build clusterState: %v",err)
		return nil, err
	}

	digitalOceanService := driver.digitalOceanFactory(clusterState.Token)

	kubeConfig, err := digitalOceanService.GetKubeConfig(clusterState.ClusterID)

	if err != nil {
		logrus.Debugf("Error get kubeConfig %v",err)
		return nil, err
	}

	if len(kubeConfig.Clusters) > 0 {
		cluster := kubeConfig.Clusters[0].Cluster
		clusterInfo.RootCaCertificate = cluster.CertificateAuthorityData
		clusterInfo.Endpoint = cluster.Server
	}else{
		return nil, errors.New("the kubeconfig file is invalid. Cluster not found")
	}

	if len(kubeConfig.Users) > 0 {
		clusterInfo.ServiceAccountToken = kubeConfig.Users[0].User.Token
	}else{
		return nil, errors.New("the kubeconfig file is invalid. Token not found")
	}

	currentNodeCount, err := digitalOceanService.GetNodeCount(ctx, clusterState.ClusterID)

	if err != nil {
		logrus.Debugf("Error get node count %v",err)
		return nil, err
	}

	clusterInfo.Version = clusterState.VersionSlug
	clusterInfo.NodeCount = int64(currentNodeCount)

	return clusterInfo, nil
}

func (*Driver) Update(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions) (*types.ClusterInfo, error) {

	return nil, nil
}

func (driver *Driver) Remove(ctx context.Context, clusterInfo *types.ClusterInfo) error {
	logrus.Debug("DigitalOcean.Driver.Remove(...) called")

	clusterState, err := driver.stateBuilder.BuildStateFromClusterInfo(clusterInfo)

	if err != nil {
		logrus.Debugf("Error build state %v",err)
		return err
	}

	digitalOceanService := driver.digitalOceanFactory(clusterState.Token)

	err = digitalOceanService.DeleteCluster(ctx, clusterState.ClusterID)

	if err != nil {
		logrus.Debugf("Error delete cluster %v",err)
		return err
	}

	err = digitalOceanService.WaitClusterDeleted(ctx, clusterState.ClusterID)

	if err != nil {
		logrus.Debugf("Error wait delete cluster %v",err)
		return err
	}

	return nil
}

func (driver *Driver) GetVersion(ctx context.Context, clusterInfo *types.ClusterInfo) (*types.KubernetesVersion, error) {

	clusterState, err := driver.stateBuilder.BuildStateFromClusterInfo(clusterInfo)

	if err != nil {
		logrus.Debugf("Error BuildStateFromClusterInfo in get version %v",err)
		return nil, err
	}

	digitalOceanService := driver.digitalOceanFactory(clusterState.Token)

	kubernetesVersion, err :=  digitalOceanService.GetKubernetesClusterVersion(ctx, clusterState.ClusterID)

	if err != nil {
		logrus.Debugf("Error digital ocean service get cluster version in get version %v",err)
		return nil, err
	}

	return &types.KubernetesVersion{Version: kubernetesVersion}, nil
}

func (driver *Driver) SetVersion(ctx context.Context, clusterInfo *types.ClusterInfo, version *types.KubernetesVersion) error {

	clusterState, err := driver.stateBuilder.BuildStateFromClusterInfo(clusterInfo)

	if err != nil {
		logrus.Debugf("Error build state from cluster info in set version %v",err)
		return err
	}

	digitalOceanService := driver.digitalOceanFactory(clusterState.Token)

	err = digitalOceanService.UpgradeKubernetesVersion(ctx, clusterState.ClusterID, version.Version)

	if err != nil {
		logrus.Debugf("Error upgrade kubernetes version %v",err)
		return err
	}

	return nil
}

func (driver *Driver) GetClusterSize(ctx context.Context, clusterInfo *types.ClusterInfo) (*types.NodeCount, error) {

	clusterState, err :=  driver.stateBuilder.BuildStateFromClusterInfo(clusterInfo)

	if err != nil {
		logrus.Debugf("Error BuildStateFromClusterInfo in GetClusterSize %v",err)
		return nil, err
	}

	digitalOceanService := driver.digitalOceanFactory(clusterState.Token)

	nodeCount, err := digitalOceanService.GetNodeCount(ctx, clusterState.ClusterID)

	if err != nil {
		logrus.Debugf("Error GetNodeCount in GetClusterSize")
		return nil, err
	}

	return &types.NodeCount{Count: int64(nodeCount)}, nil
}

func (driver *Driver) SetClusterSize(ctx context.Context, clusterInfo *types.ClusterInfo, count *types.NodeCount) error {

	clusterState, err := driver.stateBuilder.BuildStateFromClusterInfo(clusterInfo)

	if err != nil {
		logrus.Debugf("Error BuildStateFromClusterInfo in SetClusterSize")
		return err
	}

	nodePool := clusterState.NodePool

	if nodePool.AutoScale != nil && *nodePool.AutoScale {
		if int64(nodePool.MinNodes) > count.Count{
			nodePool.MinNodes = int(count.Count)
		}else if int64(nodePool.MaxNodes) < count.Count {
			nodePool.MaxNodes = int(count.Count)
		}
	}

	nodePool.Count = int(count.Count)

	digitalOceanService := driver.digitalOceanFactory(clusterState.Token)

	err = digitalOceanService.UpdateNodePool(ctx, clusterState.ClusterID, clusterState.NodePool)

	if err != nil {
		logrus.Debugf("Error UpdateNodePool in SetClusterSize")
		return err
	}

	clusterState.Save(clusterInfo)

	return errors.New("operation resize cluster size  not implemented")
}

func (*Driver) GetCapabilities(ctx context.Context) (*types.Capabilities, error) {
	return nil, nil
}

func (*Driver) RemoveLegacyServiceAccount(ctx context.Context, clusterInfo *types.ClusterInfo) error {
	return errors.New("operation remove service account not implemented")
}

func (*Driver) ETCDSave(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions, snapshotName string) error {
	return errors.New("etcd backup operations are not implemented")
}

func (*Driver) ETCDRestore(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions, snapshotName string) (*types.ClusterInfo, error) {
	return nil, errors.New("etcd backup operations are not implemented")
}

func (*Driver) GetK8SCapabilities(ctx context.Context, opts *types.DriverOptions) (*types.K8SCapabilities, error) {
	return nil, nil
}

func (*Driver) ETCDRemoveSnapshot(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions, snapshotName string) error {
	return errors.New("etcd backup operations are not implemented")
}

