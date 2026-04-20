# Profile Schema

The profile schema defines the structure for career data in `schema/profile.go`.

## FullProfile

The top-level structure containing all profile data:

```go
type FullProfile struct {
    Profile              Profile                `json:"profile"`
    Tenures              []Tenure               `json:"tenures,omitempty"`
    Skills               []Skill                `json:"skills,omitempty"`
    Education            []Education            `json:"education,omitempty"`
    Certifications       []Certification        `json:"certifications,omitempty"`
    Credentials          []VerifiableCredential `json:"credentials,omitempty"`
    Publications         []Publication          `json:"publications,omitempty"`
    Opportunities        []Opportunity          `json:"opportunities,omitempty"`
    Applications         []Application          `json:"applications,omitempty"`
    CoverLetters         []CoverLetter          `json:"cover_letters,omitempty"`
    CoverLetterTemplates []CoverLetterTemplate  `json:"cover_letter_templates,omitempty"`
    PrepSets             []InterviewPrepSet     `json:"prep_sets,omitempty"`
}
```

## Profile

Core identity and contact information:

```go
type Profile struct {
    BaseEntity

    Name     string   `json:"name"`
    Email    string   `json:"email,omitempty"`
    Phone    string   `json:"phone,omitempty"`
    Location string   `json:"location,omitempty"`
    Links    []Link   `json:"links,omitempty"`
    Summaries Summaries `json:"summaries,omitempty"`
}
```

### Links

```go
type Link struct {
    Type string `json:"type"`  // "linkedin", "github", "website", etc.
    URL  string `json:"url"`
    Text string `json:"text,omitempty"`
}
```

### Summaries

Domain-specific professional summaries:

```go
type Summaries struct {
    Default  string            `json:"default,omitempty"`
    ByDomain map[string]string `json:"by_domain,omitempty"`
}
```

## Tenure

Represents time at a company:

```go
type Tenure struct {
    BaseEntity

    Company   string     `json:"company"`
    StartDate string     `json:"start_date"`
    EndDate   string     `json:"end_date,omitempty"`
    Positions []Position `json:"positions,omitempty"`
}
```

## Position

A role held within a tenure:

```go
type Position struct {
    BaseEntity

    TenureID    string        `json:"tenure_id"`
    Title       string        `json:"title"`
    Description string        `json:"description,omitempty"`
    StartDate   string        `json:"start_date"`
    EndDate     string        `json:"end_date,omitempty"`
    Achievements []Achievement `json:"achievements,omitempty"`
}
```

## Achievement

STAR-format accomplishment:

```go
type Achievement struct {
    BaseEntity

    PositionID  string            `json:"position_id"`
    Name        string            `json:"name"`
    Situation   string            `json:"situation"`
    Task        string            `json:"task"`
    Action      string            `json:"action"`
    Result      string            `json:"result"`
    Skills      []string          `json:"skills,omitempty"`
    Tags        []string          `json:"tags,omitempty"`
    Metrics     map[string]any    `json:"metrics,omitempty"`
}
```

### STAR Format

| Field | Purpose | Example |
|-------|---------|---------|
| `situation` | Context/challenge | "Company needed developer ecosystem" |
| `task` | Goal/objective | "to drive platform adoption" |
| `action` | What you did | "by launching REST APIs with OAuth 2.0" |
| `result` | Measurable outcome | "resulting in 70K+ developers" |

## Skill

```go
type Skill struct {
    BaseEntity

    Name  string `json:"name"`
    Level string `json:"level"` // beginner, intermediate, advanced, expert
    Years int    `json:"years,omitempty"`
    Tags  []string `json:"tags,omitempty"`
}
```

## Education

```go
type Education struct {
    BaseEntity

    Institution string `json:"institution"`
    Degree      string `json:"degree"`
    Field       string `json:"field,omitempty"`
    Year        int    `json:"year,omitempty"`
    Honors      string `json:"honors,omitempty"`
}
```

## Certification

```go
type Certification struct {
    BaseEntity

    Name   string `json:"name"`
    Issuer string `json:"issuer"`
    Date   string `json:"date,omitempty"`
    URL    string `json:"url,omitempty"`
}
```

## BaseEntity

All entities include:

```go
type BaseEntity struct {
    ID        string    `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

## JSON Schema

The corresponding JSON Schema is at `schema/profile.schema.json`.

## Example

See [Profile Format](../getting-started/profile-format.md) for complete examples.
