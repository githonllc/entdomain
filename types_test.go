package entdomain

import (
	"testing"
)

func TestListRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		req     *ListRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: &ListRequest{
				Size:  10,
				Page: 0,
				SortBy: "name",
				Order:  "asc",
			},
			wantErr: false,
		},
		{
			name: "valid request with desc order",
			req: &ListRequest{
				Size:  20,
				Page: 10,
				SortBy: "created_at",
				Order:  "desc",
			},
			wantErr: false,
		},
		{
			name: "negative limit",
			req: &ListRequest{
				Size:  -1,
				Page: 0,
			},
			wantErr: true,
		},
		{
			name: "negative offset",
			req: &ListRequest{
				Size:  10,
				Page: -1,
			},
			wantErr: true,
		},
		{
			name: "limit too large",
			req: &ListRequest{
				Size:  1001,
				Page: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid order",
			req: &ListRequest{
				Size:  10,
				Page: 0,
				Order:  "invalid",
			},
			wantErr: true,
		},
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSearchRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		req     *SearchRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: &SearchRequest{
				Query:   "test",
				Filters: map[string]any{"status": "active"},
				Size:   10,
				Page:  0,
				SortBy:  "name",
				Order:   "asc",
			},
			wantErr: false,
		},
		{
			name: "empty query with filters",
			req: &SearchRequest{
				Query:   "",
				Filters: map[string]any{"status": "active"},
				Size:   10,
				Page:  0,
			},
			wantErr: false,
		},
		{
			name: "query without filters",
			req: &SearchRequest{
				Query:  "test",
				Size:  10,
				Page: 0,
			},
			wantErr: false,
		},
		{
			name: "empty query and no filters",
			req: &SearchRequest{
				Query:  "",
				Size:  10,
				Page: 0,
			},
			wantErr: true,
		},
		{
			name: "negative limit",
			req: &SearchRequest{
				Query:  "test",
				Size:  -1,
				Page: 0,
			},
			wantErr: true,
		},
		{
			name: "negative offset",
			req: &SearchRequest{
				Query:  "test",
				Size:  10,
				Page: -1,
			},
			wantErr: true,
		},
		{
			name: "limit too large",
			req: &SearchRequest{
				Query:  "test",
				Size:  1001,
				Page: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid order",
			req: &SearchRequest{
				Query: "test",
				Size: 10,
				Order: "invalid",
			},
			wantErr: true,
		},
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListRequestDefaults(t *testing.T) {
	req := &ListRequest{}
	req.SetDefaults()

	if req.Size != DefaultPageSize {
		t.Errorf("Default size should be %d, got %d", DefaultPageSize, req.Size)
	}

	// Validate should pass after defaults
	if err := req.Validate(); err != nil {
		t.Errorf("Validation should not fail after SetDefaults: %v", err)
	}
}

func TestSearchRequestDefaults(t *testing.T) {
	req := &SearchRequest{
		Query: "test",
	}
	req.SetDefaults()

	if req.Size != DefaultPageSize {
		t.Errorf("Default size should be %d, got %d", DefaultPageSize, req.Size)
	}

	// Validate should pass after defaults
	if err := req.Validate(); err != nil {
		t.Errorf("Validation should not fail after SetDefaults: %v", err)
	}
}
