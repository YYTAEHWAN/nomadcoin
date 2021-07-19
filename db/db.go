package db

import (
	"learngo/github.com/nomadcoders/utils"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

const (
	dbName      = "blockchain.db"
	dataBucket  = "data"
	blockBucket = "blocks"
)

func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open("dbName", 0600, nil)
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
