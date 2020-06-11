package digitalocean

import (
	"context"
	"errors"
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
	buildUpdateOptionsMock func()*types.DriverFlags
}

func (m *OptionsBuilderMock) BuildCreateOptions() *types.DriverFlags{
	m.Called()
	return m.buildCreateOptionsMock()
}

func (m *OptionsBuilderMock) BuildUpdateOptions() *types.DriverFlags{
	m.Called()
	return m.buildUpdateOptionsMock()
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
	getKubernetesClusterVersionMock func(ctx context.Context, clusterID string)(string, error)
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

func (m *DigitalOceanMock) GetKubernetesClusterVersion(ctx context.Context, clusterID string)(string,error){
	m.Called(ctx,clusterID)
	return m.getKubernetesClusterVersionMock(ctx, clusterID)
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
		NodePool: state.NodePool{
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
	returnState := state.State{
		DisplayName: "cluster-test",
		Name:        "my-cluster",
		RegionSlug:  "1.17.5-do.0",
		NodePool: state.NodePool{
			Name:  "node-pool-1",
			Size:  "s-2vcpu-2gb",
			Count: 5,
		},
	}

	stateBuilderMock := &StateBuilderMock{
		buildStateFromOptsMock: func(_ *types.DriverOptions) (state.State, error) {
			return returnState, nil
		},
	}

	driver := Driver{
		stateBuilder: stateBuilderMock,
	}

	options :=  &types.DriverOptions{}

	stateBuilderMock.On("BuildStateFromOpts",options).Return(returnState, nil)

	_, err := driver.Create(context.TODO(), &types.DriverOptions{}, &types.ClusterInfo{})

	stateBuilderMock.AssertExpectations(t)
	assert.Error(t,err, "Error in create cluster: not token")

}

func TestDriverCreateErrorInDigitalOceanServiceCreate(t *testing.T){
	returnState := state.State{
		Token:       "a405b7bd3e0d6193f605368102a2deafe4067ed542c92166c6d772fe7e2df019",
		DisplayName: "cluster-test",
		Name:        "my-cluster",
		RegionSlug:  "1.17.5-do.0",
		NodePool: state.NodePool{
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

	digitalOceanMock := &DigitalOceanMock{
		createClusterMock: func(_ context.Context, _ state.State) (string, error) {
			return "", errors.New("error in create cluster")
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
	digitalOceanMock.On("CreateCluster", ctx, returnState).Return("",nil)

	_, err := driver.Create(ctx, options , clusterInfo)

	digitalOceanMock.AssertExpectations(t)
	stateBuilderMock.AssertExpectations(t)

	assert.Error(t, err, "Error in create cluster")

}

func TestDriverCreateErrorInWaitClusterCreated(t *testing.T){
	returnState := state.State{
		Token:       "a405b7bd3e0d6193f605368102a2deafe4067ed542c92166c6d772fe7e2df019",
		DisplayName: "cluster-test",
		Name:        "my-cluster",
		RegionSlug:  "1.17.5-do.0",
		NodePool: state.NodePool{
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
			return errors.New("error in wait cluster")
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

	_, err := driver.Create(ctx, options , clusterInfo)

	digitalOceanMock.AssertExpectations(t)
	stateBuilderMock.AssertExpectations(t)

	assert.Error(t, err, "error in wait cluster created")

}

func TestRemoveCluster(t *testing.T){
	returnState := state.State{
		Token:       "a405b7bd3e0d6193f605368102a2deafe4067ed542c92166c6d772fe7e2df019",
		DisplayName: "cluster-test",
		Name:        "my-cluster",
		RegionSlug:  "1.17.5-do.0",
		NodePool: state.NodePool{
			Name:  "node-pool-1",
			Size:  "s-2vcpu-2gb",
			Count: 5,
		},
	}

	stateBuilderMock := &StateBuilderMock{
		buildStateFromClusterInfo: func(_ *types.ClusterInfo) (state.State, error) {
			return returnState, nil
		},
	}

	digitalOceanMock := DigitalOceanMock{
		deleteClusterMock: func(_ context.Context, _ string) error {
			return nil
		},
		waitClusterDeleted: func(_ context.Context, _ string) error {
			return nil
		},
	}

	digitalOceanFactory := func(token string)service.DigitalOcean{
		return &digitalOceanMock
	}

	driver := Driver{
		digitalOceanFactory: digitalOceanFactory,
		stateBuilder: stateBuilderMock,
	}

	clusterInfo := &types.ClusterInfo{}
	ctx := context.TODO()

	stateBuilderMock.On("BuildStateFromClusterInfo", clusterInfo).Return(returnState)
	digitalOceanMock.On("DeleteCluster",ctx, returnState.ClusterID).Return(nil)
	digitalOceanMock.On("WaitClusterDeleted",ctx, returnState.ClusterID).Return(nil)

	err := driver.Remove(ctx, clusterInfo)

	stateBuilderMock.AssertExpectations(t)
	digitalOceanMock.AssertExpectations(t)

	assert.NoError(t, err, "Not error in remove cluster")
}

func TestRemoveClusterErrorInBuildState(t *testing.T){

	returnState := state.State{}
	returnError := errors.New("error in build state")

	stateBuilderMock := &StateBuilderMock{
		buildStateFromClusterInfo: func(_ *types.ClusterInfo) (state.State, error) {
			return returnState, returnError
		},
	}

	driver := Driver{
		stateBuilder: stateBuilderMock,
	}

	clusterInfo := &types.ClusterInfo{}
	ctx := context.TODO()

	stateBuilderMock.On("BuildStateFromClusterInfo", clusterInfo).Return(returnError)

	err := driver.Remove(ctx, clusterInfo)

	stateBuilderMock.AssertExpectations(t)

	assert.Error(t, err, "Error in remove cluster")
}

func TestRemoveClusterErrorInDigitalOceanDelete(t *testing.T){

	returnState := state.State{
		Token:       "a405b7bd3e0d6193f605368102a2deafe4067ed542c92166c6d772fe7e2df019",
		DisplayName: "cluster-test",
		Name:        "my-cluster",
		RegionSlug:  "1.17.5-do.0",
		NodePool: state.NodePool{
			Name:  "node-pool-1",
			Size:  "s-2vcpu-2gb",
			Count: 5,
		},
	}

	returnError := errors.New("error in delete cluster")

	stateBuilderMock := &StateBuilderMock{
		buildStateFromClusterInfo: func(_ *types.ClusterInfo) (state.State, error) {
			return returnState, nil
		},
	}

	digitalOceanMock := DigitalOceanMock{
		deleteClusterMock: func(_ context.Context, _ string) error {
			return returnError
		},
	}

	digitalOceanFactory := func(token string)service.DigitalOcean{
		return &digitalOceanMock
	}

	driver := Driver{
		digitalOceanFactory: digitalOceanFactory,
		stateBuilder: stateBuilderMock,
	}

	clusterInfo := &types.ClusterInfo{}
	ctx := context.TODO()

	stateBuilderMock.On("BuildStateFromClusterInfo", clusterInfo).Return(returnState)
	digitalOceanMock.On("DeleteCluster",ctx, returnState.ClusterID).Return(returnError)

	err := driver.Remove(ctx, clusterInfo)

	stateBuilderMock.AssertExpectations(t)
	digitalOceanMock.AssertExpectations(t)

	assert.Error(t, err, "Error in remove cluster")
}

func TestRemoveClusterErrorInWaitDeleted(t *testing.T){

	returnState := state.State{
		Token:       "a405b7bd3e0d6193f605368102a2deafe4067ed542c92166c6d772fe7e2df019",
		DisplayName: "cluster-test",
		Name:        "my-cluster",
		RegionSlug:  "1.17.5-do.0",
		NodePool: state.NodePool{
			Name:  "node-pool-1",
			Size:  "s-2vcpu-2gb",
			Count: 5,
		},
	}

	stateBuilderMock := &StateBuilderMock{
		buildStateFromClusterInfo: func(_ *types.ClusterInfo) (state.State, error) {
			return returnState, nil
		},
	}

	returnError := errors.New("error in waht cluster deleted")

	digitalOceanMock := DigitalOceanMock{
		deleteClusterMock: func(_ context.Context, _ string) error {
			return nil
		},
		waitClusterDeleted: func(_ context.Context, _ string) error {
			return returnError
		},
	}

	digitalOceanFactory := func(token string)service.DigitalOcean{
		return &digitalOceanMock
	}

	driver := Driver{
		digitalOceanFactory: digitalOceanFactory,
		stateBuilder: stateBuilderMock,
	}

	clusterInfo := &types.ClusterInfo{}
	ctx := context.TODO()

	stateBuilderMock.On("BuildStateFromClusterInfo", clusterInfo).Return(returnState)
	digitalOceanMock.On("DeleteCluster",ctx, returnState.ClusterID).Return(nil)
	digitalOceanMock.On("WaitClusterDeleted",ctx, returnState.ClusterID).Return(returnError)

	err := driver.Remove(ctx, clusterInfo)

	stateBuilderMock.AssertExpectations(t)
	digitalOceanMock.AssertExpectations(t)

	assert.Error(t, err, "Error in remove cluster")
}

func TestGetClusterSize(t *testing.T){

	nodeCount := 5
	returnClusterID := "abcd"

	returnState := state.State{
		Token:       "a405b7bd3e0d6193f605368102a2deafe4067ed542c92166c6d772fe7e2df019",
		DisplayName: "cluster-test",
		Name:        "my-cluster",
		ClusterID: returnClusterID,
		RegionSlug:  "1.17.5-do.0",
		NodePool: state.NodePool{
			Name:  "node-pool-1",
			Size:  "s-2vcpu-2gb",
			Count: nodeCount,
		},
	}

	stateBuilderMock := &StateBuilderMock{
		buildStateFromClusterInfo: func(_ *types.ClusterInfo) (state.State, error) {
			return returnState, nil
		},
	}

	digitalOceanMock := &DigitalOceanMock{
		getNodeCountMock: func(_ context.Context, _ string) (int, error) {
			return nodeCount, nil
		},
	}

	driver := Driver{
		stateBuilder: stateBuilderMock,
		digitalOceanFactory: func(token string) service.DigitalOcean {return digitalOceanMock},
	}

	ctx := context.TODO()
	clusterInfo := &types.ClusterInfo{}

	stateBuilderMock.On("BuildStateFromClusterInfo",clusterInfo).Return(returnState)
	digitalOceanMock.On("GetNodeCount", ctx, returnClusterID).Return(nodeCount,nil)

	clusterSize, err := driver.GetClusterSize(ctx, clusterInfo)

	stateBuilderMock.AssertExpectations(t)
	digitalOceanMock.AssertExpectations(t)

	assert.NoError(t, err, "Not error in get cluster size")
	assert.Equal(t, int64(nodeCount), clusterSize.Count, "NodeCount equals")

}

func TestGetVersion(t *testing.T){

}
