package matcher

import (
	"testing"

	"github.com/grokify/structured-profile/schema"
)

func TestNewMatcher(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("expected non-nil matcher")
	}

	if m.Weights.RequiredSkill != 10.0 {
		t.Errorf("expected RequiredSkill weight 10.0, got %f", m.Weights.RequiredSkill)
	}
}

func TestMatchNilInputs(t *testing.T) {
	m := New()

	// Nil profile
	result := m.Match(nil, &schema.JobDescParsed{})
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	// Nil JD
	result = m.Match(&schema.FullProfile{}, nil)
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	// Both nil
	result = m.Match(nil, nil)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestMatchSkills(t *testing.T) {
	m := New()

	profile := &schema.FullProfile{
		Skills: []schema.Skill{
			{Name: "Go"},
			{Name: "Python"},
			{Name: "Kubernetes"},
			{Name: "AWS"},
		},
	}

	jd := &schema.JobDescParsed{
		RequiredSkills:  []string{"Go", "Kubernetes", "Docker"},
		PreferredSkills: []string{"Python", "GraphQL"},
	}

	result := m.Match(profile, jd)

	// Check matched required skills
	if len(result.MatchedRequiredSkills) != 2 {
		t.Errorf("expected 2 matched required skills, got %d: %v",
			len(result.MatchedRequiredSkills), result.MatchedRequiredSkills)
	}

	// Check missing required skills
	if len(result.MissingRequiredSkills) != 1 || result.MissingRequiredSkills[0] != "Docker" {
		t.Errorf("expected Docker as missing required skill, got %v", result.MissingRequiredSkills)
	}

	// Check matched preferred skills
	if len(result.MatchedPreferredSkills) != 1 {
		t.Errorf("expected 1 matched preferred skill, got %d", len(result.MatchedPreferredSkills))
	}

	// Check missing preferred skills
	if len(result.MissingPreferredSkills) != 1 || result.MissingPreferredSkills[0] != "GraphQL" {
		t.Errorf("expected GraphQL as missing preferred skill, got %v", result.MissingPreferredSkills)
	}
}

func TestRankAchievements(t *testing.T) {
	m := New()

	profile := &schema.FullProfile{
		Tenures: []schema.Tenure{
			{
				Positions: []schema.Position{
					{
						Achievements: []schema.Achievement{
							{
								BaseEntity: schema.NewBaseEntity(),
								Name:    "Built microservices with Go",
								Skills:     []string{"Go", "Kubernetes", "Docker"},
								Tags:       []string{"backend", "cloud"},
							},
							{
								BaseEntity: schema.NewBaseEntity(),
								Name:    "Developed frontend with React",
								Skills:     []string{"React", "TypeScript"},
								Tags:       []string{"frontend", "web"},
							},
							{
								BaseEntity: schema.NewBaseEntity(),
								Name:    "Led team of 5 engineers",
								Skills:     []string{"Leadership", "Agile"},
								Tags:       []string{"management"},
							},
						},
					},
				},
			},
		},
	}

	jd := &schema.JobDescParsed{
		RequiredSkills:  []string{"Go", "Kubernetes"},
		PreferredSkills: []string{"Docker"},
		Keywords:        []string{"microservices", "cloud"},
	}

	result := m.Match(profile, jd)

	if len(result.RankedAchievements) != 3 {
		t.Fatalf("expected 3 ranked achievements, got %d", len(result.RankedAchievements))
	}

	// First achievement should be the microservices one (highest match)
	if result.RankedAchievements[0].Achievement.Name != "Built microservices with Go" {
		t.Errorf("expected microservices achievement first, got %s",
			result.RankedAchievements[0].Achievement.Name)
	}

	// First should have highest score
	if result.RankedAchievements[0].Score <= result.RankedAchievements[1].Score {
		t.Errorf("expected first achievement to have higher score: %f vs %f",
			result.RankedAchievements[0].Score, result.RankedAchievements[1].Score)
	}
}

func TestTopAchievements(t *testing.T) {
	result := &MatchResult{
		RankedAchievements: []RankedAchievement{
			{Achievement: schema.Achievement{Name: "First"}, Score: 30},
			{Achievement: schema.Achievement{Name: "Second"}, Score: 20},
			{Achievement: schema.Achievement{Name: "Third"}, Score: 10},
		},
	}

	top2 := result.TopAchievements(2)
	if len(top2) != 2 {
		t.Errorf("expected 2 achievements, got %d", len(top2))
	}

	if top2[0].Name != "First" || top2[1].Name != "Second" {
		t.Errorf("expected First and Second, got %s and %s", top2[0].Name, top2[1].Name)
	}

	// Request more than available
	top10 := result.TopAchievements(10)
	if len(top10) != 3 {
		t.Errorf("expected 3 achievements (all available), got %d", len(top10))
	}
}

func TestFilterByMinScore(t *testing.T) {
	result := &MatchResult{
		RankedAchievements: []RankedAchievement{
			{Achievement: schema.Achievement{Name: "First"}, Score: 30},
			{Achievement: schema.Achievement{Name: "Second"}, Score: 20},
			{Achievement: schema.Achievement{Name: "Third"}, Score: 10},
			{Achievement: schema.Achievement{Name: "Fourth"}, Score: 5},
		},
	}

	filtered := result.FilterByMinScore(15)
	if len(filtered) != 2 {
		t.Errorf("expected 2 achievements with score >= 15, got %d", len(filtered))
	}

	filtered = result.FilterByMinScore(30)
	if len(filtered) != 1 {
		t.Errorf("expected 1 achievement with score >= 30, got %d", len(filtered))
	}

	filtered = result.FilterByMinScore(100)
	if len(filtered) != 0 {
		t.Errorf("expected 0 achievements with score >= 100, got %d", len(filtered))
	}
}

func TestOverallScore(t *testing.T) {
	m := New()

	// Profile with all required skills
	profile := &schema.FullProfile{
		Skills: []schema.Skill{
			{Name: "Go"},
			{Name: "Python"},
			{Name: "Kubernetes"},
		},
	}

	jd := &schema.JobDescParsed{
		RequiredSkills: []string{"Go", "Python", "Kubernetes"},
	}

	result := m.Match(profile, jd)

	// Should have high score since all required skills match
	if result.OverallScore < 50 {
		t.Errorf("expected high overall score, got %f", result.OverallScore)
	}
}

func TestMatchKeywordsInAchievements(t *testing.T) {
	m := New()

	profile := &schema.FullProfile{
		Tenures: []schema.Tenure{
			{
				Positions: []schema.Position{
					{
						Achievements: []schema.Achievement{
							{
								BaseEntity: schema.NewBaseEntity(),
								Name:       "Implemented API gateway",
								Situation:  "Company needed scalable microservices",
								Action:     "Built cloud-native services on AWS",
								Result:     "Improved latency by 50%",
								Tags:       []string{"API", "cloud"}, // Explicit tags for profile keyword matching
							},
						},
					},
				},
			},
		},
	}

	jd := &schema.JobDescParsed{
		Keywords: []string{"API", "microservices", "cloud", "AWS"},
	}

	result := m.Match(profile, jd)

	// Keywords are matched from explicit profile tags/skills
	if len(result.MatchedKeywords) == 0 {
		t.Error("expected keywords to be matched from profile tags")
	}

	// Achievement-level keyword matching comes from text search
	if len(result.RankedAchievements) == 0 {
		t.Fatal("expected at least one ranked achievement")
	}

	// Check that the achievement has matched keywords from its text
	if len(result.RankedAchievements[0].MatchedKeywords) == 0 {
		t.Error("expected achievement to have matched keywords from text")
	}
}

func TestExtractProfileSkillsFromAchievements(t *testing.T) {
	m := New()

	profile := &schema.FullProfile{
		// No explicit skills, but achievements have skills
		Tenures: []schema.Tenure{
			{
				Positions: []schema.Position{
					{
						Achievements: []schema.Achievement{
							{
								Skills: []string{"Go", "Docker"},
							},
						},
					},
				},
			},
		},
	}

	skills := m.extractProfileSkills(profile)

	if !skills["go"] || !skills["docker"] {
		t.Error("expected skills to be extracted from achievements")
	}
}

func TestMetricsBonus(t *testing.T) {
	m := New()

	profile := &schema.FullProfile{
		Tenures: []schema.Tenure{
			{
				Positions: []schema.Position{
					{
						Achievements: []schema.Achievement{
							{
								BaseEntity: schema.NewBaseEntity(),
								Name:    "With metrics",
								Metrics: schema.Metrics{
									Values: map[string]string{"Revenue": "50%"},
								},
							},
							{
								BaseEntity: schema.NewBaseEntity(),
								Name:    "Without metrics",
							},
						},
					},
				},
			},
		},
	}

	jd := &schema.JobDescParsed{} // Empty JD

	result := m.Match(profile, jd)

	// Achievement with metrics should score higher
	var withMetrics, withoutMetrics float64
	for _, ra := range result.RankedAchievements {
		if ra.Achievement.Name == "With metrics" {
			withMetrics = ra.Score
		} else {
			withoutMetrics = ra.Score
		}
	}

	if withMetrics <= withoutMetrics {
		t.Errorf("expected achievement with metrics to score higher: %f vs %f",
			withMetrics, withoutMetrics)
	}
}

func TestCustomWeights(t *testing.T) {
	weights := MatchWeights{
		RequiredSkill:  100.0, // Very high weight for required skills
		PreferredSkill: 1.0,
		Keyword:        1.0,
		Tag:            1.0,
	}

	m := NewWithWeights(weights)

	profile := &schema.FullProfile{
		Tenures: []schema.Tenure{
			{
				Positions: []schema.Position{
					{
						Achievements: []schema.Achievement{
							{
								BaseEntity: schema.NewBaseEntity(),
								Name:    "Has required skill",
								Skills:     []string{"Go"},
							},
							{
								BaseEntity: schema.NewBaseEntity(),
								Name:    "Has preferred skill",
								Skills:     []string{"Python"},
							},
						},
					},
				},
			},
		},
	}

	jd := &schema.JobDescParsed{
		RequiredSkills:  []string{"Go"},
		PreferredSkills: []string{"Python"},
	}

	result := m.Match(profile, jd)

	var goScore, pythonScore float64
	for _, ra := range result.RankedAchievements {
		if ra.Achievement.Name == "Has required skill" {
			goScore = ra.Score
		} else {
			pythonScore = ra.Score
		}
	}

	// Go (required) should score much higher than Python (preferred)
	if goScore <= pythonScore*10 {
		t.Errorf("expected Go score much higher than Python: %f vs %f", goScore, pythonScore)
	}
}

func TestToSet(t *testing.T) {
	items := []string{"Go", "PYTHON", "kubernetes"}
	set := toSet(items)

	if !set["go"] || !set["python"] || !set["kubernetes"] {
		t.Error("expected all items to be in set (lowercase)")
	}

	if set["Go"] || set["PYTHON"] {
		t.Error("expected set to be lowercase only")
	}
}
