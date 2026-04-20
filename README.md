# structured-profile

A Go library for managing structured professional profiles and generating tailored resumes and cover letters using AI-powered semantic matching.

## Features

- **Structured Profile Schema** - JSON-based career data format for work history, education, skills, and achievements
- **Job Description Parser** - Extract structured requirements from markdown job descriptions
- **Profile Matching** - Score how well a profile matches a job description
- **Resume Generation** - Create tailored resumes highlighting relevant experience
- **Cover Letter Generation** - Generate personalized cover letters
- **Dual Evaluation System** - Profile-to-JD matching + document quality assessment
- **Multi-Format Export** - Output to Markdown, PDF, and DOCX via Pandoc
- **AI Agent Pipeline** - Semantic matching using LLM-based agents

## Job Application Pipeline

```
┌─────────────────────────────────────────────────────────────────┐
│                    JOB APPLICATION PIPELINE                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1. JD Analysis     ──►  2. Profile Match  ──►  matcheval.json  │
│                              (fit check)                         │
│                                  │                               │
│                                  ▼                               │
│                          3. Generate Resume                      │
│                          4. Generate Cover Letter                │
│                                  │                               │
│                                  ▼                               │
│                          5. Doc Quality    ──►  doceval.json    │
│                              (QA check)                          │
│                                  │                               │
│                                  ▼                               │
│                          6. Refine if needed                     │
│                          7. Export PDFs                          │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

## Dual Evaluation Approach

Run **two evaluations** for each application:

| Evaluation | Schema | Purpose | When |
|------------|--------|---------|------|
| **Profile Match** | `matcheval.json` | Does candidate's experience fit the role? | Before generating documents |
| **Document Quality** | `doceval.json` | Do documents effectively present the fit? | After generating documents |

### Profile Match Categories

Standard categories defined in `schema/matcheval.go`:

| Category | Description | Match Type |
|----------|-------------|------------|
| `technical_skills` | Technical expertise | semantic |
| `domain_experience` | Industry knowledge | semantic |
| `leadership` | Team/people management | threshold |
| `years_experience` | Total relevant years | threshold |
| `platform_scale` | Scale of systems managed | threshold |
| `apis_sdks` | API/SDK experience | exact |
| `identity_security` | IAM/auth expertise | semantic |
| `compliance` | Regulatory experience | semantic |
| `thought_leadership` | Industry presence | boolean |
| `developer_relations` | DevRel/advocacy | semantic |

### Document Quality Categories

Standard categories defined in `schema/doceval.go`:

| Category | Document | Weight | Description |
|----------|----------|--------|-------------|
| `keyword_coverage` | Both | 20% | JD keywords represented |
| `achievement_relevance` | Resume | 20% | Top achievements align with JD |
| `gap_mitigation` | Cover Letter | 15% | Known gaps addressed |
| `quantification` | Resume | 15% | Achievements include metrics |
| `ats_compatibility` | Resume | 10% | Format compatible with ATS |
| `narrative_coherence` | Both | 10% | Documents tell coherent story |
| `value_proposition` | Cover Letter | 10% | Clear unique value stated |

### Severity Levels

| Severity | Description | Impact |
|----------|-------------|--------|
| `critical` | Missing must-have | Blocks |
| `high` | Significant gap | Blocks |
| `medium` | Addressable with effort | Flag |
| `low` | Nice-to-have missing | Note |
| `info` | Informational (strengths) | Highlight |

## Installation

```bash
go install github.com/grokify/structured-profile/cmd/sprofile@latest
```

## Quick Start

### 1. Create a Profile

Store your profile as JSON in `~/.sprofile/profiles/`:

```json
{
  "profile": {
    "id": "john-doe",
    "name": "John Doe",
    "email": "john@example.com",
    "summaries": {
      "default": "Platform leader with 15+ years experience..."
    }
  },
  "tenures": [
    {
      "company": "Acme Corp",
      "positions": [
        {
          "title": "VP Platform",
          "start_date": "2023-01",
          "achievements": [
            {
              "name": "API Platform Launch",
              "situation": "Company needed developer ecosystem",
              "task": "to drive platform adoption",
              "action": "by launching REST APIs with OAuth 2.0",
              "result": "resulting in 70K+ developers and $500M ARR",
              "skills": ["api-design", "oauth2", "leadership"]
            }
          ]
        }
      ]
    }
  ]
}
```

### 2. Generate Application Materials

```bash
# Analyze job description
sprofile jd analyze jobdescription.md

# Match profile to JD
sprofile match --profile john-doe --jd jobdescription.md

# Generate tailored resume
sprofile resume generate \
  --profile john-doe \
  --jd jobdescription.md \
  --output resume.md

# Generate cover letter
sprofile cover generate \
  --profile john-doe \
  --jd jobdescription.md \
  --output coverletter.md

# Export to PDF (via pandoc)
pandoc resume.md -o resume.pdf --pdf-engine=lualatex
```

### 3. Application Directory Structure

For each job application:

```
applications/app_{date}_{company}_{role}/
├── jobdescription.md    # Original JD
├── jdanalysis.json      # Output from jd-analyst
├── matcheval.json       # Profile-to-JD match evaluation
├── resume.md            # Tailored resume (markdown)
├── resume.pdf           # PDF export (via pandoc)
├── coverletter.md       # Cover letter (markdown)
├── coverletter.pdf      # PDF export (via pandoc)
└── doceval.json         # Document quality evaluation
```

## CLI Commands

```
sprofile jd analyze       Analyze a job description
sprofile match            Match profile to job description
sprofile resume generate  Generate a tailored resume
sprofile cover generate   Generate a cover letter
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--profile` | Profile ID | (required) |
| `--profile-dir` | Profile directory | `~/.sprofile/profiles/` |
| `--jd` | Job description markdown file | (required) |
| `--output` | Output file path | stdout |

Output format is determined by file extension: `.md`, `.pdf`, `.docx`

## Library Usage

```go
import (
    "github.com/grokify/structured-profile/schema"
    "github.com/grokify/structured-profile/service"
    "github.com/grokify/structured-profile/jdparser"
    "github.com/grokify/structured-profile/export"
)

// Load profile
profile, _ := schema.LoadProfile("path/to/profile.json")

// Parse job description
jd, _ := jdparser.ParseFile("path/to/jd.md")

// Generate resume
svc := service.NewResumeService()
resume, _ := svc.Generate(ctx, profile, jd)

// Export to PDF
export.ExportResumeToFile(resume, "resume.pdf", nil)
```

## AI Agent Pipeline

For semantic matching beyond keyword-based scoring, use the AI agent pipeline:

```
┌─────────────┐    ┌─────────────────┐    ┌──────────────────┐
│ jd-analyst  │───▶│ profile-matcher │───▶│ resume-generator │
└─────────────┘    └─────────────────┘    └──────────────────┘
                                                   │
                                                   ▼
┌─────────────────┐    ┌──────────────────┐    ┌──────────────────┐
│ resume-refiner  │◀───│ resume-evaluator │◀───│ (resume draft)   │
└─────────────────┘    └──────────────────┘    └──────────────────┘
```

Agent specifications are in `agents/specs/`:

| Agent | Purpose |
|-------|---------|
| `jd-analyst` | Semantic job description parsing |
| `profile-matcher` | Profile-to-JD matching with scoring |
| `resume-generator` | Tailored resume creation |
| `resume-evaluator` | LLM-as-Judge quality scoring |
| `resume-refiner` | Feedback-based improvement |
| `resume-coordinator` | Pipeline orchestration |

## Schema Reference

### Match Evaluation (`matcheval.json`)

```go
// Pass criteria defaults
PassCriteria{
    MaxCritical: 0,    // No critical gaps allowed
    MaxHigh:     0,    // No high gaps allowed
    MaxMedium:   -1,   // Unlimited medium (flag only)
    MinScore:    70.0, // 70% minimum match
}
```

### Document Evaluation (`doceval.json`)

```go
// Document quality criteria
DocEvalCriteria{
    MaxCritical:       0,     // No critical issues
    MaxHigh:           0,     // No high issues
    MaxMedium:         2,     // Max 2 medium issues
    MinScore:          75.0,  // 75% minimum quality
    MinKeywordCoverage: 70.0, // 70% JD keywords
}
```

### Decision Status Reference

| matcheval Status | Description |
|------------------|-------------|
| `pass` | Strong match, proceed |
| `conditional` | Good match with caveats |
| `fail` | Has blocking gaps |
| `human_review` | Borderline, needs review |

| doceval Status | Description |
|----------------|-------------|
| `excellent` | Ready to submit |
| `good` | Minor improvements optional |
| `needs_work` | Should address issues |
| `major_revision` | Significant rewrite needed |

## Project Structure

```
structured-profile/
├── cmd/sprofile/      # CLI entry point
├── schema/            # Profile and evaluation schemas
│   ├── profile.go     # Profile types
│   ├── matcheval.go   # Profile-to-JD match evaluation
│   └── doceval.go     # Document quality evaluation
├── jdparser/          # Job description parser
├── matcher/           # Profile-to-JD matching
├── service/           # Resume & cover letter services
├── store/             # Profile storage
├── export/            # PDF/DOCX export via Pandoc
├── migrate/           # Data migration utilities
├── agents/specs/      # AI agent specifications
├── docs/              # MkDocs documentation
└── testdata/          # Sample JDs for testing
```

## Documentation

Full documentation available at: [structured-profile docs](docs/)

- [CLAUDE.md](CLAUDE.md) - AI agent guide for resume generation

## Requirements

- Go 1.21+
- Pandoc (for PDF/DOCX export)
- LuaLaTeX (for PDF with custom fonts)

## License

MIT
