package entdomain

import (
	"context"
	"fmt"
)

// Repository defines the generic repository interface for domain model CRUD and query operations.
type Repository[T DomainModel] interface {
	Create(ctx context.Context, model T) (T, error)
	GetByID(ctx context.Context, id ID) (T, error)
	Update(ctx context.Context, model T) (T, error)
	Delete(ctx context.Context, id ID) error

	CreateBatch(ctx context.Context, models []T) ([]T, error)
	UpdateBatch(ctx context.Context, models []T) ([]T, error)
	DeleteBatch(ctx context.Context, ids []ID) error

	List(ctx context.Context, req *ListRequest) ([]T, int, error)
	Search(ctx context.Context, req *SearchRequest) ([]T, int, error)
	Count(ctx context.Context, req *SearchRequest) (int, error)
	Exists(ctx context.Context, id ID) (bool, error)

	FindBy(ctx context.Context, field string, value any) ([]T, error)
	FindOneBy(ctx context.Context, field string, value any) (T, error)
}

// GenericDomainService defines the generic domain service interface.
type GenericDomainService[
	T DomainModel,
	CR CreateRequest,
	UR UpdateRequest,
	R any,
	LR any,
	QP QueryParams,
] interface {
	Create(ctx context.Context, req CR) (R, error)
	GetByID(ctx context.Context, id ID) (R, error)
	Update(ctx context.Context, id ID, req UR) (R, error)
	Delete(ctx context.Context, id ID) error
	List(ctx context.Context, page, size int, sortBy, order string) (LR, error)
	Search(ctx context.Context, params QP) (LR, error)
}

// Converters holds the conversion functions needed by BaseGenericDomainService.
// These replace the previous reflection-based approach with compile-time type safety.
type Converters[T DomainModel, R any, LR any] struct {
	// ToResponse converts a domain model to a response DTO.
	ToResponse func(T) R
	// ToListResponse converts a slice of domain models to a list response DTO.
	ToListResponse func(models []T, total, page, size int) LR
}

// BaseGenericDomainService provides a base implementation of GenericDomainService.
// Unlike the previous version, this uses explicit converter functions instead of
// reflection, providing compile-time type safety.
type BaseGenericDomainService[
	T DomainModel,
	CR CreateRequest,
	UR UpdateRequest,
	R any,
	LR any,
	QP QueryParams,
] struct {
	repo Repository[T]
	conv Converters[T, R, LR]
}

// NewBaseGenericDomainService creates a new service with explicit converters.
func NewBaseGenericDomainService[
	T DomainModel,
	CR CreateRequest,
	UR UpdateRequest,
	R any,
	LR any,
	QP QueryParams,
](
	repo Repository[T],
	conv Converters[T, R, LR],
) *BaseGenericDomainService[T, CR, UR, R, LR, QP] {
	return &BaseGenericDomainService[T, CR, UR, R, LR, QP]{
		repo: repo,
		conv: conv,
	}
}

// Create validates the request, converts it to a domain model, persists it,
// and returns the response representation.
func (s *BaseGenericDomainService[T, CR, UR, R, LR, QP]) Create(ctx context.Context, req CR) (R, error) {
	var zero R

	if err := req.Validate(); err != nil {
		return zero, fmt.Errorf("validation failed: %w", err)
	}

	model, ok := req.ToDomainModel().(T)
	if !ok {
		return zero, fmt.Errorf("type assertion failed: cannot convert CreateRequest to domain model")
	}

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return zero, fmt.Errorf("failed to create: %w", err)
	}

	return s.conv.ToResponse(created), nil
}

// GetByID retrieves a domain model by its ID and returns the response representation.
func (s *BaseGenericDomainService[T, CR, UR, R, LR, QP]) GetByID(ctx context.Context, id ID) (R, error) {
	var zero R

	if id.IsZero() {
		return zero, fmt.Errorf("invalid ID: %s", id)
	}

	model, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return zero, fmt.Errorf("failed to get by ID: %w", err)
	}

	return s.conv.ToResponse(model), nil
}

// Update validates the request, applies changes, persists them,
// and returns the response representation.
func (s *BaseGenericDomainService[T, CR, UR, R, LR, QP]) Update(ctx context.Context, id ID, req UR) (R, error) {
	var zero R

	if id.IsZero() {
		return zero, fmt.Errorf("invalid ID: %s", id)
	}

	if err := req.Validate(); err != nil {
		return zero, fmt.Errorf("validation failed: %w", err)
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return zero, fmt.Errorf("failed to get existing model: %w", err)
	}

	updated, ok := req.ApplyToDomainModel(existing).(T)
	if !ok {
		return zero, fmt.Errorf("type assertion failed: cannot convert updated model to domain model")
	}
	updated.SetID(id)

	result, err := s.repo.Update(ctx, updated)
	if err != nil {
		return zero, fmt.Errorf("failed to update: %w", err)
	}

	return s.conv.ToResponse(result), nil
}

// Delete removes a domain model by its ID.
func (s *BaseGenericDomainService[T, CR, UR, R, LR, QP]) Delete(ctx context.Context, id ID) error {
	if id.IsZero() {
		return fmt.Errorf("invalid ID: %s", id)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}

	return nil
}

// List retrieves a paginated list of domain models.
func (s *BaseGenericDomainService[T, CR, UR, R, LR, QP]) List(ctx context.Context, page, size int, sortBy, order string) (LR, error) {
	var zero LR

	if size <= 0 || size > MaxPageSize {
		size = DefaultPageSize
	}
	if page < 0 {
		page = 0
	}

	req := &ListRequest{
		Page:   page,
		Size:   size,
		SortBy: sortBy,
		Order:  order,
	}

	models, total, err := s.repo.List(ctx, req)
	if err != nil {
		return zero, fmt.Errorf("failed to list: %w", err)
	}

	return s.conv.ToListResponse(models, total, page, size), nil
}

// Search delegates to the repository using the SearchRequest extracted from
// the typed QueryParams.
func (s *BaseGenericDomainService[T, CR, UR, R, LR, QP]) Search(ctx context.Context, params QP) (LR, error) {
	var zero LR

	if err := params.Validate(); err != nil {
		return zero, fmt.Errorf("validation failed: %w", err)
	}

	req := params.ToSearchRequest()

	models, total, err := s.repo.Search(ctx, req)
	if err != nil {
		return zero, fmt.Errorf("failed to search: %w", err)
	}

	return s.conv.ToListResponse(models, total, req.Page, req.Size), nil
}
