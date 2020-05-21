package digitalocean

import (
	"context"

	"github.com/rancher/kontainer-engine/types"
	"github.com/ribeiro-rodrigo/kontainer-engine-driver-digitalocean/digitalocean/options"
	"github.com/ribeiro-rodrigo/kontainer-engine-driver-digitalocean/digitalocean/state"
)

type DigitalOceanDriver struct {
	stateBuilder       state.StateBuilder
	optionsBuilder     options.OptionsBuilder
	driverCapabilities types.Capabilities
}

func NewDigitalOceanDriver() DigitalOceanDriver {
	driver := DigitalOceanDriver{
		stateBuilder:   state.NewStateBuilder(),
		optionsBuilder: options.NewOptionsBuilder(),
	}

	return driver
}

func (driver *DigitalOceanDriver) GetDriverCreateOptions(ctx context.Context) (*types.DriverFlags, error) {

	return driver.optionsBuilder.BuildCreateOptions(), nil
}

func (*DigitalOceanDriver) GetDriverUpdateOptions(ctx context.Context) (*types.DriverFlags, error) {
	return nil, nil
}

func (*DigitalOceanDriver) Create(ctx context.Context, opts *types.DriverOptions, clusterInfo *types.ClusterInfo) (*types.ClusterInfo, error) {

	return nil, nil
}

func (*DigitalOceanDriver) Update(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions) (*types.ClusterInfo, error) {

	return nil, nil
}

func (*DigitalOceanDriver) PostCheck(ctx context.Context, clusterInfo *types.ClusterInfo) (*types.ClusterInfo, error) {
	return nil, nil
}

func (*DigitalOceanDriver) Remove(ctx context.Context, clusterInfo *types.ClusterInfo) error {
	return nil
}

func (*DigitalOceanDriver) GetVersion(ctx context.Context, clusterInfo *types.ClusterInfo) (*types.KubernetesVersion, error) {
	return nil, nil
}

func (*DigitalOceanDriver) SetVersion(ctx context.Context, clusterInfo *types.ClusterInfo, version *types.KubernetesVersion) error {
	return nil
}

func (*DigitalOceanDriver) GetClusterSize(ctx context.Context, clusterInfo *types.ClusterInfo) (*types.NodeCount, error) {
	return nil, nil
}

func (*DigitalOceanDriver) SetClusterSize(ctx context.Context, clusterInfo *types.ClusterInfo, count *types.NodeCount) error {

	return nil
}

func (*DigitalOceanDriver) GetCapabilities(ctx context.Context) (*types.Capabilities, error) {
	return nil, nil
}

func (*DigitalOceanDriver) RemoveLegacyServiceAccount(ctx context.Context, clusterInfo *types.ClusterInfo) error {
	return nil
}

func (*DigitalOceanDriver) ETCDSave(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions, snapshotName string) error {
	return nil
}

func (*DigitalOceanDriver) ETCDRestore(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions, snapshotName string) (*types.ClusterInfo, error) {
	return nil, nil
}

func (*DigitalOceanDriver) GetK8SCapabilities(ctx context.Context, opts *types.DriverOptions) (*types.K8SCapabilities, error) {
	return nil, nil
}

func (*DigitalOceanDriver) ETCDRemoveSnapshot(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions, snapshotName string) error {
	return nil
}
