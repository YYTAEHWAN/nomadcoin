package db

import (
	"fmt"
	"os"

	"github.com/nomadcoders/utils"

	bolt "go.etcd.io/bbolt"
)

var db *bolt.DB

const (
	dbName      = "blockchain"
	dataBucket  = "data"
	blockBucket = "blocks"
	checkpoint  = "checkpoint"
)

type DB struct{}

func (DB) GetBlockHashFromDb(hash string) []byte {
	return getBlockHashFromDb(hash)
} // nico's findBlock()
func (DB) SaveBlock(hash string, data []byte) {
	saveBlock(hash, data)
} // nico's saveBlock()
func (DB) SaveCheckpoint(data []byte) {
	saveCheckpoint(data)
} // nico's saveChain()
func (DB) GetCheckPointFromDb() []byte {
	return getCheckPointFromDb()
} // nico's loadChain()
func (DB) DeleteAllBlocks() {
	deleteAllBlocks()
}

func getDbName() string {
	port := os.Args[2][7:]
	return fmt.Sprintf("%s_%s.db", dbName, port)
}

func InitDB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(getDbName(), 0600, nil)
		db = dbPointer
		utils.HandleErr(err)
		err = db.Update(func(t *bolt.Tx) error {
			_, err := t.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleErr(err)
			_, err = t.CreateBucketIfNotExists([]byte(blockBucket))
			return err
		})
		utils.HandleErr(err)
	}
	return db
}
func Close() {
	db.Close()
}
func saveBlock(hash string, data []byte) {
	err := db.Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blockBucket))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleErr(err)
}

func saveCheckpoint(data []byte) {
	err := db.Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data)
		return err
	})
	utils.HandleErr(err)
}

func getCheckPointFromDb() []byte {
	var data []byte
	db.View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	return data
}

func getBlockHashFromDb(hash string) []byte {
	var data []byte
	db.View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blockBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})
	return data
}

func deleteAllBlocks() {
	db.Update(func(t *bolt.Tx) error {
		utils.HandleErr(t.DeleteBucket([]byte(blockBucket)))
		_, err := t.CreateBucket([]byte(blockBucket))
		utils.HandleErr(err)
		return nil
	})
}
