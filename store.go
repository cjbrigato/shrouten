package main

import (
	"bytes"

	bolt "go.etcd.io/bbolt"
)

var Panic = func(v interface{}) {
	panic(v)
}

var _ Store = &DB{}

// one bucket to shroutnedn them all
var brucklet = []byte("urls")

func openDBFile(file string) *bolt.DB {
	// boltdB creates if not existing
	db, err := bolt.Open(file, 0600, nil)
	if err != nil {
		Panic(err)
	}

	buckets := [...][]byte{
		brucklet,
	}

	db.Update(func(tx *bolt.Tx) (err error) {
		for _, bucket := range buckets {
			_, err = tx.CreateBucketIfNotExists(bucket)
			if err != nil {
				Panic(err)
			}
		}
		return
	})
	return db
}

// NewDB return instances with already opened db
// TODO : implement a Store interface ? ~sortof
func NewDB(file string) *DB {
	return &DB{
		db: openDBFile(file),
	}
}

//put in.
// for now we don't update on same value, but why ? todo : at least stop telling success when not updating because value exists
//if time, fuzz with some factory/store interface to check to what extent we've lost knownledge
func (d *DB) Set(key string, value string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(brucklet)
		if err != nil {
			return err
		}
		k := []byte(key)
		valueB := []byte(value)
		c := b.Cursor() // valid for transaction duration scope....
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

// Key returns an empty string if not found, not nil.
func (d *DB) Get(key string) (value string) {
	keyB := []byte(key)
	d.db.Update(func(tx *bolt.Tx) error {
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
			u := ShortenedUrl{string(k), string(v)}
			urls = append(urls, u)
			return nil
		})
		return nil
	})
	return
}

func (d *DB) Len() (num int) {
	d.db.View(func(tx *bolt.Tx) error {
		//dirty fixed 0 total key in template conditional...
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
		Panic(err)
	}
}
