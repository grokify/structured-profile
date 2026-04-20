# Document Evaluation Schema

The document evaluation schema (`schema/doceval.go`) defines the structure for resume/cover letter quality assessment.

## DocEvaluation

Top-level structure:

```go
type DocEvaluation struct {
    Schema    string           `json:"$schema,omitempty"`
    Metadata  DocEvalMetadata  `json:"metadata"`
    Job       JobSummary       `json:"job"`
    Documents DocumentSet      `json:"documents"`
    Eval      DocEvalResult    `json:"eval"`
    Decision  DocEvalDecision  `json:"decision"`
    NextSteps []DocActionItem  `json:"nextSteps,omitempty"`
}
```

## DocEvalMetadata

```go
type DocEvalMetadata struct {
    EvaluationID  string    `json:"evaluationId"`
    GeneratedAt   time.Time `json:"generatedAt"`
    GeneratedBy   string    `json:"generatedBy"`
    Model         string    `json:"model,omitempty"`
    ModelProvider string    `json:"modelProvider,omitempty"`
    SchemaVersion string    `json:"schemaVersion"`
    MatchEvalID   string    `json:"matchEvalId,omitempty"` // Reference to matcheval
}
```

## DocumentSet

```go
type DocumentSet struct {
    ResumePath        string `json:"resumePath"`
    CoverLetterPath   string `json:"coverLetterPath,omitempty"`
    ResumeFormat      string `json:"resumeFormat,omitempty"`      // md, pdf, docx
    CoverLetterFormat string `json:"coverLetterFormat,omitempty"`
}
```

## DocEvalResult

```go
type DocEvalResult struct {
    OverallScore  float64            `json:"overallScore"`  // 0-100
    WeightedScore float64            `json:"weightedScore"` // 0-10
    Categories    []DocEvalCategory  `json:"categories"`
    Findings      []DocEvalFinding   `json:"findings"`
}
```

## DocEvalCategory

```go
type DocEvalCategory struct {
    Category      DocCategoryType `json:"category"`
    Weight        float64         `json:"weight"`
    Score         float64         `json:"score"`
    MaxScore      float64         `json:"maxScore"`
    Status        ScoreStatus     `json:"status"`
    Justification string          `json:"justification"`
    Evidence      []string        `json:"evidence"`
    Suggestions   []string        `json:"suggestions,omitempty"`
}
```

### DocCategoryType

```go
const (
    // Core categories
    DocCategoryKeywordCoverage      DocCategoryType = "keyword_coverage"
    DocCategoryAchievementRelevance DocCategoryType = "achievement_relevance"
    DocCategoryGapMitigation        DocCategoryType = "gap_mitigation"
    DocCategoryQuantification       DocCategoryType = "quantification"
    DocCategoryATSCompatibility     DocCategoryType = "ats_compatibility"
    DocCategoryNarrativeCoherence   DocCategoryType = "narrative_coherence"

    // Resume-specific
    DocCategorySummaryImpact        DocCategoryType = "summary_impact"
    DocCategoryExperienceOrdering   DocCategoryType = "experience_ordering"
    DocCategorySkillsAlignment      DocCategoryType = "skills_alignment"
    DocCategoryFormatting           DocCategoryType = "formatting"

    // Cover letter-specific
    DocCategoryOpeningHook          DocCategoryType = "opening_hook"
    DocCategoryValueProposition     DocCategoryType = "value_proposition"
    DocCategoryCompanyFit           DocCategoryType = "company_fit"
    DocCategoryCallToAction         DocCategoryType = "call_to_action"
)
```

## DocEvalFinding

```go
type DocEvalFinding struct {
    ID             string          `json:"id"`
    Document       DocumentType    `json:"document"`
    Category       DocCategoryType `json:"category"`
    FindingType    FindingType     `json:"findingType"`
    Severity       Severity        `json:"severity"`
    Title          string          `json:"title"`
    Description    string          `json:"description"`
    Location       string          `json:"location,omitempty"`
    CurrentText    string          `json:"currentText,omitempty"`
    SuggestedText  string          `json:"suggestedText,omitempty"`
    Recommendation string          `json:"recommendation,omitempty"`
    Effort         EffortLevel     `json:"effort,omitempty"`
}
```

### DocumentType

```go
const (
    DocumentResume      DocumentType = "resume"
    DocumentCoverLetter DocumentType = "cover_letter"
    DocumentBoth        DocumentType = "both"
)
```

### Finding ID Convention

| Type | Format | Example |
|------|--------|---------|
| Doc Strength | DS### | DS001, DS002 |
| Doc Gap | DG### | DG001, DG002 |

## DocEvalDecision

```go
type DocEvalDecision struct {
    Status        DocDecisionStatus `json:"status"`
    Passed        bool              `json:"passed"`
    Rationale     string            `json:"rationale"`
    FindingCounts FindingCounts     `json:"findingCounts"`
    OverallScore  float64           `json:"overallScore"`
    WeightedScore float64           `json:"weightedScore"`
    ReadyToSubmit bool              `json:"readyToSubmit"`
}
```

### DocDecisionStatus

```go
const (
    DocDecisionExcellent     DocDecisionStatus = "excellent"       // Ready to submit
    DocDecisionGood          DocDecisionStatus = "good"            // Minor improvements
    DocDecisionNeedsWork     DocDecisionStatus = "needs_work"      // Address issues
    DocDecisionMajorRevision DocDecisionStatus = "major_revision"  // Significant rewrite
)
```

## DocEvalCriteria

```go
type DocEvalCriteria struct {
    MaxCritical        int     `json:"maxCritical"`        // Default 0
    MaxHigh            int     `json:"maxHigh"`            // Default 0
    MaxMedium          int     `json:"maxMedium"`          // Default 2
    MinScore           float64 `json:"minScore"`           // Default 75.0
    MinKeywordCoverage float64 `json:"minKeywordCoverage"` // Default 70.0
}

func DefaultDocEvalCriteria() DocEvalCriteria {
    return DocEvalCriteria{
        MaxCritical:        0,
        MaxHigh:            0,
        MaxMedium:          2,
        MinScore:           75.0,
        MinKeywordCoverage: 70.0,
    }
}
```

## Standard Rubrics

```go
func StandardDocEvalRubrics() DocEvalRubricSet {
    return DocEvalRubricSet{
        ID:      "doc-eval-v1",
        Name:    "Document Quality Rubrics",
        Rubrics: []DocEvalRubric{
            {
                Category:    DocCategoryKeywordCoverage,
                Weight:      0.20,
                Document:    DocumentBoth,
                Anchors: []ScoreAnchor{
                    {MinScore: 9, MaxScore: 10, Label: "Excellent"},
                    {MinScore: 7, MaxScore: 8.9, Label: "Good"},
                    {MinScore: 5, MaxScore: 6.9, Label: "Partial"},
                    {MinScore: 0, MaxScore: 4.9, Label: "Poor"},
                },
            },
            // ... more rubrics
        },
    }
}
```

## JSON Schema

See `schema/doceval.schema.json` for the JSON Schema definition.
