---
name: entdomain
description: Work with EntDomain — the Ent DDD code generation extension. Use when creating Ent schemas with domain annotations, understanding generated code, debugging codegen issues, or asking about EntDomain annotation builders, field scopes, ID types, or the generic service layer.
---

# EntDomain — Ent DDD Code Generation Extension

You are assisting a developer working with EntDomain, an Ent Framework extension that generates DDD layer code (domain models, repositories, services) from annotated schemas.

## Core Concepts

### Field Scopes

Scopes control which **handler-layer DTOs** include a field. They do NOT restrict service or repository layer access.

| Scope | Constant | Affects |
|-------|----------|---------|
| Create | `ScopeCreate` | `{Entity}CreateRequest` struct |
| Update | `ScopeUpdate` | `{Entity}UpdateRequest` struct |
| Query | `ScopeQuery` | `{Entity}QueryParams` struct |
| Response | `ScopeResponse` | `{Entity}Response` struct |

### Annotation Builders (choose one as the base)

| Builder | Scopes | Use For |
|---------|--------|---------|
| `DefaultField()` | create, update, query, response + searchable, filterable, sortable | Most business fields (name, email, status) |
| `InputOnlyField()` | create, update + sensitive | Password, secrets |
| `OutputOnlyField()` | query, response + searchable, filterable, sortable | System fields (timestamps) |
| `CreateOnlyField()` | create, query, response + searchable, filterable, sortable | Immutable after creation (creator_id) |
| `IdField()` | query, response + read-only | Primary key |
| `AuditLogField()` | query, response + read-only | Audit fields |
| `DomainFieldWithScopes(...)` | custom | Any custom combination |

### Fluent Methods (chain after base builder)

**Scope & validation:**
- `.WithRequired(scope)` — required in that scope's DTO
- `.WithValidation(rules)` — validation rules map

**Capabilities:**
- `.AsSearchable()` — enable text search
- `.AsFilterable()` — enable filtering
- `.AsSortable()` — enable sorting
- `.AsUniqueLookup()` — generate `FindByX` returning single result
- `.AsRangeLookup()` — generate `FindByXRange` for time/numeric fields
- `.AsSensitive()` — mark as sensitive (excluded from responses)

**Metadata (for future OpenAPI generation):**
- `.WithDescription(desc)`, `.WithExample(val)`, `.WithFormat(fmt)`
- `.WithPattern(regex)`, `.WithRange(min, max *float64)`, `.WithLength(min, max *int)`
- `.WithTitle(title)`, `.WithEnum(values...)`, `.WithTags(tags...)`
- `.AsReadOnly()`, `.AsWriteOnly()`, `.AsDeprecated()`

## Schema Pattern

```go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "github.com/githonllc/entdomain"
)

type User struct { ent.Schema }

func (User) Fields() []ent.Field {
    return []ent.Field{
        field.Int64("id").Unique().Immutable().
            Annotations(entdomain.IdField()),

        field.String("name").NotEmpty().
            Annotations(entdomain.DefaultField().
                WithRequired(entdomain.ScopeCreate).
                AsSearchable()),

        field.String("email").Optional().
            Annotations(entdomain.DefaultField().
                WithFormat("email").
                AsUniqueLookup()),

        field.String("password").Sensitive().
            Annotations(entdomain.InputOnlyField()),

        field.Time("created_at").Default(time.Now).Immutable().
            Annotations(entdomain.OutputOnlyField().
                AsRangeLookup()),
    }
}
```

## Generated Files Per Entity

For entity `User`, three files are generated:

### `user_domain_model.go`
- `UserDomainModel` — full domain model implementing `entdomain.DomainModel`
- `UserCreateRequest` — fields with `ScopeCreate`, implements `entdomain.CreateRequest`
- `UserUpdateRequest` — fields with `ScopeUpdate` as pointer types (partial updates), implements `entdomain.UpdateRequest`
- `UserResponse` — fields with `ScopeResponse`
- `UserQueryParams` — fields with `ScopeQuery`, includes `form` tags for gin `ShouldBindQuery`
- `UserListResponse` — paginated response wrapper with optional `PageInfo` for cursor pagination
- Conversion methods: `ToDomainModel()`, `ApplyToDomainModel()`, `ToResponse()` (uses `Ptr`/`PtrOrNil`/`PtrTimeOrNil`), `Clone()`

### `user_domain_repository.go`
- `UserRepository` interface — CRUD, batch ops, List, ListByCursor, Search, Count, Exists, WithTx
- `UserRepositoryImpl` — Ent ORM implementation with typed error wrapping
- Field-specific methods from annotations:
  - `AsUniqueLookup()` → `FindByEmail(ctx, email string) (*UserDomainModel, error)`
  - `AsRangeLookup()` → `FindByCreatedAtRange(ctx, start, end time.Time) ([]*UserDomainModel, error)`
- **Error wrapping**: `GetByID`/`FindOneBy` → `ErrNotFound`, `Create`/`Update` → `ErrAlreadyExists`
- **Create skips Default fields**: fields with `Default(time.Now)` etc. are not set in Create (Ent's default hook runs)
- **Keyset pagination**: `ListByCursor` uses `sql.CompositeGT`/`CompositeLT` for O(log n) index seek

### `user_domain_service.go`
- `UserService` interface — mirrors repository + `ListByCursor`
- `DefaultUserService` — pure delegation to repository, embeddable for selective override
- NO generated implementation with TODO stubs (removed — users compose their own logic)

## Typed Errors

Repositories wrap Ent errors with standard sentinels:

```go
// Sentinels in entdomain package
entdomain.ErrNotFound      // entity not found (wrapped from ent.IsNotFound)
entdomain.ErrAlreadyExists // uniqueness constraint (wrapped from ent.IsConstraintError)
entdomain.ErrValidation    // validation failed

// Check with errors.Is or helpers
if entdomain.IsNotFound(err) { ... }
if entdomain.IsAlreadyExists(err) { ... }
```

**Important**: `FindOneBy` and `FindByX` return `ErrNotFound` when no results found (NOT `nil, nil`). Always check `err` first.

## DefaultService Pattern

Generated `Default{Entity}Service` delegates all methods to the repository. Embed it and override only what needs custom logic:

```go
// Simple entity: use DefaultService directly
type auditLogService struct {
    ent.DefaultAuditLogService
}

// Complex entity: embed and override
type personService struct {
    ent.DefaultPersonService
}

func (s *personService) Create(ctx context.Context, model *ent.PersonDomainModel) (*ent.PersonDomainModel, error) {
    if model.Email == "" {
        return nil, entdomain.ErrValidation
    }
    return s.DefaultPersonService.Create(ctx, model)
}
```

## Pagination

Two modes supported simultaneously:

### Offset Pagination (traditional)
```go
items, total, err := repo.List(ctx, &entdomain.ListRequest{
    Page: 2, Size: 20, SortBy: "created_at", Order: "desc",
})
```

### Cursor Pagination (keyset, O(log n))
```go
// First page
items, pageInfo, err := repo.ListByCursor(ctx, &entdomain.ListRequest{
    Size: 20, SortBy: "created_at", Order: "desc",
})
// pageInfo.HasNextPage, pageInfo.EndCursor

// Next page
items, pageInfo, err = repo.ListByCursor(ctx, &entdomain.ListRequest{
    Cursor: pageInfo.EndCursor, Size: 20, SortBy: "created_at", Order: "desc",
})
```

Cursor is opaque (base64 JSON), contains sort field value + entity ID for tie-breaking.

## QueryParams Binding

Generated QueryParams include `form` tags for automatic gin binding:

```go
// Instead of manual parsing:
var params ent.PersonQueryParams
if err := c.ShouldBindQuery(&params); err != nil { ... }
// params.Query, params.Page, params.Size are auto-parsed
```

`ListRequest` also has `form` tags for direct binding.

## Pointer Helpers (Runtime)

Used in generated `ToResponse()` for optional field conversion:

```go
entdomain.Ptr(value)          // generic: returns *T
entdomain.PtrOrNil(s)         // string: returns nil if empty
entdomain.PtrTimeOrNil(t)     // time.Time: returns nil if zero
```

## Extension Setup (`entc.go`)

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
        // entdomain.WithRepository(false),  // skip repository generation
        // entdomain.WithService(false),      // skip service generation
        // entdomain.WithEntDomainPackage("custom/import/path"),
    )
    if err := entc.Generate("./schema", &gen.Config{}, entc.Extensions(ext)); err != nil {
        log.Fatal(err)
    }
}
```

## ID System

```go
// Two backing types
id := entdomain.NewIDFromInt64(12345)   // Int64ID
id := entdomain.NewIDFromString("abc")  // StringID

// Interface methods
id.String()    // string representation
id.Int64()     // (int64, error)
id.IsZero()    // bool
```

## Generic Base Service

For handler-layer CRUD with automatic DTO validation and conversion:

```go
svc := entdomain.NewBaseGenericDomainService[
    *UserDomainModel,
    *UserCreateRequest,
    *UserUpdateRequest,
    *UserResponse,
    *UserListResponse,
    *UserQueryParams,
](repo, entdomain.Converters[*UserDomainModel, *UserResponse, *UserListResponse]{
    ToResponse:     func(m *UserDomainModel) *UserResponse { return m.ToResponse() },
    ToListResponse: func(models []*UserDomainModel, total, page, size int) *UserListResponse { ... },
})
```

Note: `BaseGenericDomainService` operates on DTOs (CreateRequest in / Response out) — different from the generated `{Entity}Service` interface which operates on DomainModel. Use `BaseGenericDomainService` in the handler layer; use `DefaultService` in the domain layer.

## Common Issues

1. **After schema changes, always run `make gen`** (or `go generate ./...`) to regenerate code
2. **Never edit generated files** — they are overwritten on each generation
3. **ID type mismatch** — ensure schema ID type matches generated code expectations (int64 vs string)
4. **Missing scope** — if a field doesn't appear in a request/response DTO, check its annotation scopes
5. **Immutable fields** — fields with `Immutable()` are excluded from `updateableFields` in repository layer
6. **Default fields** — fields with `Default(time.Now)` are skipped in repository `Create` (Ent's default hook applies)
7. **FindOneBy returns ErrNotFound** — not `nil, nil`. Always check `err` before using the result
8. **Cursor pagination cursor** — opaque string, do not parse client-side. Pass `pageInfo.EndCursor` as-is

## Architecture Reference

For implementation details, see these source files:
- [annotations.go](annotations.go) — all annotation types and builders
- [types.go](types.go) — ID, DomainModel, Request/Response interfaces, ListRequest, Ptr helpers
- [errors.go](errors.go) — ErrNotFound, ErrAlreadyExists, ErrValidation sentinels
- [cursor.go](cursor.go) — Cursor, PageInfo, EncodeCursor/DecodeCursor
- [service.go](service.go) — Repository and GenericDomainService interfaces
- [extension.go](extension.go) — Extension configuration and options
- [templates/](templates/) — Go templates for code generation
