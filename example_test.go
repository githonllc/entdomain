package entdomain_test

import (
	"fmt"

	"github.com/githonllc/entdomain"
)

func ExampleDefaultField() {
	df := entdomain.DefaultField()
	fmt.Println(df.Searchable)
	fmt.Println(df.Filterable)
	fmt.Println(df.Sortable)
	fmt.Println(len(df.Scopes))
	// Output:
	// true
	// true
	// true
	// 4
}

func ExampleNewIDFromInt64() {
	id := entdomain.NewIDFromInt64(42)
	fmt.Println(id.String())
	fmt.Println(id.IsZero())

	val, _ := id.Int64()
	fmt.Println(val)
	// Output:
	// 42
	// false
	// 42
}

func ExampleNewIDFromString() {
	id := entdomain.NewIDFromString("user-abc")
	fmt.Println(id.String())
	fmt.Println(id.IsZero())
	// Output:
	// user-abc
	// false
}

func ExampleListRequest_SetDefaults() {
	req := &entdomain.ListRequest{}
	req.SetDefaults()
	fmt.Println(req.Size)
	// Output:
	// 20
}

func ExampleDomainFieldWithScopes() {
	df := entdomain.DomainFieldWithScopes(
		entdomain.ScopeCreate,
		entdomain.ScopeResponse,
	)
	fmt.Println(len(df.Scopes))
	// Output:
	// 2
}

func ExampleDomainField_chaining() {
	df := entdomain.DefaultField().
		WithRequired(entdomain.ScopeCreate).
		WithDescription("User email address").
		WithFormat("email").
		AsUniqueLookup()

	fmt.Println(df.Description)
	fmt.Println(df.UniqueLookup)
	fmt.Println(df.Required[entdomain.ScopeCreate])
	// Output:
	// User email address
	// true
	// true
}
