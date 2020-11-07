package device_registry

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Device struct {
	Name     string `json:"name,omitempty"`
	PollTime int    `json:"pollTime,omitempty"`
}

type Registry struct {
	db *bolt.DB
}

func Open(dbFileName string) (*Registry, error) {
	db, err := bolt.Open(dbFileName, 0600, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "error opening database file '%v'", dbFileName)
	}

	err = initializeDevicesBucket(db)
	if err != nil {
		return nil, err
	}

	return &Registry{db: db}, nil
}

func (r *Registry) Close() error {
	return r.db.Close()
}

func (r *Registry) GetOrCreate(id string) (Device, error) {
	d := Device{}

	err := r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Devices"))

		log.Debugf("Getting device '%v'", id)

		buf := b.Get([]byte(id))
		if buf == nil {
			return putDevice(tx, id, d)
		} else {
			err := json.Unmarshal(buf, &d)
			if err != nil {
				return errors.Wrapf(err, "failed to unmarshal device from db, data: %v", string(buf))
			}
			log.Debugf("Device '%v' found: %+v", id, d)
		}
		return nil
	})

	return d, err
}

func (r *Registry) Update(id string, dev Device) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		return putDevice(tx, id, dev)
	})
}

func putDevice(tx *bolt.Tx, id string, dev Device) error {
	b := tx.Bucket([]byte("Devices"))

	log.Debugf("Putting device '%v': %+v", id, dev)
	buf, err := json.Marshal(dev)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal: %v", dev)
	}

	err = b.Put([]byte(id), buf)
	if err != nil {
		return errors.Wrapf(err, "failed to put: %v", dev)
	}

	return nil
}