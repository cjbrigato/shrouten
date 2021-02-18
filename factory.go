package main

import (
	"net/url"

	"github.com/google/uuid"
)

var DefaultGenerator = func(unused string) string {
	id, _ := uuid.NewRandom()
	return id.String()
}
var CustomNamed = func(id string) string {
	return id
}

// Gen key, abstraction
func (f *Factory) Gen(uri string, opt string) (key string, err error) {
	_, err = url.ParseRequestURI(uri)
	if err != nil {
		return "", err
	}
	key = f.generator(opt)
	//uniqueness
	for {
		if v := f.store.Get(key); v == "" {
			break
		}
		key = f.generator(opt)
	}
	return key, nil
}

//lazy
func NewFactory(generator Generator, store Store) *Factory {
	return &Factory{
		store:     store,
		generator: generator,
	}
}
