package digitalocean

import (
	"context"

	"github.com/rancher/kontainer-engine/types"
)

type DigitalOceanDriver struct {
	driverCapabilities types.Capabilities
}

func (*DigitalOceanDriver) GetDriverOptions(ctx context.Context) (*types.DriverFlags, error) {

	return nil, nil
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

func test() {

}
