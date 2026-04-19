package json

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/grokify/structured-profile/schema"
	"github.com/grokify/structured-profile/store"
)

func TestNewStore(t *testing.T) {
	tmpDir := t.TempDir()

	s, err := New(Config{BaseDir: tmpDir})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	if s.baseDir != tmpDir {
		t.Errorf("Expected baseDir %s, got %s", tmpDir, s.baseDir)
	}
}

func TestSaveAndGetFullProfile(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	s, err := New(Config{BaseDir: tmpDir})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	// Create a profile
	fp := schema.NewFullProfile("John Doe")
	fp.Profile.Email = "john@example.com"
	fp.Profile.Summaries = schema.Summaries{Default: "A software engineer"}

	// Add experience
	tenure := schema.NewTenure("Acme Corp", schema.NewDate(2022, 1))
	position := schema.NewPosition("Engineer", schema.NewDate(2022, 1))
	achievement := schema.NewSTARAchievement("feat-1", "Context", "Goal", "Action", "Outcome")
	achievement.Tags = []string{"api", "backend"}
	achievement.Skills = []string{"Go", "REST"}
	position.AddAchievement(*achievement)
	tenure.AddPosition(*position)
	fp.Tenures = append(fp.Tenures, *tenure)

	// Add skill
	fp.Skills = append(fp.Skills, *schema.NewSkillWithCategory("Go", "technical", "expert"))

	// Save
	if err := s.SaveFullProfile(ctx, fp); err != nil {
		t.Fatalf("Failed to save profile: %v", err)
	}

	// Verify file exists
	path := s.ProfilePath(fp.Profile.ID)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("Expected profile file to exist")
	}

	// Get
	loaded, err := s.GetFullProfile(ctx, fp.Profile.ID)
	if err != nil {
		t.Fatalf("Failed to get profile: %v", err)
	}

	if loaded.Profile.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", loaded.Profile.Name)
	}
	if loaded.Profile.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", loaded.Profile.Email)
	}
	if len(loaded.Tenures) != 1 {
		t.Errorf("Expected 1 tenure, got %d", len(loaded.Tenures))
	}
	if len(loaded.Skills) != 1 {
		t.Errorf("Expected 1 skill, got %d", len(loaded.Skills))
	}
}

func TestGetFullProfileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	s, err := New(Config{BaseDir: tmpDir})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	_, err = s.GetFullProfile(ctx, "nonexistent")
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestDeleteProfile(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	s, err := New(Config{BaseDir: tmpDir})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	// Create and save a profile
	fp := schema.NewFullProfile("John Doe")
	if err := s.SaveFullProfile(ctx, fp); err != nil {
		t.Fatalf("Failed to save profile: %v", err)
	}

	// Delete
	if err := s.DeleteProfile(ctx, fp.Profile.ID); err != nil {
		t.Fatalf("Failed to delete profile: %v", err)
	}

	// Verify deleted
	_, err = s.GetFullProfile(ctx, fp.Profile.ID)
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("Expected ErrNotFound after delete, got %v", err)
	}
}

func TestListProfiles(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	s, err := New(Config{BaseDir: tmpDir})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	// Create and save multiple profiles
	fp1 := schema.NewFullProfile("John Doe")
	fp2 := schema.NewFullProfile("Jane Smith")

	if err := s.SaveFullProfile(ctx, fp1); err != nil {
		t.Fatalf("Failed to save profile 1: %v", err)
	}
	if err := s.SaveFullProfile(ctx, fp2); err != nil {
		t.Fatalf("Failed to save profile 2: %v", err)
	}

	// List
	profiles, err := s.ListProfiles(ctx)
	if err != nil {
		t.Fatalf("Failed to list profiles: %v", err)
	}

	if len(profiles) != 2 {
		t.Errorf("Expected 2 profiles, got %d", len(profiles))
	}
}

func TestSearchAchievementsByTags(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	s, err := New(Config{BaseDir: tmpDir})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	// Create profile with achievements
	fp := schema.NewFullProfile("John Doe")
	tenure := schema.NewTenure("Acme", schema.NewDate(2022, 1))
	position := schema.NewPosition("Engineer", schema.NewDate(2022, 1))

	a1 := schema.NewSTARAchievement("a1", "S1", "T1", "A1", "R1")
	a1.Tags = []string{"api", "backend"}
	a2 := schema.NewSTARAchievement("a2", "S2", "T2", "A2", "R2")
	a2.Tags = []string{"frontend", "react"}

	position.AddAchievement(*a1)
	position.AddAchievement(*a2)
	tenure.AddPosition(*position)
	fp.Tenures = append(fp.Tenures, *tenure)

	if err := s.SaveFullProfile(ctx, fp); err != nil {
		t.Fatalf("Failed to save profile: %v", err)
	}

	// Search by tag
	matches, err := s.SearchAchievementsByTags(ctx, fp.Profile.ID, []string{"api"})
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}

	if len(matches) != 1 {
		t.Errorf("Expected 1 match, got %d", len(matches))
	}
	if matches[0].Name != "a1" {
		t.Errorf("Expected achievement 'a1', got '%s'", matches[0].Name)
	}
}

func TestSearchAchievementsBySkills(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	s, err := New(Config{BaseDir: tmpDir})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	// Create profile with achievements
	fp := schema.NewFullProfile("John Doe")
	tenure := schema.NewTenure("Acme", schema.NewDate(2022, 1))
	position := schema.NewPosition("Engineer", schema.NewDate(2022, 1))

	a1 := schema.NewSTARAchievement("a1", "S1", "T1", "A1", "R1")
	a1.Skills = []string{"Go", "REST"}
	a2 := schema.NewSTARAchievement("a2", "S2", "T2", "A2", "R2")
	a2.Skills = []string{"Python", "Django"}

	position.AddAchievement(*a1)
	position.AddAchievement(*a2)
	tenure.AddPosition(*position)
	fp.Tenures = append(fp.Tenures, *tenure)

	if err := s.SaveFullProfile(ctx, fp); err != nil {
		t.Fatalf("Failed to save profile: %v", err)
	}

	// Search by skill
	matches, err := s.SearchAchievementsBySkills(ctx, fp.Profile.ID, []string{"Go"})
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}

	if len(matches) != 1 {
		t.Errorf("Expected 1 match, got %d", len(matches))
	}
	if matches[0].Name != "a1" {
		t.Errorf("Expected achievement 'a1', got '%s'", matches[0].Name)
	}
}

func TestStoreWithCache(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	s, err := New(Config{BaseDir: tmpDir, UseCache: true})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	// Create and save a profile
	fp := schema.NewFullProfile("John Doe")
	if err := s.SaveFullProfile(ctx, fp); err != nil {
		t.Fatalf("Failed to save profile: %v", err)
	}

	// First read (populates cache)
	_, err = s.GetFullProfile(ctx, fp.Profile.ID)
	if err != nil {
		t.Fatalf("Failed to get profile: %v", err)
	}

	// Second read (from cache)
	cached, err := s.GetFullProfile(ctx, fp.Profile.ID)
	if err != nil {
		t.Fatalf("Failed to get cached profile: %v", err)
	}

	if cached.Profile.Name != "John Doe" {
		t.Errorf("Expected cached name 'John Doe', got '%s'", cached.Profile.Name)
	}

	// Clear cache
	s.ClearCache()

	// Should still work (reads from disk)
	_, err = s.GetFullProfile(ctx, fp.Profile.ID)
	if err != nil {
		t.Fatalf("Failed to get profile after cache clear: %v", err)
	}
}

func TestProfilePath(t *testing.T) {
	s := &Store{baseDir: "/profiles"}
	path := s.ProfilePath("test-id")
	expected := filepath.Join("/profiles", "test-id.json")
	if path != expected {
		t.Errorf("Expected path %s, got %s", expected, path)
	}
}

func TestGetOpportunity(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	s, err := New(Config{BaseDir: tmpDir})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	// Create profile with opportunity
	fp := schema.NewFullProfile("John Doe")
	opp := schema.NewOpportunity("Google", "Software Engineer")
	fp.Opportunities = append(fp.Opportunities, *opp)

	if err := s.SaveFullProfile(ctx, fp); err != nil {
		t.Fatalf("Failed to save profile: %v", err)
	}

	// Get opportunity
	found, err := s.GetOpportunity(ctx, fp.Profile.ID, opp.ID)
	if err != nil {
		t.Fatalf("Failed to get opportunity: %v", err)
	}

	if found.Company != "Google" {
		t.Errorf("Expected company 'Google', got '%s'", found.Company)
	}
}

func TestGetApplication(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	s, err := New(Config{BaseDir: tmpDir})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	// Create profile with application
	fp := schema.NewFullProfile("John Doe")
	app := schema.NewApplication("opp-123")
	fp.Applications = append(fp.Applications, *app)

	if err := s.SaveFullProfile(ctx, fp); err != nil {
		t.Fatalf("Failed to save profile: %v", err)
	}

	// Get application
	found, err := s.GetApplication(ctx, fp.Profile.ID, app.ID)
	if err != nil {
		t.Fatalf("Failed to get application: %v", err)
	}

	if found.OpportunityID != "opp-123" {
		t.Errorf("Expected opportunityID 'opp-123', got '%s'", found.OpportunityID)
	}
}

func TestInvalidID(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	s, err := New(Config{BaseDir: tmpDir})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	// Test empty ID
	_, err = s.GetFullProfile(ctx, "")
	if !errors.Is(err, store.ErrInvalidID) {
		t.Errorf("Expected ErrInvalidID for empty ID, got %v", err)
	}

	// Test delete with empty ID
	err = s.DeleteProfile(ctx, "")
	if !errors.Is(err, store.ErrInvalidID) {
		t.Errorf("Expected ErrInvalidID for delete with empty ID, got %v", err)
	}
}

func TestOpportunityCRUD(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	s, err := New(Config{BaseDir: tmpDir})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	// Create profile
	fp := schema.NewFullProfile("John Doe")
	if err := s.SaveFullProfile(ctx, fp); err != nil {
		t.Fatalf("Failed to save profile: %v", err)
	}

	// Create opportunity
	opp := schema.NewOpportunity("Google", "Software Engineer")
	opp.JobDescRaw = "We are looking for a software engineer..."
	opp.SalaryMin = 150000
	opp.SalaryMax = 200000

	// Save opportunity
	if err := s.SaveOpportunity(ctx, fp.Profile.ID, opp); err != nil {
		t.Fatalf("Failed to save opportunity: %v", err)
	}

	// List opportunities
	opps, err := s.ListOpportunities(ctx, fp.Profile.ID)
	if err != nil {
		t.Fatalf("Failed to list opportunities: %v", err)
	}
	if len(opps) != 1 {
		t.Errorf("Expected 1 opportunity, got %d", len(opps))
	}

	// Get opportunity
	found, err := s.GetOpportunity(ctx, fp.Profile.ID, opp.ID)
	if err != nil {
		t.Fatalf("Failed to get opportunity: %v", err)
	}
	if found.Company != "Google" {
		t.Errorf("Expected company 'Google', got '%s'", found.Company)
	}
	if found.SalaryMin != 150000 {
		t.Errorf("Expected salary min 150000, got %d", found.SalaryMin)
	}

	// Update opportunity
	opp.Position = "Senior Software Engineer"
	if err := s.SaveOpportunity(ctx, fp.Profile.ID, opp); err != nil {
		t.Fatalf("Failed to update opportunity: %v", err)
	}

	// Verify update
	updated, err := s.GetOpportunity(ctx, fp.Profile.ID, opp.ID)
	if err != nil {
		t.Fatalf("Failed to get updated opportunity: %v", err)
	}
	if updated.Position != "Senior Software Engineer" {
		t.Errorf("Expected position 'Senior Software Engineer', got '%s'", updated.Position)
	}

	// Delete opportunity
	if err := s.DeleteOpportunity(ctx, fp.Profile.ID, opp.ID); err != nil {
		t.Fatalf("Failed to delete opportunity: %v", err)
	}

	// Verify deletion
	_, err = s.GetOpportunity(ctx, fp.Profile.ID, opp.ID)
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("Expected ErrNotFound after delete, got %v", err)
	}
}

func TestApplicationCRUD(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	s, err := New(Config{BaseDir: tmpDir})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	// Create profile with opportunity
	fp := schema.NewFullProfile("John Doe")
	opp := schema.NewOpportunity("Google", "Software Engineer")
	fp.Opportunities = append(fp.Opportunities, *opp)
	if err := s.SaveFullProfile(ctx, fp); err != nil {
		t.Fatalf("Failed to save profile: %v", err)
	}

	// Create application
	app := schema.NewApplication(opp.ID)
	app.Notes = "Applied via referral"

	// Save application
	if err := s.SaveApplication(ctx, fp.Profile.ID, app); err != nil {
		t.Fatalf("Failed to save application: %v", err)
	}

	// List applications
	apps, err := s.ListApplications(ctx, fp.Profile.ID)
	if err != nil {
		t.Fatalf("Failed to list applications: %v", err)
	}
	if len(apps) != 1 {
		t.Errorf("Expected 1 application, got %d", len(apps))
	}

	// Get application
	found, err := s.GetApplication(ctx, fp.Profile.ID, app.ID)
	if err != nil {
		t.Fatalf("Failed to get application: %v", err)
	}
	if found.OpportunityID != opp.ID {
		t.Errorf("Expected opportunity ID '%s', got '%s'", opp.ID, found.OpportunityID)
	}
	if found.Status != schema.ApplicationStatusDraft {
		t.Errorf("Expected status 'draft', got '%s'", found.Status)
	}

	// Update application (submit)
	app.Submit()
	if err := s.SaveApplication(ctx, fp.Profile.ID, app); err != nil {
		t.Fatalf("Failed to update application: %v", err)
	}

	// Verify update
	updated, err := s.GetApplication(ctx, fp.Profile.ID, app.ID)
	if err != nil {
		t.Fatalf("Failed to get updated application: %v", err)
	}
	if updated.Status != schema.ApplicationStatusSubmitted {
		t.Errorf("Expected status 'submitted', got '%s'", updated.Status)
	}

	// Delete application
	if err := s.DeleteApplication(ctx, fp.Profile.ID, app.ID); err != nil {
		t.Fatalf("Failed to delete application: %v", err)
	}

	// Verify deletion
	_, err = s.GetApplication(ctx, fp.Profile.ID, app.ID)
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("Expected ErrNotFound after delete, got %v", err)
	}
}
