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

## Quick Start

```bash
# Install
go install github.com/grokify/structured-profile/cmd/sprofile@latest

# Generate resume
sprofile resume generate --profile john-doe --jd jobdescription.md -o resume.md

# Generate cover letter
sprofile cover generate --profile john-doe --jd jobdescription.md -o coverletter.md

# Export to PDF
pandoc resume.md -o resume.pdf --pdf-engine=lualatex
```

## Application Directory Structure

```
applications/app_{date}_{company}_{role}/
├── jobdescription.md    # Original JD
├── jdanalysis.json      # Output from jd-analyst
├── matcheval.json       # Profile-to-JD match evaluation
├── resume.md            # Tailored resume (markdown)
├── resume.pdf           # PDF export
├── coverletter.md       # Cover letter (markdown)
├── coverletter.pdf      # PDF export
└── doceval.json         # Document quality evaluation
```

## Getting Started

- [Installation](getting-started/installation.md)
- [Quick Start](getting-started/quickstart.md)
- [Profile Format](getting-started/profile-format.md)

## Learn More

- [Pipeline Overview](pipeline/overview.md)
- [Dual Evaluation](pipeline/dual-evaluation.md)
- [AI Agents](pipeline/agents.md)
- [Schema Reference](schema/profile.md)
