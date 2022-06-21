package database

import (
	"elize/elizeutils"
	"fmt"
	"os"
	"strings"

	bolt "go.etcd.io/bbolt"
)

const (
	dbname       = "ElizeDB"
	chainBucket  = "chain"
	blocksBucket = "blocks"
	lastPoint    = "lastPoint"
)

var db *bolt.DB

type DB struct{}

func (DB) SaveBlock(hash string, data []byte) {
	saveBlock(hash, data)
}
func (DB) SaveBlockchain(data []byte) {
	saveBlockchain(data)
}
func (DB) LastBlockPoint() []byte {
	return RestoreChain()
}
func (DB) FindBlock(hash string) []byte {
	return FindOneBlock(hash)
}
func (DB) DeleteBlocks() {
	deleteBlocks()
}

func GetDBname() string {
	var portNum string
	for _, arg := range os.Args {
		if strings.Contains(arg, "port") {
			portNum = strings.Split(arg, "=")[1]
		}
	}
	return fmt.Sprintf("%s_%s.db", dbname, portNum)
}

func InitDB() {
	if db == nil {
		newDBpointer, err := bolt.Open(GetDBname(), 0644, nil)
		db = newDBpointer
		elizeutils.Errchk(err)
		elizeutils.Errchk(
			db.Update(func(t *bolt.Tx) error {
				_, err := t.CreateBucketIfNotExists([]byte(chainBucket))
				elizeutils.Errchk(err)
				_, err = t.CreateBucketIfNotExists([]byte(blocksBucket))
				return err
			}))
	}
}

func Close() {
	db.Close()
}

func saveBlock(hash string, data []byte) {
	err := db.Update(func(t *bolt.Tx) error {
		blockbucket := t.Bucket([]byte(blocksBucket))
		return blockbucket.Put([]byte(hash), data)
	})
	elizeutils.Errchk(err)
}

func saveBlockchain(data []byte) {
	err := db.Update(func(t *bolt.Tx) error {
		blockbucket := t.Bucket([]byte(chainBucket))
		return blockbucket.Put([]byte(lastPoint), data)
	})
	elizeutils.Errchk(err)
}

func RestoreChain() []byte {
	var data []byte
	db.View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(chainBucket))
		data = bucket.Get([]byte(lastPoint))
		return nil
	})
	return data
}

func FindOneBlock(hash string) []byte {
	var data []byte
	db.View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})
	return data
}

func deleteBlocks() {
	db.Update(func(t *bolt.Tx) error {
		elizeutils.Errchk(t.DeleteBucket([]byte(blocksBucket)))
		_, err := t.CreateBucket([]byte(blocksBucket))
		elizeutils.Errchk(err)
		return nil
	})
}
