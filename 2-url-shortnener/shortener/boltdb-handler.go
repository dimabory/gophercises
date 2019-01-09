package shortener

import (
	"github.com/boltdb/bolt"
	"log"
)

func Open(filepath string) (*bolt.DB, error) {
	db, err := bolt.Open(filepath, 0600, nil)
	if err != nil {
		panic(err)
	}

	return db, nil
}

func NewBoltDbUrlMapper(filename, bucket string) (func(string) (string, bool), error) {
	db, err := Open(filename)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	populateBucket(db, bucket)

	var mapping = make(map[string]string)

	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if err := b.ForEach(func(k, v []byte) error {
			mapping[string(k)] = string(v)
			return nil
		}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Fatal(err)
	}

	return NewBaseUrlMapper(mapping), nil
}

func populateBucket(db *bolt.DB, bucket string) {
	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}

		if err := putDefaultValuesIntoTheBucket(b); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Fatal(err)
	}
}

func putDefaultValuesIntoTheBucket(bucket *bolt.Bucket) error {
	if err := bucket.Put([]byte("/1"), []byte("https://google.com")); err != nil {
		return err
	}
	if err := bucket.Put([]byte("/2"), []byte("https://github.com")); err != nil {
		return err
	}
	if err := bucket.Put([]byte("/3"), []byte("https://bbc.com")); err != nil {
		return err
	}

	return nil
}
