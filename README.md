# Chain - Fluent Method Chaining Library

## Project Introduction

Chain is a Go language library that provides fluent method chaining capabilities with automatic parameter injection and error handling. It allows you to build complex operation sequences in a declarative manner without manually handling intermediate results and errors.

## Features

- **Fluent Method Chaining**: Use the `.Then()` method to chain multiple function calls
- **Automatic Parameter Injection**: Automatically retrieve required values from the chain based on function parameter types
- **Error Handling**: Automatically stop execution and propagate errors when they occur in the chain
- **Type Safety**: Provide generic helper functions to ensure type-safe operations
- **Flexible Data Storage**: Store and retrieve values through type mapping
- **Rich Helper Functions**: Including `Get`, `GetOrDefault`, `MustGet`, `Map`, `FlatMap`, etc.

## Installation

```bash
go get github.com/shunshouda/chain
```

## Basic Usage

### Creating a Chain

```go
import "github.com/shunshouda/chain"

c := chain.NewChain()
```

### Basic Chaining

```go
c.Then(func() (string, error) {
    return "42", nil
}).Then(func(s string) (int, error) {
    return strconv.Atoi(s)
}).Then(func(i int) float64 {
    return float64(i) * 1.5
})
```

### Error Handling

```go
c := chain.NewChain()

c.OnError(func(ctx *ErrorContext) {
    fmt.Printf("Error: %v\n", ctx.Error)
    
    // Access any values that were available when the error occurred
    if user, err := GetFromError[User](ctx); err == nil {
        fmt.Printf("User at error time: %+v\n", user)
    }
    
    if config, err := GetFromError[Config](ctx); err == nil {
        fmt.Printf("Config at error time: %+v\n", config)
    }
    
    // Check if specific types exist
    if HasInError[Transaction](ctx) {
        fmt.Println("Transaction data available for recovery")
    }
    
    // Access the chain itself for potential recovery
    fmt.Printf("Chain state: %v\n", ctx.Chain.IsSuccess())
})

c.Then(func() (User, error) {
    return User{ID: 1, Name: "Alice"}, nil
}).Then(func(u User) (Config, error) {
    return Config{MaxRetries: 3}, nil
}).Then(func(u User, cfg Config) (string, error) {
    return "", fmt.Errorf("processing failed")
})
```

### Type-Safe Value Retrieval

```go
// Get value, returns error if not found or chain has error
if str, err := chain.Get[string](c); err == nil {
    fmt.Printf("String value: %s\n", str)
}

// Get value, returns default value if not found
defaultVal := chain.GetOrDefault(c, "default")

// Get value, panics if not found
str := chain.MustGet[string](c)

// Check if value exists
hasString := chain.Has[string](c)

// Clear value of specific type
chain.Clear[string](c)
```

### Data Transformation

```go
// Use Map to transform value
chain.Map(c, func(i int) string {
    return fmt.Sprintf("Number: %d", i)
})

// Use FlatMap to transform value (supports error return)
chain.FlatMap(c, func(s string) (int, error) {
    return len(s), nil
})
```

### Complex Business Logic

```go
type Customer struct {
    ID    string
    Email string
}

type Product struct {
    Name  string
    Price float64
}

type Invoice struct {
    CustomerID string
    Total      float64
    Items      []Product
}

c := chain.NewChain()
c.Then(func() (Customer, error) {
    return Customer{ID: "C001", Email: "customer@example.com"}, nil
}).Then(func(c Customer) ([]Product, error) {
    products := []Product{
        {Name: "Laptop", Price: 999.99},
        {Name: "Mouse", Price: 29.99},
    }
    return products, nil
}).Then(func(c Customer, products []Product) (Invoice, error) {
    total := 0.0
    for _, p := range products {
        total += p.Price
    }
    return Invoice{
        CustomerID: c.ID,
        Total:      total,
        Items:      products,
    }, nil
})

if invoice, err := chain.Get[Invoice](c); err == nil {
    fmt.Printf("Invoice total: $%.2f\n", invoice.Total)
}
```

## Core API

### Main Methods

- `NewChain()`: Create a new chain instance
- `Then(fn interface{}) *Chain`: Add a function to the chain and execute it
- `OnError(callback func(error)) *Chain`: Set a custom error handling function
- `GetAll() map[reflect.Type]interface{}`: Return all stored values
- `GetError() error`: Return the current error
- `Reset() *Chain`: Reset the chain's state
- `IsSuccess() bool`: Check if the chain is successful
- `IsError() bool`: Check if the chain has an error

### Helper Functions

- `Get[T any](c *Chain) (T, error)`: Get value of specified type
- `GetOrDefault[T any](c *Chain, defaultValue T) T`: Get value of specified type or default value
- `MustGet[T any](c *Chain) T`: Get value of specified type or panic
- `Has[T any](c *Chain) bool`: Check if value of specified type exists
- `Clear[T any](c *Chain) *Chain`: Clear value of specified type
- `Map[T1, T2 any](c *Chain, fn func(T1) T2) *Chain`: Transform value
- `FlatMap[T1, T2 any](c *Chain, fn func(T1) (T2, error)) *Chain`: Transform value with error handling

## Working Principle

1. **Value Storage**: The chain uses `map[reflect.Type]interface{}` to store values, with types as keys
2. **Parameter Injection**: When executing `Then()`, the chain automatically looks up and injects required values based on function parameter types
3. **Error Handling**: If a function returns an error (last return value is error type), the chain captures and stores the error, and subsequent `Then()` calls are skipped
4. **Value Passing**: Function return values (except errors) are stored in the chain for use by subsequent functions

## Application Scenarios

- **Data Processing Pipeline**: Build complex data transformation and processing flows
- **Business Logic Orchestration**: Organize business logic steps in a declarative manner
- **Dependency Injection**: Automatically inject function dependencies
- **Error Handling**: Centralize error handling and avoid nested error checks

## Examples

Check the `chain_test.go` file for more detailed usage examples.

## Contribution

Welcome to submit Issues and Pull Requests to improve this library.
