package device_registry

import (
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func initializeDevicesBucket(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Devices"))
		if err != nil {
			return errors.WithStack(err)
		}
		log.Info("Initialized bucket Devices")
		return nil
	})
}
