package device_registry

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

var ip = net.ParseIP("ffff::1")

func TestRegistry_CRUD(t *testing.T) {
	reg := CreateTestRegistry(t)

	dev, err := reg.GetOrCreate("12345")
	require.NoError(t, err)
	assert.Equal(t, Device{}, dev)

	d := update(t, reg, "12345", Device{"D100", -4, 5000, []net.IP{ip}})

	dev, err = reg.GetOrCreate("12345")
	require.NoError(t, err)
	assert.Equal(t, d, dev)

	err = reg.Delete("12345")
	require.NoError(t, err)

	dev, err = reg.GetOrCreate("12345")
	require.NoError(t, err)
	assert.Equal(t, Device{}, dev) // Should create new empty Device here

	err = reg.UpdateAddresses("12345", []net.IP{ip})
	require.NoError(t, err)

	dev, err = reg.GetOrCreate("12345")
	require.NoError(t, err)
	assert.Equal(t, Device{Addresses: []net.IP{ip}}, dev)

	err = reg.UpdateAddresses("12345", []net.IP{})
	require.NoError(t, err)

	dev, err = reg.GetOrCreate("12345")
	require.NoError(t, err)
	assert.Equal(t, Device{}, dev)
}

func TestRegistry_GetAll(t *testing.T) {
	reg := CreateTestRegistry(t)

	expected := map[string]Device{}
	assert.Equal(t, expected, getAll(t, reg))

	expected["12345"] = update(t, reg, "12345", Device{"D100", -4, 5000, []net.IP{ip}})
	assert.Equal(t, expected, getAll(t, reg))

	expected["AABBCCDD"] = update(t, reg, "AABBCCDD", Device{"D100", -4, 5000, nil})
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
