package main

import (
	"bytes"

	bolt "go.etcd.io/bbolt"
)

type Store interface {
	All() ShortenedUrls
	Set(key string, value string) error
	Get(key string) string
	Len() int
	Clear() error
	Close()
}

//implements store
type DB struct {
	db *bolt.DB
}

// Was trick to ensure improper implementation of Store interface would come up at Compile time.
// This is unneeded as api function now properly takes a Store as parameter.
//var _ Store = &DB{}

// The bucket (in boltdb sense) used to keep shroutened urls.
var brucklet = []byte("urls")

func openDBFile(file string) (*bolt.DB, error) {

	// boltdB creates a db if it doesn't exists
	// here we considered that combined failure to either open or create proper database makes persistence failure and that such failure is fatal.
	// But this was not only argable but also a failure to provide agnosticism. So such panic is deported to NewDB func.
	db, err := bolt.Open(file, 0600, nil)
	if err != nil {
		return nil, err
	}

	// We only use one bucket, but the "..." denotes you can add more as buckets lenght is then specified
	// as equals to the number of elements in buckets array literal.
	buckets := [...][]byte{
		brucklet,
	}

	// bucket inavailability is not fatal and can be recovered
	db.Update(func(tx *bolt.Tx) (err error) {
		for _, bucket := range buckets {
			_, err = tx.CreateBucketIfNotExists(bucket)
			if err != nil {
				return err
			}
		}
		return
	})
	return db, nil
}

// NewDB return instances with already opened db
func NewDB(file string) *DB {
	db, err := openDBFile(file)
	if err != nil {
		panic(err)
	}
	return &DB{db}
}

//put in
func (d *DB) Set(key string, value string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(brucklet)
		if err != nil {
			return err
		}
		k := []byte(key)
		valueB := []byte(value)
		c := b.Cursor()
		exists := false
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if bytes.Equal(valueB, v) {
				exists = true
				break
			}
		}
		if exists {
			return nil
		}
		return b.Put(k, []byte(value))
	})
}

// Poof!
func (d *DB) Clear() error {
	return d.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(brucklet)
	})
}

// Warning: boltdb will return an empty string if key is not found, not nil.
func (d *DB) Get(key string) (value string) {
	keyB := []byte(key)
	d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(brucklet)
		if b == nil {
			return nil
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if bytes.Equal(keyB, k) {
				value = string(v)
				break
			}
		}

		return nil
	})
	return
}

func (d *DB) All() (urls ShortenedUrls) {
	d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(brucklet)
		if b == nil {
			return nil
		}
		b.ForEach(func(k, v []byte) error {
			urls.addInPlace(string(k), string(v))
			return nil
		})
		return nil
	})
	return
}

func (d *DB) Len() (num int) {
	d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(brucklet)
		if b == nil {
			return nil
		}
		b.ForEach(func([]byte, []byte) error {
			num++
			return nil
		})
		return nil
	})
	return
}

// Close-shutdowns db handle.
func (d *DB) Close() {
	if err := d.db.Close(); err != nil {
		panic(err)
	}
}
