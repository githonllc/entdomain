package entdomain

import (
	"fmt"
	"strconv"
	"time"
)

const (
	// DefaultPageSize is the default number of items per page when not specified or invalid.
	DefaultPageSize = 20

	// MaxPageSize is the maximum allowed number of items per page.
	MaxPageSize = 1000
)

// ID is the entity identifier interface. Concrete implementations (StringID,
// Int64ID) allow domain code to remain agnostic of the underlying key type.
type ID interface {
	// String returns a human-readable representation of the identifier.
	String() string
	// IsZero reports whether the ID is the zero value for its type.
	IsZero() bool
	// Int64 returns the numeric value of the ID, or an error if conversion is not possible.
	Int64() (int64, error)
}

// StringID is an ID backed by a string value.
type StringID string

// String returns the string representation of the ID
func (id StringID) String() string {
	return string(id)
}

// IsZero returns true if the ID is empty/zero
func (id StringID) IsZero() bool {
	return string(id) == ""
}

// Int64 converts the ID to int64 (if possible)
func (id StringID) Int64() (int64, error) {
	return strconv.ParseInt(string(id), 10, 64)
}

// Int64ID is an ID backed by an int64 value (e.g., snowflake IDs).
type Int64ID int64

// String returns the string representation of the ID
func (id Int64ID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

// IsZero returns true if the ID is empty/zero
func (id Int64ID) IsZero() bool {
	return int64(id) == 0
}

// Int64 converts the ID to int64
func (id Int64ID) Int64() (int64, error) {
	return int64(id), nil
}

// NewIDFromString creates a StringID from the given string value.
func NewIDFromString(s string) ID {
	return StringID(s)
}

// NewIDFromInt64 creates an Int64ID from the given int64 value.
func NewIDFromInt64(i int64) ID {
	return Int64ID(i)
}

// Entity defines the identity contract for domain entities. Every domain
// entity must be identifiable and its ID must be gettable and settable.
type Entity interface {
	GetID() ID
	SetID(id ID)
}

// Cloner defines deep-copy capability for domain models.
type Cloner interface {
	Clone() DomainModel
}

// DomainModel combines Entity and Cloner — the base contract for all
// domain models operated on by the service and repository layers.
type DomainModel interface {
	Entity
	Cloner
}

// Validatable defines the validation contract. Types that embed this
// interface must return a non-nil error when their state is invalid.
type Validatable interface {
	Validate() error
}

// DomainConverter converts a request or response DTO to a DomainModel.
type DomainConverter interface {
	ToDomainModel() DomainModel
}

// DomainApplier applies partial updates from a request DTO onto an
// existing DomainModel, returning the modified model.
type DomainApplier interface {
	ApplyToDomainModel(domain DomainModel) DomainModel
}

// CreateRequest is the contract for handler-layer create request DTOs.
// Implementations validate input and convert to a domain model.
type CreateRequest interface {
	Validatable
	DomainConverter
}

// UpdateRequest is the contract for handler-layer update request DTOs.
// In addition to validation and conversion, it can apply partial updates
// to an existing domain model.
type UpdateRequest interface {
	Validatable
	DomainConverter
	DomainApplier
}

// QueryParams defines the interface for query parameter models.
// Implementations must provide ToSearchRequest to bridge typed query parameters
// to the generic SearchRequest used by the repository layer.
type QueryParams interface {
	Validatable
	ToSearchRequest() *SearchRequest
}

// ListRequest represents a paginated list request with optional sorting.
// Supports both offset-based (Page/Size) and cursor-based (Cursor/Size) pagination.
// When Cursor is set, keyset pagination is used; otherwise offset pagination applies.
type ListRequest struct {
	Size   int    `json:"size,omitempty" form:"size" validate:"omitempty,min=1,max=100"`
	Page   int    `json:"page,omitempty" form:"page" validate:"omitempty,min=0"`
	SortBy string `json:"sort_by,omitempty" form:"sort_by"`
	Order  string `json:"order,omitempty" form:"order" validate:"omitempty,oneof=asc desc"`
	Cursor string `json:"cursor,omitempty" form:"cursor"` // opaque cursor for keyset pagination
}

// SetDefaults fills in zero-valued fields with sensible defaults.
// Call this before using the request to ensure pagination works correctly.
func (r *ListRequest) SetDefaults() {
	if r.Size == 0 {
		r.Size = DefaultPageSize
	}
}

// Validate checks that all fields are within acceptable bounds.
// It does NOT modify the receiver — call SetDefaults first if needed.
func (r *ListRequest) Validate() error {
	if r == nil {
		return fmt.Errorf("list request cannot be nil")
	}

	if r.Size < 0 {
		return fmt.Errorf("size cannot be negative")
	}
	if r.Size > MaxPageSize {
		return fmt.Errorf("size cannot exceed %d", MaxPageSize)
	}

	if r.Page < 0 {
		return fmt.Errorf("page cannot be negative")
	}

	if r.Order != "" && r.Order != "asc" && r.Order != "desc" {
		return fmt.Errorf("order must be 'asc' or 'desc'")
	}

	return nil
}

// SearchRequest represents a search request with a free-text query, typed
// filters, pagination, and optional sorting.
type SearchRequest struct {
	Query   string         `json:"query,omitempty"`
	Filters map[string]any `json:"filters,omitempty"`
	Size    int            `json:"size,omitempty" validate:"omitempty,min=1,max=100"`
	Page    int            `json:"page,omitempty" validate:"omitempty,min=0"`
	SortBy  string         `json:"sort_by,omitempty"`
	Order   string         `json:"order,omitempty" validate:"omitempty,oneof=asc desc"`
}

// SetDefaults fills in zero-valued fields with sensible defaults.
func (r *SearchRequest) SetDefaults() {
	if r.Size == 0 {
		r.Size = DefaultPageSize
	}
}

// Validate checks that the search request has at least a query or filters,
// and that pagination parameters are within bounds.
// It does NOT modify the receiver — call SetDefaults first if needed.
func (r *SearchRequest) Validate() error {
	if r == nil {
		return fmt.Errorf("search request cannot be nil")
	}

	if r.Query == "" && len(r.Filters) == 0 {
		return fmt.Errorf("either query or filters must be provided")
	}

	if r.Size < 0 {
		return fmt.Errorf("size cannot be negative")
	}
	if r.Size > MaxPageSize {
		return fmt.Errorf("size cannot exceed %d", MaxPageSize)
	}

	if r.Page < 0 {
		return fmt.Errorf("page cannot be negative")
	}

	if r.Order != "" && r.Order != "asc" && r.Order != "desc" {
		return fmt.Errorf("order must be 'asc' or 'desc'")
	}

	return nil
}

// Ptr returns a pointer to the given value.
func Ptr[T any](v T) *T { return &v }

// PtrOrNil returns a pointer to s, or nil if s is empty.
func PtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// PtrTimeOrNil returns a pointer to t, or nil if t is zero.
func PtrTimeOrNil(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}
