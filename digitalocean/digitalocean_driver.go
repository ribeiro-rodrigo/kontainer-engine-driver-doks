package digitalocean

import (
	"context"

	"github.com/rancher/kontainer-engine/types"
	"github.com/ribeiro-rodrigo/kontainer-engine-driver-digitalocean/digitalocean/options"
	"github.com/ribeiro-rodrigo/kontainer-engine-driver-digitalocean/digitalocean/state"
)

type Driver struct {
	stateBuilder       state.StateBuilder
	optionsBuilder     options.OptionsBuilder
	driverCapabilities types.Capabilities
}

func NewDriver() Driver {
	driver := Driver{
		stateBuilder:   state.NewStateBuilder(),
		optionsBuilder: options.NewOptionsBuilder(),
	}

	return driver
}

func (driver *Driver) GetDriverCreateOptions(ctx context.Context) (*types.DriverFlags, error) {

	return driver.optionsBuilder.BuildCreateOptions(), nil
}

func (*Driver) GetDriverUpdateOptions(ctx context.Context) (*types.DriverFlags, error) {
	return nil, nil
}

func (*Driver) Create(ctx context.Context, opts *types.DriverOptions, _ *types.ClusterInfo) (*types.ClusterInfo, error) {

	return nil, nil
}

func (*Driver) Update(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions) (*types.ClusterInfo, error) {

	return nil, nil
}

func (*Driver) PostCheck(ctx context.Context, clusterInfo *types.ClusterInfo) (*types.ClusterInfo, error) {
	return nil, nil
}

func (*Driver) Remove(ctx context.Context, clusterInfo *types.ClusterInfo) error {
	return nil
}

func (*Driver) GetVersion(ctx context.Context, clusterInfo *types.ClusterInfo) (*types.KubernetesVersion, error) {
	return nil, nil
}

func (*Driver) SetVersion(ctx context.Context, clusterInfo *types.ClusterInfo, version *types.KubernetesVersion) error {
	return nil
}

func (*Driver) GetClusterSize(ctx context.Context, clusterInfo *types.ClusterInfo) (*types.NodeCount, error) {
	return nil, nil
}

func (*Driver) SetClusterSize(ctx context.Context, clusterInfo *types.ClusterInfo, count *types.NodeCount) error {

	return nil
}

func (*Driver) GetCapabilities(ctx context.Context) (*types.Capabilities, error) {
	return nil, nil
}

func (*Driver) RemoveLegacyServiceAccount(ctx context.Context, clusterInfo *types.ClusterInfo) error {
	return nil
}

func (*Driver) ETCDSave(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions, snapshotName string) error {
	return nil
}

func (*Driver) ETCDRestore(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions, snapshotName string) (*types.ClusterInfo, error) {
	return nil, nil
}

func (*Driver) GetK8SCapabilities(ctx context.Context, opts *types.DriverOptions) (*types.K8SCapabilities, error) {
	return nil, nil
}

func (*Driver) ETCDRemoveSnapshot(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions, snapshotName string) error {
	return nil
}
