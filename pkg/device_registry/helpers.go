package device_registry

import (
	"github.com/boltdb/bolt"
	"github.com/chacal/thread-mgmt-server/pkg/test"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"testing"
)

func initializeBucket(db *bolt.DB, bucketName string) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return errors.WithStack(err)
		}
		log.Infof("Initialized bucket %v", bucketName)
		return nil
	})
}

func CreateTestRegistry(t *testing.T) *Registry {
	reg, err := Open(test.Tempfile())
	require.NoError(t, err)
	t.Cleanup(func() {
		reg.Close()
	})
	return reg
}
