package structs

import (
	"fmt"
	"sync"
)

var (
	_ Singleton[string] = &OnceSingleton[string]{}
)

type (
	Provider[T any] interface {
		Singleton[T]
		Set(singleton Singleton[T])
	}

	Singleton[T any] interface {
		Get() (T, error)
		MustGet() T
	}

	OnceSingleton[T any] struct {
		object  T
		err     error
		once    sync.Once
		factory func() (T, error)
	}

	DefaultProvider[T any] struct {
		singleton Singleton[T]
	}
)

func NewProvider[T any]() Provider[T] {
	return &DefaultProvider[T]{}
}

func NewSingleton[T any](factory func() (T, error)) Singleton[T] {
	return &OnceSingleton[T]{once: sync.Once{}, factory: factory}
}

func (o *OnceSingleton[T]) Get() (T, error) {
	o.once.Do(func() {
		object, err := o.factory()
		o.object = object
		o.err = err
	})

	return o.object, o.err
}

func (o *OnceSingleton[T]) MustGet() T {
	object, err := o.Get()
	if err != nil {
		panic(fmt.Errorf("failed to get object: %w", err))
	}

	return object
}

func (d *DefaultProvider[T]) Set(singleton Singleton[T]) {
	d.singleton = singleton
}

func (d *DefaultProvider[T]) Get() (T, error) {
	var t T
	if d.singleton == nil {
		return t, fmt.Errorf("empty object")
	}

	return d.singleton.Get()
}

func (d *DefaultProvider[T]) MustGet() T {
	get, err := d.Get()
	if err != nil {
		panic(err)
	}

	return get
}
