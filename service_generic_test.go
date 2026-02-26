package entdomain

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

// --- Mock types ---

// mockModel implements DomainModel (Entity + Cloner)
type mockModel struct {
	ID   ID
	Name string
}

func (m *mockModel) GetID() ID   { return m.ID }
func (m *mockModel) SetID(id ID) { m.ID = id }
func (m *mockModel) Clone() DomainModel {
	return &mockModel{ID: m.ID, Name: m.Name}
}

// ToResponse converts mockModel to *mockResponseDTO (used by domainModelToResponse via reflection)
func (m *mockModel) ToResponse() *mockResponseDTO {
	return &mockResponseDTO{ID: m.ID, Name: m.Name}
}

// mockResponseDTO is the response DTO returned by mockModel.ToResponse()
type mockResponseDTO struct {
	ID   ID
	Name string
}

// ToDomainModel converts response back to domain model (for testing)
func (r *mockResponseDTO) ToDomainModel() DomainModel {
	return &mockModel{ID: r.ID, Name: r.Name}
}

// mockCreateReq implements CreateRequest (Validate + ToDomainModel)
type mockCreateReq struct {
	Name       string
	shouldFail bool
}

func (r *mockCreateReq) Validate() error {
	if r.shouldFail {
		return fmt.Errorf("create validation failed")
	}
	return nil
}

func (r *mockCreateReq) ToDomainModel() DomainModel {
	return &mockModel{Name: r.Name}
}

// mockUpdateReq implements UpdateRequest (Validate + ToDomainModel + ApplyToDomainModel)
type mockUpdateReq struct {
	Name       string
	shouldFail bool
}

func (r *mockUpdateReq) Validate() error {
	if r.shouldFail {
		return fmt.Errorf("update validation failed")
	}
	return nil
}

func (r *mockUpdateReq) ToDomainModel() DomainModel {
	return &mockModel{Name: r.Name}
}

func (r *mockUpdateReq) ApplyToDomainModel(domain DomainModel) DomainModel {
	m := domain.(*mockModel)
	return &mockModel{ID: m.ID, Name: r.Name}
}

// mockQueryParams implements QueryParams (Validate + ToSearchRequest)
type mockQueryParams struct {
	query      string
	shouldFail bool
}

func (p *mockQueryParams) Validate() error {
	if p.shouldFail {
		return fmt.Errorf("query params validation failed")
	}
	return nil
}

func (p *mockQueryParams) ToSearchRequest() *SearchRequest {
	return &SearchRequest{
		Query:   p.query,
		Size:    DefaultPageSize,
		Page:    0,
		Filters: make(map[string]any),
	}
}

// mockListResponseDTO is the list response DTO for testing
type mockListResponseDTO struct {
	Data   []*mockResponseDTO
	Total  int
	Limit  int
	Offset int
}

func (r *mockListResponseDTO) GetData() []*mockResponseDTO { return r.Data }
func (r *mockListResponseDTO) GetTotal() int               { return r.Total }
func (r *mockListResponseDTO) GetLimit() int               { return r.Limit }
func (r *mockListResponseDTO) GetOffset() int              { return r.Offset }

// mockRepo implements Repository[*mockModel] with configurable behavior
type mockRepo struct {
	createFn      func(ctx context.Context, model *mockModel) (*mockModel, error)
	getByIDFn     func(ctx context.Context, id ID) (*mockModel, error)
	updateFn      func(ctx context.Context, model *mockModel) (*mockModel, error)
	deleteFn      func(ctx context.Context, id ID) error
	listFn        func(ctx context.Context, req *ListRequest) ([]*mockModel, int, error)
	searchFn      func(ctx context.Context, req *SearchRequest) ([]*mockModel, int, error)
	createBatchFn func(ctx context.Context, models []*mockModel) ([]*mockModel, error)
	updateBatchFn func(ctx context.Context, models []*mockModel) ([]*mockModel, error)
	deleteBatchFn func(ctx context.Context, ids []ID) error
	countFn       func(ctx context.Context, req *SearchRequest) (int, error)
	existsFn      func(ctx context.Context, id ID) (bool, error)
	findByFn      func(ctx context.Context, field string, value any) ([]*mockModel, error)
	findOneByFn   func(ctx context.Context, field string, value any) (*mockModel, error)
}

func (r *mockRepo) Create(ctx context.Context, model *mockModel) (*mockModel, error) {
	if r.createFn != nil {
		return r.createFn(ctx, model)
	}
	return model, nil
}

func (r *mockRepo) GetByID(ctx context.Context, id ID) (*mockModel, error) {
	if r.getByIDFn != nil {
		return r.getByIDFn(ctx, id)
	}
	return &mockModel{ID: id, Name: "found"}, nil
}

func (r *mockRepo) Update(ctx context.Context, model *mockModel) (*mockModel, error) {
	if r.updateFn != nil {
		return r.updateFn(ctx, model)
	}
	return model, nil
}

func (r *mockRepo) Delete(ctx context.Context, id ID) error {
	if r.deleteFn != nil {
		return r.deleteFn(ctx, id)
	}
	return nil
}

func (r *mockRepo) CreateBatch(ctx context.Context, models []*mockModel) ([]*mockModel, error) {
	if r.createBatchFn != nil {
		return r.createBatchFn(ctx, models)
	}
	return models, nil
}

func (r *mockRepo) UpdateBatch(ctx context.Context, models []*mockModel) ([]*mockModel, error) {
	if r.updateBatchFn != nil {
		return r.updateBatchFn(ctx, models)
	}
	return models, nil
}

func (r *mockRepo) DeleteBatch(ctx context.Context, ids []ID) error {
	if r.deleteBatchFn != nil {
		return r.deleteBatchFn(ctx, ids)
	}
	return nil
}

func (r *mockRepo) List(ctx context.Context, req *ListRequest) ([]*mockModel, int, error) {
	if r.listFn != nil {
		return r.listFn(ctx, req)
	}
	return []*mockModel{}, 0, nil
}

func (r *mockRepo) Search(ctx context.Context, req *SearchRequest) ([]*mockModel, int, error) {
	if r.searchFn != nil {
		return r.searchFn(ctx, req)
	}
	return []*mockModel{}, 0, nil
}

func (r *mockRepo) Count(ctx context.Context, req *SearchRequest) (int, error) {
	if r.countFn != nil {
		return r.countFn(ctx, req)
	}
	return 0, nil
}

func (r *mockRepo) Exists(ctx context.Context, id ID) (bool, error) {
	if r.existsFn != nil {
		return r.existsFn(ctx, id)
	}
	return true, nil
}

func (r *mockRepo) FindBy(ctx context.Context, field string, value any) ([]*mockModel, error) {
	if r.findByFn != nil {
		return r.findByFn(ctx, field, value)
	}
	return []*mockModel{}, nil
}

func (r *mockRepo) FindOneBy(ctx context.Context, field string, value any) (*mockModel, error) {
	if r.findOneByFn != nil {
		return r.findOneByFn(ctx, field, value)
	}
	return &mockModel{}, nil
}

// --- Type aliases for readability ---

type testService = BaseGenericDomainService[
	*mockModel,
	*mockCreateReq,
	*mockUpdateReq,
	*mockResponseDTO,
	*mockListResponseDTO,
	*mockQueryParams,
]

func newTestService(repo *mockRepo) *testService {
	return NewBaseGenericDomainService[
		*mockModel,
		*mockCreateReq,
		*mockUpdateReq,
		*mockResponseDTO,
		*mockListResponseDTO,
		*mockQueryParams,
	](repo, Converters[*mockModel, *mockResponseDTO, *mockListResponseDTO]{
		ToResponse: func(m *mockModel) *mockResponseDTO {
			return m.ToResponse()
		},
		ToListResponse: func(models []*mockModel, total, page, size int) *mockListResponseDTO {
			data := make([]*mockResponseDTO, len(models))
			for i, m := range models {
				data[i] = m.ToResponse()
			}
			return &mockListResponseDTO{Data: data, Total: total, Limit: size, Offset: page * size}
		},
	})
}

// --- Tests ---

func TestBaseGenericDomainService_Create(t *testing.T) {
	tests := []struct {
		name    string
		req     *mockCreateReq
		repo    *mockRepo
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request creates successfully",
			req:  &mockCreateReq{Name: "new-entity"},
			repo: &mockRepo{
				createFn: func(ctx context.Context, model *mockModel) (*mockModel, error) {
					model.ID = Int64ID(1)
					return model, nil
				},
			},
			wantErr: false,
		},
		{
			name:    "validation failure returns error",
			req:     &mockCreateReq{Name: "bad", shouldFail: true},
			repo:    &mockRepo{},
			wantErr: true,
			errMsg:  "validation failed",
		},
		{
			name: "repo error returns error",
			req:  &mockCreateReq{Name: "entity"},
			repo: &mockRepo{
				createFn: func(ctx context.Context, model *mockModel) (*mockModel, error) {
					return nil, fmt.Errorf("db connection refused")
				},
			},
			wantErr: true,
			errMsg:  "failed to create",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.repo)
			result, err := svc.Create(context.Background(), tt.req)

			if tt.wantErr {
				if err == nil {
					t.Fatal("Create() expected error, got nil")
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Create() error = %q, want to contain %q", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Fatalf("Create() unexpected error: %v", err)
			}

			if result == nil {
				t.Fatal("Create() returned nil result")
			}

			if result.Name != tt.req.Name {
				t.Errorf("Create() result.Name = %q, want %q", result.Name, tt.req.Name)
			}
		})
	}
}

func TestBaseGenericDomainService_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		id      ID
		repo    *mockRepo
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid ID returns model",
			id:   Int64ID(1),
			repo: &mockRepo{
				getByIDFn: func(ctx context.Context, id ID) (*mockModel, error) {
					return &mockModel{ID: id, Name: "found"}, nil
				},
			},
			wantErr: false,
		},
		{
			name:    "zero ID returns error",
			id:      Int64ID(0),
			repo:    &mockRepo{},
			wantErr: true,
			errMsg:  "invalid ID",
		},
		{
			name: "not found returns error",
			id:   Int64ID(999),
			repo: &mockRepo{
				getByIDFn: func(ctx context.Context, id ID) (*mockModel, error) {
					return nil, fmt.Errorf("not found")
				},
			},
			wantErr: true,
			errMsg:  "failed to get by ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.repo)
			result, err := svc.GetByID(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Fatal("GetByID() expected error, got nil")
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("GetByID() error = %q, want to contain %q", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetByID() unexpected error: %v", err)
			}

			if result == nil {
				t.Fatal("GetByID() returned nil result")
			}
		})
	}
}

func TestBaseGenericDomainService_Update(t *testing.T) {
	tests := []struct {
		name    string
		id      ID
		req     *mockUpdateReq
		repo    *mockRepo
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid update succeeds",
			id:   Int64ID(1),
			req:  &mockUpdateReq{Name: "updated"},
			repo: &mockRepo{
				getByIDFn: func(ctx context.Context, id ID) (*mockModel, error) {
					return &mockModel{ID: id, Name: "original"}, nil
				},
				updateFn: func(ctx context.Context, model *mockModel) (*mockModel, error) {
					return model, nil
				},
			},
			wantErr: false,
		},
		{
			name:    "zero ID returns error",
			id:      Int64ID(0),
			req:     &mockUpdateReq{Name: "updated"},
			repo:    &mockRepo{},
			wantErr: true,
			errMsg:  "invalid ID",
		},
		{
			name:    "validation failure returns error",
			id:      Int64ID(1),
			req:     &mockUpdateReq{Name: "bad", shouldFail: true},
			repo:    &mockRepo{},
			wantErr: true,
			errMsg:  "validation failed",
		},
		{
			name: "not found returns error",
			id:   Int64ID(999),
			req:  &mockUpdateReq{Name: "updated"},
			repo: &mockRepo{
				getByIDFn: func(ctx context.Context, id ID) (*mockModel, error) {
					return nil, fmt.Errorf("entity not found")
				},
			},
			wantErr: true,
			errMsg:  "failed to get existing model",
		},
		{
			name: "repo update error returns error",
			id:   Int64ID(1),
			req:  &mockUpdateReq{Name: "updated"},
			repo: &mockRepo{
				getByIDFn: func(ctx context.Context, id ID) (*mockModel, error) {
					return &mockModel{ID: id, Name: "original"}, nil
				},
				updateFn: func(ctx context.Context, model *mockModel) (*mockModel, error) {
					return nil, fmt.Errorf("db write error")
				},
			},
			wantErr: true,
			errMsg:  "failed to update",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.repo)
			result, err := svc.Update(context.Background(), tt.id, tt.req)

			if tt.wantErr {
				if err == nil {
					t.Fatal("Update() expected error, got nil")
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Update() error = %q, want to contain %q", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Fatalf("Update() unexpected error: %v", err)
			}

			if result == nil {
				t.Fatal("Update() returned nil result")
			}

			if result.Name != tt.req.Name {
				t.Errorf("Update() result.Name = %q, want %q", result.Name, tt.req.Name)
			}
		})
	}
}

func TestBaseGenericDomainService_Delete(t *testing.T) {
	tests := []struct {
		name    string
		id      ID
		repo    *mockRepo
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid delete succeeds",
			id:   Int64ID(1),
			repo: &mockRepo{
				deleteFn: func(ctx context.Context, id ID) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name:    "zero ID returns error",
			id:      Int64ID(0),
			repo:    &mockRepo{},
			wantErr: true,
			errMsg:  "invalid ID",
		},
		{
			name: "repo error returns error",
			id:   Int64ID(1),
			repo: &mockRepo{
				deleteFn: func(ctx context.Context, id ID) error {
					return fmt.Errorf("foreign key constraint")
				},
			},
			wantErr: true,
			errMsg:  "failed to delete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.repo)
			err := svc.Delete(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Fatal("Delete() expected error, got nil")
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Delete() error = %q, want to contain %q", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Fatalf("Delete() unexpected error: %v", err)
			}
		})
	}
}

func TestBaseGenericDomainService_List(t *testing.T) {
	t.Run("valid params calls repo with correct request", func(t *testing.T) {
		var capturedReq *ListRequest
		repo := &mockRepo{
			listFn: func(ctx context.Context, req *ListRequest) ([]*mockModel, int, error) {
				capturedReq = req
				return []*mockModel{
					{ID: Int64ID(1), Name: "first"},
					{ID: Int64ID(2), Name: "second"},
				}, 2, nil
			},
		}
		svc := newTestService(repo)

		result, err := svc.List(context.Background(), 0, 10, "name", "asc")
		if err != nil {
			t.Fatalf("List() unexpected error: %v", err)
		}

		// Verify the repo was called with correct parameters
		if capturedReq == nil {
			t.Fatal("List() did not call repo.List")
		}
		if capturedReq.Page != 0 {
			t.Errorf("List() repo received Page = %d, want 0", capturedReq.Page)
		}
		if capturedReq.Size != 10 {
			t.Errorf("List() repo received Size = %d, want 10", capturedReq.Size)
		}
		if capturedReq.SortBy != "name" {
			t.Errorf("List() repo received SortBy = %q, want %q", capturedReq.SortBy, "name")
		}
		if capturedReq.Order != "asc" {
			t.Errorf("List() repo received Order = %q, want %q", capturedReq.Order, "asc")
		}

		// Verify response was properly converted
		if result == nil {
			t.Fatal("List() returned nil result")
		}
		if result.Total != 2 {
			t.Errorf("List() total = %d, want 2", result.Total)
		}
	})

	t.Run("default size when size is zero", func(t *testing.T) {
		var capturedReq *ListRequest
		repo := &mockRepo{
			listFn: func(ctx context.Context, req *ListRequest) ([]*mockModel, int, error) {
				capturedReq = req
				return []*mockModel{}, 0, nil
			},
		}
		svc := newTestService(repo)

		_, _ = svc.List(context.Background(), 0, 0, "", "")
		if capturedReq == nil {
			t.Fatal("List() did not call repo.List")
		}
		if capturedReq.Size != DefaultPageSize {
			t.Errorf("List() default size = %d, want %d", capturedReq.Size, DefaultPageSize)
		}
	})

	t.Run("default size when size exceeds MaxPageSize", func(t *testing.T) {
		var capturedReq *ListRequest
		repo := &mockRepo{
			listFn: func(ctx context.Context, req *ListRequest) ([]*mockModel, int, error) {
				capturedReq = req
				return []*mockModel{}, 0, nil
			},
		}
		svc := newTestService(repo)

		_, _ = svc.List(context.Background(), 0, MaxPageSize+100, "", "")
		if capturedReq == nil {
			t.Fatal("List() did not call repo.List")
		}
		if capturedReq.Size != DefaultPageSize {
			t.Errorf("List() size for oversized input = %d, want %d", capturedReq.Size, DefaultPageSize)
		}
	})

	t.Run("negative page is clamped to 0", func(t *testing.T) {
		var capturedReq *ListRequest
		repo := &mockRepo{
			listFn: func(ctx context.Context, req *ListRequest) ([]*mockModel, int, error) {
				capturedReq = req
				return []*mockModel{}, 0, nil
			},
		}
		svc := newTestService(repo)

		_, _ = svc.List(context.Background(), -5, 10, "", "")
		if capturedReq == nil {
			t.Fatal("List() did not call repo.List")
		}
		if capturedReq.Page != 0 {
			t.Errorf("List() page for negative input = %d, want 0", capturedReq.Page)
		}
	})

	t.Run("repo error returns error", func(t *testing.T) {
		repo := &mockRepo{
			listFn: func(ctx context.Context, req *ListRequest) ([]*mockModel, int, error) {
				return nil, 0, fmt.Errorf("db timeout")
			},
		}
		svc := newTestService(repo)
		_, err := svc.List(context.Background(), 0, 10, "", "")

		if err == nil {
			t.Fatal("List() expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to list") {
			t.Errorf("List() error = %q, want to contain %q", err.Error(), "failed to list")
		}
	})
}

func TestBaseGenericDomainService_Search(t *testing.T) {
	t.Run("valid search calls repo", func(t *testing.T) {
		repoCalled := false
		repo := &mockRepo{
			searchFn: func(ctx context.Context, req *SearchRequest) ([]*mockModel, int, error) {
				repoCalled = true
				return []*mockModel{
					{ID: Int64ID(1), Name: "result"},
				}, 1, nil
			},
		}
		svc := newTestService(repo)

		result, err := svc.Search(context.Background(), &mockQueryParams{shouldFail: false})
		if err != nil {
			t.Fatalf("Search() unexpected error: %v", err)
		}
		if !repoCalled {
			t.Fatal("Search() did not call repo.Search")
		}
		if result == nil {
			t.Fatal("Search() returned nil result")
		}
		if result.Total != 1 {
			t.Errorf("Search() total = %d, want 1", result.Total)
		}
	})

	t.Run("validation failure returns error", func(t *testing.T) {
		svc := newTestService(&mockRepo{})
		_, err := svc.Search(context.Background(), &mockQueryParams{shouldFail: true})

		if err == nil {
			t.Fatal("Search() expected error, got nil")
		}
		if !strings.Contains(err.Error(), "validation failed") {
			t.Errorf("Search() error = %q, want to contain %q", err.Error(), "validation failed")
		}
	})

	t.Run("repo error returns error", func(t *testing.T) {
		repo := &mockRepo{
			searchFn: func(ctx context.Context, req *SearchRequest) ([]*mockModel, int, error) {
				return nil, 0, fmt.Errorf("search index unavailable")
			},
		}
		svc := newTestService(repo)
		_, err := svc.Search(context.Background(), &mockQueryParams{shouldFail: false})

		if err == nil {
			t.Fatal("Search() expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to search") {
			t.Errorf("Search() error = %q, want to contain %q", err.Error(), "failed to search")
		}
	})
}
