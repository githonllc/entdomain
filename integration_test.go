package entdomain

import (
	"fmt"
	"testing"
	"time"
)

// MockDomainModel for testing
type MockPersonDomainModel struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Implement DomainModel interface
func (m *MockPersonDomainModel) GetID() ID {
	return NewIDFromInt64(m.ID)
}

func (m *MockPersonDomainModel) SetID(id ID) {
	if i, err := id.Int64(); err == nil {
		m.ID = i
	}
}

func (m *MockPersonDomainModel) Clone() DomainModel {
	if m == nil {
		return nil
	}
	clone := *m
	return &clone
}

// MockCreateRequest for testing
type MockPersonCreateRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email,omitempty"`
}

func (r *MockPersonCreateRequest) Validate() error {
	if r == nil {
		return fmt.Errorf("create request cannot be nil")
	}
	if r.FirstName == "" {
		return fmt.Errorf("first name is required")
	}
	if r.LastName == "" {
		return fmt.Errorf("last name is required")
	}
	return nil
}

func (r *MockPersonCreateRequest) ToDomainModel() DomainModel {
	if r == nil {
		return nil
	}
	return &MockPersonDomainModel{
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// MockUpdateRequest for testing
type MockPersonUpdateRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Email     *string `json:"email,omitempty"`
}

func (r *MockPersonUpdateRequest) Validate() error {
	if r == nil {
		return fmt.Errorf("update request cannot be nil")
	}
	// At least one field should be provided for update
	if r.FirstName == nil && r.LastName == nil && r.Email == nil {
		return fmt.Errorf("at least one field must be provided for update")
	}
	return nil
}

func (r *MockPersonUpdateRequest) ApplyToDomainModel(model DomainModel) DomainModel {
	if r == nil || model == nil {
		return model
	}

	person, ok := model.(*MockPersonDomainModel)
	if !ok {
		return model
	}

	if r.FirstName != nil {
		person.FirstName = *r.FirstName
	}
	if r.LastName != nil {
		person.LastName = *r.LastName
	}
	if r.Email != nil {
		person.Email = *r.Email
	}

	person.UpdatedAt = time.Now()
	return person
}

// Integration tests
func TestDomainModelInterface(t *testing.T) {
	person := &MockPersonDomainModel{
		ID:        123,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test GetID
	id := person.GetID()
	if id.String() != "123" {
		t.Errorf("Expected ID '123', got '%s'", id.String())
	}

	// Test SetID
	newID := NewIDFromInt64(456)
	person.SetID(newID)
	if person.ID != 456 {
		t.Errorf("Expected ID 456, got %d", person.ID)
	}

	// Test Clone
	cloned := person.Clone().(*MockPersonDomainModel)
	if cloned.ID != person.ID {
		t.Errorf("Cloned ID should match original")
	}
	if cloned.FirstName != person.FirstName {
		t.Errorf("Cloned FirstName should match original")
	}

	// Modify clone to ensure it's independent
	cloned.FirstName = "Jane"
	if person.FirstName == "Jane" {
		t.Error("Original should not be affected by clone modification")
	}
}

func TestCreateRequestInterface(t *testing.T) {
	// Test valid request
	createReq := &MockPersonCreateRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	if err := createReq.Validate(); err != nil {
		t.Errorf("Valid request should not fail validation: %v", err)
	}

	domainModel := createReq.ToDomainModel().(*MockPersonDomainModel)
	if domainModel.FirstName != "John" {
		t.Errorf("Expected FirstName 'John', got '%s'", domainModel.FirstName)
	}

	// Test invalid request
	invalidReq := &MockPersonCreateRequest{
		FirstName: "",
		LastName:  "Doe",
	}

	if err := invalidReq.Validate(); err == nil {
		t.Error("Invalid request should fail validation")
	}

	// Test nil request
	var nilReq *MockPersonCreateRequest
	if err := nilReq.Validate(); err == nil {
		t.Error("Nil request should fail validation")
	}
}

func TestUpdateRequestInterface(t *testing.T) {
	person := &MockPersonDomainModel{
		ID:        123,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test valid update
	newFirstName := "Jane"
	updateReq := &MockPersonUpdateRequest{
		FirstName: &newFirstName,
	}

	if err := updateReq.Validate(); err != nil {
		t.Errorf("Valid update request should not fail validation: %v", err)
	}

	originalUpdatedAt := person.UpdatedAt
	time.Sleep(1 * time.Millisecond) // Ensure time difference

	result := updateReq.ApplyToDomainModel(person)
	updatedPerson, ok := result.(*MockPersonDomainModel)
	if !ok {
		t.Fatal("ApplyToDomainModel should return *MockPersonDomainModel")
	}

	if updatedPerson.FirstName != "Jane" {
		t.Errorf("Expected FirstName 'Jane', got '%s'", updatedPerson.FirstName)
	}

	if !updatedPerson.UpdatedAt.After(originalUpdatedAt) {
		t.Error("UpdatedAt should be updated")
	}

	// Test empty update request
	emptyReq := &MockPersonUpdateRequest{}
	if err := emptyReq.Validate(); err == nil {
		t.Error("Empty update request should fail validation")
	}

	// Test nil request
	var nilReq *MockPersonUpdateRequest
	if err := nilReq.Validate(); err == nil {
		t.Error("Nil update request should fail validation")
	}
}

func TestRequestValidationWorkflow(t *testing.T) {
	// Test complete workflow: Create -> Update

	// 1. Create
	createReq := &MockPersonCreateRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	if err := createReq.Validate(); err != nil {
		t.Fatalf("Create request validation failed: %v", err)
	}

	person := createReq.ToDomainModel().(*MockPersonDomainModel)
	person.SetID(NewIDFromInt64(123)) // Simulate DB assignment

	// 2. Update
	newEmail := "john.doe.updated@example.com"
	updateReq := &MockPersonUpdateRequest{
		Email: &newEmail,
	}

	if err := updateReq.Validate(); err != nil {
		t.Fatalf("Update request validation failed: %v", err)
	}

	updateReq.ApplyToDomainModel(person)

	if person.Email != newEmail {
		t.Errorf("Expected email '%s', got '%s'", newEmail, person.Email)
	}

	// 3. Verify ID preservation
	if person.GetID().String() != "123" {
		t.Errorf("ID should be preserved through updates")
	}
}

func TestListAndSearchRequests(t *testing.T) {
	// Test ListRequest
	listReq := &ListRequest{
		Size:  10,
		Page: 0,
		SortBy: "created_at",
		Order:  "desc",
	}

	if err := listReq.Validate(); err != nil {
		t.Errorf("Valid list request should not fail: %v", err)
	}

	// Test SearchRequest
	searchReq := &SearchRequest{
		Query:   "John",
		Filters: map[string]any{"status": "active"},
		Size:   20,
		Page:  0,
		SortBy:  "name",
		Order:   "asc",
	}

	if err := searchReq.Validate(); err != nil {
		t.Errorf("Valid search request should not fail: %v", err)
	}

	// Test invalid requests
	invalidListReq := &ListRequest{
		Size:  -1,
		Page: 0,
	}

	if err := invalidListReq.Validate(); err == nil {
		t.Error("Invalid list request should fail validation")
	}

	invalidSearchReq := &SearchRequest{
		Query:  "",
		Size:  10,
		Page: 0,
	}

	if err := invalidSearchReq.Validate(); err == nil {
		t.Error("Invalid search request should fail validation")
	}
}

// Benchmark tests
func BenchmarkDomainModelClone(b *testing.B) {
	person := &MockPersonDomainModel{
		ID:        123,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = person.Clone()
	}
}

func BenchmarkIDConversion(b *testing.B) {
	id := NewIDFromInt64(123456789)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = id.Int64()
		_ = id.String()
	}
}

func BenchmarkRequestValidation(b *testing.B) {
	createReq := &MockPersonCreateRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = createReq.Validate()
	}
}
