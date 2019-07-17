package database

import (
	"log"
	"encoding/json"
	"github.com/pkg/errors"

	edb "github.com/Varunram/essentials/database"
	globals "github.com/YaleOpenLab/openclimate/globals"
)

// Our definition of "City" includes states,
// regions, provinces, prefectures, etc. The
// following struct defines the relevant fields.
type Region struct {

	// Identifying info
	Index   int
	Name    string
	Country string

	// Contextual data
	Area        float64
	Iso         string
	Population  int
	Latitude    float64
	Longitude   float64
	Revenue     float64
	CompanySize int
	HQ          string
	// EntityType		string

	//	For regions: children = companies (divided by region)
	Children 			[]string

	// Data that is reported (through self-reporting, databases, IoT, etc.)
	// as opposed to data that is aggregated from its parts/children. Data
	// is stored on IPFS, so Reports holds the IPFS hashes.
	Reports				[]RepData

	AggEmissions 		AggEmiData
	AggMitigation		AggMitData
	AggAdaptation 		AggAdptData
}

/*
	ISSUE: edb.Save() asks for an key argument of type INT,
	but currently we are passing in a key argument of type string.
	This issue needs to be resolved. Could maybe just use a hash.

	RESOLVED: currently using solution previously implemented in OpenX;
	incrementing index for each new region, so the key is of type int.
*/

// Function that creates a new region object given its name and country
// and saves the object in the regions bucket.
func NewRegion(name string, country string) (Region, error) {
	var new Region
	var err error
	var lenRegions int
	// naive implementation of assigning keys to bucket items (simple indexing)
	regions, err := RetrieveAllRegions()
	if err != nil {
		// regions doesn't exist yet
		lenRegions = 0
	} else {
		lenRegions = len(regions)
	}

	new.Index = lenRegions + 1
	new.Name = name
	new.Country = country

	return new, new.Save()
}

// Saves region object in regions bucket. Called by NewRegion
func (region *Region) Save() error {
	return edb.Save(globals.DbPath, RegionBucket, region, region.Index)
}

// Given a key of type int, retrieves the corresponding region object
// from the database regions bucket.
func RetrieveRegion(key int) (Region, error) {
	var region Region
	regionBytes, err := edb.Retrieve(globals.DbPath, RegionBucket, key)
	if err != nil {
		return region, errors.Wrap(err, "error while retrieving key from bucket")
	}
	err = json.Unmarshal(regionBytes, &region)
	return region, err
}

// Given the name and country of the region, retrieves the
// corresponding region object from the database regions bucket.
func RetrieveRegionByName(name string, country string) (Region, error) {
	var region Region
	allRegions, err := RetrieveAllRegions()
	if err != nil {
		return region, errors.Wrap(err, "Error while retrieving all regions from database")
	}

	for _, val := range allRegions {
		if val.Name == name && val.Country == country {
			region = val
			return region, nil
		}
	}

	return region, errors.New("could not find regions")
}

// Retrieves all regions from the regions bucket.
func RetrieveAllRegions() ([]Region, error) {
	var regions []Region
	keys, err := edb.RetrieveAllKeys(globals.DbPath, RegionBucket)
	if err != nil {
		log.Println("couldn't retrieve regions")
		return regions, errors.Wrap(err, "error while retrieving all regions")
	}

	for _, val := range keys {
		var region Region
		err = json.Unmarshal(val, &region)
		if err != nil {
			return regions, errors.Wrap(err, "could not unmarshal json")
		}
		regions = append(regions, region)
	}

	return regions, nil
}
