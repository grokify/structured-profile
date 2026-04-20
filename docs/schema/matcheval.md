# Match Evaluation Schema

The match evaluation schema (`schema/matcheval.go`) defines the structure for profile-to-JD fit assessment.

## MatchEvaluation

Top-level structure:

```go
type MatchEvaluation struct {
    Schema    string          `json:"$schema,omitempty"`
    Metadata  MatchMetadata   `json:"metadata"`
    Job       JobSummary      `json:"job"`
    Profile   ProfileSummary  `json:"profile"`
    Match     MatchResult     `json:"match"`
    Decision  MatchDecision   `json:"decision"`
    NextSteps []ActionItem    `json:"nextSteps,omitempty"`
}
```

## MatchMetadata

```go
type MatchMetadata struct {
    EvaluationID  string    `json:"evaluationId"`
    GeneratedAt   time.Time `json:"generatedAt"`
    GeneratedBy   string    `json:"generatedBy"`
    Model         string    `json:"model,omitempty"`
    ModelProvider string    `json:"modelProvider,omitempty"`
    SchemaVersion string    `json:"schemaVersion"`
}
```

## JobSummary

```go
type JobSummary struct {
    DocumentPath string   `json:"documentPath"`
    Title        string   `json:"title"`
    Company      string   `json:"company"`
    Level        string   `json:"level"`
    Location     []string `json:"location,omitempty"`
    Compensation string   `json:"compensation,omitempty"`
}
```

## MatchResult

```go
type MatchResult struct {
    OverallScore  float64         `json:"overallScore"`  // 0-100
    WeightedScore float64         `json:"weightedScore"` // 0-10
    Categories    []MatchCategory `json:"categories"`
    Findings      []MatchFinding  `json:"findings"`
}
```

## MatchCategory

```go
type MatchCategory struct {
    Category      MatchCategoryType `json:"category"`
    Weight        float64           `json:"weight"`        // 0.0-1.0
    Score         float64           `json:"score"`         // 0-10
    MaxScore      float64           `json:"maxScore"`
    Status        ScoreStatus       `json:"status"`
    MatchType     MatchType         `json:"matchType"`
    Justification string            `json:"justification"`
    Evidence      []string          `json:"evidence"`
    JDRequirement string            `json:"jdRequirement,omitempty"`
    ProfileMatch  string            `json:"profileMatch,omitempty"`
}
```

### MatchCategoryType

```go
const (
    CategoryTechnicalSkills    MatchCategoryType = "technical_skills"
    CategoryDomainExperience   MatchCategoryType = "domain_experience"
    CategoryLeadership         MatchCategoryType = "leadership"
    CategoryYearsExperience    MatchCategoryType = "years_experience"
    CategoryEducation          MatchCategoryType = "education"
    CategoryIndustryKnowledge  MatchCategoryType = "industry_knowledge"
    CategoryCommunication      MatchCategoryType = "communication"
    CategoryCulturalFit        MatchCategoryType = "cultural_fit"
    CategoryPlatformScale      MatchCategoryType = "platform_scale"
    CategoryAPIsSDKs           MatchCategoryType = "apis_sdks"
    CategoryDeveloperRelations MatchCategoryType = "developer_relations"
    CategoryIdentitySecurity   MatchCategoryType = "identity_security"
    CategoryCompliance         MatchCategoryType = "compliance"
    CategoryThoughtLeadership  MatchCategoryType = "thought_leadership"
    CategoryTeamManagement     MatchCategoryType = "team_management"
    CategoryCrossFunctional    MatchCategoryType = "cross_functional"
)
```

### MatchType

```go
const (
    MatchTypeExact     MatchType = "exact"     // Specific skill required
    MatchTypeSemantic  MatchType = "semantic"  // Equivalent skills accepted
    MatchTypePartial   MatchType = "partial"   // Related but not identical
    MatchTypeThreshold MatchType = "threshold" // Meets minimum level
    MatchTypeBoolean   MatchType = "boolean"   // Present or absent
    MatchTypeWeighted  MatchType = "weighted"  // Combined sub-factors
)
```

## MatchFinding

```go
type MatchFinding struct {
    ID             string            `json:"id"`
    Category       MatchCategoryType `json:"category"`
    FindingType    FindingType       `json:"findingType"` // strength or gap
    Severity       Severity          `json:"severity"`
    Title          string            `json:"title"`
    Description    string            `json:"description"`
    Evidence       string            `json:"evidence,omitempty"`
    Recommendation string            `json:"recommendation,omitempty"`
    Effort         EffortLevel       `json:"effort,omitempty"`
}
```

### Finding ID Convention

| Type | Format | Example |
|------|--------|---------|
| Strength | S### | S001, S002 |
| Gap | G### | G001, G002 |

## MatchDecision

```go
type MatchDecision struct {
    Status        DecisionStatus `json:"status"`
    Passed        bool           `json:"passed"`
    Rationale     string         `json:"rationale"`
    FindingCounts FindingCounts  `json:"findingCounts"`
    OverallScore  float64        `json:"overallScore"`
    WeightedScore float64        `json:"weightedScore"`
}
```

### DecisionStatus

```go
const (
    DecisionPass        DecisionStatus = "pass"         // Strong match
    DecisionConditional DecisionStatus = "conditional"  // Good with caveats
    DecisionFail        DecisionStatus = "fail"         // Blocking gaps
    DecisionHumanReview DecisionStatus = "human_review" // Borderline
)
```

## PassCriteria

```go
type PassCriteria struct {
    MaxCritical int     `json:"maxCritical"` // Default 0
    MaxHigh     int     `json:"maxHigh"`     // Default 0
    MaxMedium   int     `json:"maxMedium"`   // Default -1 (unlimited)
    MinScore    float64 `json:"minScore"`    // Default 70.0
}

func DefaultPassCriteria() PassCriteria {
    return PassCriteria{
        MaxCritical: 0,
        MaxHigh:     0,
        MaxMedium:   -1,
        MinScore:    70.0,
    }
}

func StrictPassCriteria() PassCriteria {
    return PassCriteria{
        MaxCritical: 0,
        MaxHigh:     0,
        MaxMedium:   2,
        MinScore:    80.0,
    }
}
```

## JSON Schema

See `schema/matcheval.schema.json` for the JSON Schema definition.
