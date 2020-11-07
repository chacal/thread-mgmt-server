package device_registry

import (
	"github.com/chacal/thread-mgmt-server/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegistry(t *testing.T) {
	reg, err := Open(test.Tempfile())
	require.NoError(t, err)
	defer reg.Close()

	dev, err := reg.GetOrCreate("12345")
	require.NoError(t, err)
	assert.Equal(t, Device{}, dev)

	d := Device{"D100", 5000}
	err = reg.Update("12345", d)
	require.NoError(t, err)

	dev, err = reg.GetOrCreate("12345")
	require.NoError(t, err)
	assert.Equal(t, d, dev)
}
