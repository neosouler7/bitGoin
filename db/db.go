package db

import (
	"fmt"
	"os"
	"sync"

	"github.com/neosouler7/bitGoin/utils"
	bolt "go.etcd.io/bbolt"
)

var (
	db   *bolt.DB
	once sync.Once
)

type DB struct{}

func (DB) FindBlock(hash string) []byte {
	return findBlock(hash)
}

func (DB) LoadChain() []byte {
	return loadChain()
}

func (DB) SaveBlock(hash string, data []byte) {
	saveBlock(hash, data)
}

func (DB) SaveChain(data []byte) {
	saveChain(data)
}

func (DB) DeleteAllBlocks() {
	deleteAllBlocks()
}

const (
	dbName       = "blockchain"
	dataBucket   = "data"
	blocksBucket = "blocks"
	checkpoint   = "checkpoint"
)

func getDbName() string {
	port := os.Args[2][6:]
	return fmt.Sprintf("%s_%s.db", dbName, port)
}

func Close() {
	utils.HandleErr(db.Close())
}

func InitDB() {
	if db == nil {
		once.Do(func() {
			dbPointer, err := bolt.Open(getDbName(), 0600, nil)
			db = dbPointer
			utils.HandleErr(err)

			err = db.Update(func(tx *bolt.Tx) error {
				_, err = tx.CreateBucketIfNotExists([]byte(dataBucket))
				utils.HandleErr(err)
				_, err = tx.CreateBucketIfNotExists([]byte(blocksBucket))
				return err
			})
			utils.HandleErr(err)
		})
	}
}

func saveBlock(hash string, data []byte) {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleErr(err)
}

func saveChain(data []byte) {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data)
		return err
	})
	utils.HandleErr(err)
}

func loadChain() []byte {
	var data []byte
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	utils.HandleErr(err)
	return data
}

func findBlock(hash string) []byte {
	var data []byte
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})
	utils.HandleErr(err)
	return data
}

func deleteAllBlocks() {
	err := db.Update(func(tx *bolt.Tx) error {
		utils.HandleErr(tx.DeleteBucket([]byte(blocksBucket)))
		_, err := tx.CreateBucket([]byte(blocksBucket))
		utils.HandleErr(err)
		return nil
	})
	utils.HandleErr(err)
}
