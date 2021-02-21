package main

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"
)

// a generator is is the function responsible for creating proper "shortened" name of url.
// As Generator can be complex beast, we expect unknown (if any) arguments to them as long as they ultimately return a string and feedback error
// Identity generator (for user enforced names) and uuid generator are provided as exemples.
type Generator func(arg ...interface{}) (string, error)

type Factory struct {
	store     Store
	generator Generator
}

//used for user-defined short urls
var IdentityGenerator = func(id ...interface{}) (string, error) {
	str := fmt.Sprintf("%s", id[0].(string))
	return str, nil
}

//used for UUID-made short urls
var DefaultGenerator = func(...interface{}) (string, error) {
	id, err := uuid.NewRandom()
	return id.String(), err
}

// Takes uri and any optional argument that generator needs.
// 1. Check URL is valid
// 2.a Call Generator it was constructed with to generate key accordingly
// 2.b Ensure uniqueness of said key
// 3. sends back resulting key if all went well, error if any
func (f *Factory) Gen(uri string, opt ...interface{}) (key string, err error) {
	_, err = url.ParseRequestURI(uri)
	if err != nil {
		return "", err
	}
	key, err = f.generator(opt...)
	if err != nil {
		return "", err
	}
	for {
		if v := f.store.Get(key); v == "" {
			break
		}
		key, err = f.generator(opt...)
		if err != nil {
			return "", err
		}
	}
	return key, nil
}

func NewFactory(generator Generator, store Store) *Factory {
	return &Factory{
		store:     store,
		generator: generator,
	}
}
