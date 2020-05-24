package digitalocean

import (
	"context"
	"github.com/digitalocean/godo"
	"github.com/rancher/kontainer-engine/store"
	"github.com/rancher/kontainer-engine/types"
	"github.com/ribeiro-rodrigo/kontainer-engine-driver-digitalocean/digitalocean/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

/*************** Defining Mocks *************/
type OptionsBuilderMock struct{
	mock.Mock
	buildCreateOptionsMock func()*types.DriverFlags
}

func (m *OptionsBuilderMock) BuildCreateOptions() *types.DriverFlags{
	m.Called()
	return m.buildCreateOptionsMock()
}

type StateBuilderMock struct {
	mock.Mock
	buildStateFromOptsMock func (driverOptions *types.DriverOptions) (state.State, error)
	buildStateFromClusterInfo func (clusterInfo *types.ClusterInfo)(state.State,error)
}

func (m *StateBuilderMock) BuildStateFromOpts(driverOptions *types.DriverOptions) (state.State, error){
	m.Called(driverOptions)
	return m.buildStateFromOptsMock(driverOptions)
}

func (m *StateBuilderMock) BuildStateFromClusterInfo(clusterInfo *types.ClusterInfo)(state.State,error){
	m.Called(clusterInfo)
	return m.buildStateFromClusterInfo(clusterInfo)
}

type DigitalOceanMock struct {
	mock.Mock
	createClusterMock func(ctx context.Context, state state.State) (string, error)
	deleteClusterMock func (ctx context.Context, clusterID string)error
	getNodeCountMock func (ctx context.Context, clusterID string) (int,error)
	getKubeConfigMock func (clusterID string)(*store.KubeConfig,error)
	waitClusterCreated func (ctx context.Context, clusterID string)error
	waitClusterDeleted func (ctx context.Context, clusterID string)error
}

func (m *DigitalOceanMock) CreateCluster(ctx context.Context, state state.State) (string, error){
	m.Called(ctx,state)
	return m.createClusterMock(ctx, state)
}

func (m *DigitalOceanMock) DeleteCluster(ctx context.Context, clusterID string)error {
	m.Called(ctx,clusterID)
	return m.deleteClusterMock(ctx,clusterID)
}
func (m *DigitalOceanMock) GetNodeCount(ctx context.Context, clusterID string) (int,error){
	m.Called(ctx,clusterID)
	return m.getNodeCountMock(ctx,clusterID)
}
func (m *DigitalOceanMock) GetKubeConfig(clusterID string)(*store.KubeConfig,error){
	m.Called(clusterID)
	return m.getKubeConfigMock(clusterID)
}
func (m *DigitalOceanMock) WaitClusterCreated(ctx context.Context, clusterID string)error{
	m.Called(ctx,clusterID)
	return m.waitClusterCreated(ctx,clusterID)
}
func (m *DigitalOceanMock) WaitClusterDeleted(ctx context.Context, clusterID string)error{
	m.Called(ctx,clusterID)
	return m.WaitClusterDeleted(ctx,clusterID)
}

/*************** Defining Tests *************/

func TestGetDriverCreateOptions(t *testing.T) {

	returnBuildCreateOptions := &types.DriverFlags{}

	builderMock := &OptionsBuilderMock{
		buildCreateOptionsMock: func() *types.DriverFlags {
			return returnBuildCreateOptions
		},
	}

	driver := Driver{
		optionsBuilder: builderMock,
	}

	builderMock.On("BuildCreateOptions").Return(returnBuildCreateOptions)

	flags, err := driver.GetDriverCreateOptions(context.TODO())

	assert.NoError(t, err, "GetDriverCreateOptions not error")
	assert.Equal(t, returnBuildCreateOptions, flags, "Flags mock equals")
	builderMock.AssertExpectations(t)
}

func TestDriverCreate(t *testing.T) {



}
