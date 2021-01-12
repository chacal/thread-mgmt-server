package device_registry

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

type Device struct {
	Defaults Defaults `json:"defaults"`
	State    *State   `json:"state,omitempty"`
	Config   Config   `json:"config"`
}

const (
	GOOD_DISPLAY_1_54IN      = "GOOD_DISPLAY_1_54IN"
	GOOD_DISPLAY_2_13IN      = "GOOD_DISPLAY_2_13IN"
	GOOD_DISPLAY_2_9IN       = "GOOD_DISPLAY_2_9IN"
	GOOD_DISPLAY_2_9IN_4GRAY = "GOOD_DISPLAY_2_9IN_4GRAY"
)

const (
	E73          = "E73"
	MS88SF2_V1_0 = "MS88SF2_V1_0"
)

type Defaults struct {
	Instance    string `json:"instance"`
	TxPower     int    `json:"txPower"`
	PollPeriod  int    `json:"pollPeriod"`
	DisplayType string `json:"displayType"`
	HwVersion   string `json:"hwVersion"`
}

type ParentInfo struct {
	Rloc16         string `json:"rloc16"`
	LinkQualityIn  int    `json:"linkQualityIn"`
	LinkQualityOut int    `json:"linkQualityOut"`
	AvgRssi        int    `json:"avgRssi"`
	LatestRssi     int    `json:"latestRssi"`
}

type State struct {
	Addresses  []net.IP   `json:"addresses"`
	Vcc        int        `json:"vcc"`
	Instance   string     `json:"instance"`
	TxPower    int        `json:"txPower"`
	PollPeriod int        `json:"pollPeriod"`
	Parent     ParentInfo `json:"parent"`
}

type Config struct {
	MainIp                  net.IP `json:"mainIp"`
	StatePollingEnabled     bool   `json:"statePollingEnabled"`
	StatePollingIntervalSec int    `json:"statePollingIntervalSec"`
}

const DevicesBucket = "Devices"
const DefaultsBucket = "Defaults"
const StateBucket = "State"
const ConfigBucket = "Config"

var DefaultDefaults = Defaults{Instance: "0000", TxPower: 0, PollPeriod: 1000, DisplayType: "", HwVersion: ""}
var DefaultConfig = Config{MainIp: nil, StatePollingEnabled: false, StatePollingIntervalSec: 600}
var DefaultDevice = Device{DefaultDefaults, nil, DefaultConfig}

type Registry struct {
	db *bolt.DB
}

func Open(dbFileName string) (*Registry, error) {
	db, err := bolt.Open(dbFileName, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		return nil, errors.Wrapf(err, "error opening database file '%v'", dbFileName)
	}

	err = initializeBucket(db, DevicesBucket)
	if err != nil {
		return nil, err
	}

	return &Registry{db: db}, nil
}

func (r *Registry) Close() error {
	return r.db.Close()
}

func (r *Registry) Get(id string) (*Device, error) {
	var d *Device = nil
	var err error

	err = r.db.View(func(tx *bolt.Tx) error {
		err = assertDeviceExistsInTx(tx, id)
		if err != nil {
			return err
		}
		d, err = getDeviceInTx(tx, id)
		return err
	})

	return d, err
}

func (r *Registry) Create(id string) (*Device, error) {
	var d *Device = nil
	var err error

	err = r.db.Update(func(tx *bolt.Tx) error {
		devices := tx.Bucket([]byte(DevicesBucket))
		device := devices.Bucket([]byte(id))
		if device != nil {
			return errors.Errorf("device with id '%v' alredy exists", id)
		}
		_, err = devices.CreateBucket([]byte(id))
		if err != nil {
			return errors.WithStack(err)
		}
		err = putToDeviceBucket(tx, DefaultsBucket, id, DefaultDefaults)
		if err != nil {
			return err
		}
		err = putToDeviceBucket(tx, ConfigBucket, id, DefaultConfig)
		if err != nil {
			return err
		}

		d, err = getDeviceInTx(tx, id)
		return err
	})

	return d, err
}

func (r *Registry) Contains(id string) (bool, error) {
	var found bool

	err := r.db.View(func(tx *bolt.Tx) error {
		devices := tx.Bucket([]byte(DevicesBucket))
		device := devices.Bucket([]byte(id))
		found = device != nil
		return nil
	})

	return found, err
}

func (r *Registry) UpdateDefaults(id string, defaults Defaults) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		err := assertDeviceExistsInTx(tx, id)
		if err != nil {
			return err
		}
		return putToDeviceBucket(tx, DefaultsBucket, id, defaults)
	})
}

func (r *Registry) UpdateState(id string, state State) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		err := assertDeviceExistsInTx(tx, id)
		if err != nil {
			return err
		}
		return putToDeviceBucket(tx, StateBucket, id, state)
	})
}

func (r *Registry) UpdateConfig(id string, config Config) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		err := assertDeviceExistsInTx(tx, id)
		if err != nil {
			return err
		}
		return putToDeviceBucket(tx, ConfigBucket, id, config)
	})
}

func (r *Registry) GetDevices() (map[string]Device, error) {
	devices := make(map[string]Device)
	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DevicesBucket))
		return b.ForEach(func(k []byte, v []byte) error {
			id := string(k)
			device, err := getDeviceInTx(tx, id)
			if err != nil {
				return err
			}

			devices[string(k)] = *device
			return nil
		})
	})
	return devices, err
}

func (r *Registry) DeleteDevice(id string) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		err := assertDeviceExistsInTx(tx, id)
		if err != nil {
			return err
		}

		b := tx.Bucket([]byte(DevicesBucket))
		log.Debugf("Deleting device '%v'", id)
		err = b.DeleteBucket([]byte(id))
		if err != nil {
			return errors.Wrapf(err, "failed to delete device, id: '%v'", id)
		}

		return nil
	})
}

func assertDeviceExistsInTx(tx *bolt.Tx, id string) error {
	devices := tx.Bucket([]byte(DevicesBucket))
	device := devices.Bucket([]byte(id))
	if device == nil {
		return errors.Errorf("device with id '%v' not found", id)
	} else {
		return nil
	}
}

func getDeviceInTx(tx *bolt.Tx, id string) (*Device, error) {
	devices := tx.Bucket([]byte(DevicesBucket))
	device := devices.Bucket([]byte(id))
	if device == nil {
		return nil, nil
	}

	d := Device{}

	defaults, err := getDefaultsInTx(tx, id)
	if err != nil {
		return nil, err
	}
	if defaults != nil {
		d.Defaults = *defaults
	}

	state, err := getStateInTx(tx, id)
	if err != nil {
		return nil, err
	}
	if state != nil {
		d.State = state
	}

	config, err := getConfigInTx(tx, id)
	if err != nil {
		return nil, err
	}
	if config != nil {
		d.Config = *config
	}

	return &d, nil
}

func getDefaultsInTx(tx *bolt.Tx, id string) (*Defaults, error) {
	buf := getFromDeviceBucket(tx, DefaultsBucket, id)
	if buf == nil {
		return nil, nil
	}

	d, err := defaultsFromJSON(buf)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func getStateInTx(tx *bolt.Tx, id string) (*State, error) {
	buf := getFromDeviceBucket(tx, StateBucket, id)
	if buf == nil {
		return nil, nil
	}

	state, err := StateFromJSON(buf)
	if err != nil {
		return nil, err
	}

	return &state, nil
}

func getConfigInTx(tx *bolt.Tx, id string) (*Config, error) {
	buf := getFromDeviceBucket(tx, ConfigBucket, id)
	if buf == nil {
		return nil, nil
	}

	config, err := configFromJSON(buf)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func getFromDeviceBucket(tx *bolt.Tx, bucketName string, id string) []byte {
	devices := tx.Bucket([]byte(DevicesBucket))
	device := devices.Bucket([]byte(id))
	if device == nil {
		return nil
	}

	bucket := device.Bucket([]byte(bucketName))
	if bucket == nil {
		return nil
	}

	return bucket.Get([]byte(id))
}

func putToDeviceBucket(tx *bolt.Tx, bucketName string, id string, obj interface{}) error {
	devices := tx.Bucket([]byte(DevicesBucket))
	device, err := devices.CreateBucketIfNotExists([]byte(id))
	if err != nil {
		return errors.WithStack(err)
	}
	b, err := device.CreateBucketIfNotExists([]byte(bucketName))
	if err != nil {
		return errors.WithStack(err)
	}

	log.Debugf("Putting to bucket %v '%v': %+v", bucketName, id, obj)
	buf, err := json.Marshal(obj)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal: %v", obj)
	}

	err = b.Put([]byte(id), buf)
	if err != nil {
		return errors.Wrapf(err, "failed to put: %+v", obj)
	}

	return nil
}

func defaultsFromJSON(buf []byte) (Defaults, error) {
	d := Defaults{}
	err := json.Unmarshal(buf, &d)
	if err != nil {
		return d, errors.Wrapf(err, "failed to unmarshal defaults from db, data: %v", string(buf))
	}
	return d, nil
}

func StateFromJSON(buf []byte) (State, error) {
	state := State{}
	err := json.Unmarshal(buf, &state)
	if err != nil {
		return state, errors.Wrapf(err, "failed to unmarshal state from db, data: %v", string(buf))
	}
	return state, nil
}

func configFromJSON(buf []byte) (Config, error) {
	config := Config{}
	err := json.Unmarshal(buf, &config)
	if err != nil {
		return config, errors.Wrapf(err, "failed to unmarshal config from db, data: %v", string(buf))
	}
	return config, nil
}
