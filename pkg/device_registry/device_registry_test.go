package device_registry

import (
	"github.com/chacal/thread-mgmt-server/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegistry_GetOrCreate(t *testing.T) {
	reg, err := Open(test.Tempfile())
	require.NoError(t, err)
	defer reg.Close()

	dev, err := reg.GetOrCreate("12345")
	require.NoError(t, err)
	assert.Equal(t, Device{}, dev)

	d := update(t, reg, "12345", Device{"D100", 5000})

	dev, err = reg.GetOrCreate("12345")
	require.NoError(t, err)
	assert.Equal(t, d, dev)
}

func TestRegistry_GetAll(t *testing.T) {
	reg, err := Open(test.Tempfile())
	require.NoError(t, err)
	defer reg.Close()

	expected := map[string]Device{}
	assert.Equal(t, expected, getAll(t, reg))

	expected["12345"] = update(t, reg, "12345", Device{"D100", 5000})
	assert.Equal(t, expected, getAll(t, reg))

	expected["AABBCCDD"] = update(t, reg, "AABBCCDD", Device{"D100", 5000})
	assert.Equal(t, expected, getAll(t, reg))
}

func update(t *testing.T, reg *Registry, id string, d Device) Device {
	err := reg.Update(id, d)
	require.NoError(t, err)
	return d
}

func getAll(t *testing.T, reg *Registry) map[string]Device {
	devices, err := reg.GetAll()
	require.NoError(t, err)
	return devices
}
