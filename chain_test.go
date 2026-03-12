package chain

import (
	"fmt"
	"strconv"
	"testing"
)

func TestNewChain(t *testing.T) {
	// Example 1: Basic chain operations
	fmt.Println("=== Example 1: Basic Chain Operations ===")
	c1 := NewChain()

	c1.Then(func() (string, error) {
		return "42", nil
	}).Then(func(s string) (int, error) {
		return strconv.Atoi(s)
	}).Then(func(i int) float64 {
		return float64(i) * 1.5
	}).Then(func(f float64) string {
		return fmt.Sprintf("Result: %.2f", f)
	})

	// Retrieve values using helper functions
	if str, err := Get[string](c1); err == nil {
		fmt.Printf("String value: %s\n", str)
	}

	if num, err := Get[int](c1); err == nil {
		fmt.Printf("Integer value: %d\n", num)
	}

	if f, err := Get[float64](c1); err == nil {
		fmt.Printf("Float value: %.2f\n", f)
	}

	// Example 2: Error handling
	fmt.Println("\n=== Example 2: Error Handling ===")
	c2 := NewChain()
	c2.OnError(func(err error) {
		fmt.Printf("Custom error handler: %v\n", err)
	})

	c2.Then(func() (int, error) {
		return 0, fmt.Errorf("simulated error")
	}).Then(func(i int) string {
		fmt.Println("This won't execute")
		return "won't execute"
	})

	if c2.IsError() {
		fmt.Printf("Chain has error: %v\n", c2.GetError())
	}

	// Example 3: Working with custom types
	fmt.Println("\n=== Example 3: Custom Types ===")
	type User struct {
		Name string
		Age  int
	}

	type Order struct {
		ID     string
		Amount float64
	}

	c3 := NewChain()
	c3.Then(func() (User, error) {
		return User{Name: "Alice", Age: 30}, nil
	}).Then(func(u User) (Order, error) {
		return Order{ID: "ORD001", Amount: 99.99}, nil
	}).Then(func(u User, o Order) string {
		return fmt.Sprintf("User %s has order %s for $%.2f",
			u.Name, o.ID, o.Amount)
	})

	if user, err := Get[User](c3); err == nil {
		fmt.Printf("User: %+v\n", user)
	}

	if order, err := Get[Order](c3); err == nil {
		fmt.Printf("Order: %+v\n", order)
	}

	if msg, err := Get[string](c3); err == nil {
		fmt.Printf("Message: %s\n", msg)
	}

	// Example 4: Data transformation pipeline
	fmt.Println("\n=== Example 4: Data Transformation ===")
	c4 := NewChain()
	c4.Then(func() []int {
		return []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	}).Then(func(nums []int) ([]int, error) {
		// Filter even numbers
		var result []int
		for _, n := range nums {
			if n%2 == 0 {
				result = append(result, n)
			}
		}
		return result, nil
	}).Then(func(nums []int) int {
		// Calculate sum
		sum := 0
		for _, n := range nums {
			sum += n
		}
		return sum
	}).Then(func(sum int) float64 {
		// Calculate average
		return float64(sum) / 5.0 // 5 even numbers from 1-10
	})

	if avg, err := Get[float64](c4); err == nil {
		fmt.Printf("Average of even numbers: %.2f\n", avg)
	}

	// Example 5: Multiple return values
	fmt.Println("\n=== Example 5: Multiple Return Values ===")
	c5 := NewChain()
	c5.Then(func() (string, int, error) {
		return "score", 100, nil
	}).Then(func(name string, score int) (bool, string) {
		passed := score >= 60
		return passed, fmt.Sprintf("%s: %d - %v", name, score, passed)
	})

	if passed, err := Get[bool](c5); err == nil {
		fmt.Printf("Passed: %v\n", passed)
	}

	if result, err := Get[string](c5); err == nil {
		fmt.Printf("Result: %s\n", result)
	}

	// Example 6: Utility methods
	fmt.Println("\n=== Example 6: Utility Methods ===")
	c6 := NewChain()
	c6.Then(func() string {
		return "hello"
	}).Then(func(s string) int {
		return len(s)
	})

	// Check if types exist using helper functions
	fmt.Printf("Has string? %v\n", Has[string](c6))
	fmt.Printf("Has int? %v\n", Has[int](c6))
	fmt.Printf("Has bool? %v\n", Has[bool](c6))

	// Get with default value
	defaultVal := GetOrDefault(c6, "default")
	fmt.Printf("GetOrDefault: %s\n", defaultVal)

	// Clear a type using helper function
	Clear[string](c6)
	fmt.Printf("After clearing, has string? %v\n", Has[string](c6))

	// Reset everything
	c6.Reset()
	fmt.Printf("After reset, has int? %v\n", Has[int](c6))

	// Example 7: Using Map and FlatMap
	fmt.Println("\n=== Example 7: Map and FlatMap ===")
	c7 := NewChain()
	c7.Then(func() int {
		return 5
	})

	// Use Map helper
	Map(c7, func(i int) string {
		return fmt.Sprintf("Number: %d", i)
	})

	if str, err := Get[string](c7); err == nil {
		fmt.Printf("Mapped value: %s\n", str)
	}

	// Use FlatMap helper
	FlatMap(c7, func(s string) (int, error) {
		return len(s), nil
	})

	if length, err := Get[int](c7); err == nil {
		fmt.Printf("FlatMapped value: %d\n", length)
	}

	// Example 8: Complex business logic
	fmt.Println("\n=== Example 8: Complex Business Logic ===")
	type Product struct {
		Name  string
		Price float64
	}

	type Customer struct {
		ID    string
		Email string
	}

	type Invoice struct {
		CustomerID string
		Total      float64
		Items      []Product
	}

	c8 := NewChain()
	c8.Then(func() (Customer, error) {
		return Customer{ID: "C001", Email: "customer@example.com"}, nil
	}).Then(func(c Customer) ([]Product, error) {
		products := []Product{
			{Name: "Laptop", Price: 999.99},
			{Name: "Mouse", Price: 29.99},
			{Name: "Keyboard", Price: 89.99},
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
	}).Then(func(inv Invoice) (string, error) {
		return fmt.Sprintf("Invoice total: $%.2f with %d items",
			inv.Total, len(inv.Items)), nil
	})

	if summary, err := Get[string](c8); err == nil {
		fmt.Printf("Business Result: %s\n", summary)
	}

	if invoice, err := Get[Invoice](c8); err == nil {
		fmt.Printf("Invoice Details: %+v\n", invoice)
	}
}
