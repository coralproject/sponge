/*
Package report

CSV errors file with this fields:
table, id, row, "what went wrong"

*/

package report

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/log"
	"github.com/boltdb/bolt"
)

const (
	dbnameDefault = "sponge.db"
)

var (
	dbname string
	uuid   string
)

// Note is the schema for each row in the db
type Note struct {
	Row     string
	Details string
	Error   string
}

// Init initialize the records to be recorded
func Init(u string, dbn string) {

	uuid = u

	// 	{"table", "id", "row", "note", "error"},

	if dbn == "" {
		dbname = dbnameDefault
	} else {
		dbname = dbn
	}

}

// Record adds a new record to the report on failed imports
func Record(model string, id interface{}, row map[string]interface{}, n string, e error) {

	var srow string
	for k, v := range row {
		srow = fmt.Sprintf("%v/%s:%s", srow, k, v)
	}

	reflect.TypeOf(id)

	var key []byte

	switch v := id.(type) {
	case string:
		key = []byte(v)
	case bson.ObjectId:
		key, _ = v.MarshalJSON()
		//key = []byte(v.Hex())
	default:
		log.Error(uuid, "report.record", fmt.Errorf("Error on assertion. Type is %v", reflect.TypeOf(id)), "Asserting the ID to string")
	}

	note := &Note{
		Row:     srow,
		Details: n,
		Error:   e.Error(),
	}

	value, err := json.Marshal(note)
	if err != nil {
		log.Error(uuid, "report.record", err, "Marshaling data.")
	}

	db, err := bolt.Open(dbname, 0600, nil)
	if err != nil {
		log.Error(uuid, "report.record", err, "Initializing database.")
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(model))
		if err != nil {
			log.Error(uuid, "report.record", err, "Creating bucket.")
			return err
		}

		err = b.Put(key, value)
		if err != nil {
			log.Error(uuid, "report.record", err, "Recording data.")
		}

		return err
	})

	if err != nil {
		log.Error(uuid, "report.record", e, "Commiting to Bolt DB.")
	}
}

// GetRecords returns the actual records that I'm recording
func GetRecords(model string) (map[string]interface{}, error) {

	m := make(map[string]interface{})
	var err error

	db, err := bolt.Open(dbname, 0600, nil)
	if err != nil {
		log.Error(uuid, "report.record", err, "Initializing database.")
	}
	defer db.Close()

	var vj Note
	var kj int

	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(model))

		err := b.ForEach(func(k []byte, v []byte) error {
			err := json.Unmarshal(v, &vj)
			if err != nil {
				log.Error(uuid, "report.getrecords", err, "Unmarshalling.")
			}
			err = json.Unmarshal(k, &kj)
			if err != nil {
				reflect.TypeOf(kj)
				log.Error(uuid, "report.getrecords", err, "Unmarshalling.")
				return err
			}

			m[strconv.Itoa(kj)] = vj

			return err
		})
		return err
	})
	if err != nil {
		log.Error(uuid, "report.getrecords", err, "Reading bolt database.")
	}

	return m, err
}

// GetRecordsForBucket use the bucket to look all the records
func GetRecordsForBucket(b *bolt.Bucket) (map[string]interface{}, error) {
	m := make(map[string]interface{})

	var vj Note
	var kj string

	err := b.ForEach(func(k []byte, v []byte) error {
		err := json.Unmarshal(v, &vj)
		if err != nil {
			log.Error(uuid, "report.getrecordsforbucket", err, "Unmarshalling.")
		}
		err = json.Unmarshal(k, &kj)
		if err != nil {
			log.Error(uuid, "report.getrecordsforbucket", err, "Unmarshalling.")
		}

		m[kj] = vj

		return err
	})

	return m, err

}

// ReadReport gets the ids that needs to be imported from the report already in disk and return the records in a way that can be easily read it
func ReadReport(dbname string) (map[string][]string, error) { //(map[string]map[string]interface{}, error) {
	maa := make(map[string][]string) //make(map[string]map[string]interface{})

	db, err := bolt.Open(dbname, 0600, nil)
	if err != nil {
		log.Error(uuid, "report.record", err, "Opening bolt database.")
	}

	err = db.View(func(tx *bolt.Tx) error {

		var err error
		err = tx.ForEach(func(bucketname []byte, b *bolt.Bucket) error {

			m, err := GetRecordsForBucket(b)
			if err != nil {
				return err
			}
			s := make([]string, 0, len(m))
			for i := range m {
				s = append(s, i)
			}
			maa[string(bucketname)] = s

			return err
		})

		return err
	})

	return maa, err
}

// Print all the reports
func Print() {

	//GetRecords(model string) (map[string]interface{}, error) {
	m, err := ReadReport(dbname)
	if err != nil {
		fmt.Printf("Error on reading %s. Error: %v.", dbname, err)
		return
	}

	for i := range m {
		fmt.Println("Model: ", i)

		records, err := GetRecords(i)
		if err != nil {
			fmt.Println("Error on getting records: ", err)
		}
		for j, k := range records {
			fmt.Println("ID: ", j)
			fmt.Println("Details: ", k.(Note).Details)
			fmt.Println("Error: ", k.(Note).Error)
		}
	}
}

// SetImportDate
// GetLastImportDate
