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

	defaults, err := reg.GetDefaults("12345")
	require.NoError(t, err)
	assert.Equal(t, (*Defaults)(nil), defaults)

	dev, err := reg.GetOrCreate("12345")
	require.NoError(t, err)
	assert.Equal(t, Device{}, dev)

	expectedDefaults := updateDefaults(t, reg, "12345", Defaults{"D100", -4, 500})

	dev, err = reg.GetOrCreate("12345")
	require.NoError(t, err)
	assert.Equal(t, Device{Defaults: expectedDefaults}, dev)

	defaults, err = reg.GetDefaults("12345")
	require.NoError(t, err)
	assert.Equal(t, &expectedDefaults, defaults)

	expectedState := updateState(t, reg, "12345", State{[]net.IP{ip}, 2970})

	dev, err = reg.GetOrCreate("12345")
	require.NoError(t, err)
	assert.Equal(t, Device{Defaults: expectedDefaults, State: expectedState}, dev)

	expectedConfig := updateConfig(t, reg, "12345", Config{ip})

	dev, err = reg.GetOrCreate("12345")
	require.NoError(t, err)
	assert.Equal(t, Device{expectedDefaults, expectedState, expectedConfig}, dev)

	err = reg.DeleteDevice("12345")
	require.NoError(t, err)

	dev, err = reg.GetOrCreate("12345")
	require.NoError(t, err)
	assert.Equal(t, Device{}, dev) // Should create new empty Device here

	// Updating defaults of non-exising device should create it
	expectedDefaults = updateDefaults(t, reg, "AABBCC", Defaults{"D100", -4, 500})

	dev, err = reg.GetOrCreate("AABBCC")
	require.NoError(t, err)
	assert.Equal(t, Device{Defaults: expectedDefaults}, dev)

	err = reg.DeleteDevice("AABBCC")
	require.NoError(t, err)

	// Updating state of non-exising device should create it
	expectedState = updateState(t, reg, "AABBCC", State{[]net.IP{ip}, 2970})

	dev, err = reg.GetOrCreate("AABBCC")
	require.NoError(t, err)
	assert.Equal(t, Device{State: expectedState}, dev)

	_ = updateState(t, reg, "AABBCC", State{})

	dev, err = reg.GetOrCreate("AABBCC")
	require.NoError(t, err)
	assert.Equal(t, Device{}, dev)

	err = reg.DeleteDevice("AABBCC")
	require.NoError(t, err)

	// Updating config of non-exising device should create it
	expectedConfig = updateConfig(t, reg, "AABBCC", Config{ip})

	dev, err = reg.GetOrCreate("AABBCC")
	require.NoError(t, err)
	assert.Equal(t, Device{Config: expectedConfig}, dev)

	_ = updateConfig(t, reg, "AABBCC", Config{})

	dev, err = reg.GetOrCreate("AABBCC")
	require.NoError(t, err)
	assert.Equal(t, Device{}, dev)
}

func TestRegistry_GetAll(t *testing.T) {
	reg := CreateTestRegistry(t)

	expected := map[string]Device{}
	assert.Equal(t, expected, getAll(t, reg))

	_, err := reg.GetOrCreate("EMPTY")
	assert.NoError(t, err)
	expected["EMPTY"] = Device{}
	assert.Equal(t, expected, getAll(t, reg))

	expectedDefaults := updateDefaults(t, reg, "12345", Defaults{"D100", -4, 5000})
	expected["12345"] = Device{Defaults: expectedDefaults}
	assert.Equal(t, expected, getAll(t, reg))

	expectedState := updateState(t, reg, "AABBCC", State{[]net.IP{ip}, 2970})
	expected["AABBCC"] = Device{State: expectedState}
	assert.Equal(t, expected, getAll(t, reg))

	expectedState = updateState(t, reg, "12345", State{[]net.IP{ip}, 2970})
	expected["12345"] = Device{Defaults: expectedDefaults, State: expectedState}
	assert.Equal(t, expected, getAll(t, reg))

	err = reg.DeleteDevice("12345")
	assert.NoError(t, err)
	delete(expected, "12345")
	assert.Equal(t, expected, getAll(t, reg))
}

func updateDefaults(t *testing.T, reg *Registry, id string, d Defaults) Defaults {
	err := reg.UpdateDefaults(id, d)
	require.NoError(t, err)
	return d
}

func updateState(t *testing.T, reg *Registry, id string, state State) State {
	err := reg.UpdateState(id, state)
	require.NoError(t, err)
	return state
}

func updateConfig(t *testing.T, reg *Registry, id string, config Config) Config {
	err := reg.UpdateConfig(id, config)
	require.NoError(t, err)
	return config
}

func getAll(t *testing.T, reg *Registry) map[string]Device {
	devices, err := reg.GetDevices()
	require.NoError(t, err)
	return devices
}
