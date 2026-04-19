package schema

import (
	"encoding/json"
	"testing"
	"time"
)

func TestBaseEntity(t *testing.T) {
	entity := NewBaseEntity()

	if entity.ID == "" {
		t.Error("Expected ID to be set")
	}
	if entity.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
	if entity.IsDeleted() {
		t.Error("Expected entity to not be deleted")
	}

	entity.SoftDelete()
	if !entity.IsDeleted() {
		t.Error("Expected entity to be deleted after SoftDelete")
	}

	entity.Restore()
	if entity.IsDeleted() {
		t.Error("Expected entity to not be deleted after Restore")
	}
}

func TestDate(t *testing.T) {
	tests := []struct {
		name     string
		date     Date
		expected string
	}{
		{"YYYYMM", NewDate(2024, 6), "2024-06"},
		{"YYYYMMDD", NewDateFull(2024, 6, 15), "2024-06-15"},
		{"FromDT6", NewDateFromDT6(202406), "2024-06"},
		{"FromDT8", NewDateFromDT8(20240615), "2024-06-15"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.date.String(); got != tt.expected {
				t.Errorf("Date.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDateJSON(t *testing.T) {
	date := NewDateFull(2024, 6, 15)
	data, err := json.Marshal(date)
	if err != nil {
		t.Fatalf("Failed to marshal date: %v", err)
	}

	var parsed Date
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal date: %v", err)
	}

	if !date.Equal(parsed) {
		t.Errorf("Date mismatch: got %v, want %v", parsed, date)
	}
}

func TestDateRange(t *testing.T) {
	start := NewDate(2022, 1)
	end := NewDate(2024, 6)
	r := DateRange{Start: start, End: &end}

	if r.IsCurrent() {
		t.Error("Expected DateRange to not be current")
	}

	duration := r.Duration()
	expected := 29 // 2 years and 5 months
	if duration != expected {
		t.Errorf("Duration() = %v, want %v", duration, expected)
	}
}

func TestProfile(t *testing.T) {
	profile := NewProfile("John Doe")

	if profile.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", profile.Name)
	}
	if profile.ID == "" {
		t.Error("Expected ID to be set")
	}
}

func TestSummaries(t *testing.T) {
	s := Summaries{
		Default:  "Default summary",
		ByDomain: map[string]string{"devx": "DevX summary"},
	}

	if got := s.ForDomain("devx"); got != "DevX summary" {
		t.Errorf("ForDomain('devx') = %v, want 'DevX summary'", got)
	}
	if got := s.ForDomain("unknown"); got != "Default summary" {
		t.Errorf("ForDomain('unknown') = %v, want 'Default summary'", got)
	}
}

func TestTenure(t *testing.T) {
	tenure := NewTenure("Acme Corp", NewDate(2022, 1))

	if tenure.Company != "Acme Corp" {
		t.Errorf("Expected company 'Acme Corp', got '%s'", tenure.Company)
	}
	if !tenure.IsCurrent() {
		t.Error("Expected tenure to be current (no end date)")
	}

	position := NewPosition("Software Engineer", NewDate(2022, 1))
	tenure.AddPosition(*position)

	if len(tenure.Positions) != 1 {
		t.Errorf("Expected 1 position, got %d", len(tenure.Positions))
	}
}

func TestAchievementSTAR(t *testing.T) {
	achievement := NewSTARAchievement(
		"api-redesign",
		"Legacy API had 40% churn",
		"Redesign API for better DX",
		"Implemented OpenAPI-first design with SDK generation",
		"Reduced churn by 60%, increased adoption by 300%",
	)

	star := achievement.STARString()
	if star == "" {
		t.Error("Expected STAR string to not be empty")
	}

	// Verify all parts are included
	if !contains(star, "Legacy API") {
		t.Error("STAR string missing Situation")
	}
	if !contains(star, "Redesign API") {
		t.Error("STAR string missing Task")
	}
	if !contains(star, "OpenAPI-first") {
		t.Error("STAR string missing Action")
	}
	if !contains(star, "60%") {
		t.Error("STAR string missing Result")
	}
}

func TestPositionDomainConfig(t *testing.T) {
	position := NewPosition("VP Product", NewDate(2022, 1))

	// Add achievements
	a1 := NewSTARAchievement("achievement1", "S1", "T1", "A1", "R1")
	a2 := NewSTARAchievement("achievement2", "S2", "T2", "A2", "R2")
	position.AddAchievement(*a1)
	position.AddAchievement(*a2)

	// Add domain config with specific order
	config := NewPositionDomainConfig("devx")
	config.Skills = []string{"API Design", "SDK Development"}
	config.AchievementOrder = []string{"achievement2", "achievement1"}
	position.DomainConfigs = append(position.DomainConfigs, config)

	// Test skills for domain
	skills := position.SkillsForDomain("devx")
	if len(skills) != 2 {
		t.Errorf("Expected 2 skills for devx, got %d", len(skills))
	}

	// Test achievements for domain
	achievements := position.AchievementsForDomain("devx")
	if len(achievements) != 2 {
		t.Errorf("Expected 2 achievements, got %d", len(achievements))
	}
	if achievements[0].Name != "achievement2" {
		t.Error("Expected achievement2 first (per domain config order)")
	}
}

func TestSkill(t *testing.T) {
	skill := NewSkillWithCategory("Go", "technical", "expert")

	if skill.Name != "Go" {
		t.Errorf("Expected name 'Go', got '%s'", skill.Name)
	}
	if skill.Category != "technical" {
		t.Errorf("Expected category 'technical', got '%s'", skill.Category)
	}
	if skill.Level != "expert" {
		t.Errorf("Expected level 'expert', got '%s'", skill.Level)
	}
}

func TestEducation(t *testing.T) {
	edu := NewEducation("MIT", "Bachelor of Science")
	edu.Field = "Computer Science"
	edu.Honors = "Magna Cum Laude"

	if edu.Institution != "MIT" {
		t.Errorf("Expected institution 'MIT', got '%s'", edu.Institution)
	}
	if !edu.Display {
		t.Error("Expected Display to be true by default")
	}
}

func TestCertification(t *testing.T) {
	cert := NewCertificationWithIssuer("AWS Solutions Architect", "Amazon", NewDate(2023, 6))

	if cert.Name != "AWS Solutions Architect" {
		t.Errorf("Expected name 'AWS Solutions Architect', got '%s'", cert.Name)
	}
	if cert.IsExpired() {
		t.Error("Expected certification to not be expired (no expiration date)")
	}
}

func TestOpportunity(t *testing.T) {
	opp := NewOpportunity("Google", "Senior Software Engineer")

	if opp.Company != "Google" {
		t.Errorf("Expected company 'Google', got '%s'", opp.Company)
	}
	if opp.Position != "Senior Software Engineer" {
		t.Errorf("Expected position 'Senior Software Engineer', got '%s'", opp.Position)
	}
}

func TestApplication(t *testing.T) {
	app := NewApplication("opp-123")

	if app.Status != ApplicationStatusDraft {
		t.Errorf("Expected status 'draft', got '%s'", app.Status)
	}

	app.Submit()
	if app.Status != ApplicationStatusSubmitted {
		t.Errorf("Expected status 'submitted', got '%s'", app.Status)
	}
	if app.AppliedAt == nil {
		t.Error("Expected AppliedAt to be set after Submit")
	}

	interview := NewInterview(1, InterviewTypePhone)
	app.AddInterview(*interview)
	if len(app.Interviews) != 1 {
		t.Errorf("Expected 1 interview, got %d", len(app.Interviews))
	}
	if app.Status != ApplicationStatusInterview {
		t.Errorf("Expected status 'interview', got '%s'", app.Status)
	}
}

func TestInterviewPrepSet(t *testing.T) {
	prepSet := NewInterviewPrepSet("Technical Interview Prep")

	section := NewPrepSection("Technical", "System Design")
	question := NewPrepQuestion("Design a URL shortener")
	question.Difficulty = "medium"
	question.AddAnswer(PrepAnswer{Text: "Use hash function", Correct: true})
	question.AddAnswer(PrepAnswer{Text: "Use random string", Correct: false})

	section.AddQuestion(*question)
	prepSet.AddSection(section)

	if prepSet.TotalQuestions() != 1 {
		t.Errorf("Expected 1 total question, got %d", prepSet.TotalQuestions())
	}
}

func TestFullProfile(t *testing.T) {
	fp := NewFullProfile("John Doe")

	// Add a tenure with position and achievement
	tenure := NewTenure("Acme Corp", NewDate(2022, 1))
	position := NewPosition("Engineer", NewDate(2022, 1))
	achievement := NewSTARAchievement("feat-1", "S", "T", "A", "R")
	position.AddAchievement(*achievement)
	tenure.AddPosition(*position)
	fp.Tenures = append(fp.Tenures, *tenure)

	// Test finder methods
	foundTenure := fp.FindTenure(tenure.ID)
	if foundTenure == nil {
		t.Error("Expected to find tenure by ID")
	}

	foundPosition := fp.FindPosition(position.ID)
	if foundPosition == nil {
		t.Error("Expected to find position by ID")
	}

	foundAchievement := fp.FindAchievement(achievement.ID)
	if foundAchievement == nil {
		t.Error("Expected to find achievement by ID")
	}

	allAchievements := fp.AllAchievements()
	if len(allAchievements) != 1 {
		t.Errorf("Expected 1 achievement, got %d", len(allAchievements))
	}
}

func TestFullProfileJSON(t *testing.T) {
	fp := NewFullProfile("John Doe")
	fp.Profile.Email = "john@example.com"
	fp.Profile.Summaries = Summaries{Default: "A software engineer"}

	tenure := NewTenure("Acme Corp", NewDate(2022, 1))
	position := NewPosition("Engineer", NewDate(2022, 1))
	achievement := NewSTARAchievement("feat-1", "Context", "Goal", "Action", "Outcome")
	position.AddAchievement(*achievement)
	tenure.AddPosition(*position)
	fp.Tenures = append(fp.Tenures, *tenure)

	fp.Skills = append(fp.Skills, *NewSkillWithCategory("Go", "technical", "expert"))

	// Marshal to JSON
	data, err := json.MarshalIndent(fp, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal FullProfile: %v", err)
	}

	// Unmarshal back
	var parsed FullProfile
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal FullProfile: %v", err)
	}

	if parsed.Profile.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", parsed.Profile.Name)
	}
	if len(parsed.Tenures) != 1 {
		t.Errorf("Expected 1 tenure, got %d", len(parsed.Tenures))
	}
	if len(parsed.Skills) != 1 {
		t.Errorf("Expected 1 skill, got %d", len(parsed.Skills))
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Ensure time package is used (for certification expiry test setup)
var _ = time.Now
