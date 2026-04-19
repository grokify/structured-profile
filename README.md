# structured-profile

A Go library for managing structured professional profiles and generating tailored resumes and cover letters using AI-powered semantic matching.

## Features

- **Structured Profile Schema** - JSON-based career data format for work history, education, skills, and achievements
- **Job Description Parser** - Extract structured requirements from markdown job descriptions
- **Profile Matching** - Score how well a profile matches a job description
- **Resume Generation** - Create tailored resumes highlighting relevant experience
- **Cover Letter Generation** - Generate personalized cover letters
- **Multi-Format Export** - Output to Markdown, PDF, and DOCX via Pandoc
- **AI Agent Pipeline** - Semantic matching using LLM-based agents

## Installation

```bash
go install github.com/grokify/structured-profile/cmd/sprofile@latest
```

## Quick Start

### 1. Create a Profile

Store your profile as JSON in `~/.sprofile/profiles/`:

```json
{
  "id": "john-doe",
  "personal": {
    "name": "John Doe",
    "email": "john@example.com"
  },
  "experience": [
    {
      "company": "Acme Corp",
      "title": "Senior Engineer",
      "start_date": "2020-01",
      "achievements": [
        "Led team of 5 engineers building microservices platform",
        "Reduced deployment time by 80% through CI/CD automation"
      ]
    }
  ]
}
```

### 2. Generate a Resume

```bash
# Generate tailored resume from profile + job description
sprofile resume generate \
  --profile john-doe \
  --jd path/to/job-description.md \
  --output resume.pdf
```

### 3. Generate a Cover Letter

```bash
sprofile cover generate \
  --profile john-doe \
  --jd path/to/job-description.md \
  --output coverletter.pdf
```

## CLI Commands

```
sprofile resume generate   Generate a tailored resume
sprofile cover generate    Generate a cover letter
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
JD Analysis → Profile Matching → Resume Generation → Evaluation → Refinement
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

## Project Structure

```
structured-profile/
├── cmd/sprofile/      # CLI entry point
├── schema/            # Profile JSON schema types
├── jdparser/          # Job description parser
├── matcher/           # Profile-to-JD matching
├── service/           # Resume & cover letter services
├── store/             # Profile storage
├── export/            # PDF/DOCX export via Pandoc
├── migrate/           # Data migration utilities
├── agents/specs/      # AI agent specifications
└── testdata/          # Sample profiles and JDs
```

## Requirements

- Go 1.21+
- Pandoc (for PDF/DOCX export)
- LuaLaTeX (for PDF with custom fonts)

## License

MIT
