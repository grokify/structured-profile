package service

import (
	"context"
	"strings"
	"testing"

	"github.com/grokify/structured-profile/schema"
)

func TestNewCoverLetterService(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	svc := NewCoverLetterService(store)
	if svc == nil {
		t.Fatal("expected non-nil service")
	}

	if svc.store != store {
		t.Error("expected store to be set")
	}

	if svc.matcher == nil {
		t.Error("expected matcher to be set")
	}
}

func TestCoverLetterGenerateBasic(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	// Save test profile
	profile := createTestProfile()
	if err := store.SaveFullProfile(ctx, profile); err != nil {
		t.Fatalf("failed to save profile: %v", err)
	}

	// Save opportunity
	opp := schema.NewOpportunity("TechCorp", "Staff Engineer")
	opp.HiringManager = "Jane Smith"
	opp.JobDescParsed = &schema.JobDescParsed{
		RequiredSkills:  []string{"Go", "Kubernetes"},
		PreferredSkills: []string{"Python", "AWS"},
		Keywords:        []string{"microservices", "API", "backend"},
	}

	if err := store.SaveOpportunity(ctx, profile.Profile.ID, opp); err != nil {
		t.Fatalf("failed to save opportunity: %v", err)
	}

	// Generate cover letter
	svc := NewCoverLetterService(store)
	result, err := svc.Generate(ctx, GenerateCoverLetterInput{
		ProfileID:     profile.Profile.ID,
		OpportunityID: opp.ID,
	})

	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.CoverLetter == nil {
		t.Fatal("expected non-nil cover letter")
	}

	cl := result.CoverLetter

	// Check basic fields
	if cl.TargetCompany != "TechCorp" {
		t.Errorf("expected company 'TechCorp', got %q", cl.TargetCompany)
	}

	if cl.TargetPosition != "Staff Engineer" {
		t.Errorf("expected position 'Staff Engineer', got %q", cl.TargetPosition)
	}

	if cl.HiringManager != "Jane Smith" {
		t.Errorf("expected hiring manager 'Jane Smith', got %q", cl.HiringManager)
	}

	// Check content was generated
	if cl.Opening == "" {
		t.Error("expected opening to be generated")
	}

	if cl.Body == "" {
		t.Error("expected body to be generated")
	}

	if cl.Closing == "" {
		t.Error("expected closing to be generated")
	}

	// Check match result
	if result.MatchResult == nil {
		t.Error("expected match result when opportunity has JD")
	}

	// Check STAR refs were captured
	if len(cl.STARRefs) == 0 {
		t.Error("expected STAR refs to be set")
	}
}

func TestCoverLetterGenerateWithTemplate(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	// Save test profile
	profile := createTestProfile()
	if err := store.SaveFullProfile(ctx, profile); err != nil {
		t.Fatalf("failed to save profile: %v", err)
	}

	// Save opportunity
	opp := schema.NewOpportunity("StartupInc", "Senior Engineer")
	if err := store.SaveOpportunity(ctx, profile.Profile.ID, opp); err != nil {
		t.Fatalf("failed to save opportunity: %v", err)
	}

	// Create custom template
	tmpl := &schema.CoverLetterTemplate{
		BaseEntity:      schema.NewBaseEntity(),
		Name:            "Custom",
		OpeningTemplate: "Dear {{.HiringManager}}, I am very interested in the {{.Position}} role at {{.Company}}.",
		BodyTemplate:    "My top achievement: {{.STAR1}}",
		ClosingTemplate: "Best regards, {{.Name}}",
	}

	// Generate with template
	svc := NewCoverLetterService(store)
	result, err := svc.Generate(ctx, GenerateCoverLetterInput{
		ProfileID:     profile.Profile.ID,
		OpportunityID: opp.ID,
		Template:      tmpl,
	})

	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	cl := result.CoverLetter

	// Check template was used
	if !strings.Contains(cl.Opening, "StartupInc") {
		t.Error("expected opening to contain company name from template")
	}

	if !strings.Contains(cl.Body, "My top achievement") {
		t.Error("expected body to use template format")
	}

	if !strings.Contains(cl.Closing, "Test User") {
		t.Error("expected closing to contain profile name")
	}
}

func TestCoverLetterGenerateWithSavedTemplate(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	// Save test profile
	profile := createTestProfile()
	if err := store.SaveFullProfile(ctx, profile); err != nil {
		t.Fatalf("failed to save profile: %v", err)
	}

	// Save template
	tmpl := schema.NewCoverLetterTemplate("Saved Template")
	tmpl.OpeningTemplate = "Custom opening for {{.Company}}"
	tmpl.BodyTemplate = "Custom body"
	tmpl.ClosingTemplate = "Custom closing"

	if err := store.SaveCoverLetterTemplate(ctx, profile.Profile.ID, tmpl); err != nil {
		t.Fatalf("failed to save template: %v", err)
	}

	// Save opportunity
	opp := schema.NewOpportunity("BigCorp", "Lead Engineer")
	if err := store.SaveOpportunity(ctx, profile.Profile.ID, opp); err != nil {
		t.Fatalf("failed to save opportunity: %v", err)
	}

	// Generate with saved template
	svc := NewCoverLetterService(store)
	result, err := svc.Generate(ctx, GenerateCoverLetterInput{
		ProfileID:     profile.Profile.ID,
		OpportunityID: opp.ID,
		TemplateID:    tmpl.ID,
	})

	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	if !strings.Contains(result.CoverLetter.Opening, "BigCorp") {
		t.Error("expected saved template to be used")
	}
}

func TestCoverLetterGenerateNumSTAR(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	// Save test profile with multiple achievements
	profile := createTestProfile()
	if err := store.SaveFullProfile(ctx, profile); err != nil {
		t.Fatalf("failed to save profile: %v", err)
	}

	// Save opportunity
	opp := schema.NewOpportunity("TestCo", "Engineer")
	if err := store.SaveOpportunity(ctx, profile.Profile.ID, opp); err != nil {
		t.Fatalf("failed to save opportunity: %v", err)
	}

	// Generate with 1 STAR
	svc := NewCoverLetterService(store)
	result, err := svc.Generate(ctx, GenerateCoverLetterInput{
		ProfileID:     profile.Profile.ID,
		OpportunityID: opp.ID,
		NumSTAR:       1,
	})

	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	if len(result.CoverLetter.STARRefs) > 1 {
		t.Errorf("expected max 1 STAR ref, got %d", len(result.CoverLetter.STARRefs))
	}
}

func TestCoverLetterSaveCRUD(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	// Save test profile
	profile := createTestProfile()
	if err := store.SaveFullProfile(ctx, profile); err != nil {
		t.Fatalf("failed to save profile: %v", err)
	}

	svc := NewCoverLetterService(store)

	// Create and save a cover letter
	cl := schema.NewCoverLetter(profile.Profile.ID, "opp-123")
	cl.TargetCompany = "TestCorp"
	cl.TargetPosition = "Engineer"
	cl.Opening = "Hello"
	cl.Body = "Body content"
	cl.Closing = "Thanks"

	if err := svc.SaveCoverLetter(ctx, cl); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	// Retrieve
	retrieved, err := svc.GetCoverLetter(ctx, profile.Profile.ID, cl.ID)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if retrieved.TargetCompany != "TestCorp" {
		t.Errorf("expected company 'TestCorp', got %q", retrieved.TargetCompany)
	}

	// List
	letters, err := svc.ListCoverLetters(ctx, profile.Profile.ID)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}

	if len(letters) != 1 {
		t.Errorf("expected 1 cover letter, got %d", len(letters))
	}

	// Delete
	if err := svc.DeleteCoverLetter(ctx, profile.Profile.ID, cl.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	// Verify deleted
	letters, err = svc.ListCoverLetters(ctx, profile.Profile.ID)
	if err != nil {
		t.Fatalf("list after delete failed: %v", err)
	}

	if len(letters) != 0 {
		t.Errorf("expected 0 cover letters after delete, got %d", len(letters))
	}
}

func TestTemplateData(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	// Save test profile
	profile := createTestProfile()
	if err := store.SaveFullProfile(ctx, profile); err != nil {
		t.Fatalf("failed to save profile: %v", err)
	}

	// Save opportunity with JD
	opp := schema.NewOpportunity("DataCorp", "Data Engineer")
	opp.HiringManager = "Bob Manager"
	opp.JobDescParsed = &schema.JobDescParsed{
		RequiredSkills:  []string{"Go", "Python"},
		PreferredSkills: []string{"Kubernetes"},
	}

	if err := store.SaveOpportunity(ctx, profile.Profile.ID, opp); err != nil {
		t.Fatalf("failed to save opportunity: %v", err)
	}

	// Create template that uses all data fields
	tmpl := &schema.CoverLetterTemplate{
		BaseEntity: schema.NewBaseEntity(),
		OpeningTemplate: `Company: {{.Company}}
Position: {{.Position}}
Manager: {{.HiringManager}}
Name: {{.Name}}
Years: {{.YearsExperience}}`,
		BodyTemplate: `Skills: {{.TopSkills}}
STAR1: {{.STAR1}}
{{if .STAR2}}STAR2: {{.STAR2}}{{end}}`,
		ClosingTemplate: "End",
	}

	svc := NewCoverLetterService(store)
	result, err := svc.Generate(ctx, GenerateCoverLetterInput{
		ProfileID:     profile.Profile.ID,
		OpportunityID: opp.ID,
		Template:      tmpl,
		NumSTAR:       2,
	})

	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	// Check all fields populated correctly
	if !strings.Contains(result.CoverLetter.Opening, "DataCorp") {
		t.Error("expected company in opening")
	}

	if !strings.Contains(result.CoverLetter.Opening, "Data Engineer") {
		t.Error("expected position in opening")
	}

	if !strings.Contains(result.CoverLetter.Opening, "Bob Manager") {
		t.Error("expected manager in opening")
	}

	if !strings.Contains(result.CoverLetter.Opening, "Test User") {
		t.Error("expected name in opening")
	}

	if !strings.Contains(result.CoverLetter.Body, "Go") {
		t.Error("expected matched skills in body")
	}
}

func TestDefaultCoverLetterTemplate(t *testing.T) {
	tmpl := DefaultCoverLetterTemplate()

	if tmpl == nil {
		t.Fatal("expected non-nil template")
	}

	if tmpl.Name != "Default" {
		t.Errorf("expected name 'Default', got %q", tmpl.Name)
	}

	if tmpl.OpeningTemplate == "" {
		t.Error("expected opening template")
	}

	if tmpl.BodyTemplate == "" {
		t.Error("expected body template")
	}

	if tmpl.ClosingTemplate == "" {
		t.Error("expected closing template")
	}
}

func TestCoverLetterStoreOperations(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	// Save test profile
	profile := createTestProfile()
	if err := store.SaveFullProfile(ctx, profile); err != nil {
		t.Fatalf("failed to save profile: %v", err)
	}

	// Test CoverLetter CRUD via store directly
	cl := schema.NewCoverLetter(profile.Profile.ID, "opp-1")
	cl.TargetCompany = "StoreCorp"

	// Save
	if err := store.SaveCoverLetter(ctx, profile.Profile.ID, cl); err != nil {
		t.Fatalf("save cover letter failed: %v", err)
	}

	// Get
	retrieved, err := store.GetCoverLetter(ctx, profile.Profile.ID, cl.ID)
	if err != nil {
		t.Fatalf("get cover letter failed: %v", err)
	}
	if retrieved.TargetCompany != "StoreCorp" {
		t.Errorf("expected 'StoreCorp', got %q", retrieved.TargetCompany)
	}

	// List
	list, err := store.ListCoverLetters(ctx, profile.Profile.ID)
	if err != nil {
		t.Fatalf("list cover letters failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 cover letter, got %d", len(list))
	}

	// Delete
	if err := store.DeleteCoverLetter(ctx, profile.Profile.ID, cl.ID); err != nil {
		t.Fatalf("delete cover letter failed: %v", err)
	}

	// Verify deleted
	_, err = store.GetCoverLetter(ctx, profile.Profile.ID, cl.ID)
	if err == nil {
		t.Error("expected error getting deleted cover letter")
	}
}

func TestCoverLetterTemplateStoreOperations(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	// Save test profile
	profile := createTestProfile()
	if err := store.SaveFullProfile(ctx, profile); err != nil {
		t.Fatalf("failed to save profile: %v", err)
	}

	// Test CoverLetterTemplate CRUD
	tmpl := schema.NewCoverLetterTemplate("Test Template")
	tmpl.OpeningTemplate = "Hello {{.Company}}"
	tmpl.Description = "A test template"

	// Save
	if err := store.SaveCoverLetterTemplate(ctx, profile.Profile.ID, tmpl); err != nil {
		t.Fatalf("save template failed: %v", err)
	}

	// Get
	retrieved, err := store.GetCoverLetterTemplate(ctx, profile.Profile.ID, tmpl.ID)
	if err != nil {
		t.Fatalf("get template failed: %v", err)
	}
	if retrieved.Name != "Test Template" {
		t.Errorf("expected 'Test Template', got %q", retrieved.Name)
	}
	if retrieved.OpeningTemplate != "Hello {{.Company}}" {
		t.Errorf("expected template content, got %q", retrieved.OpeningTemplate)
	}

	// List
	list, err := store.ListCoverLetterTemplates(ctx, profile.Profile.ID)
	if err != nil {
		t.Fatalf("list templates failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 template, got %d", len(list))
	}

	// Update
	retrieved.Description = "Updated description"
	if err := store.SaveCoverLetterTemplate(ctx, profile.Profile.ID, retrieved); err != nil {
		t.Fatalf("update template failed: %v", err)
	}

	updated, err := store.GetCoverLetterTemplate(ctx, profile.Profile.ID, tmpl.ID)
	if err != nil {
		t.Fatalf("get updated template failed: %v", err)
	}
	if updated.Description != "Updated description" {
		t.Errorf("expected 'Updated description', got %q", updated.Description)
	}

	// Delete
	if err := store.DeleteCoverLetterTemplate(ctx, profile.Profile.ID, tmpl.ID); err != nil {
		t.Fatalf("delete template failed: %v", err)
	}

	// Verify deleted
	_, err = store.GetCoverLetterTemplate(ctx, profile.Profile.ID, tmpl.ID)
	if err == nil {
		t.Error("expected error getting deleted template")
	}
}

func TestYearsExperienceCalculation(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	svc := NewCoverLetterService(store)

	// Profile with tenure starting in 2020
	profile := createTestProfile()
	years := svc.calculateYearsExperience(profile)

	// Should be approximately 4 years (2024 - 2020)
	if years != 4 {
		t.Errorf("expected 4 years experience, got %d", years)
	}

	// Empty profile
	emptyProfile := &schema.FullProfile{}
	years = svc.calculateYearsExperience(emptyProfile)
	if years != 0 {
		t.Errorf("expected 0 years for empty profile, got %d", years)
	}
}

