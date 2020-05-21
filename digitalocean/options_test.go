package digitalocean

import (
	"testing"

	"github.com/rancher/kontainer-engine/types"
	"github.com/stretchr/testify/assert"
)

func TestGetCreateOptions(t *testing.T) {
	options := getCreateOptions()

	DisplayNameFlag, ok := options.Options["display-name"]

	assert.True(t, ok, "DisplayName flag is present")
	assert.Equal(t, types.StringType, DisplayNameFlag.GetType(), "DisplayName type is string")
}
