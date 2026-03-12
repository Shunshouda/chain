package chain

import (
	"fmt"
	"reflect"
)

// Chain is a fluent chain call handler that enables method chaining
// with automatic parameter injection and error handling
type Chain struct {
	values  map[reflect.Type]interface{}
	err     error
	errCtx  *ErrorContext
	onError func(*ErrorContext)
}

// ErrorContext contains the error and the values available at the time of error
type ErrorContext struct {
	Error  error
	Values map[reflect.Type]interface{}
	Chain  *Chain
}

// NewChain creates a new Chain instance
func NewChain() *Chain {
	return &Chain{
		values: make(map[reflect.Type]interface{}),
		err:    nil,
		onError: func(ctx *ErrorContext) {
			fmt.Printf("Error: %v\n", ctx.Error)
		},
	}
}

// OnError sets a custom error handling callback that receives an ErrorContext
func (c *Chain) OnError(callback func(*ErrorContext)) *Chain {
	if c.err == nil {
		c.onError = callback
	}
	return c
}

// Then adds a function to the chain and executes it
// The function can have multiple parameters and return values
// If the last return value is of type error, it will be treated as error
func (c *Chain) Then(fn interface{}) *Chain {
	// Skip execution if there's already an error
	if c.err != nil {
		return c
	}

	// Validate that fn is a function
	fnValue := reflect.ValueOf(fn)
	if fnValue.Kind() != reflect.Func {
		c.err = fmt.Errorf("provided value is not a function")
		c.triggerError(c.err)
		return c
	}

	// Prepare arguments for the function call
	args, err := c.prepareArguments(fnValue.Type())
	if err != nil {
		c.err = err
		c.triggerError(c.err)
		return c
	}

	// Execute the function
	results := fnValue.Call(args)

	// Process the return values
	c.processResults(results)

	return c
}

// prepareArguments prepares arguments for function call by matching types
// from previously stored values
func (c *Chain) prepareArguments(fnType reflect.Type) ([]reflect.Value, error) {
	numIn := fnType.NumIn()
	if numIn == 0 {
		return []reflect.Value{}, nil
	}

	args := make([]reflect.Value, numIn)

	// Match each parameter with stored values by type
	for i := 0; i < numIn; i++ {
		paramType := fnType.In(i)
		found := false

		for t, val := range c.values {
			if t.AssignableTo(paramType) {
				args[i] = reflect.ValueOf(val)
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("no value found for parameter type: %v at position %d", paramType, i)
		}
	}

	return args, nil
}

// processResults handles function return values
// If the last value is error, it's treated specially
func (c *Chain) processResults(results []reflect.Value) {
	if len(results) == 0 {
		return
	}

	lastIdx := len(results) - 1
	lastResult := results[lastIdx]

	// Check if the last return value is an error
	if lastResult.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		if !lastResult.IsNil() {
			c.err = lastResult.Interface().(error)
			c.triggerError(c.err)
			return
		}
		// Store all non-error return values
		for i := 0; i < len(results)-1; i++ {
			c.storeValue(results[i].Interface())
		}
	} else {
		// Store all return values
		for i := 0; i < len(results); i++ {
			c.storeValue(results[i].Interface())
		}
	}
}

// triggerError triggers the error callback with context
func (c *Chain) triggerError(err error) {
	c.errCtx = &ErrorContext{
		Error:  err,
		Values: c.copyValues(),
		Chain:  c,
	}
	c.onError(c.errCtx)
}

// copyValues creates a copy of the values map
func (c *Chain) copyValues() map[reflect.Type]interface{} {
	copy := make(map[reflect.Type]interface{})
	for k, v := range c.values {
		copy[k] = v
	}
	return copy
}

// storeValue stores a value by its type
// Values of the same type will be overridden
func (c *Chain) storeValue(val interface{}) {
	if val == nil {
		return
	}
	t := reflect.TypeOf(val)
	c.values[t] = val
}

// GetAll returns all stored values mapped by their types
func (c *Chain) GetAll() map[reflect.Type]interface{} {
	return c.values
}

// GetError returns the current error if any
func (c *Chain) GetError() error {
	return c.err
}

// GetErrorContext returns the error context if any
func (c *Chain) GetErrorContext() *ErrorContext {
	return c.errCtx
}

// Reset clears all values and error from the chain
func (c *Chain) Reset() *Chain {
	c.values = make(map[reflect.Type]interface{})
	c.err = nil
	c.errCtx = nil
	return c
}

// IsSuccess returns true if no error has occurred
func (c *Chain) IsSuccess() bool {
	return c.err == nil
}

// IsError returns true if an error has occurred
func (c *Chain) IsError() bool {
	return c.err != nil
}

// Value returns a value by its reflect.Type (for dynamic type access)
func (c *Chain) Value(t reflect.Type) interface{} {
	return c.values[t]
}
