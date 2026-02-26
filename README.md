# EntDomain

[![Go Reference](https://pkg.go.dev/badge/github.com/githonllc/entdomain.svg)](https://pkg.go.dev/github.com/githonllc/entdomain)
[![Go Report Card](https://goreportcard.com/badge/github.com/githonllc/entdomain)](https://goreportcard.com/report/github.com/githonllc/entdomain)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

An [Ent](https://entgo.io) extension that generates DDD (Domain-Driven Design) layer code from annotated schemas — domain models, repositories, and services.

## Features

- **Annotation-driven** — mark field scopes with concise builders (`DefaultField`, `InputOnlyField`, `OutputOnlyField`, etc.)
- **Full DDD stack** — generates domain models, create/update requests, responses, query params, repository interfaces, and service interfaces
- **Type-safe ID system** — `ID` interface supporting string and int64 backing types
- **Generic base service** — `BaseGenericDomainService` with CRUD, pagination, and search out of the box
- **Deep copy & partial updates** — generated `Clone()` and `ApplyToDomainModel()` methods

## Requirements

- Go 1.23+
- [Ent](https://entgo.io) v0.14+

## Installation

```bash
go get github.com/githonllc/entdomain
```

## Setup

Wire the extension in your `entc.go`:

```go
//go:build ignore

package main

import (
    "log"

    "entgo.io/ent/entc"
    "entgo.io/ent/entc/gen"
    "github.com/githonllc/entdomain"
)

func main() {
    ext := entdomain.NewExtensionWithOptions(
        entdomain.WithOutputDir("./ent"),
        entdomain.WithPackageName("ent"),
    )

    if err := entc.Generate("./schema", &gen.Config{}, entc.Extensions(ext)); err != nil {
        log.Fatal(err)
    }
}
```

Then run:

```bash
go generate ./...
```

## Annotation Builders

### Base Builders

```go
entdomain.DefaultField()                      // all scopes: create, update, query, response
entdomain.InputOnlyField()                    // create + update only (e.g., password)
entdomain.OutputOnlyField()                   // query + response only (e.g., timestamps)
entdomain.CreateOnlyField()                   // create + query + response (immutable after creation)
entdomain.IdField()                           // response + query, read-only
entdomain.AuditLogField()                     // response + query, read-only
entdomain.DomainFieldWithScopes(scopes...)    // custom scope combination
```

### Fluent Builder API

```go
field.String("email").
    Annotations(
        entdomain.DefaultField().
            WithRequired(entdomain.ScopeCreate).
            WithDescription("Email address").
            WithFormat("email").
            AsSearchable().
            AsFilterable().
            AsUniqueLookup(),
    )
```

### Available Methods

| Method | Description |
|--------|-------------|
| `WithRequired(scope)` | Mark field as required in a scope |
| `WithDescription(desc)` | Set field description |
| `WithExample(val)` | Set example value |
| `WithFormat(fmt)` | Set format (email, date-time, uuid, etc.) |
| `WithPattern(regex)` | Set validation pattern |
| `WithRange(min, max *float64)` | Set numeric range (nil = unbounded) |
| `WithLength(min, max *int)` | Set string length (nil = unbounded) |
| `AsSearchable()` | Enable text search |
| `AsFilterable()` | Enable filtering |
| `AsSortable()` | Enable sorting |
| `AsUniqueLookup()` | Generate `FindByX` method |
| `AsRangeLookup()` | Generate `FindByXRange` method |
| `AsSensitive()` | Mark as sensitive |
| `AsReadOnly()` | Mark as read-only |

## Schema Example

```go
package schema

import (
    "time"

    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "github.com/githonllc/entdomain"
)

type User struct {
    ent.Schema
}

func (User) Fields() []ent.Field {
    return []ent.Field{
        field.Int64("id").
            Unique().
            Immutable().
            Annotations(entdomain.IdField()),

        field.String("name").
            NotEmpty().
            Annotations(
                entdomain.DefaultField().
                    WithRequired(entdomain.ScopeCreate).
                    AsSearchable(),
            ),

        field.String("email").
            Optional().
            Annotations(
                entdomain.DefaultField().
                    WithFormat("email").
                    AsUniqueLookup(),
            ),

        field.Time("created_at").
            Default(time.Now).
            Immutable().
            Annotations(
                entdomain.OutputOnlyField().
                    AsRangeLookup(),
            ),
    }
}
```

## Generated Code

For each annotated schema, three files are generated:

### Domain Model (`user_domain_model.go`)

```go
type UserDomainModel struct { ... }     // full domain model
type UserCreateRequest struct { ... }   // create request DTO
type UserUpdateRequest struct { ... }   // update request DTO (pointer fields for partial updates)
type UserResponse struct { ... }        // response DTO
type UserQueryParams struct { ... }     // search/filter parameters
type UserListResponse struct { ... }    // paginated response
```

Generated conversion methods:
- `ToDomainModel()` — request → domain model
- `ApplyToDomainModel()` — partial update on existing model
- `ToResponse()` — domain model → response (uses `Ptr`/`PtrOrNil`/`PtrTimeOrNil` helpers for optional fields)
- `Clone()` — deep copy

### Repository (`user_domain_repository.go`)

```go
type UserRepository interface {
    Create(ctx, model) (*UserDomainModel, error)
    GetByID(ctx, id) (*UserDomainModel, error)
    Update(ctx, model) (*UserDomainModel, error)
    Delete(ctx, id) error
    List(ctx, req) ([]*UserDomainModel, int, error)
    Search(ctx, req) ([]*UserDomainModel, int, error)
    FindByEmail(ctx, email) (*UserDomainModel, error)  // from AsUniqueLookup
    FindByCreatedAtRange(ctx, start, end) ([]*UserDomainModel, error)  // from AsRangeLookup
    // ... batch operations, count, exists
}
```

### Service (`user_domain_service.go`)

The service template generates an **interface only** — no implementation. This lets you compose `BaseGenericDomainService` with your own business logic:

```go
// Generated interface (lean, no boilerplate impl):
type UserService interface {
    Create(ctx, model) (*UserDomainModel, error)
    GetByID(ctx, id) (*UserDomainModel, error)
    Update(ctx, model) (*UserDomainModel, error)
    Delete(ctx, id) error
    List(ctx, req) ([]*UserDomainModel, int, error)
    Search(ctx, req) ([]*UserDomainModel, int, error)
    Count(ctx, req) (int, error)
    Exists(ctx, id) (bool, error)
    CreateBatch(ctx, models) ([]*UserDomainModel, error)
    UpdateBatch(ctx, models) ([]*UserDomainModel, error)
    DeleteBatch(ctx, ids) error
    FindByEmail(ctx, email) (*UserDomainModel, error)       // from AsUniqueLookup
    FindByCreatedAtRange(ctx, start, end) ([]*UserDomainModel, error)  // from AsRangeLookup
}
```

A `DefaultUserService` is also generated — embed it and override only what needs custom logic:

```go
type userServiceImpl struct {
    ent.DefaultUserService  // embed for default CRUD delegation
}

// Override Create to add custom validation:
func (s *userServiceImpl) Create(ctx context.Context, model *ent.UserDomainModel) (*ent.UserDomainModel, error) {
    if model.Email == "" {
        return nil, entdomain.ErrValidation
    }
    return s.DefaultUserService.Create(ctx, model)
}
```

## Typed Errors

Repositories wrap Ent errors with standard sentinel values:

```go
var (
    entdomain.ErrNotFound      // entity not found
    entdomain.ErrAlreadyExists // uniqueness constraint violation
    entdomain.ErrValidation    // validation failed
)

// Check errors with errors.Is or helpers:
if entdomain.IsNotFound(err) { ... }
if entdomain.IsAlreadyExists(err) { ... }
```

## QueryParams Binding

Generated `QueryParams` structs include `form` tags for automatic gin binding:

```go
var params ent.UserQueryParams
if err := c.ShouldBindQuery(&params); err != nil { ... }
// params.Query, params.Page, params.Size are automatically parsed
```

## Cursor Pagination

For large datasets, use cursor-based (keyset) pagination instead of offset:

```go
// First page:
items, pageInfo, err := repo.ListByCursor(ctx, &entdomain.ListRequest{
    Size: 20, SortBy: "created_at", Order: "desc",
})
// pageInfo.HasNextPage == true
// pageInfo.EndCursor == "eyJpZ..." (opaque)

// Next page:
items, pageInfo, err = repo.ListByCursor(ctx, &entdomain.ListRequest{
    Cursor: pageInfo.EndCursor,
    Size: 20, SortBy: "created_at", Order: "desc",
})
```

Cursor pagination is O(log n) vs offset's O(n) — it seeks directly via index.
Offset pagination (`List()`) remains available for use cases that need page jumping.

## ID System

```go
// Create IDs
id := entdomain.NewIDFromInt64(12345)
id := entdomain.NewIDFromString("user-abc")

// Use IDs
str := id.String()       // "12345"
num, err := id.Int64()   // 12345, nil
zero := id.IsZero()      // false
```

## Generic Base Service

For handler-layer CRUD with automatic validation and conversion:

```go
service := entdomain.NewBaseGenericDomainService[
    *UserDomainModel,
    *UserCreateRequest,
    *UserUpdateRequest,
    *UserResponse,
    *UserListResponse,
    *UserQueryParams,
](repo, converters)

resp, err := service.Create(ctx, createReq)
resp, err := service.GetByID(ctx, id)
resp, err := service.Update(ctx, id, updateReq)
err := service.Delete(ctx, id)
```

## Field Scopes

Scopes control which handler-layer DTOs include a field. They do **not** restrict service or repository layer access.

| Scope | Description |
|-------|-------------|
| `ScopeCreate` | Field appears in create request |
| `ScopeUpdate` | Field appears in update request |
| `ScopeQuery` | Field appears in query params |
| `ScopeResponse` | Field appears in response |

## Extension Options

```go
entdomain.WithOutputDir("./ent")              // output directory
entdomain.WithPackageName("ent")              // Go package name
entdomain.WithRepository(true)                // generate repository (default: true)
entdomain.WithService(true)                   // generate service (default: true)
entdomain.WithEntDomainPackage("custom/path") // override entdomain import path
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

## License

[MIT](LICENSE)
