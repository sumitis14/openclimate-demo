package database

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"os"

	edb "github.com/Varunram/essentials/database"
	"github.com/YaleOpenLab/openclimate/globals"
	"github.com/boltdb/bolt"
)

var UserBucket = []byte("Users")
var CompanyBucket = []byte("Companies")
var RegionBucket = []byte("Regions")
var CityBucket = []byte("Cities")

// CreateHomeDir creates a home directory
func CreateHomeDir() {
	if _, err := os.Stat(globals.HomeDir); os.IsNotExist(err) {
		// directory does not exist, create one
		log.Println("Creating home directory")
		os.MkdirAll(globals.HomeDir, os.ModePerm)
	}

	if _, err := os.Stat(globals.DbDir); os.IsNotExist(err) {
		os.MkdirAll(globals.DbDir, os.ModePerm)
	}
}

// don't lock since boltdb can only process one operation at a time. As the application
// grows bigger, this would be a major reason to search for a new db system

// OpenDB opens the db
func OpenDB() (*bolt.DB, error) {
	return edb.CreateDB(globals.DbDir+"/openclimate.db", UserBucket, CompanyBucket, RegionBucket, CityBucket)
}

// DeleteKeyFromBucket deletes a given key from the bucket bucketName but doesn
// not shift indices of elements succeeding the deleted element's index
func DeleteKeyFromBucket(key int, bucketName []byte) error {
	return edb.DeleteKeyFromBucket(globals.DbDir+"/openclimate.db", key, bucketName)
}

func RetrieveKey(bucketName []byte, key int) ([]byte, error) {
	x, err := edb.Retrieve(globals.DbDir+"/openclimate.db", CompanyBucket, key)
	if err != nil {
		log.Println(x)
		return nil, errors.Wrap(err, "error while retrieving key from bucket")
	}

	xBytes, err := json.Marshal(x)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal json, quitting")
	}

	return xBytes, nil
}
