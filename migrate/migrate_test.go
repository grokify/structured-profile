package migrate

import (
	"testing"
)

func TestMigrateJohnWangProfile(t *testing.T) {
	fp, err := MigrateJohnWangProfile("")
	if err != nil {
		t.Fatalf("Failed to migrate profile: %v", err)
	}

	// Basic checks
	if fp.Profile.Name != "John Wang" {
		t.Errorf("Expected name 'John Wang', got '%s'", fp.Profile.Name)
	}
	if fp.Profile.Email != "johncwang@gmail.com" {
		t.Errorf("Expected email 'johncwang@gmail.com', got '%s'", fp.Profile.Email)
	}
	if fp.Profile.ID == "" {
		t.Error("Expected profile ID to be set")
	}

	// Check tenures
	if len(fp.Tenures) == 0 {
		t.Error("Expected at least one tenure")
	}
	// Check for Saviynt tenure (most recent)
	found := false
	for _, t := range fp.Tenures {
		if t.Company == "Saviynt" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find Saviynt tenure")
	}

	// Check education
	if len(fp.Education) != 4 {
		t.Errorf("Expected 4 education entries, got %d", len(fp.Education))
	}

	// Check certifications
	if len(fp.Certifications) == 0 {
		t.Error("Expected at least one certification")
	}

	// Check publications
	if len(fp.Publications) == 0 {
		t.Error("Expected at least one publication")
	}

	// Check credentials
	if len(fp.Credentials) == 0 {
		t.Error("Expected at least one credential")
	}

	// Check skills
	if len(fp.Skills) == 0 {
		t.Error("Expected at least one skill")
	}

	// Check summaries
	if fp.Profile.Summaries.Default == "" {
		t.Error("Expected default summary to be set")
	}
	if len(fp.Profile.Summaries.ByDomain) == 0 {
		t.Error("Expected domain-specific summaries")
	}
}

func TestMigrateWithCustomID(t *testing.T) {
	customID := "custom-profile-id"
	fp, err := MigrateJohnWangProfile(customID)
	if err != nil {
		t.Fatalf("Failed to migrate profile: %v", err)
	}

	if fp.Profile.ID != customID {
		t.Errorf("Expected profile ID '%s', got '%s'", customID, fp.Profile.ID)
	}
}

func TestSTARAchievementCreation(t *testing.T) {
	a := createSTARAchievement(
		"test-achievement",
		[]string{"tag1", "tag2"},
		[]string{"skill1", "skill2"},
		"Test situation",
		"Test task",
		"Test action",
		"Test result",
	)

	if a.Name != "test-achievement" {
		t.Errorf("Expected name 'test-achievement', got '%s'", a.Name)
	}
	if len(a.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(a.Tags))
	}
	if len(a.Skills) != 2 {
		t.Errorf("Expected 2 skills, got %d", len(a.Skills))
	}
	if a.Situation != "Test situation" {
		t.Errorf("Expected situation 'Test situation', got '%s'", a.Situation)
	}
}

func TestTenureMigration(t *testing.T) {
	fp, err := MigrateJohnWangProfile("")
	if err != nil {
		t.Fatalf("Failed to migrate profile: %v", err)
	}

	// Check that each tenure has at least one position
	for _, tenure := range fp.Tenures {
		if len(tenure.Positions) == 0 {
			t.Errorf("Tenure '%s' has no positions", tenure.Company)
		}

		// Check that each position has achievements
		for _, pos := range tenure.Positions {
			if len(pos.Achievements) == 0 {
				t.Errorf("Position '%s' at '%s' has no achievements", pos.Title, tenure.Company)
			}
		}
	}
}

func TestAllAchievements(t *testing.T) {
	fp, err := MigrateJohnWangProfile("")
	if err != nil {
		t.Fatalf("Failed to migrate profile: %v", err)
	}

	achievements := fp.AllAchievements()
	if len(achievements) == 0 {
		t.Error("Expected at least one achievement")
	}

	// Check that achievements have STAR components
	for _, a := range achievements {
		if a.Situation == "" {
			t.Errorf("Achievement '%s' missing Situation", a.Name)
		}
	}
}
