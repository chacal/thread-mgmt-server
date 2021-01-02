package device_registry

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

var ip = net.ParseIP("ffff::1")

var testState = State{[]net.IP{ip}, 2970, "A100", -4, 1000,
	ParentInfo{"0x4400", 3, 2, -65, -63},
}

func TestRegistry_GetAndCreate(t *testing.T) {
	reg := CreateTestRegistry(t)

	dev, err := reg.Get("12345")
	assert.Error(t, err)
	assert.Equal(t, (*Device)(nil), dev)

	dev, err = reg.Create("12345")
	require.NoError(t, err)
	assert.Equal(t, &DefaultDevice, dev)
}

func TestRegistry_UpdateDefaults(t *testing.T) {
	reg := CreateTestRegistry(t)

	err := reg.UpdateDefaults("12345", DefaultDefaults)
	assert.Error(t, err)

	dev, _ := reg.Create("12345")

	expectedDefaults := updateDefaults(t, reg, "12345", Defaults{"D100", -4, 500, GOOD_DISPLAY_1_54IN})
	dev, _ = reg.Get("12345")
	assert.Equal(t, &Device{Defaults: expectedDefaults, Config: DefaultConfig}, dev)

	_ = updateDefaults(t, reg, "12345", Defaults{})
	dev, _ = reg.Get("12345")
	assert.Equal(t, &Device{Defaults: Defaults{"", 0, 0, ""}, Config: DefaultConfig}, dev)
}

func TestRegistry_UpdateConfig(t *testing.T) {
	reg := CreateTestRegistry(t)

	err := reg.UpdateConfig("12345", DefaultConfig)
	assert.Error(t, err)

	dev, _ := reg.Create("12345")

	expectedConfig := updateConfig(t, reg, "12345", Config{ip, true, 300})
	dev, _ = reg.Get("12345")
	assert.Equal(t, &Device{Defaults: DefaultDefaults, Config: expectedConfig}, dev)

	expectedConfig = updateConfig(t, reg, "12345", Config{})
	dev, _ = reg.Get("12345")
	assert.Equal(t, &Device{Defaults: DefaultDefaults, Config: Config{nil, false, 0}}, dev)
}

func TestRegistry_UpdateState(t *testing.T) {
	reg := CreateTestRegistry(t)

	err := reg.UpdateState("12345", testState)
	assert.Error(t, err)

	dev, _ := reg.Create("12345")

	expectedState := updateState(t, reg, "12345", testState)
	dev, _ = reg.Get("12345")
	assert.Equal(t, &Device{Defaults: DefaultDefaults, Config: DefaultConfig, State: expectedState}, dev)
}

func TestRegistry_GetDevices(t *testing.T) {
	reg := CreateTestRegistry(t)

	expected := map[string]Device{}
	assert.Equal(t, expected, getAll(t, reg))

	_, err := reg.Create("EMPTY")
	assert.NoError(t, err)
	expected["EMPTY"] = DefaultDevice
	assert.Equal(t, expected, getAll(t, reg))

	_, err = reg.Create("12345")
	assert.NoError(t, err)
	expectedDefaults := updateDefaults(t, reg, "12345", Defaults{"D100", -4, 5000, GOOD_DISPLAY_1_54IN})
	expected["12345"] = Device{Defaults: expectedDefaults, Config: DefaultConfig}
	assert.Equal(t, expected, getAll(t, reg))

	expectedConfig := updateConfig(t, reg, "12345", Config{ip, true, 100})
	expected["12345"] = Device{Defaults: expectedDefaults, Config: expectedConfig}
	assert.Equal(t, expected, getAll(t, reg))

	expectedState := updateState(t, reg, "12345", testState)
	expected["12345"] = Device{Defaults: expectedDefaults, State: expectedState, Config: expectedConfig}
	assert.Equal(t, expected, getAll(t, reg))

	_, err = reg.Create("AABBCC")

	expected["AABBCC"] = DefaultDevice
	assert.Equal(t, expected, getAll(t, reg))

	expectedState = updateState(t, reg, "AABBCC", testState)
	expected["AABBCC"] = Device{Defaults: DefaultDefaults, State: expectedState, Config: DefaultConfig}
	assert.Equal(t, expected, getAll(t, reg))

	err = reg.DeleteDevice("12345")
	assert.NoError(t, err)
	delete(expected, "12345")
	assert.Equal(t, expected, getAll(t, reg))
}

func TestRegistry_DeleteDevice(t *testing.T) {
	reg := CreateTestRegistry(t)

	err := reg.DeleteDevice("12345")
	assert.Error(t, err)

	_, err = reg.Create("12345")
	assert.NoError(t, err)

	dev, err := reg.Get("12345")
	require.NoError(t, err)
	assert.Equal(t, &DefaultDevice, dev)

	err = reg.DeleteDevice("12345")
	assert.NoError(t, err)

	_, err = reg.Get("12345")
	assert.Error(t, err)
}

func TestRegistry_Contains(t *testing.T) {
	reg := CreateTestRegistry(t)

	contains, err := reg.Contains("12345")
	require.NoError(t, err)
	assert.Equal(t, false, contains)

	_, err = reg.Create("12345")
	require.NoError(t, err)

	contains, err = reg.Contains("12345")
	require.NoError(t, err)
	assert.Equal(t, true, contains)
}

func updateDefaults(t *testing.T, reg *Registry, id string, d Defaults) Defaults {
	err := reg.UpdateDefaults(id, d)
	require.NoError(t, err)
	return d
}

func updateState(t *testing.T, reg *Registry, id string, state State) *State {
	err := reg.UpdateState(id, state)
	require.NoError(t, err)
	return &state
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
