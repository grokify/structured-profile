// Package schema provides types for structured profile data and job matching.
package schema

import (
	"encoding/json"
	"time"
)

// DocEvaluation represents a resume/cover letter quality assessment
// against job requirements.
type DocEvaluation struct {
	Schema    string           `json:"$schema,omitempty"`
	Metadata  DocEvalMetadata  `json:"metadata"`
	Job       JobSummary       `json:"job"`
	Documents DocumentSet      `json:"documents"`
	Eval      DocEvalResult    `json:"eval"`
	Decision  DocEvalDecision  `json:"decision"`
	NextSteps []DocActionItem  `json:"nextSteps,omitempty"`
}

// DocEvalMetadata contains evaluation metadata.
type DocEvalMetadata struct {
	EvaluationID  string    `json:"evaluationId"`
	GeneratedAt   time.Time `json:"generatedAt"`
	GeneratedBy   string    `json:"generatedBy"`
	Model         string    `json:"model,omitempty"`
	ModelProvider string    `json:"modelProvider,omitempty"`
	SchemaVersion string    `json:"schemaVersion"`
	// Reference to the profile match evaluation
	MatchEvalID string `json:"matchEvalId,omitempty"`
}

// DocumentSet contains paths to the documents being evaluated.
type DocumentSet struct {
	ResumePath      string `json:"resumePath"`
	CoverLetterPath string `json:"coverLetterPath,omitempty"`
	ResumeFormat    string `json:"resumeFormat,omitempty"`    // md, pdf, docx
	CoverLetterFormat string `json:"coverLetterFormat,omitempty"`
}

// DocEvalResult contains the document quality assessment.
type DocEvalResult struct {
	OverallScore  float64              `json:"overallScore"`  // 0-100 percentage
	WeightedScore float64              `json:"weightedScore"` // 0-10 scale
	Categories    []DocEvalCategory    `json:"categories"`
	Findings      []DocEvalFinding     `json:"findings"`
}

// DocEvalCategory represents a scored dimension of document quality.
type DocEvalCategory struct {
	Category      DocCategoryType `json:"category"`
	Weight        float64         `json:"weight"`        // 0.0-1.0, should sum to 1.0
	Score         float64         `json:"score"`         // 0-10 scale
	MaxScore      float64         `json:"maxScore"`      // Default 10.0
	Status        ScoreStatus     `json:"status"`        // pass/warn/fail
	Justification string          `json:"justification"` // Why this score
	Evidence      []string        `json:"evidence"`      // Supporting examples
	Suggestions   []string        `json:"suggestions,omitempty"` // Improvement suggestions
}

// DocCategoryType defines document quality dimensions.
type DocCategoryType string

const (
	// Core document quality categories
	DocCategoryKeywordCoverage     DocCategoryType = "keyword_coverage"
	DocCategoryAchievementRelevance DocCategoryType = "achievement_relevance"
	DocCategoryGapMitigation       DocCategoryType = "gap_mitigation"
	DocCategoryQuantification      DocCategoryType = "quantification"
	DocCategoryATSCompatibility    DocCategoryType = "ats_compatibility"
	DocCategoryNarrativeCoherence  DocCategoryType = "narrative_coherence"

	// Resume-specific categories
	DocCategorySummaryImpact       DocCategoryType = "summary_impact"
	DocCategoryExperienceOrdering  DocCategoryType = "experience_ordering"
	DocCategorySkillsAlignment     DocCategoryType = "skills_alignment"
	DocCategoryFormatting          DocCategoryType = "formatting"

	// Cover letter-specific categories
	DocCategoryOpeningHook         DocCategoryType = "opening_hook"
	DocCategoryValueProposition    DocCategoryType = "value_proposition"
	DocCategoryCompanyFit          DocCategoryType = "company_fit"
	DocCategoryCallToAction        DocCategoryType = "call_to_action"
)

// DocEvalFinding represents a strength or issue in the documents.
type DocEvalFinding struct {
	ID             string          `json:"id"`
	Document       DocumentType    `json:"document"`       // resume, cover_letter, both
	Category       DocCategoryType `json:"category"`
	FindingType    FindingType     `json:"findingType"`    // strength or gap
	Severity       Severity        `json:"severity"`
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	Location       string          `json:"location,omitempty"`       // Section or line reference
	CurrentText    string          `json:"currentText,omitempty"`    // What it currently says
	SuggestedText  string          `json:"suggestedText,omitempty"`  // What it should say
	Recommendation string          `json:"recommendation,omitempty"`
	Effort         EffortLevel     `json:"effort,omitempty"`
}

// DocumentType indicates which document a finding applies to.
type DocumentType string

const (
	DocumentResume      DocumentType = "resume"
	DocumentCoverLetter DocumentType = "cover_letter"
	DocumentBoth        DocumentType = "both"
)

// DocEvalDecision contains the final quality assessment.
type DocEvalDecision struct {
	Status        DocDecisionStatus `json:"status"`
	Passed        bool              `json:"passed"`
	Rationale     string            `json:"rationale"`
	FindingCounts FindingCounts     `json:"findingCounts"`
	OverallScore  float64           `json:"overallScore"`  // 0-100
	WeightedScore float64           `json:"weightedScore"` // 0-10
	ReadyToSubmit bool              `json:"readyToSubmit"` // Can submit as-is
}

// DocDecisionStatus indicates the document quality recommendation.
type DocDecisionStatus string

const (
	DocDecisionExcellent   DocDecisionStatus = "excellent"    // Ready to submit
	DocDecisionGood        DocDecisionStatus = "good"         // Minor improvements optional
	DocDecisionNeedsWork   DocDecisionStatus = "needs_work"   // Should address issues
	DocDecisionMajorRevision DocDecisionStatus = "major_revision" // Significant rewrite needed
)

// DocActionItem represents a recommended document improvement.
type DocActionItem struct {
	Action     string          `json:"action"`
	Document   DocumentType    `json:"document"`
	Category   DocCategoryType `json:"category"`
	Severity   Severity        `json:"severity"`
	Effort     EffortLevel     `json:"effort"`
	Priority   int             `json:"priority,omitempty"` // 1 = highest
}

// DocEvalCriteria defines thresholds for document quality decisions.
type DocEvalCriteria struct {
	MaxCritical       int     `json:"maxCritical"`       // Default 0
	MaxHigh           int     `json:"maxHigh"`           // Default 0
	MaxMedium         int     `json:"maxMedium"`         // Default 2
	MinScore          float64 `json:"minScore"`          // Default 75.0
	MinKeywordCoverage float64 `json:"minKeywordCoverage"` // Default 70.0 (percentage)
}

// DefaultDocEvalCriteria returns standard document quality criteria.
func DefaultDocEvalCriteria() DocEvalCriteria {
	return DocEvalCriteria{
		MaxCritical:       0,
		MaxHigh:           0,
		MaxMedium:         2,
		MinScore:          75.0,
		MinKeywordCoverage: 70.0,
	}
}

// StrictDocEvalCriteria returns stricter criteria for senior roles.
func StrictDocEvalCriteria() DocEvalCriteria {
	return DocEvalCriteria{
		MaxCritical:       0,
		MaxHigh:           0,
		MaxMedium:         1,
		MinScore:          85.0,
		MinKeywordCoverage: 80.0,
	}
}

// DocEvalRubric defines standard scoring criteria for a document category.
type DocEvalRubric struct {
	Category    DocCategoryType `json:"category"`
	Description string          `json:"description"`
	Weight      float64         `json:"weight"`
	Document    DocumentType    `json:"document"` // Which document this applies to
	Anchors     []ScoreAnchor   `json:"anchors"`
}

// DocEvalRubricSet defines rubrics for document evaluation.
type DocEvalRubricSet struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Version     string          `json:"version"`
	Rubrics     []DocEvalRubric `json:"rubrics"`
}

// StandardDocEvalRubrics returns rubrics for resume/cover letter evaluation.
func StandardDocEvalRubrics() DocEvalRubricSet {
	return DocEvalRubricSet{
		ID:          "doc-eval-v1",
		Name:        "Document Quality Rubrics",
		Description: "Standard rubrics for evaluating resume and cover letter quality",
		Version:     "1.0",
		Rubrics: []DocEvalRubric{
			{
				Category:    DocCategoryKeywordCoverage,
				Description: "JD keywords and terminology represented in documents",
				Weight:      0.20,
				Document:    DocumentBoth,
				Anchors: []ScoreAnchor{
					{MinScore: 9, MaxScore: 10, Label: "Excellent", Description: "90%+ JD keywords present with natural usage"},
					{MinScore: 7, MaxScore: 8.9, Label: "Good", Description: "70-90% keywords, mostly natural"},
					{MinScore: 5, MaxScore: 6.9, Label: "Partial", Description: "50-70% keywords, some forced"},
					{MinScore: 0, MaxScore: 4.9, Label: "Poor", Description: "<50% keywords missing"},
				},
			},
			{
				Category:    DocCategoryAchievementRelevance,
				Description: "Top achievements align with JD requirements",
				Weight:      0.20,
				Document:    DocumentResume,
				Anchors: []ScoreAnchor{
					{MinScore: 9, MaxScore: 10, Label: "Excellent", Description: "Top 3 achievements directly address JD requirements"},
					{MinScore: 7, MaxScore: 8.9, Label: "Good", Description: "Most achievements relevant, good ordering"},
					{MinScore: 5, MaxScore: 6.9, Label: "Partial", Description: "Some relevant achievements buried or missing"},
					{MinScore: 0, MaxScore: 4.9, Label: "Poor", Description: "Achievements don't align with JD"},
				},
			},
			{
				Category:    DocCategoryGapMitigation,
				Description: "Known gaps addressed or mitigated in documents",
				Weight:      0.15,
				Document:    DocumentCoverLetter,
				Anchors: []ScoreAnchor{
					{MinScore: 9, MaxScore: 10, Label: "Excellent", Description: "Gaps proactively addressed with transferable skills"},
					{MinScore: 7, MaxScore: 8.9, Label: "Good", Description: "Major gaps addressed, minor ones acceptable"},
					{MinScore: 5, MaxScore: 6.9, Label: "Partial", Description: "Some gaps addressed, others ignored"},
					{MinScore: 0, MaxScore: 4.9, Label: "Poor", Description: "Gaps not addressed, obvious omissions"},
				},
			},
			{
				Category:    DocCategoryQuantification,
				Description: "Achievements include metrics and measurable impact",
				Weight:      0.15,
				Document:    DocumentResume,
				Anchors: []ScoreAnchor{
					{MinScore: 9, MaxScore: 10, Label: "Excellent", Description: "80%+ achievements have metrics (%, $, scale)"},
					{MinScore: 7, MaxScore: 8.9, Label: "Good", Description: "60-80% quantified, meaningful metrics"},
					{MinScore: 5, MaxScore: 6.9, Label: "Partial", Description: "40-60% quantified, some vague"},
					{MinScore: 0, MaxScore: 4.9, Label: "Poor", Description: "<40% quantified, mostly generic"},
				},
			},
			{
				Category:    DocCategoryATSCompatibility,
				Description: "Format and structure compatible with ATS systems",
				Weight:      0.10,
				Document:    DocumentResume,
				Anchors: []ScoreAnchor{
					{MinScore: 9, MaxScore: 10, Label: "Excellent", Description: "Clean format, standard sections, no tables/graphics"},
					{MinScore: 7, MaxScore: 8.9, Label: "Good", Description: "Mostly ATS-friendly, minor issues"},
					{MinScore: 5, MaxScore: 6.9, Label: "Partial", Description: "Some ATS issues (columns, headers)"},
					{MinScore: 0, MaxScore: 4.9, Label: "Poor", Description: "Heavy formatting, likely ATS problems"},
				},
			},
			{
				Category:    DocCategoryNarrativeCoherence,
				Description: "Documents tell a coherent career story",
				Weight:      0.10,
				Document:    DocumentBoth,
				Anchors: []ScoreAnchor{
					{MinScore: 9, MaxScore: 10, Label: "Excellent", Description: "Clear progression, consistent theme, compelling story"},
					{MinScore: 7, MaxScore: 8.9, Label: "Good", Description: "Logical flow, minor gaps in narrative"},
					{MinScore: 5, MaxScore: 6.9, Label: "Partial", Description: "Disjointed, unclear progression"},
					{MinScore: 0, MaxScore: 4.9, Label: "Poor", Description: "No coherent story, random listing"},
				},
			},
			{
				Category:    DocCategoryValueProposition,
				Description: "Cover letter clearly states unique value for role",
				Weight:      0.10,
				Document:    DocumentCoverLetter,
				Anchors: []ScoreAnchor{
					{MinScore: 9, MaxScore: 10, Label: "Excellent", Description: "Clear, specific value prop tied to JD requirements"},
					{MinScore: 7, MaxScore: 8.9, Label: "Good", Description: "Value stated, mostly specific"},
					{MinScore: 5, MaxScore: 6.9, Label: "Partial", Description: "Generic value statements"},
					{MinScore: 0, MaxScore: 4.9, Label: "Poor", Description: "No clear value proposition"},
				},
			},
		},
	}
}

// ToJSON serializes the evaluation to JSON.
func (e *DocEvaluation) ToJSON() ([]byte, error) {
	return json.MarshalIndent(e, "", "  ")
}

// FromJSON deserializes an evaluation from JSON.
func (e *DocEvaluation) FromJSON(data []byte) error {
	return json.Unmarshal(data, e)
}

// ComputeOverallScore calculates weighted score from categories.
func (r *DocEvalResult) ComputeOverallScore() {
	var totalWeight, weightedSum float64
	for _, cat := range r.Categories {
		totalWeight += cat.Weight
		weightedSum += cat.Score * cat.Weight
	}
	if totalWeight > 0 {
		r.WeightedScore = weightedSum / totalWeight
		r.OverallScore = r.WeightedScore * 10 // Convert to percentage
	}
}

// ComputeDecision determines document quality status based on criteria.
func (e *DocEvaluation) ComputeDecision(criteria DocEvalCriteria) {
	counts := FindingCounts{}
	for _, f := range e.Eval.Findings {
		if f.FindingType == FindingTypeGap {
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
		}
		counts.Total++
	}

	e.Decision.FindingCounts = counts
	e.Decision.OverallScore = e.Eval.OverallScore
	e.Decision.WeightedScore = e.Eval.WeightedScore

	// Determine status
	if counts.Critical > criteria.MaxCritical {
		e.Decision.Status = DocDecisionMajorRevision
		e.Decision.Passed = false
		e.Decision.ReadyToSubmit = false
		e.Decision.Rationale = "Critical document issues require major revision"
	} else if counts.High > criteria.MaxHigh {
		e.Decision.Status = DocDecisionNeedsWork
		e.Decision.Passed = false
		e.Decision.ReadyToSubmit = false
		e.Decision.Rationale = "High severity issues should be addressed before submission"
	} else if e.Eval.OverallScore < criteria.MinScore {
		e.Decision.Status = DocDecisionNeedsWork
		e.Decision.Passed = false
		e.Decision.ReadyToSubmit = false
		e.Decision.Rationale = "Document quality score below minimum threshold"
	} else if counts.Medium > criteria.MaxMedium {
		e.Decision.Status = DocDecisionGood
		e.Decision.Passed = true
		e.Decision.ReadyToSubmit = true
		e.Decision.Rationale = "Good quality with some improvements recommended"
	} else {
		e.Decision.Status = DocDecisionExcellent
		e.Decision.Passed = true
		e.Decision.ReadyToSubmit = true
		e.Decision.Rationale = "Excellent document quality, ready to submit"
	}
}
