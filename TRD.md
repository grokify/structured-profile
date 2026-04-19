# Technical Requirements Document: structured-profile

**Version:** 1.0.0
**Status:** Draft
**Last Updated:** 2026-02-08

## 1. System Architecture

### 1.1 High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              structured-profile                              │
├─────────────────────────────────────────────────────────────────────────────┤
│  CLI Layer                                                                   │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐   │
│  │ profile │ │ resume  │ │ cover   │ │  eval   │ │  apply  │ │  prep   │   │
│  │   cmd   │ │   cmd   │ │   cmd   │ │   cmd   │ │   cmd   │ │   cmd   │   │
│  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘   │
├───────┴──────────┴──────────┴──────────┴──────────┴──────────┴─────────────┤
│  Service Layer                                                               │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │  Profile    │ │   Resume    │ │ Evaluation  │ │  Interview  │           │
│  │  Service    │ │   Service   │ │   Service   │ │   Service   │           │
│  └──────┬──────┘ └──────┬──────┘ └──────┬──────┘ └──────┬──────┘           │
├─────────┴────────────────┴──────────────┴────────────────┴─────────────────┤
│  Domain Layer                                                                │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │  schema/  - Core domain types (Profile, Experience, Skill, etc.)     │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
├─────────────────────────────────────────────────────────────────────────────┤
│  Infrastructure Layer                                                        │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │    Store    │ │   Export    │ │   LLM       │ │  External   │           │
│  │ (JSON/PG)   │ │ (MD/PDF)    │ │  Provider   │ │    APIs     │           │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘           │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Package Structure

```
structured-profile/
├── cmd/
│   └── sprofile/              # Main CLI application
│       └── main.go
├── schema/                    # Core domain types
│   ├── profile.go             # Profile, Contact
│   ├── experience.go          # Tenure, Position, Achievement (STAR)
│   ├── skill.go               # Skill, SkillCategory
│   ├── education.go           # Education, Certification
│   ├── credential.go          # VerifiableCredential (GitHub, SO, etc.)
│   ├── opportunity.go         # JobOpportunity, JobDescription
│   ├── application.go         # Application, Interview, Feedback
│   ├── resume.go              # Resume, CoverLetter
│   └── base.go                # BaseEntity (ID, timestamps)
├── store/                     # Data persistence
│   ├── store.go               # Store interface
│   ├── json/                  # JSON file backend
│   │   └── store.go
│   └── postgres/              # PostgreSQL backend (future)
│       ├── store.go
│       └── migrations/
├── service/                   # Business logic
│   ├── profile.go             # Profile CRUD operations
│   ├── resume.go              # Resume generation
│   ├── coverletter.go         # Cover letter generation
│   ├── evaluation.go          # LLM-as-a-Judge integration
│   ├── application.go         # Application tracking
│   └── interview.go           # Interview prep generation
├── agents/                    # Multi-agent definitions (multi-agent-spec)
│   ├── resume-pipeline/       # Resume generation agent team
│   │   ├── jd-analyst.md      # JD parsing agent
│   │   ├── resume-generator.md # Resume generation agent
│   │   ├── cover-letter-generator.md # Cover letter agent
│   │   ├── resume-evaluator.md # LLM-as-a-Judge agent
│   │   └── resume-refiner.md  # Iterative refinement agent
│   └── teams/
│       └── resume-pipeline.json # Team definition
├── export/                    # Output formats
│   ├── markdown.go            # Markdown generation
│   ├── pandoc.go              # PDF/DOCX via Pandoc
│   ├── json.go                # JSON IR
│   └── h5p.go                 # H5P quiz export
├── llm/                       # LLM provider abstraction
│   ├── provider.go            # Provider interface
│   ├── anthropic.go           # Anthropic Claude
│   ├── openai.go              # OpenAI GPT
│   └── bedrock.go             # AWS Bedrock
├── orchestrator/              # Agent team orchestration
│   ├── orchestrator.go        # Workflow execution
│   ├── runner.go              # Step runner
│   └── context.go             # Shared context management
├── jdparser/                  # Job description parsing (via LLM)
│   └── parser.go              # LLM-powered JD analysis
├── matcher/                   # Profile-JD matching
│   └── matcher.go             # Relevance scoring
├── internal/                  # Internal utilities
│   └── testutil/              # Test helpers
├── examples/                  # Example data and configs
│   ├── profiles/
│   └── opportunities/
├── docs/                      # Documentation (MkDocs)
├── go.mod
├── go.sum
├── PRD.md
├── TRD.md
├── PLAN.md
├── ROADMAP.md
└── README.md
```

## 2. Data Models

### 2.1 Core Entity Design

All entities extend `BaseEntity` for database compatibility:

```go
package schema

import "time"

type BaseEntity struct {
    ID        string     `json:"id" db:"id"`
    CreatedAt time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

func NewBaseEntity() BaseEntity {
    now := time.Now().UTC()
    return BaseEntity{
        ID:        uuid.New().String(),
        CreatedAt: now,
        UpdatedAt: now,
    }
}
```

### 2.2 Profile Schema

```go
package schema

type Profile struct {
    BaseEntity

    // Personal Information
    Name        string       `json:"name"`
    Email       string       `json:"email"`
    Phone       string       `json:"phone,omitempty"`
    Location    string       `json:"location,omitempty"`

    // Online Presence
    Links       []Link       `json:"links,omitempty"`

    // Summaries by domain/filter
    Summaries   Summaries    `json:"summaries,omitempty"`

    // Professional Data (stored separately, linked by ProfileID)
    // Tenures      []Tenure      - via store.ListTenures(profileID)
    // Skills       []Skill       - via store.ListSkills(profileID)
    // Education    []Education   - via store.ListEducation(profileID)
    // etc.
}

type Link struct {
    Type string `json:"type"` // "linkedin", "github", "website", "portfolio"
    URL  string `json:"url"`
    Text string `json:"text,omitempty"`
}

type Summaries struct {
    Default   string            `json:"default"`
    ByDomain  map[string]string `json:"by_domain,omitempty"` // "devx" -> "..."
}
```

### 2.3 Experience Schema (STAR Format)

```go
package schema

type Tenure struct {
    BaseEntity
    ProfileID   string     `json:"profile_id" db:"profile_id"`
    Company     string     `json:"company"`
    StartDate   Date       `json:"start_date"`
    EndDate     *Date      `json:"end_date,omitempty"` // nil = current
    Description string     `json:"description,omitempty"`
    Location    string     `json:"location,omitempty"`

    // Positions within this tenure
    // Positions []Position - via store.ListPositions(tenureID)
}

type Position struct {
    BaseEntity
    TenureID    string     `json:"tenure_id" db:"tenure_id"`
    Title       string     `json:"title"`
    StartDate   Date       `json:"start_date"`
    EndDate     *Date      `json:"end_date,omitempty"`
    Description string     `json:"description,omitempty"`

    // Domain-specific configurations
    DomainConfigs []PositionDomainConfig `json:"domain_configs,omitempty"`

    // Achievements within this position
    // Achievements []Achievement - via store.ListAchievements(positionID)
}

type PositionDomainConfig struct {
    BaseEntity
    PositionID       string   `json:"position_id" db:"position_id"`
    Domain           string   `json:"domain"`           // "devx", "iam", "platform"
    Skills           []string `json:"skills"`           // Domain-specific skill list
    AchievementOrder []string `json:"achievement_order"` // Achievement IDs in display order
}

type Achievement struct {
    BaseEntity
    PositionID  string   `json:"position_id" db:"position_id"`
    Name        string   `json:"name"`        // Slug for ordering/reference

    // STAR Format
    Situation   string   `json:"situation"`
    Task        string   `json:"task"`
    Action      string   `json:"action"`
    Result      string   `json:"result"`

    // Metadata
    Skills      []string `json:"skills,omitempty"`
    Tags        []string `json:"tags,omitempty"`
    Metrics     Metrics  `json:"metrics,omitempty"`

    // Display control
    SkipDisplay bool     `json:"skip_display,omitempty"`
}

type Metrics struct {
    Values map[string]string `json:"values,omitempty"` // "revenue_increase": "40%"
}

// STAR string generation
func (a Achievement) STARString() string {
    return strings.TrimSpace(strings.Join([]string{
        a.Situation, a.Task, a.Action, a.Result,
    }, " "))
}
```

### 2.4 Skills Schema

```go
package schema

type Skill struct {
    BaseEntity
    ProfileID   string `json:"profile_id" db:"profile_id"`
    Name        string `json:"name"`
    Category    string `json:"category"`    // "technical", "domain", "soft"
    Subcategory string `json:"subcategory"` // "cloud", "iam", "leadership"
    Level       string `json:"level"`       // "expert", "proficient", "familiar"

    // Verification
    Verified    bool   `json:"verified,omitempty"`
    VerifyURL   string `json:"verify_url,omitempty"`
}

type SkillCategory struct {
    Name        string   `json:"name"`
    Skills      []string `json:"skills"`
    Description string   `json:"description,omitempty"`
}
```

### 2.5 Education and Certifications

```go
package schema

type Education struct {
    BaseEntity
    ProfileID   string `json:"profile_id" db:"profile_id"`
    Institution string `json:"institution"`
    Degree      string `json:"degree"`
    Field       string `json:"field,omitempty"`
    StartDate   Date   `json:"start_date,omitempty"`
    EndDate     *Date  `json:"end_date,omitempty"`
    Honors      string `json:"honors,omitempty"`
    GPA         string `json:"gpa,omitempty"`
    Display     bool   `json:"display"`
}

type Certification struct {
    BaseEntity
    ProfileID      string `json:"profile_id" db:"profile_id"`
    Name           string `json:"name"`
    Issuer         string `json:"issuer"`
    IssueDate      Date   `json:"issue_date"`
    ExpirationDate *Date  `json:"expiration_date,omitempty"`
    CredentialID   string `json:"credential_id,omitempty"`
    CredentialURL  string `json:"credential_url,omitempty"`
    Display        bool   `json:"display"`
}
```

### 2.6 Verifiable Credentials

```go
package schema

type VerifiableCredential struct {
    BaseEntity
    ProfileID   string         `json:"profile_id" db:"profile_id"`
    Type        string         `json:"type"`     // "github", "stackoverflow", "linkedin"
    Username    string         `json:"username"`
    ProfileURL  string         `json:"profile_url"`

    // Platform-specific data
    Data        CredentialData `json:"data,omitempty"`

    // Verification
    VerifiedAt  *time.Time     `json:"verified_at,omitempty"`
}

type CredentialData struct {
    // GitHub
    Repositories    int      `json:"repositories,omitempty"`
    Stars           int      `json:"stars,omitempty"`
    Contributions   int      `json:"contributions,omitempty"`
    Languages       []string `json:"languages,omitempty"`

    // StackOverflow
    Reputation      int      `json:"reputation,omitempty"`
    GoldBadges      int      `json:"gold_badges,omitempty"`
    SilverBadges    int      `json:"silver_badges,omitempty"`
    BronzeBadges    int      `json:"bronze_badges,omitempty"`
    TopTags         []string `json:"top_tags,omitempty"`

    // LinkedIn
    Connections     int      `json:"connections,omitempty"`
    Endorsements    int      `json:"endorsements,omitempty"`
}
```

### 2.7 Opportunity and Application

```go
package schema

type Opportunity struct {
    BaseEntity
    ProfileID       string          `json:"profile_id" db:"profile_id"`

    // Job Details
    Company         string          `json:"company"`
    Position        string          `json:"position"`
    Location        string          `json:"location,omitempty"`
    Remote          bool            `json:"remote,omitempty"`
    URL             string          `json:"url,omitempty"`

    // Job Description
    JobDescRaw      string          `json:"job_desc_raw"`
    JobDescParsed   *JobDescParsed  `json:"job_desc_parsed,omitempty"`

    // Company Research
    CompanyInfo     *CompanyInfo    `json:"company_info,omitempty"`

    // Generated Documents
    ResumeID        string          `json:"resume_id,omitempty"`
    CoverLetterID   string          `json:"cover_letter_id,omitempty"`
    EvaluationID    string          `json:"evaluation_id,omitempty"`
}

type JobDescParsed struct {
    RequiredSkills   []string `json:"required_skills"`
    PreferredSkills  []string `json:"preferred_skills"`
    ExperienceYears  int      `json:"experience_years"`
    Keywords         []string `json:"keywords"`
    Responsibilities []string `json:"responsibilities"`
}

type CompanyInfo struct {
    Name        string   `json:"name"`
    Industry    string   `json:"industry,omitempty"`
    Size        string   `json:"size,omitempty"`        // "startup", "mid", "enterprise"
    Values      []string `json:"values,omitempty"`
    Culture     []string `json:"culture,omitempty"`
    TechStack   []string `json:"tech_stack,omitempty"`
    Links       []Link   `json:"links,omitempty"`
}

type Application struct {
    BaseEntity
    OpportunityID string            `json:"opportunity_id" db:"opportunity_id"`

    // Status Tracking
    Status        ApplicationStatus `json:"status"`
    AppliedAt     *time.Time        `json:"applied_at,omitempty"`

    // Documents Used
    ResumeVersion   string          `json:"resume_version,omitempty"`
    CoverLetterVersion string       `json:"cover_letter_version,omitempty"`

    // Outcome
    Outcome       *ApplicationOutcome `json:"outcome,omitempty"`

    // Interviews are stored separately
    // Interviews []Interview - via store.ListInterviews(applicationID)
}

type ApplicationStatus string

const (
    StatusDraft      ApplicationStatus = "draft"
    StatusSubmitted  ApplicationStatus = "submitted"
    StatusScreening  ApplicationStatus = "screening"
    StatusInterview  ApplicationStatus = "interview"
    StatusOffer      ApplicationStatus = "offer"
    StatusAccepted   ApplicationStatus = "accepted"
    StatusRejected   ApplicationStatus = "rejected"
    StatusWithdrawn  ApplicationStatus = "withdrawn"
)

type ApplicationOutcome struct {
    Result        string   `json:"result"`       // "offer", "rejected", "withdrawn"
    OfferDetails  string   `json:"offer_details,omitempty"`
    RejectionReason string `json:"rejection_reason,omitempty"`
    LessonsLearned []string `json:"lessons_learned,omitempty"`
    Strengths     []string `json:"strengths,omitempty"`
    Improvements  []string `json:"improvements,omitempty"`
}
```

### 2.8 Interview and Feedback

```go
package schema

type Interview struct {
    BaseEntity
    ApplicationID string        `json:"application_id" db:"application_id"`

    // Interview Details
    Round         int           `json:"round"`
    Type          InterviewType `json:"type"`
    ScheduledAt   time.Time     `json:"scheduled_at"`
    DurationMins  int           `json:"duration_mins,omitempty"`

    // Participants
    Interviewers  []string      `json:"interviewers,omitempty"`

    // Content
    Questions     []InterviewQuestion `json:"questions,omitempty"`

    // Assessment
    SelfAssessment *SelfAssessment   `json:"self_assessment,omitempty"`
    Feedback       *InterviewFeedback `json:"feedback,omitempty"`
}

type InterviewType string

const (
    InterviewPhone      InterviewType = "phone"
    InterviewTechnical  InterviewType = "technical"
    InterviewBehavioral InterviewType = "behavioral"
    InterviewSystem     InterviewType = "system_design"
    InterviewCoding     InterviewType = "coding"
    InterviewOnsite     InterviewType = "onsite"
    InterviewFinal      InterviewType = "final"
)

type InterviewQuestion struct {
    Question      string   `json:"question"`
    Category      string   `json:"category"`      // "behavioral", "technical", etc.
    MyAnswer      string   `json:"my_answer,omitempty"`
    IdealAnswer   string   `json:"ideal_answer,omitempty"`
    RelatedSTAR   string   `json:"related_star,omitempty"` // Achievement ID
    Tags          []string `json:"tags,omitempty"`
    AddToPrep     bool     `json:"add_to_prep,omitempty"`  // Flag to add to prep bank
}

type SelfAssessment struct {
    OverallRating  int      `json:"overall_rating"`  // 1-5
    WentWell       []string `json:"went_well"`
    CouldImprove   []string `json:"could_improve"`
    QuestionsToAdd []string `json:"questions_to_add"`
    Notes          string   `json:"notes,omitempty"`
}

type InterviewFeedback struct {
    Source        string   `json:"source"`        // "recruiter", "interviewer", "self"
    Rating        string   `json:"rating"`        // "strong_yes", "yes", "no", "strong_no"
    Strengths     []string `json:"strengths"`
    Concerns      []string `json:"concerns"`
    Recommendation string  `json:"recommendation,omitempty"`
    Notes         string   `json:"notes,omitempty"`
}
```

### 2.9 Interview Prep (H5P Integration)

```go
package schema

type InterviewPrepSet struct {
    BaseEntity
    ProfileID     string              `json:"profile_id" db:"profile_id"`
    OpportunityID string              `json:"opportunity_id,omitempty"` // Optional: role-specific

    Title         string              `json:"title"`
    Description   string              `json:"description,omitempty"`

    Sections      []PrepSection       `json:"sections"`

    // H5P Export Settings
    PassPercentage int               `json:"pass_percentage"`
    FeedbackRanges []FeedbackRange   `json:"feedback_ranges,omitempty"`
}

type PrepSection struct {
    Name       string         `json:"name"`       // "Behavioral", "Technical"
    Topic      string         `json:"topic"`      // "Leadership", "System Design"
    Questions  []PrepQuestion `json:"questions"`
}

type PrepQuestion struct {
    BaseEntity
    Question          string   `json:"question"`
    Difficulty        string   `json:"difficulty"` // "easy", "medium", "hard"
    Category          string   `json:"category"`

    // Answer Options (for multiple choice)
    Answers           []PrepAnswer `json:"answers"`

    // Linked Content
    RelatedSTARID     string   `json:"related_star_id,omitempty"`
    RelatedSkills     []string `json:"related_skills,omitempty"`

    // Learning
    LearningObjective string   `json:"learning_objective,omitempty"`

    // Source Tracking
    Source            string   `json:"source,omitempty"` // "interview:app-123", "book:xyz"
    SourceInterview   string   `json:"source_interview,omitempty"`
}

type PrepAnswer struct {
    Text     string `json:"text"`
    Correct  bool   `json:"correct"`
    Feedback string `json:"feedback,omitempty"`
    Tip      string `json:"tip,omitempty"`
}

type FeedbackRange struct {
    From int    `json:"from"`
    To   int    `json:"to"`
    Text string `json:"text"`
}
```

## 3. Store Interface

```go
package store

import (
    "context"
    "github.com/grokify/structured-profile/schema"
)

type Store interface {
    // Profile
    GetProfile(ctx context.Context, id string) (*schema.Profile, error)
    ListProfiles(ctx context.Context) ([]schema.Profile, error)
    SaveProfile(ctx context.Context, p *schema.Profile) error
    DeleteProfile(ctx context.Context, id string) error

    // Tenures
    GetTenure(ctx context.Context, id string) (*schema.Tenure, error)
    ListTenures(ctx context.Context, profileID string) ([]schema.Tenure, error)
    SaveTenure(ctx context.Context, t *schema.Tenure) error
    DeleteTenure(ctx context.Context, id string) error

    // Positions
    GetPosition(ctx context.Context, id string) (*schema.Position, error)
    ListPositions(ctx context.Context, tenureID string) ([]schema.Position, error)
    SavePosition(ctx context.Context, p *schema.Position) error
    DeletePosition(ctx context.Context, id string) error

    // Achievements
    GetAchievement(ctx context.Context, id string) (*schema.Achievement, error)
    ListAchievements(ctx context.Context, positionID string) ([]schema.Achievement, error)
    SearchAchievements(ctx context.Context, profileID string, tags []string) ([]schema.Achievement, error)
    SaveAchievement(ctx context.Context, a *schema.Achievement) error
    DeleteAchievement(ctx context.Context, id string) error

    // Skills
    GetSkill(ctx context.Context, id string) (*schema.Skill, error)
    ListSkills(ctx context.Context, profileID string) ([]schema.Skill, error)
    ListSkillsByCategory(ctx context.Context, profileID, category string) ([]schema.Skill, error)
    SaveSkill(ctx context.Context, s *schema.Skill) error
    DeleteSkill(ctx context.Context, id string) error

    // Education
    ListEducation(ctx context.Context, profileID string) ([]schema.Education, error)
    SaveEducation(ctx context.Context, e *schema.Education) error

    // Certifications
    ListCertifications(ctx context.Context, profileID string) ([]schema.Certification, error)
    SaveCertification(ctx context.Context, c *schema.Certification) error

    // Credentials
    ListCredentials(ctx context.Context, profileID string) ([]schema.VerifiableCredential, error)
    SaveCredential(ctx context.Context, c *schema.VerifiableCredential) error

    // Opportunities
    GetOpportunity(ctx context.Context, id string) (*schema.Opportunity, error)
    ListOpportunities(ctx context.Context, profileID string) ([]schema.Opportunity, error)
    SaveOpportunity(ctx context.Context, o *schema.Opportunity) error

    // Applications
    GetApplication(ctx context.Context, id string) (*schema.Application, error)
    ListApplications(ctx context.Context, opportunityID string) ([]schema.Application, error)
    ListApplicationsByStatus(ctx context.Context, profileID string, status schema.ApplicationStatus) ([]schema.Application, error)
    SaveApplication(ctx context.Context, a *schema.Application) error

    // Interviews
    ListInterviews(ctx context.Context, applicationID string) ([]schema.Interview, error)
    SaveInterview(ctx context.Context, i *schema.Interview) error

    // Interview Prep
    GetInterviewPrepSet(ctx context.Context, id string) (*schema.InterviewPrepSet, error)
    ListInterviewPrepSets(ctx context.Context, profileID string) ([]schema.InterviewPrepSet, error)
    SaveInterviewPrepSet(ctx context.Context, ps *schema.InterviewPrepSet) error

    // Transactions (for database backend)
    BeginTx(ctx context.Context) (Store, error)
    Commit() error
    Rollback() error
}
```

## 4. LLM Integration

### 4.1 Provider Interface

```go
package llm

import "context"

type Provider interface {
    // Generate text completion
    Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)

    // Generate structured output (JSON)
    CompleteJSON(ctx context.Context, req CompletionRequest, schema any) (any, error)

    // Get provider info
    Info() ProviderInfo
}

type CompletionRequest struct {
    SystemPrompt string
    UserPrompt   string
    Temperature  float64
    MaxTokens    int
    Model        string
}

type CompletionResponse struct {
    Content      string
    Model        string
    InputTokens  int
    OutputTokens int
    Latency      time.Duration
}

type ProviderInfo struct {
    Name    string
    Models  []string
    Default string
}
```

### 4.2 Evaluation Integration

```go
package service

import (
    "context"
    "github.com/agentplexus/structured-evaluation/evaluation"
    "github.com/grokify/structured-profile/schema"
)

type EvaluationService struct {
    llm    llm.Provider
    store  store.Store
}

func (s *EvaluationService) EvaluateResume(
    ctx context.Context,
    resume *schema.Resume,
    opportunity *schema.Opportunity,
) (*evaluation.EvaluationReport, error) {
    // 1. Create evaluation report
    report := evaluation.NewEvaluationReport("resume-jd-match", resume.ID)

    // 2. Set up judge metadata
    judge := evaluation.NewJudgeMetadata(s.llm.Info().Default).
        WithProvider(s.llm.Info().Name).
        WithPrompt("resume-eval-v1", "1.0").
        WithRubric("resume-jd-match-v1", "1.0")
    report.Judge = judge

    // 3. Evaluate each category via LLM
    categories := []string{
        "skills_alignment",
        "experience_level",
        "achievement_clarity",
        "industry_fit",
        "communication",
    }

    for _, cat := range categories {
        score, justification := s.evaluateCategory(ctx, resume, opportunity, cat)
        report.AddCategory(evaluation.NewCategoryScore(
            cat, categoryWeights[cat], score, justification,
        ))
    }

    // 4. Generate findings
    findings := s.generateFindings(ctx, resume, opportunity)
    for _, f := range findings {
        report.AddFinding(f)
    }

    // 5. Finalize and return
    report.Finalize("sprofile eval resume")
    return report, nil
}
```

### 4.3 Agent Team Architecture (multi-agent-spec)

Resume generation, cover letter creation, and evaluation are orchestrated by a multi-agent team defined using the [multi-agent-spec](https://github.com/agentplexus/multi-agent-spec) framework.

#### Agent Definitions

Agent definitions are stored as Markdown files with YAML frontmatter in `agents/resume-pipeline/`:

**jd-analyst.md** - Job Description Analysis Agent
```markdown
---
name: jd-analyst
namespace: resume-pipeline
description: Analyzes job descriptions to extract requirements and keywords
model: sonnet
tools: [Read, WebFetch, WebSearch]
---

# JD Analyst Agent

You are a job description analyst specializing in extracting structured
information from job postings.

## Instructions

1. Parse the raw job description text
2. Extract required skills, preferred skills, and experience levels
3. Identify company culture signals and keywords
4. Output structured JobDescParsed JSON
...
```

**resume-generator.md** - Resume Generation Agent
```markdown
---
name: resume-generator
namespace: resume-pipeline
description: Generates tailored resumes from master profile
model: sonnet
tools: [Read]
dependencies: [jd-analyst]
---

# Resume Generator Agent

You are a resume optimization specialist. Given a master profile and
parsed job description, you select and tailor relevant experiences.

## Instructions

1. Match master profile skills to JD requirements
2. Select most relevant STAR achievements
3. Reorder experiences by relevance score
4. Generate customized summary for this opportunity
5. Output tailored resume in Markdown format
...
```

**resume-evaluator.md** - LLM-as-a-Judge Agent
```markdown
---
name: resume-evaluator
namespace: resume-pipeline
description: Evaluates resume quality against job descriptions
model: opus
tools: [Read]
dependencies: [resume-generator]
---

# Resume Evaluator Agent

You are an expert resume evaluator using LLM-as-a-Judge methodology.
Evaluate resumes against job descriptions using the structured-evaluation
framework.

## Evaluation Categories

1. **Skills Alignment (25%)**: Match between resume skills and JD requirements
2. **Experience Level (25%)**: Seniority and domain experience match
3. **Achievement Clarity (20%)**: STAR format quality and quantification
4. **Industry Fit (15%)**: Domain and industry experience relevance
5. **Communication (15%)**: Writing quality and presentation

## Output Format

Return a structured-evaluation EvaluationReport JSON with:
- Category scores with justifications
- Findings with severity levels
- Pass/fail decision with rationale
...
```

#### Team Definition

**teams/resume-pipeline.json**
```json
{
  "$schema": "https://agentplexus.github.io/multi-agent-spec/schemas/team.json",
  "name": "resume-pipeline",
  "version": "1.0.0",
  "description": "Multi-agent team for resume generation and evaluation",
  "agents": [
    "resume-pipeline/jd-analyst",
    "resume-pipeline/resume-generator",
    "resume-pipeline/cover-letter-generator",
    "resume-pipeline/resume-evaluator",
    "resume-pipeline/resume-refiner"
  ],
  "workflow": {
    "type": "dag",
    "steps": [
      {
        "name": "analyze-jd",
        "agent": "resume-pipeline/jd-analyst",
        "inputs": [
          {"name": "job_description", "type": "string", "required": true}
        ],
        "outputs": [
          {"name": "parsed_jd", "type": "object"}
        ]
      },
      {
        "name": "generate-resume",
        "agent": "resume-pipeline/resume-generator",
        "depends_on": ["analyze-jd"],
        "inputs": [
          {"name": "master_profile", "type": "object", "required": true},
          {"name": "parsed_jd", "type": "object", "from": "analyze-jd.parsed_jd"}
        ],
        "outputs": [
          {"name": "tailored_resume", "type": "string"}
        ]
      },
      {
        "name": "generate-cover-letter",
        "agent": "resume-pipeline/cover-letter-generator",
        "depends_on": ["analyze-jd", "generate-resume"],
        "inputs": [
          {"name": "master_profile", "type": "object", "required": true},
          {"name": "parsed_jd", "type": "object", "from": "analyze-jd.parsed_jd"},
          {"name": "resume", "type": "string", "from": "generate-resume.tailored_resume"}
        ],
        "outputs": [
          {"name": "cover_letter", "type": "string"}
        ]
      },
      {
        "name": "evaluate-resume",
        "agent": "resume-pipeline/resume-evaluator",
        "depends_on": ["analyze-jd", "generate-resume"],
        "inputs": [
          {"name": "resume", "type": "string", "from": "generate-resume.tailored_resume"},
          {"name": "parsed_jd", "type": "object", "from": "analyze-jd.parsed_jd"}
        ],
        "outputs": [
          {"name": "evaluation_report", "type": "object"}
        ]
      },
      {
        "name": "refine-resume",
        "agent": "resume-pipeline/resume-refiner",
        "depends_on": ["generate-resume", "evaluate-resume"],
        "inputs": [
          {"name": "resume", "type": "string", "from": "generate-resume.tailored_resume"},
          {"name": "evaluation", "type": "object", "from": "evaluate-resume.evaluation_report"}
        ],
        "outputs": [
          {"name": "refined_resume", "type": "string"}
        ]
      }
    ]
  },
  "context": "This team generates tailored resumes and cover letters from a master profile, evaluates them against job descriptions using LLM-as-a-Judge, and iteratively refines based on feedback."
}
```

#### Orchestrator Implementation

```go
package orchestrator

import (
    "context"

    multiagent "github.com/agentplexus/multi-agent-spec/sdk/go"
    "github.com/grokify/structured-profile/llm"
    "github.com/grokify/structured-profile/schema"
)

type ResumePipelineOrchestrator struct {
    team     *multiagent.Team
    agents   map[string]*multiagent.Agent
    provider llm.Provider
}

func NewResumePipelineOrchestrator(provider llm.Provider) (*ResumePipelineOrchestrator, error) {
    // Load team definition
    team, err := multiagent.LoadTeamFromFile("agents/teams/resume-pipeline.json")
    if err != nil {
        return nil, err
    }

    // Load agent definitions
    agents, err := multiagent.LoadAgentsFromDir("agents/resume-pipeline")
    if err != nil {
        return nil, err
    }

    return &ResumePipelineOrchestrator{
        team:     team,
        agents:   agents,
        provider: provider,
    }, nil
}

func (o *ResumePipelineOrchestrator) GenerateResume(
    ctx context.Context,
    profile *schema.Profile,
    jobDescription string,
) (*PipelineResult, error) {
    // Execute DAG workflow
    runner := NewWorkflowRunner(o.team.Workflow, o.agents, o.provider)

    inputs := map[string]any{
        "master_profile":  profile,
        "job_description": jobDescription,
    }

    results, err := runner.Execute(ctx, inputs)
    if err != nil {
        return nil, err
    }

    return &PipelineResult{
        ParsedJD:         results["analyze-jd.parsed_jd"],
        TailoredResume:   results["generate-resume.tailored_resume"].(string),
        CoverLetter:      results["generate-cover-letter.cover_letter"].(string),
        EvaluationReport: results["evaluate-resume.evaluation_report"],
        RefinedResume:    results["refine-resume.refined_resume"].(string),
    }, nil
}

type PipelineResult struct {
    ParsedJD         any
    TailoredResume   string
    CoverLetter      string
    EvaluationReport any
    RefinedResume    string
}
```

#### Workflow Runner

```go
package orchestrator

import (
    "context"
    "sort"

    multiagent "github.com/agentplexus/multi-agent-spec/sdk/go"
    "github.com/grokify/structured-profile/llm"
)

type WorkflowRunner struct {
    workflow *multiagent.Workflow
    agents   map[string]*multiagent.Agent
    provider llm.Provider
}

func (r *WorkflowRunner) Execute(ctx context.Context, inputs map[string]any) (map[string]any, error) {
    results := make(map[string]any)

    // Copy initial inputs to results
    for k, v := range inputs {
        results[k] = v
    }

    // Topological sort steps by dependencies
    ordered := r.topologicalSort()

    // Execute steps in order
    for _, step := range ordered {
        // Resolve input references
        stepInputs := r.resolveInputs(step, results)

        // Get agent for this step
        agent := r.agents[step.Agent]

        // Execute agent
        output, err := r.executeAgent(ctx, agent, stepInputs)
        if err != nil {
            return nil, err
        }

        // Store outputs with step prefix
        for _, port := range step.Outputs {
            results[step.Name+"."+port.Name] = output[port.Name]
        }
    }

    return results, nil
}

func (r *WorkflowRunner) executeAgent(
    ctx context.Context,
    agent *multiagent.Agent,
    inputs map[string]any,
) (map[string]any, error) {
    // Build prompt from agent instructions and inputs
    prompt := r.buildPrompt(agent, inputs)

    // Map canonical model to provider model
    model := r.provider.MapModel(string(agent.Model))

    // Execute LLM call
    resp, err := r.provider.CompleteJSON(ctx, llm.CompletionRequest{
        SystemPrompt: agent.Instructions,
        UserPrompt:   prompt,
        Model:        model,
        Temperature:  0.0, // Deterministic for evaluation
    }, nil)

    if err != nil {
        return nil, err
    }

    return resp.(map[string]any), nil
}
```

## 5. Export Pipeline

### 5.1 Markdown Generation

```go
package export

import (
    "strings"
    "github.com/grokify/structured-profile/schema"
)

type MarkdownExporter struct {
    opts MarkdownOptions
}

type MarkdownOptions struct {
    IncludeContact      bool
    IncludeSummary      bool
    IncludeExperience   bool
    IncludeEducation    bool
    IncludeCertifications bool
    IncludeSkills       bool
    Domain              string // Filter domain
    Template            string // Optional template name
}

func (e *MarkdownExporter) Export(p *FullProfile, opts MarkdownOptions) (string, error) {
    var lines []string

    // Header
    lines = append(lines, "# "+p.Profile.Name)

    // Contact
    if opts.IncludeContact {
        lines = append(lines, e.renderContact(p.Profile)...)
    }

    // Summary
    if opts.IncludeSummary {
        summary := p.Profile.Summaries.ForDomain(opts.Domain)
        lines = append(lines, "", "## Summary", "", summary)
    }

    // Experience
    if opts.IncludeExperience {
        lines = append(lines, e.renderExperience(p.Tenures, opts.Domain)...)
    }

    // ... other sections

    return strings.Join(lines, "\n"), nil
}
```

### 5.2 H5P Export

```go
package export

import (
    "github.com/grokify/h5p-go"
    "github.com/grokify/h5p-go/schemas"
    "github.com/grokify/structured-profile/schema"
)

type H5PExporter struct{}

func (e *H5PExporter) ExportPrepSet(ps *schema.InterviewPrepSet) (*h5p.H5PPackage, error) {
    builder := h5p.NewQuestionSetBuilder().
        SetTitle(ps.Title).
        SetPassPercentage(ps.PassPercentage)

    for _, section := range ps.Sections {
        for i, q := range section.Questions {
            // Convert to H5P answers
            answers := make([]schemas.AnswerOption, len(q.Answers))
            for j, a := range q.Answers {
                answers[j] = schemas.AnswerOption{
                    Text:    a.Text,
                    Correct: a.Correct,
                    TipsAndFeedback: &schemas.AnswerTipsAndFeedback{
                        Tip:            a.Tip,
                        ChosenFeedback: a.Feedback,
                    },
                }
            }

            // Create MultiChoice question with extensions
            mcParams := &schemas.MultiChoiceParams{
                Question: q.Question,
                Answers:  answers,
            }

            mcQuestion := h5p.NewMultiChoiceQuestion(mcParams).
                WithExtensions(
                    h5p.NewH5PGoExtension(section.Name, i+1).
                        WithTopic(section.Topic).
                        WithDifficulty(q.Difficulty).
                        WithLearningObjective(q.LearningObjective),
                )

            builder.AddTypedQuestion(mcQuestion)
        }
    }

    // Add feedback ranges
    for _, fr := range ps.FeedbackRanges {
        builder.AddOverallFeedback([]h5p.FeedbackRange{
            {From: fr.From, To: fr.To, Text: fr.Text},
        })
    }

    qs, err := builder.Build()
    if err != nil {
        return nil, err
    }

    return h5p.NewH5PPackageFromQuestionSet(qs)
}
```

## 6. Database Schema (PostgreSQL)

```sql
-- Profiles
CREATE TABLE profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT,
    phone TEXT,
    location TEXT,
    links JSONB DEFAULT '[]',
    summaries JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Tenures
CREATE TABLE tenures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID REFERENCES profiles(id) ON DELETE CASCADE,
    company TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    description TEXT,
    location TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_tenures_profile ON tenures(profile_id);

-- Positions
CREATE TABLE positions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenure_id UUID REFERENCES tenures(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    description TEXT,
    domain_configs JSONB DEFAULT '[]',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_positions_tenure ON positions(tenure_id);

-- Achievements (STAR)
CREATE TABLE achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    position_id UUID REFERENCES positions(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    situation TEXT,
    task TEXT,
    action TEXT,
    result TEXT,
    skills TEXT[] DEFAULT '{}',
    tags TEXT[] DEFAULT '{}',
    metrics JSONB DEFAULT '{}',
    skip_display BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_achievements_position ON achievements(position_id);
CREATE INDEX idx_achievements_tags ON achievements USING GIN(tags);
CREATE INDEX idx_achievements_skills ON achievements USING GIN(skills);

-- Skills
CREATE TABLE skills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID REFERENCES profiles(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    category TEXT,
    subcategory TEXT,
    level TEXT,
    verified BOOLEAN DEFAULT FALSE,
    verify_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE(profile_id, name)
);

CREATE INDEX idx_skills_profile ON skills(profile_id);
CREATE INDEX idx_skills_category ON skills(category);

-- Education
CREATE TABLE education (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID REFERENCES profiles(id) ON DELETE CASCADE,
    institution TEXT NOT NULL,
    degree TEXT NOT NULL,
    field TEXT,
    start_date DATE,
    end_date DATE,
    honors TEXT,
    gpa TEXT,
    display BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Certifications
CREATE TABLE certifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID REFERENCES profiles(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    issuer TEXT,
    issue_date DATE,
    expiration_date DATE,
    credential_id TEXT,
    credential_url TEXT,
    display BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Verifiable Credentials
CREATE TABLE verifiable_credentials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID REFERENCES profiles(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    username TEXT,
    profile_url TEXT,
    data JSONB DEFAULT '{}',
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Opportunities
CREATE TABLE opportunities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID REFERENCES profiles(id) ON DELETE CASCADE,
    company TEXT NOT NULL,
    position TEXT NOT NULL,
    location TEXT,
    remote BOOLEAN DEFAULT FALSE,
    url TEXT,
    job_desc_raw TEXT,
    job_desc_parsed JSONB,
    company_info JSONB,
    resume_id UUID,
    cover_letter_id UUID,
    evaluation_id UUID,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_opportunities_profile ON opportunities(profile_id);

-- Applications
CREATE TABLE applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    opportunity_id UUID REFERENCES opportunities(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'draft',
    applied_at TIMESTAMPTZ,
    resume_version TEXT,
    cover_letter_version TEXT,
    outcome JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_applications_opportunity ON applications(opportunity_id);
CREATE INDEX idx_applications_status ON applications(status);

-- Interviews
CREATE TABLE interviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID REFERENCES applications(id) ON DELETE CASCADE,
    round INT NOT NULL,
    type TEXT NOT NULL,
    scheduled_at TIMESTAMPTZ,
    duration_mins INT,
    interviewers TEXT[] DEFAULT '{}',
    questions JSONB DEFAULT '[]',
    self_assessment JSONB,
    feedback JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_interviews_application ON interviews(application_id);

-- Interview Prep Sets
CREATE TABLE interview_prep_sets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID REFERENCES profiles(id) ON DELETE CASCADE,
    opportunity_id UUID REFERENCES opportunities(id),
    title TEXT NOT NULL,
    description TEXT,
    sections JSONB NOT NULL DEFAULT '[]',
    pass_percentage INT DEFAULT 70,
    feedback_ranges JSONB DEFAULT '[]',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_prep_sets_profile ON interview_prep_sets(profile_id);
```

## 7. CLI Commands

```
sprofile - Career profile and resume management

Usage:
  sprofile [command]

Profile Commands:
  profile init          Initialize a new profile
  profile show          Display profile summary
  profile export        Export profile to JSON

Experience Commands:
  tenure list           List all tenures
  tenure add            Add a new tenure
  position list         List positions for a tenure
  position add          Add a new position
  achievement list      List achievements for a position
  achievement add       Add a new achievement (STAR)

Resume Commands:
  resume generate       Generate a tailored resume
  resume list           List generated resumes
  resume export         Export resume to format (md, pdf, docx)

Cover Letter Commands:
  cover generate        Generate a cover letter
  cover list            List cover letters
  cover export          Export cover letter

Evaluation Commands:
  eval resume           Evaluate resume against job description
  eval report           Show evaluation report
  eval history          Show evaluation history

Application Commands:
  apply create          Create a new application
  apply list            List applications
  apply status          Update application status
  apply outcome         Record application outcome

Interview Commands:
  interview add         Add interview to application
  interview feedback    Record interview feedback
  interview questions   List questions from interviews

Prep Commands:
  prep generate         Generate interview prep set
  prep export           Export prep set to H5P
  prep practice         Practice with prep questions

Config Commands:
  config init           Initialize configuration
  config show           Show current configuration
```

## 8. Security Considerations

1. **Data Encryption**: Sensitive fields (email, phone) encrypted at rest
2. **API Key Storage**: LLM API keys stored in environment variables or secure vault
3. **No PII in Logs**: Logging excludes personal information
4. **Access Control**: Database backend supports row-level security
5. **Export Sanitization**: Exported files exclude sensitive data by default

## 9. Testing Strategy

1. **Unit Tests**: All domain types and business logic
2. **Integration Tests**: Store implementations, LLM providers
3. **Golden Tests**: Markdown/PDF output comparison
4. **E2E Tests**: Full workflow from profile to exported resume

## 10. Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| github.com/agentplexus/multi-agent-spec | latest | Agent team definitions and orchestration |
| github.com/agentplexus/structured-evaluation | latest | LLM evaluation framework |
| github.com/grokify/h5p-go | latest | H5P quiz generation |
| github.com/grokify/gopandoc | latest | Markdown to PDF/DOCX |
| github.com/grokify/mogo | latest | Go utilities |
| github.com/spf13/cobra | v1.8+ | CLI framework |
| github.com/google/uuid | v1.6+ | UUID generation |
| github.com/jackc/pgx/v5 | v5.5+ | PostgreSQL driver |
| sigs.k8s.io/yaml | latest | YAML parsing |
