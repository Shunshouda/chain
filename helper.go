package chain

import (
	"fmt"
	"reflect"
)

// Get retrieves a value of type T from the chain
// Returns error if value not found or chain has error
func Get[T any](c *Chain) (T, error) {
	var zero T
	if c.err != nil {
		return zero, c.err
	}

	t := reflect.TypeOf(zero)
	if val, ok := c.values[t]; ok {
		return val.(T), nil
	}

	return zero, fmt.Errorf("no value found for type: %v", t)
}

// GetOrDefault retrieves a value of type T or returns the default value if not found
func GetOrDefault[T any](c *Chain, defaultValue T) T {
	if c.err != nil {
		return defaultValue
	}

	t := reflect.TypeOf(defaultValue)
	if val, ok := c.values[t]; ok {
		return val.(T)
	}

	return defaultValue
}

// MustGet retrieves a value of type T or panics if not found
func MustGet[T any](c *Chain) T {
	val, err := Get[T](c)
	if err != nil {
		panic(err)
	}
	return val
}

// Has checks if a value of type T exists in the chain
func Has[T any](c *Chain) bool {
	var zero T
	_, ok := c.values[reflect.TypeOf(zero)]
	return ok
}

// Clear removes a value of type T from the chain
func Clear[T any](c *Chain) *Chain {
	var zero T
	delete(c.values, reflect.TypeOf(zero))
	return c
}

// GetAs retrieves a value and converts it to the target type
func GetAs[T any](c *Chain) (T, error) {
	return Get[T](c)
}

// Map transforms a value of type T1 to T2 using the provided function
func Map[T1, T2 any](c *Chain, fn func(T1) T2) *Chain {
	return c.Then(func(t1 T1) T2 {
		return fn(t1)
	})
}

// FlatMap transforms a value of type T1 to T2 with possible error
func FlatMap[T1, T2 any](c *Chain, fn func(T1) (T2, error)) *Chain {
	return c.Then(func(t1 T1) (T2, error) {
		return fn(t1)
	})
}
