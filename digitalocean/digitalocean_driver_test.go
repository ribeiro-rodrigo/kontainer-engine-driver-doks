package digitalocean

import (
	"context"
	"errors"
	"github.com/digitalocean/godo"
	"github.com/rancher/kontainer-engine/store"
	"github.com/rancher/kontainer-engine/types"
	"github.com/ribeiro-rodrigo/kontainer-engine-driver-digitalocean/digitalocean/service"
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
	return m.waitClusterDeleted(ctx,clusterID)
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

	returnState := state.State{
		Token:       "a405b7bd3e0d6193f605368102a2deafe4067ed542c92166c6d772fe7e2df019",
		DisplayName: "cluster-test",
		Name:        "my-cluster",
		RegionSlug:  "1.17.5-do.0",
		NodePool: &godo.KubernetesNodePoolCreateRequest{
			Name:  "node-pool-1",
			Size:  "s-2vcpu-2gb",
			Count: 5,
		},
	}

	stateBuilderMock := &StateBuilderMock{
		buildStateFromOptsMock: func(do *types.DriverOptions) (state.State, error) {
			return returnState,nil
		},
	}

	returnClusterID := "abcd"

	digitalOceanMock := &DigitalOceanMock{
		createClusterMock: func(_ context.Context, _ state.State) (string, error) {
			return returnClusterID, nil
		},
		waitClusterCreated: func(_ context.Context, _ string) error {
			return nil
		},
	}

	driver := Driver{
		stateBuilder: stateBuilderMock,
		digitalOceanFactory: func(token string) service.DigitalOcean {return digitalOceanMock},
	}

	options := &types.DriverOptions{}
	ctx := context.TODO()
	clusterInfo := &types.ClusterInfo{}

	stateBuilderMock.On("BuildStateFromOpts",options).Return(returnState)
	digitalOceanMock.On("CreateCluster", ctx, returnState).Return(returnClusterID,nil)
	digitalOceanMock.On("WaitClusterCreated",ctx,returnClusterID).Return(nil)

	info, err := driver.Create(ctx, options , clusterInfo)

	digitalOceanMock.AssertExpectations(t)
	stateBuilderMock.AssertExpectations(t)

	assert.NoError(t, err, "Not error in create cluster")

	_, ok  := info.Metadata["state"]

	assert.True(t, ok, "State serialized in info")
}

func TestDriverCreateErrorInBuildStateFromOpts(t *testing.T) {

	returnState := state.State{}
	returnError := errors.New("error")

	stateBuilderMock := &StateBuilderMock{
		buildStateFromOptsMock: func(do *types.DriverOptions) (state.State, error) {
			return returnState, returnError
		},
	}

	driver := Driver{
		stateBuilder: stateBuilderMock,
	}

	options := &types.DriverOptions{}

	stateBuilderMock.On("BuildStateFromOpts",options).Return(returnState, returnError)

	_, err := driver.Create(context.TODO(), options ,&types.ClusterInfo{})

	stateBuilderMock.AssertExpectations(t)
	assert.Error(t, err, "Error in create cluster")
}

func TestDriverCreateWithoutToken(t *testing.T){

}

func TestDriverCreateErrorInDigitalOceanServiceCreate(t *testing.T){

}

func TestDriverCreateErrorInSaveState(t *testing.T){

}

func TestDriverCreateErrorInWaitClusterCreated(t *testing.T){

}

