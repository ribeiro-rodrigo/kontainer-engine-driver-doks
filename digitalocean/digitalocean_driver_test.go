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
}

var returnBuildCreateOptionsMock = &types.DriverFlags{}

func (m *OptionsBuilderMock) BuildCreateOptions() *types.DriverFlags{
	m.Called()
	return returnBuildCreateOptionsMock
}

func TestGetDriverCreateOptions(t *testing.T) {
	builderMock := new(OptionsBuilderMock)

	driver := Driver{
		optionsBuilder: builderMock,
	}

	builderMock.On("BuildCreateOptions").Return(returnBuildCreateOptionsMock)

	flags, err := driver.GetDriverCreateOptions(context.TODO())

	assert.NoError(t, err, "GetDriverCreateOptions not error")
	assert.Equal(t, returnBuildCreateOptionsMock, flags, "Flags mock equals")
	builderMock.AssertExpectations(t)
}