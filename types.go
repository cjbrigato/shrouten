package main

import bolt "go.etcd.io/bbolt"

type ShortenedUrls []ShortenedUrl

type ShortenedUrl struct {
	Short string `json:"short"`
	Uri   string `json:"uri"`
}

type Store interface {
	All() ShortenedUrls
	Set(key string, value string) error
	Get(key string) string
	Len() int
	Close()
}

type DB struct {
	db *bolt.DB
}

type Generator func(arg string) string

type Factory struct {
	store     Store
	generator Generator
}
