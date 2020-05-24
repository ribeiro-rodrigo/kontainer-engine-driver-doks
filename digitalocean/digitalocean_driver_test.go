package digitalocean

import (
	"context"
	"github.com/rancher/kontainer-engine/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type OptionsBuilderMock struct{
	mock.Mock
	buildCreateOptionsMock func()*types.DriverFlags
}

func (m *OptionsBuilderMock) BuildCreateOptions() *types.DriverFlags{
	m.Called()
	return m.buildCreateOptionsMock()
}

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

