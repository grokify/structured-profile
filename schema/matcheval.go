// Package schema provides types for structured profile data and job matching.
package schema

import (
	"encoding/json"
	"time"
)

// MatchEvaluation represents a job-to-profile match assessment
// built on top of structured-evaluation format.
type MatchEvaluation struct {
	Schema   string          `json:"$schema,omitempty"`
	Metadata MatchMetadata   `json:"metadata"`
	Job      JobSummary      `json:"job"`
	Profile  ProfileSummary  `json:"profile"`
	Match    MatchResult     `json:"match"`
	Decision MatchDecision   `json:"decision"`
	NextSteps []ActionItem   `json:"nextSteps,omitempty"`
}

// MatchMetadata contains evaluation metadata.
type MatchMetadata struct {
	EvaluationID  string    `json:"evaluationId"`
	GeneratedAt   time.Time `json:"generatedAt"`
	GeneratedBy   string    `json:"generatedBy"`
	Model         string    `json:"model,omitempty"`
	ModelProvider string    `json:"modelProvider,omitempty"`
	SchemaVersion string    `json:"schemaVersion"`
}

// JobSummary contains key job description details.
type JobSummary struct {
	DocumentPath string   `json:"documentPath"`
	Title        string   `json:"title"`
	Company      string   `json:"company"`
	Level        string   `json:"level"`
	Location     []string `json:"location,omitempty"`
	Compensation string   `json:"compensation,omitempty"`
}

// ProfileSummary contains key profile details.
type ProfileSummary struct {
	ProfileID   string `json:"profileId"`
	ProfilePath string `json:"profilePath"`
	Name        string `json:"name"`
	CurrentRole string `json:"currentRole,omitempty"`
}

// MatchResult contains the overall match assessment.
type MatchResult struct {
	OverallScore  float64           `json:"overallScore"`  // 0-100 percentage
	WeightedScore float64           `json:"weightedScore"` // 0-10 scale (structured-evaluation compatible)
	Categories    []MatchCategory   `json:"categories"`
	Findings      []MatchFinding    `json:"findings"`
}

// MatchCategory represents a scored dimension of the match.
type MatchCategory struct {
	Category      MatchCategoryType `json:"category"`
	Weight        float64           `json:"weight"`        // 0.0-1.0, should sum to 1.0
	Score         float64           `json:"score"`         // 0-10 scale
	MaxScore      float64           `json:"maxScore"`      // Default 10.0
	Status        ScoreStatus       `json:"status"`        // pass/warn/fail
	MatchType     MatchType         `json:"matchType"`     // How this was evaluated
	Justification string            `json:"justification"` // Why this score
	Evidence      []string          `json:"evidence"`      // Supporting proof
	JDRequirement string            `json:"jdRequirement,omitempty"` // What JD asked for
	ProfileMatch  string            `json:"profileMatch,omitempty"`  // What profile has
}

// MatchCategoryType defines standard match dimensions.
type MatchCategoryType string

const (
	// Core match categories
	CategoryTechnicalSkills    MatchCategoryType = "technical_skills"
	CategoryDomainExperience   MatchCategoryType = "domain_experience"
	CategoryLeadership         MatchCategoryType = "leadership"
	CategoryYearsExperience    MatchCategoryType = "years_experience"
	CategoryEducation          MatchCategoryType = "education"
	CategoryIndustryKnowledge  MatchCategoryType = "industry_knowledge"
	CategoryCommunication      MatchCategoryType = "communication"
	CategoryCulturalFit        MatchCategoryType = "cultural_fit"

	// Extended categories for specific roles
	CategoryPlatformScale      MatchCategoryType = "platform_scale"
	CategoryAPIsSDKs           MatchCategoryType = "apis_sdks"
	CategoryDeveloperRelations MatchCategoryType = "developer_relations"
	CategoryIdentitySecurity   MatchCategoryType = "identity_security"
	CategoryCompliance         MatchCategoryType = "compliance"
	CategoryThoughtLeadership  MatchCategoryType = "thought_leadership"
	CategoryTeamManagement     MatchCategoryType = "team_management"
	CategoryCrossFunctional    MatchCategoryType = "cross_functional"
)

// MatchType defines how a category was evaluated.
type MatchType string

const (
	// Exact requires the specific skill/experience
	MatchTypeExact MatchType = "exact"

	// Semantic matches equivalent/synonymous skills
	MatchTypeSemantic MatchType = "semantic"

	// Partial matches related but not identical experience
	MatchTypePartial MatchType = "partial"

	// Threshold checks if value meets minimum (e.g., years >= 10)
	MatchTypeThreshold MatchType = "threshold"

	// Boolean checks presence/absence
	MatchTypeBoolean MatchType = "boolean"

	// Weighted combines multiple sub-factors
	MatchTypeWeighted MatchType = "weighted"
)

// ScoreStatus indicates pass/warn/fail for a category.
type ScoreStatus string

const (
	StatusPass            ScoreStatus = "pass"             // >= 7.0
	StatusWarn            ScoreStatus = "warn"             // >= 5.0 && < 7.0
	StatusFail            ScoreStatus = "fail"             // < 5.0
	StatusNeedsImprovement ScoreStatus = "needs_improvement" // Alias for warn
)

// MatchFinding represents a strength or gap identified.
type MatchFinding struct {
	ID             string        `json:"id"`
	Category       MatchCategoryType `json:"category"`
	FindingType    FindingType   `json:"findingType"` // strength or gap
	Severity       Severity      `json:"severity"`
	Title          string        `json:"title"`
	Description    string        `json:"description"`
	Evidence       string        `json:"evidence,omitempty"`
	Recommendation string        `json:"recommendation,omitempty"`
	Effort         EffortLevel   `json:"effort,omitempty"`
}

// FindingType distinguishes strengths from gaps.
type FindingType string

const (
	FindingTypeStrength FindingType = "strength"
	FindingTypeGap      FindingType = "gap"
)

// Severity levels following structured-evaluation conventions.
type Severity string

const (
	SeverityCritical Severity = "critical" // Blocks - missing must-have
	SeverityHigh     Severity = "high"     // Blocks - significant gap
	SeverityMedium   Severity = "medium"   // Addressable with training
	SeverityLow      Severity = "low"      // Nice-to-have missing
	SeverityInfo     Severity = "info"     // Informational (often strengths)
)

// EffortLevel indicates effort to address a gap.
type EffortLevel string

const (
	EffortLow    EffortLevel = "low"
	EffortMedium EffortLevel = "medium"
	EffortHigh   EffortLevel = "high"
)

// MatchDecision contains the final recommendation.
type MatchDecision struct {
	Status        DecisionStatus `json:"status"`
	Passed        bool           `json:"passed"`
	Rationale     string         `json:"rationale"`
	FindingCounts FindingCounts  `json:"findingCounts"`
	OverallScore  float64        `json:"overallScore"`  // 0-100
	WeightedScore float64        `json:"weightedScore"` // 0-10
}

// DecisionStatus indicates the match recommendation.
type DecisionStatus string

const (
	DecisionPass        DecisionStatus = "pass"         // Strong match, proceed
	DecisionConditional DecisionStatus = "conditional"  // Good match with caveats
	DecisionFail        DecisionStatus = "fail"         // Has blocking gaps
	DecisionHumanReview DecisionStatus = "human_review" // Borderline, needs review
)

// FindingCounts summarizes findings by severity.
type FindingCounts struct {
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
	Info     int `json:"info"`
	Total    int `json:"total"`
}

// ActionItem represents a recommended next step.
type ActionItem struct {
	Action   string            `json:"action"`
	Category MatchCategoryType `json:"category"`
	Severity Severity          `json:"severity"`
	Effort   EffortLevel       `json:"effort"`
	Owner    string            `json:"owner,omitempty"`
}

// PassCriteria defines thresholds for pass/fail decisions.
type PassCriteria struct {
	MaxCritical int     `json:"maxCritical"` // Default 0
	MaxHigh     int     `json:"maxHigh"`     // Default 0
	MaxMedium   int     `json:"maxMedium"`   // Default -1 (unlimited)
	MinScore    float64 `json:"minScore"`    // Default 70.0 (percentage)
}

// DefaultPassCriteria returns standard pass criteria.
func DefaultPassCriteria() PassCriteria {
	return PassCriteria{
		MaxCritical: 0,
		MaxHigh:     0,
		MaxMedium:   -1, // unlimited
		MinScore:    70.0,
	}
}

// StrictPassCriteria returns stricter pass criteria for senior roles.
func StrictPassCriteria() PassCriteria {
	return PassCriteria{
		MaxCritical: 0,
		MaxHigh:     0,
		MaxMedium:   2,
		MinScore:    80.0,
	}
}

// MatchRubric defines standard scoring criteria for a category.
type MatchRubric struct {
	Category    MatchCategoryType `json:"category"`
	Description string            `json:"description"`
	MatchType   MatchType         `json:"matchType"`
	Weight      float64           `json:"weight"`
	Anchors     []ScoreAnchor     `json:"anchors"`
}

// ScoreAnchor defines what qualifies for a score range.
type ScoreAnchor struct {
	MinScore    float64  `json:"minScore"`
	MaxScore    float64  `json:"maxScore"`
	Label       string   `json:"label"`
	Description string   `json:"description"`
	Criteria    []string `json:"criteria,omitempty"`
}

// MatchRubricSet defines rubrics for a job family.
type MatchRubricSet struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	JobFamily   string        `json:"jobFamily"` // e.g., "product_management", "engineering"
	Version     string        `json:"version"`
	Rubrics     []MatchRubric `json:"rubrics"`
}

// StandardProductManagementRubrics returns rubrics for PM roles.
func StandardProductManagementRubrics() MatchRubricSet {
	return MatchRubricSet{
		ID:          "pm-match-v1",
		Name:        "Product Management Match Rubrics",
		Description: "Standard rubrics for evaluating PM candidate-job fit",
		JobFamily:   "product_management",
		Version:     "1.0",
		Rubrics: []MatchRubric{
			{
				Category:    CategoryTechnicalSkills,
				Description: "Technical depth and platform/API expertise",
				MatchType:   MatchTypeSemantic,
				Weight:      0.20,
				Anchors: []ScoreAnchor{
					{MinScore: 9, MaxScore: 10, Label: "Expert", Description: "Deep expertise, can lead technical discussions"},
					{MinScore: 7, MaxScore: 8.9, Label: "Strong", Description: "Solid skills, occasional gaps"},
					{MinScore: 5, MaxScore: 6.9, Label: "Adequate", Description: "Functional, needs growth"},
					{MinScore: 0, MaxScore: 4.9, Label: "Gap", Description: "Significant skills missing"},
				},
			},
			{
				Category:    CategoryLeadership,
				Description: "Team leadership and people management",
				MatchType:   MatchTypeThreshold,
				Weight:      0.20,
				Anchors: []ScoreAnchor{
					{MinScore: 9, MaxScore: 10, Label: "Exceeds", Description: "Larger scope than required"},
					{MinScore: 7, MaxScore: 8.9, Label: "Meets", Description: "Matches requirements"},
					{MinScore: 5, MaxScore: 6.9, Label: "Partial", Description: "Some leadership, less scope"},
					{MinScore: 0, MaxScore: 4.9, Label: "Gap", Description: "Insufficient leadership experience"},
				},
			},
			{
				Category:    CategoryDomainExperience,
				Description: "Industry and domain knowledge",
				MatchType:   MatchTypeSemantic,
				Weight:      0.15,
				Anchors: []ScoreAnchor{
					{MinScore: 9, MaxScore: 10, Label: "Direct", Description: "Same industry, same domain"},
					{MinScore: 7, MaxScore: 8.9, Label: "Adjacent", Description: "Related industry or domain"},
					{MinScore: 5, MaxScore: 6.9, Label: "Transferable", Description: "Different but applicable"},
					{MinScore: 0, MaxScore: 4.9, Label: "New", Description: "No relevant domain experience"},
				},
			},
			{
				Category:    CategoryYearsExperience,
				Description: "Total years of relevant experience",
				MatchType:   MatchTypeThreshold,
				Weight:      0.15,
				Anchors: []ScoreAnchor{
					{MinScore: 9, MaxScore: 10, Label: "Exceeds", Description: "Significantly exceeds requirement"},
					{MinScore: 7, MaxScore: 8.9, Label: "Meets", Description: "Meets or slightly exceeds"},
					{MinScore: 5, MaxScore: 6.9, Label: "Close", Description: "Slightly below requirement"},
					{MinScore: 0, MaxScore: 4.9, Label: "Below", Description: "Significantly below requirement"},
				},
			},
			{
				Category:    CategoryCrossFunctional,
				Description: "Cross-functional collaboration and influence",
				MatchType:   MatchTypeSemantic,
				Weight:      0.15,
				Anchors: []ScoreAnchor{
					{MinScore: 9, MaxScore: 10, Label: "Enterprise", Description: "C-suite and cross-org influence"},
					{MinScore: 7, MaxScore: 8.9, Label: "Organization", Description: "Multi-team collaboration"},
					{MinScore: 5, MaxScore: 6.9, Label: "Team", Description: "Team-level collaboration"},
					{MinScore: 0, MaxScore: 4.9, Label: "Limited", Description: "Minimal cross-functional work"},
				},
			},
			{
				Category:    CategoryThoughtLeadership,
				Description: "External visibility and industry presence",
				MatchType:   MatchTypeBoolean,
				Weight:      0.15,
				Anchors: []ScoreAnchor{
					{MinScore: 9, MaxScore: 10, Label: "Industry Leader", Description: "Known speaker, author, standards contributor"},
					{MinScore: 7, MaxScore: 8.9, Label: "Visible", Description: "Conference talks, published articles"},
					{MinScore: 5, MaxScore: 6.9, Label: "Emerging", Description: "Some external presence"},
					{MinScore: 0, MaxScore: 4.9, Label: "Internal", Description: "No external visibility"},
				},
			},
		},
	}
}

// ToJSON serializes the evaluation to JSON.
func (e *MatchEvaluation) ToJSON() ([]byte, error) {
	return json.MarshalIndent(e, "", "  ")
}

// FromJSON deserializes an evaluation from JSON.
func (e *MatchEvaluation) FromJSON(data []byte) error {
	return json.Unmarshal(data, e)
}

// ComputeOverallScore calculates weighted score from categories.
func (m *MatchResult) ComputeOverallScore() {
	var totalWeight, weightedSum float64
	for _, cat := range m.Categories {
		totalWeight += cat.Weight
		weightedSum += cat.Score * cat.Weight
	}
	if totalWeight > 0 {
		m.WeightedScore = weightedSum / totalWeight
		m.OverallScore = m.WeightedScore * 10 // Convert to percentage
	}
}

// ComputeDecision determines pass/fail based on criteria.
func (e *MatchEvaluation) ComputeDecision(criteria PassCriteria) {
	counts := FindingCounts{}
	for _, f := range e.Match.Findings {
		switch f.Severity {
		case SeverityCritical:
			counts.Critical++
		case SeverityHigh:
			counts.High++
		case SeverityMedium:
			counts.Medium++
		case SeverityLow:
			counts.Low++
		case SeverityInfo:
			counts.Info++
		}
		counts.Total++
	}

	e.Decision.FindingCounts = counts
	e.Decision.OverallScore = e.Match.OverallScore
	e.Decision.WeightedScore = e.Match.WeightedScore

	// Determine status
	if counts.Critical > criteria.MaxCritical {
		e.Decision.Status = DecisionFail
		e.Decision.Passed = false
		e.Decision.Rationale = "Blocked: critical findings exceed threshold"
	} else if counts.High > criteria.MaxHigh {
		e.Decision.Status = DecisionFail
		e.Decision.Passed = false
		e.Decision.Rationale = "Blocked: high severity findings exceed threshold"
	} else if e.Match.OverallScore < criteria.MinScore {
		e.Decision.Status = DecisionHumanReview
		e.Decision.Passed = false
		e.Decision.Rationale = "Score below minimum threshold"
	} else if criteria.MaxMedium >= 0 && counts.Medium > criteria.MaxMedium {
		e.Decision.Status = DecisionConditional
		e.Decision.Passed = true
		e.Decision.Rationale = "Pass with conditions: medium findings need attention"
	} else if counts.Medium > 0 || counts.Low > 0 {
		e.Decision.Status = DecisionConditional
		e.Decision.Passed = true
		e.Decision.Rationale = "Strong match with minor gaps"
	} else {
		e.Decision.Status = DecisionPass
		e.Decision.Passed = true
		e.Decision.Rationale = "Strong match, no significant gaps"
	}
}
