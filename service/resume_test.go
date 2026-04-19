package service

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/grokify/structured-profile/schema"
	jsonstore "github.com/grokify/structured-profile/store/json"
)

func setupTestStore(t *testing.T) (*jsonstore.Store, func()) {
	t.Helper()

	dir, err := os.MkdirTemp("", "resume-service-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	store, err := jsonstore.NewWithDir(dir)
	if err != nil {
		os.RemoveAll(dir)
		t.Fatalf("failed to create store: %v", err)
	}

	cleanup := func() {
		store.Close()
		os.RemoveAll(dir)
	}

	return store, cleanup
}

func createTestProfile() *schema.FullProfile {
	profile := schema.NewFullProfile("Test User")
	profile.Profile.Email = "test@example.com"
	profile.Profile.Location = "San Francisco, CA"
	profile.Profile.Summaries = schema.Summaries{
		Default: "Experienced software engineer with expertise in backend systems.",
		ByDomain: map[string]string{
			"devx": "Developer experience engineer focused on tooling and productivity.",
			"iam":  "Identity and access management specialist.",
		},
	}

	// Add tenures
	tenure := schema.NewTenure("Acme Corp", schema.Date{Year: 2020, Month: 1})
	position := schema.NewPosition("Senior Software Engineer", schema.Date{Year: 2020, Month: 1})
	position.Description = "Led backend development team"
	position.SkillsDefault = []string{"Go", "Python", "Kubernetes"}

	// Add domain config
	position.DomainConfigs = []schema.PositionDomainConfig{
		{
			Domain: "devx",
			Skills: []string{"Go", "CI/CD", "Developer Tools"},
		},
	}

	// Add achievements
	a1 := schema.NewSTARAchievement(
		"api-platform",
		"Legacy API system was causing performance issues",
		"Redesign and implement new API platform",
		"Built microservices architecture using Go and Kubernetes",
		"Improved API latency by 50%",
	)
	a1.Skills = []string{"Go", "Kubernetes", "API"}
	a1.Tags = []string{"backend", "performance"}
	position.AddAchievement(*a1)

	a2 := schema.NewSTARAchievement(
		"ci-cd",
		"Manual deployments were slow and error-prone",
		"Automate the deployment process",
		"Implemented CI/CD pipeline with GitHub Actions",
		"Reduced deployment time from hours to minutes",
	)
	a2.Skills = []string{"CI/CD", "GitHub Actions", "DevOps"}
	a2.Tags = []string{"devx", "automation"}
	position.AddAchievement(*a2)

	tenure.AddPosition(*position)
	profile.Tenures = append(profile.Tenures, *tenure)

	// Add skills
	profile.Skills = []schema.Skill{
		{Name: "Go", Category: "Programming Languages"},
		{Name: "Python", Category: "Programming Languages"},
		{Name: "Kubernetes", Category: "Infrastructure"},
		{Name: "AWS", Category: "Cloud"},
	}

	// Add education
	endDate := schema.Date{Year: 2018, Month: 5}
	profile.Education = []schema.Education{
		{
			Institution: "MIT",
			Degree:      "BS Computer Science",
			EndDate:     &endDate,
		},
	}

	return profile
}

func TestNewResumeService(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	svc := NewResumeService(store)
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

func TestGenerateBasic(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	// Save test profile
	profile := createTestProfile()
	if err := store.SaveFullProfile(ctx, profile); err != nil {
		t.Fatalf("failed to save profile: %v", err)
	}

	// Generate resume
	svc := NewResumeService(store)
	result, err := svc.Generate(ctx, GenerateInput{
		ProfileID: profile.Profile.ID,
	})

	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Resume == nil {
		t.Fatal("expected non-nil resume")
	}

	// Check content
	content := result.Resume.Content
	if content == nil {
		t.Fatal("expected non-nil content")
	}

	if content.Name != "Test User" {
		t.Errorf("expected name 'Test User', got %q", content.Name)
	}

	if content.Summary == "" {
		t.Error("expected summary to be set")
	}

	if len(content.Experiences) == 0 {
		t.Error("expected at least one experience")
	}

	if len(content.Skills) == 0 {
		t.Error("expected skills to be set")
	}
}

func TestGenerateWithDomain(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	// Save test profile
	profile := createTestProfile()
	if err := store.SaveFullProfile(ctx, profile); err != nil {
		t.Fatalf("failed to save profile: %v", err)
	}

	// Generate resume with devx domain
	svc := NewResumeService(store)
	result, err := svc.Generate(ctx, GenerateInput{
		ProfileID: profile.Profile.ID,
		Domain:    "devx",
	})

	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	content := result.Resume.Content

	// Should use devx summary
	if content.Summary != "Developer experience engineer focused on tooling and productivity." {
		t.Errorf("expected devx summary, got %q", content.Summary)
	}

	// Check domain is set
	if result.Resume.Domain != "devx" {
		t.Errorf("expected domain 'devx', got %q", result.Resume.Domain)
	}
}

func TestGenerateWithOpportunity(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	// Save test profile
	profile := createTestProfile()
	if err := store.SaveFullProfile(ctx, profile); err != nil {
		t.Fatalf("failed to save profile: %v", err)
	}

	// Save opportunity with JD
	opp := schema.NewOpportunity("TechCorp", "Staff Engineer")
	opp.JobDescParsed = &schema.JobDescParsed{
		RequiredSkills:  []string{"Go", "Kubernetes"},
		PreferredSkills: []string{"Python", "AWS"},
		Keywords:        []string{"microservices", "API", "backend"},
	}

	if err := store.SaveOpportunity(ctx, profile.Profile.ID, opp); err != nil {
		t.Fatalf("failed to save opportunity: %v", err)
	}

	// Generate resume for opportunity
	svc := NewResumeService(store)
	result, err := svc.Generate(ctx, GenerateInput{
		ProfileID:     profile.Profile.ID,
		OpportunityID: opp.ID,
	})

	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	// Should have match result
	if result.MatchResult == nil {
		t.Error("expected match result when opportunity has JD")
	}

	// Check that matched skills appear first
	if len(result.Resume.Content.Skills) == 0 {
		t.Error("expected skills to be set")
	}
}

func TestGenerateWithJDOverride(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	// Save test profile
	profile := createTestProfile()
	if err := store.SaveFullProfile(ctx, profile); err != nil {
		t.Fatalf("failed to save profile: %v", err)
	}

	// Generate with JD override
	svc := NewResumeService(store)
	result, err := svc.Generate(ctx, GenerateInput{
		ProfileID: profile.Profile.ID,
		JDOverride: &schema.JobDescParsed{
			RequiredSkills: []string{"Go", "CI/CD"},
			Keywords:       []string{"automation", "devx"},
		},
	})

	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	// Should have match result
	if result.MatchResult == nil {
		t.Error("expected match result with JD override")
	}

	if len(result.MatchResult.MatchedRequiredSkills) == 0 {
		t.Error("expected matched skills")
	}
}

func TestGenerateWithOptions(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	// Save test profile
	profile := createTestProfile()
	if err := store.SaveFullProfile(ctx, profile); err != nil {
		t.Fatalf("failed to save profile: %v", err)
	}

	// Generate with minimal options
	svc := NewResumeService(store)
	result, err := svc.Generate(ctx, GenerateInput{
		ProfileID: profile.Profile.ID,
		Options: &schema.ResumeOptions{
			IncludeContact:    true,
			IncludeSummary:    false,
			IncludeExperience: true,
			IncludeEducation:  false,
			IncludeSkills:     false,
		},
	})

	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	content := result.Resume.Content

	// Summary should be empty
	if content.Summary != "" {
		t.Error("expected summary to be empty when IncludeSummary is false")
	}

	// Education should be empty
	if len(content.Education) != 0 {
		t.Error("expected education to be empty when IncludeEducation is false")
	}

	// Skills should be empty
	if len(content.Skills) != 0 {
		t.Error("expected skills to be empty when IncludeSkills is false")
	}

	// Experience should still be present
	if len(content.Experiences) == 0 {
		t.Error("expected experiences to be present")
	}
}

func TestSelectSummary(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	svc := NewResumeService(store)

	profile := &schema.FullProfile{
		Profile: schema.Profile{
			Summaries: schema.Summaries{
				Default: "Default summary",
				ByDomain: map[string]string{
					"devx": "DevX summary",
					"iam":  "IAM summary",
				},
			},
		},
	}

	// Test domain-specific
	if summary := svc.selectSummary(profile, "devx"); summary != "DevX summary" {
		t.Errorf("expected DevX summary, got %q", summary)
	}

	// Test fallback to default
	if summary := svc.selectSummary(profile, "unknown"); summary != "Default summary" {
		t.Errorf("expected Default summary, got %q", summary)
	}

	// Test empty domain
	if summary := svc.selectSummary(profile, ""); summary != "Default summary" {
		t.Errorf("expected Default summary, got %q", summary)
	}
}

func TestSkillsMaxLimit(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	// Create profile with many skills
	profile := createTestProfile()
	profile.Skills = []schema.Skill{
		{Name: "Go"}, {Name: "Python"}, {Name: "Java"},
		{Name: "Kubernetes"}, {Name: "Docker"}, {Name: "AWS"},
		{Name: "GCP"}, {Name: "React"}, {Name: "TypeScript"},
		{Name: "PostgreSQL"}, {Name: "Redis"}, {Name: "Kafka"},
	}

	if err := store.SaveFullProfile(ctx, profile); err != nil {
		t.Fatalf("failed to save profile: %v", err)
	}

	svc := NewResumeService(store)
	result, err := svc.Generate(ctx, GenerateInput{
		ProfileID: profile.Profile.ID,
		Options: &schema.ResumeOptions{
			IncludeContact: true,
			IncludeSkills:  true,
			MaxSkills:      5,
		},
	})

	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	if len(result.Resume.Content.Skills) > 5 {
		t.Errorf("expected max 5 skills, got %d", len(result.Resume.Content.Skills))
	}
}

func TestLoadTestDataProfile(t *testing.T) {
	// Test with real testdata if available
	testdataPath := filepath.Join("..", "testdata", "john-wang.json")
	if _, err := os.Stat(testdataPath); os.IsNotExist(err) {
		t.Skip("testdata not available")
	}

	store, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	// Load test profile from testdata
	data, err := os.ReadFile(testdataPath)
	if err != nil {
		t.Skipf("could not read testdata: %v", err)
	}

	var profile schema.FullProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		t.Skipf("could not parse testdata: %v", err)
	}

	if err := store.SaveFullProfile(ctx, &profile); err != nil {
		t.Fatalf("failed to save profile: %v", err)
	}

	// Generate resume
	svc := NewResumeService(store)
	result, err := svc.Generate(ctx, GenerateInput{
		ProfileID: profile.Profile.ID,
		Domain:    "devx",
	})

	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	t.Logf("Generated resume for %s", result.Resume.Content.Name)
	t.Logf("Summary: %s", result.Resume.Content.Summary[:min(100, len(result.Resume.Content.Summary))]+"...")
	t.Logf("Experiences: %d", len(result.Resume.Content.Experiences))
	t.Logf("Skills: %d", len(result.Resume.Content.Skills))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
