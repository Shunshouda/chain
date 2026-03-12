package chain

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewChain(t *testing.T) {
	// Example 1: Basic error handling with context
	fmt.Println("=== Example 1: Error Handling with Context ===")
	c1 := NewChain()

	c1.OnError(func(ctx *ErrorContext) {
		fmt.Printf("Error occurred: %v\n", ctx.Error)

		// Access values available at the time of error
		if user, err := GetFromError[string](ctx); err == nil {
			fmt.Printf("User at error time: %s\n", user)
		}

		if count, err := GetFromError[int](ctx); err == nil {
			fmt.Printf("Count at error time: %d\n", count)
		}

		// You can also access the chain itself for recovery
		fmt.Println("Chain can be recovered or logged")
	})

	c1.Then(func() (string, error) {
		return "john_doe", nil
	}).Then(func(username string) (int, error) {
		fmt.Printf("Processing user: %s\n", username)
		return len(username), nil
	}).Then(func(length int) (bool, error) {
		if length < 5 {
			return false, fmt.Errorf("username too short: %d", length)
		}
		return true, nil
	})

	// Example 2: Different error types with context
	fmt.Println("\n=== Example 2: Different Error Types ===")
	c2 := NewChain()

	type User struct {
		ID    int
		Name  string
		Email string
	}

	type Config struct {
		MaxRetries int
		Timeout    int
	}

	c2.OnError(func(ctx *ErrorContext) {
		fmt.Printf("Error: %v\n", ctx.Error)

		// Check what values were available
		if user, err := GetFromError[User](ctx); err == nil {
			fmt.Printf("User at error: %+v\n", user)
		}

		if config, err := GetFromError[Config](ctx); err == nil {
			fmt.Printf("Config at error: %+v\n", config)
		}

		// Make decisions based on available values
		if HasInError[User](ctx) && HasInError[Config](ctx) {
			fmt.Println("Both user and config available - can implement retry logic")
		}
	})

	c2.Then(func() (User, error) {
		return User{ID: 1, Name: "Alice", Email: "alice@example.com"}, nil
	}).Then(func(u User) (Config, error) {
		return Config{MaxRetries: 3, Timeout: 30}, nil
	}).Then(func(u User, cfg Config) (string, error) {
		// Simulate an error
		return "", fmt.Errorf("failed to process user %s with config %+v", u.Name, cfg)
	})

	// Example 3: Recovery strategy based on context
	fmt.Println("\n=== Example 3: Recovery Strategy ===")
	c3 := NewChain()

	type Transaction struct {
		ID     string
		Amount float64
		Status string
	}

	type RetryInfo struct {
		Attempts  int
		LastError error
	}

	c3.OnError(func(ctx *ErrorContext) {
		fmt.Printf("Transaction failed: %v\n", ctx.Error)

		// Get transaction details
		if tx, err := GetFromError[Transaction](ctx); err == nil {
			fmt.Printf("Failed transaction: %+v\n", tx)

			// Check retry info
			if retry, err := GetFromError[RetryInfo](ctx); err == nil {
				if retry.Attempts < 3 {
					fmt.Printf("Retrying transaction %s, attempt %d\n", tx.ID, retry.Attempts+1)
					// Here you could implement retry logic
				}
			} else {
				// First attempt
				fmt.Printf("First attempt failed for transaction %s\n", tx.ID)
			}
		}
	})

	c3.Then(func() (Transaction, error) {
		return Transaction{ID: "TXN123", Amount: 99.99, Status: "pending"}, nil
	}).Then(func(tx Transaction) (RetryInfo, error) {
		return RetryInfo{Attempts: 1, LastError: nil}, nil
	}).Then(func(tx Transaction, retry RetryInfo) (string, error) {
		// Simulate failure
		return "", fmt.Errorf("payment gateway timeout")
	})

	// Example 4: Complex error handling with multiple values
	fmt.Println("\n=== Example 4: Complex Error Handling ===")
	c4 := NewChain()

	type Request struct {
		Method string
		Path   string
		Body   []byte
	}

	type Response struct {
		StatusCode int
		Body       string
	}

	type Metrics struct {
		StartTime int64
		Duration  int64
	}

	c4.OnError(func(ctx *ErrorContext) {
		fmt.Println("=== Error Report ===")
		fmt.Printf("Error: %v\n", ctx.Error)

		// Log all available values for debugging
		fmt.Println("Available context at error time:")
		for typ, val := range ctx.Values {
			fmt.Printf("  Type: %v, Value: %+v\n", typ, val)
		}

		// Make decisions based on available data
		if req, ok := ctx.Values[reflect.TypeOf(Request{})]; ok {
			fmt.Printf("Request that caused error: %+v\n", req)
		}

		if metrics, ok := ctx.Values[reflect.TypeOf(Metrics{})]; ok {
			fmt.Printf("Metrics at error: %+v\n", metrics)
		}

		fmt.Println("===================")
	})

	c4.Then(func() (Request, error) {
		return Request{
			Method: "POST",
			Path:   "/api/users",
			Body:   []byte(`{"name":"John"}`),
		}, nil
	}).Then(func(r Request) (Metrics, error) {
		return Metrics{StartTime: 123456789, Duration: 0}, nil
	}).Then(func(r Request, m Metrics) (Response, error) {
		// Simulate validation error
		return Response{}, fmt.Errorf("invalid request body: missing required field 'email'")
	})

	// Example 5: Error handling in data pipeline
	fmt.Println("\n=== Example 5: Data Pipeline Error Handling ===")
	c5 := NewChain()

	type DataSource struct {
		Name    string
		Records int
	}

	type ProcessedData struct {
		ValidCount   int
		InvalidCount int
		Errors       []string
	}

	c5.OnError(func(ctx *ErrorContext) {
		source, _ := GetFromError[DataSource](ctx)
		processed, _ := GetFromError[ProcessedData](ctx)

		fmt.Printf("Pipeline failed at source: %s\n", source.Name)
		fmt.Printf("Processed before failure: %+v\n", processed)
		fmt.Printf("Error: %v\n", ctx.Error)

		// Implement partial success handling
		if processed.ValidCount > 0 {
			fmt.Printf("Partially successful: %d records processed\n", processed.ValidCount)
		}
	})

	c5.Then(func() (DataSource, error) {
		return DataSource{Name: "users.csv", Records: 1000}, nil
	}).Then(func(ds DataSource) (ProcessedData, error) {
		return ProcessedData{ValidCount: 500, InvalidCount: 0, Errors: []string{}}, nil
	}).Then(func(ds DataSource, pd ProcessedData) (int, error) {
		// Simulate failure at 50% processing
		return 0, fmt.Errorf("disk full at record 500")
	})

	// Example 6: Using error context for monitoring
	fmt.Println("\n=== Example 6: Monitoring Integration ===")
	c6 := NewChain()

	type ServiceInfo struct {
		Name     string
		Version  string
		Endpoint string
	}

	type Trace struct {
		TraceID    string
		SpanID     string
		ParentSpan string
	}

	c6.OnError(func(ctx *ErrorContext) {
		// Extract monitoring data
		service, _ := GetFromError[ServiceInfo](ctx)
		trace, _ := GetFromError[Trace](ctx)

		// Send to monitoring system
		fmt.Printf("Monitoring Alert:\n")
		fmt.Printf("  Service: %s v%s\n", service.Name, service.Version)
		fmt.Printf("  Trace: %s/%s\n", trace.TraceID, trace.SpanID)
		fmt.Printf("  Error: %v\n", ctx.Error)
		fmt.Printf("  Timestamp: %d\n", getTimestamp())
	})

	c6.Then(func() (ServiceInfo, error) {
		return ServiceInfo{Name: "auth-service", Version: "1.2.3", Endpoint: "/login"}, nil
	}).Then(func(s ServiceInfo) (Trace, error) {
		return Trace{TraceID: "trace-123", SpanID: "span-456", ParentSpan: "span-123"}, nil
	}).Then(func(s ServiceInfo, t Trace) (bool, error) {
		return false, fmt.Errorf("rate limit exceeded for endpoint %s", s.Endpoint)
	})
}

func getTimestamp() int64 {
	return 1234567890
}
