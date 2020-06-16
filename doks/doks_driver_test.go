package doks

import (
	"context"
	"errors"
	"github.com/rancher/kontainer-engine/store"
	"github.com/rancher/kontainer-engine/types"
	"github.com/ribeiro-rodrigo/kontainer-engine-driver-doks/doks/service"
	"github.com/ribeiro-rodrigo/kontainer-engine-driver-doks/doks/state"
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
	buildStatesFromOptsMock func (driverOptions *types.DriverOptions) (state.Cluster, state.NodePool ,error)
	buildStateFromClusterInfo func (clusterInfo *types.ClusterInfo)(state.Cluster,error)
}

func (m *StateBuilderMock) BuildStatesFromOpts(driverOptions *types.DriverOptions) (state.Cluster, state.NodePool , error){
	m.Called(driverOptions)
	return m.buildStatesFromOptsMock(driverOptions)
}

func (m *StateBuilderMock) BuildClusterStateFromClusterInfo(clusterInfo *types.ClusterInfo)(state.Cluster,error){
	m.Called(clusterInfo)
	return m.buildStateFromClusterInfo(clusterInfo)
}

type DigitalOceanMock struct {
	mock.Mock
	createClusterMock func(ctx context.Context, state state.Cluster, pool state.NodePool ) (string, string, error)
	deleteClusterMock func (ctx context.Context, clusterID string)error
	getKubeConfigMock func (clusterID string)(*store.KubeConfig,error)
	waitClusterCreated func (ctx context.Context, clusterID string)error
	waitClusterDeleted func (ctx context.Context, clusterID string)error
	getKubernetesClusterVersionMock func(ctx context.Context, clusterID string)(string, error)
	upgradeKubernetesVersionMock func(ctx context.Context, clusterID, version string)error
	updateNodePoolMock func (ctx context.Context, clusterID, nodePoolID string, nodePool state.NodePool) error
	getNodePoolMock func(ctx context.Context, clusterID, nodePoolID string) (*state.NodePool,error)
}

func (m *DigitalOceanMock) CreateCluster(ctx context.Context, clusterState state.Cluster,
	pool state.NodePool) (string, string, error){
	m.Called(ctx, clusterState, pool)
	return m.createClusterMock(ctx, clusterState, pool)
}

func (m *DigitalOceanMock) DeleteCluster(ctx context.Context, clusterID string)error {
	m.Called(ctx,clusterID)
	return m.deleteClusterMock(ctx,clusterID)
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

func (m *DigitalOceanMock) UpgradeKubernetesVersion(ctx context.Context, clusterID, version string)error{
	m.Called(ctx, clusterID, version)
	return m.upgradeKubernetesVersionMock(ctx, clusterID, version)
}

func (m *DigitalOceanMock) UpdateNodePool(ctx context.Context, clusterID, nodePoolID string,
	nodePool state.NodePool) error{
	m.Called(ctx, clusterID)
	return m.updateNodePoolMock(ctx, clusterID, nodePoolID, nodePool)
}

func (m *DigitalOceanMock) GetNodePool(ctx context.Context, clusterID, nodePoolID string) (*state.NodePool,error){
	m.Called(ctx, clusterID, nodePoolID)
	return m.getNodePoolMock(ctx,clusterID,nodePoolID)
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

	returnNodePoolState := state.NodePool{
		Name:  "node-pool-1",
		Size:  "s-2vcpu-2gb",
		Count: 5,
	}

	returnClusterState := state.Cluster{
		Token:       "a405b7bd3e0d6193f605368102a2deafe4067ed542c92166c6d772fe7e2df019",
		DisplayName: "cluster-test",
		Name:        "my-cluster",
		RegionSlug:  "1.17.5-do.0",
	}

	stateBuilderMock := &StateBuilderMock{
		buildStatesFromOptsMock: func(do *types.DriverOptions) (state.Cluster, state.NodePool, error) {
			return returnClusterState, returnNodePoolState, nil
		},
	}

	returnClusterID := "abcd"
	returnNodePoolID := "zzz"

	digitalOceanMock := &DigitalOceanMock{
		createClusterMock: func(_ context.Context, _ state.Cluster, _ state.NodePool) (string, string ,error) {
			return returnClusterID, returnNodePoolID, nil
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

	stateBuilderMock.On("BuildStatesFromOpts",options).Return(returnClusterState,
		returnNodePoolState)

	digitalOceanMock.On("CreateCluster", ctx, returnClusterState,
		returnNodePoolState).Return(returnClusterID,returnNodePoolID,nil)

	digitalOceanMock.On("WaitClusterCreated",ctx,returnClusterID).Return(nil)

	info, err := driver.Create(ctx, options , nil)

	digitalOceanMock.AssertExpectations(t)
	stateBuilderMock.AssertExpectations(t)

	assert.NoError(t, err, "Not error in create cluster")

	_, ok  := info.Metadata["state"]

	assert.True(t, ok, "Cluster serialized in info")
}

func TestDriverCreateErrorInBuildStateFromOpts(t *testing.T) {

	returnClusterState := state.Cluster{}
	returnNodePoolState := state.NodePool{}
	returnError := errors.New("error")

	stateBuilderMock := &StateBuilderMock{
		buildStatesFromOptsMock: func(do *types.DriverOptions) (state.Cluster, state.NodePool ,error) {
			return returnClusterState, returnNodePoolState,returnError
		},
	}

	driver := Driver{
		stateBuilder: stateBuilderMock,
	}

	options := &types.DriverOptions{}

	stateBuilderMock.On("BuildStatesFromOpts",
		options).Return(returnClusterState, returnNodePoolState, returnError)

	_, err := driver.Create(context.TODO(), options ,nil)

	stateBuilderMock.AssertExpectations(t)
	assert.Error(t, err, "Error in create cluster")
}

func TestDriverCreateWithoutToken(t *testing.T){

	returnNodePoolState := state.NodePool{
		Name:  "node-pool-1",
		Size:  "s-2vcpu-2gb",
		Count: 5,
	}

	returnClusterState := state.Cluster{
		DisplayName: "cluster-test",
		Name:        "my-cluster",
		RegionSlug:  "1.17.5-do.0",
	}

	stateBuilderMock := &StateBuilderMock{
		buildStatesFromOptsMock: func(_ *types.DriverOptions) (state.Cluster, state.NodePool,  error) {
			return returnClusterState, returnNodePoolState, nil
		},
	}

	driver := Driver{
		stateBuilder: stateBuilderMock,
	}

	options :=  &types.DriverOptions{}

	stateBuilderMock.On("BuildStatesFromOpts",
		options).Return(returnClusterState, returnNodePoolState, nil)

	_, err := driver.Create(context.TODO(), &types.DriverOptions{}, nil)

	stateBuilderMock.AssertExpectations(t)
	assert.Error(t,err, "Error in create cluster: not token")

}

func TestDriverCreateErrorInDigitalOceanServiceCreate(t *testing.T){

	returnNodePoolState := state.NodePool{
		Name:  "node-pool-1",
		Size:  "s-2vcpu-2gb",
		Count: 5,
	}

	returnClusterState := state.Cluster{
		Token:       "a405b7bd3e0d6193f605368102a2deafe4067ed542c92166c6d772fe7e2df019",
		DisplayName: "cluster-test",
		Name:        "my-cluster",
		RegionSlug:  "1.17.5-do.0",
	}

	stateBuilderMock := &StateBuilderMock{
		buildStatesFromOptsMock: func(do *types.DriverOptions) (state.Cluster, state.NodePool, error) {
			return returnClusterState, returnNodePoolState,nil
		},
	}

	digitalOceanMock := &DigitalOceanMock{
		createClusterMock: func(_ context.Context, _ state.Cluster, _ state.NodePool) (string, string, error) {
			return "", "", errors.New("error in create cluster")
		},
	}

	driver := Driver{
		stateBuilder: stateBuilderMock,
		digitalOceanFactory: func(token string) service.DigitalOcean {return digitalOceanMock},
	}

	options := &types.DriverOptions{}
	ctx := context.TODO()

	stateBuilderMock.On("BuildStatesFromOpts",
		options).Return(returnClusterState, returnNodePoolState)

	digitalOceanMock.On("CreateCluster", ctx, returnClusterState,
		returnNodePoolState).Return("",nil)

	_, err := driver.Create(ctx, options , nil)

	digitalOceanMock.AssertExpectations(t)
	stateBuilderMock.AssertExpectations(t)

	assert.Error(t, err, "Error in create cluster")

}

func TestDriverCreateErrorInWaitClusterCreated(t *testing.T){

	returnNodePoolState := state.NodePool{
		Name:  "node-pool-1",
		Size:  "s-2vcpu-2gb",
		Count: 5,
	}

	returnClusterState := state.Cluster{
		Token:       "a405b7bd3e0d6193f605368102a2deafe4067ed542c92166c6d772fe7e2df019",
		DisplayName: "cluster-test",
		Name:        "my-cluster",
		RegionSlug:  "1.17.5-do.0",
	}

	stateBuilderMock := &StateBuilderMock{
		buildStatesFromOptsMock: func(do *types.DriverOptions) (state.Cluster, state.NodePool ,error) {
			return returnClusterState, returnNodePoolState,nil
		},
	}

	returnClusterID := "abcd"
	returnNodePoolID := "zzz"

	digitalOceanMock := &DigitalOceanMock{
		createClusterMock: func(_ context.Context, _ state.Cluster, _ state.NodePool) (string, string, error) {
			return returnClusterID, returnNodePoolID, nil
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

	stateBuilderMock.On("BuildStatesFromOpts",
		options).Return(returnClusterState, returnNodePoolState)

	digitalOceanMock.On("CreateCluster", ctx,
		returnClusterState, returnNodePoolState).Return(returnClusterID,returnNodePoolID,nil)

	digitalOceanMock.On("WaitClusterCreated",ctx,returnClusterID).Return(nil)

	_, err := driver.Create(ctx, options , nil)

	digitalOceanMock.AssertExpectations(t)
	stateBuilderMock.AssertExpectations(t)

	assert.Error(t, err, "error in wait cluster created")

}

func TestRemoveCluster(t *testing.T){

	returnState := state.Cluster{
		Token:       "a405b7bd3e0d6193f605368102a2deafe4067ed542c92166c6d772fe7e2df019",
		DisplayName: "cluster-test",
		Name:        "my-cluster",
		RegionSlug:  "1.17.5-do.0",
	}

	stateBuilderMock := &StateBuilderMock{
		buildStateFromClusterInfo: func(_ *types.ClusterInfo) (state.Cluster, error) {
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

	stateBuilderMock.On("BuildClusterStateFromClusterInfo", clusterInfo).Return(returnState)
	digitalOceanMock.On("DeleteCluster",ctx, returnState.ClusterID).Return(nil)
	digitalOceanMock.On("WaitClusterDeleted",ctx, returnState.ClusterID).Return(nil)

	err := driver.Remove(ctx, clusterInfo)

	stateBuilderMock.AssertExpectations(t)
	digitalOceanMock.AssertExpectations(t)

	assert.NoError(t, err, "Not error in remove cluster")
}

func TestRemoveClusterErrorInBuildState(t *testing.T){

	returnState := state.Cluster{}
	returnError := errors.New("error in build state")

	stateBuilderMock := &StateBuilderMock{
		buildStateFromClusterInfo: func(_ *types.ClusterInfo) (state.Cluster, error) {
			return returnState, returnError
		},
	}

	driver := Driver{
		stateBuilder: stateBuilderMock,
	}

	clusterInfo := &types.ClusterInfo{}
	ctx := context.TODO()

	stateBuilderMock.On("BuildClusterStateFromClusterInfo", clusterInfo).Return(returnError)

	err := driver.Remove(ctx, clusterInfo)

	stateBuilderMock.AssertExpectations(t)

	assert.Error(t, err, "Error in remove cluster")
}

func TestRemoveClusterErrorInDigitalOceanDelete(t *testing.T){

	returnState := state.Cluster{
		Token:       "a405b7bd3e0d6193f605368102a2deafe4067ed542c92166c6d772fe7e2df019",
		DisplayName: "cluster-test",
		Name:        "my-cluster",
		RegionSlug:  "1.17.5-do.0",
	}

	returnError := errors.New("error in delete cluster")

	stateBuilderMock := &StateBuilderMock{
		buildStateFromClusterInfo: func(_ *types.ClusterInfo) (state.Cluster, error) {
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

	stateBuilderMock.On("BuildClusterStateFromClusterInfo", clusterInfo).Return(returnState)
	digitalOceanMock.On("DeleteCluster",ctx, returnState.ClusterID).Return(returnError)

	err := driver.Remove(ctx, clusterInfo)

	stateBuilderMock.AssertExpectations(t)
	digitalOceanMock.AssertExpectations(t)

	assert.Error(t, err, "Error in remove cluster")
}

func TestRemoveClusterErrorInWaitDeleted(t *testing.T){

	returnState := state.Cluster{
		Token:       "a405b7bd3e0d6193f605368102a2deafe4067ed542c92166c6d772fe7e2df019",
		DisplayName: "cluster-test",
		Name:        "my-cluster",
		RegionSlug:  "1.17.5-do.0",
	}

	stateBuilderMock := &StateBuilderMock{
		buildStateFromClusterInfo: func(_ *types.ClusterInfo) (state.Cluster, error) {
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

	stateBuilderMock.On("BuildClusterStateFromClusterInfo", clusterInfo).Return(returnState)
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
	returnNodePoolID := "aaas"

	returnNodePool := &state.NodePool{
		Name:  "node-pool-1",
		Size:  "s-2vcpu-2gb",
		Count: nodeCount,
	}

	returnState := state.Cluster{
		Token:       "a405b7bd3e0d6193f605368102a2deafe4067ed542c92166c6d772fe7e2df019",
		DisplayName: "cluster-test",
		Name:        "my-cluster",
		ClusterID: returnClusterID,
		RegionSlug:  "1.17.5-do.0",
		NodePoolID: returnNodePoolID,
	}

	stateBuilderMock := &StateBuilderMock{
		buildStateFromClusterInfo: func(_ *types.ClusterInfo) (state.Cluster, error) {
			return returnState, nil
		},
	}

	digitalOceanMock := &DigitalOceanMock{
		getNodePoolMock: func(_ context.Context, _, _ string) (*state.NodePool, error) {
			return returnNodePool, nil
		},
	}

	driver := Driver{
		stateBuilder: stateBuilderMock,
		digitalOceanFactory: func(token string) service.DigitalOcean {return digitalOceanMock},
	}

	ctx := context.TODO()
	clusterInfo := &types.ClusterInfo{}

	stateBuilderMock.On("BuildClusterStateFromClusterInfo",clusterInfo).Return(returnState)

	digitalOceanMock.On("GetNodePool", ctx, returnClusterID,
		returnNodePoolID).Return(returnNodePool,nil)

	clusterSize, err := driver.GetClusterSize(ctx, clusterInfo)

	stateBuilderMock.AssertExpectations(t)
	digitalOceanMock.AssertExpectations(t)

	assert.NoError(t, err, "Not error in get cluster size")
	assert.Equal(t, int64(nodeCount), clusterSize.Count, "NodeCount equals")

}

func TestGetVersion(t *testing.T){

}

func TestSetVersion(t *testing.T){

}

func TestSetClusterSize(t *testing.T){

}
