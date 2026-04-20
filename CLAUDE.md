# CLAUDE.md - structured-profile AI Agent Guide

Instructions for AI agents using structured-profile to generate resumes, cover letters, and match evaluations.

## Profile Data Format

Profile data uses **structured-profile JSON format** defined in `schema/profile.go`:

```json
{
  "profile": {
    "name": "John Doe",
    "email": "john@example.com",
    "links": [
      {"type": "linkedin", "url": "https://linkedin.com/in/johndoe"},
      {"type": "github", "url": "https://github.com/johndoe"}
    ],
    "summaries": {
      "default": "Platform leader with 15+ years...",
      "by_domain": {
        "devx": "Developer experience leader...",
        "iam": "Identity and access management expert..."
      }
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
              "name": "OAuth 2.0 Token Exchange",
              "description": "Launched RFC 8693 implementation...",
              "skills": ["oauth2", "oidc", "api-security"],
              "metrics": {"adoption_rate": "40%"}
            }
          ]
        }
      ]
    }
  ],
  "skills": [...],
  "education": [...],
  "certifications": [...]
}
```

## AI Agent Pipeline

The resume/cover letter generation uses a 6-step agent pipeline:

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

### Step 1: Analyze Job Description (jd-analyst)

Extract structured requirements with semantic understanding:

```
Analyze the job description at {jd_path}.
Extract: role, level, required/preferred skills, experience years,
semantic mappings (e.g., "technical enablement" → "developer advocacy"),
company context, and compensation.
```

### Step 2: Match Profile to JD (profile-matcher)

Score profile achievements against JD requirements:

```
Match profile achievements to JD requirements.
Score using: direct skill match (40pts), semantic match (25pts),
impact alignment (20pts), recency (10pts), quantification (5pts).
Identify gaps and recommend positioning.
```

### Step 3: Generate Resume (resume-generator)

Create tailored resume highlighting relevant experience:

```
Generate a tailored resume using:
- Profile: structured-profile JSON
- JD Analysis: extracted requirements
- Match Result: ranked achievements

Include: tailored summary, relevant achievements (STAR format),
skills matching JD terminology, education/certifications.
```

### Step 4: Evaluate Resume (resume-evaluator)

Score resume quality and job fit:

```
Evaluate resume for:
- JD alignment (keywords, requirements coverage)
- Impact clarity (quantified achievements)
- Professional presentation
- ATS compatibility

Return score 0-100 with actionable feedback.
```

### Step 5: Refine if Needed (resume-refiner)

If score < 85, improve based on feedback:

```
Refine resume addressing evaluation feedback.
Priority: high-impact issues first.
Maintain authenticity - reframe, don't fabricate.
```

### Step 6: Output

Write final resume to specified format (`.md`, `.pdf` via pandoc).

## Generating a Match Evaluation

Use the **matcheval schema** (`schema/matcheval.go`) to produce structured evaluations:

```
Evaluate profile-to-job match and output matcheval.json:

1. Score each category (0-10 scale):
   - technical_skills, domain_experience, leadership
   - years_experience, platform_scale, compliance
   - thought_leadership, apis_sdks, etc.

2. Record findings (strengths and gaps):
   - id: S001/G001 format
   - severity: critical/high/medium/low/info
   - evidence, recommendation, effort level

3. Compute decision:
   - Pass criteria: 0 critical, 0 high, score >= 70%
   - Status: pass/conditional/fail/human_review

4. List next steps for candidate
```

Output follows `schema/matcheval.schema.json`.

## Generating a Cover Letter

After profile matching, generate a cover letter:

```
Generate a cover letter using:
- Profile: structured-profile JSON
- JD Analysis: role, company, culture signals
- Match Result: top 3-5 achievements, positioning

Structure:
1. Opening: Express interest, mention specific role
2. Value proposition: 2-3 relevant achievements with impact
3. Company fit: Connect experience to company needs
4. Closing: Call to action

Tone: Professional, confident, authentic
Length: 3-4 paragraphs, under 400 words
```

## Dual Evaluation Approach

Run **two evaluations** for each application:

| Evaluation | Schema | Purpose | When |
|------------|--------|---------|------|
| **Profile Match** | `matcheval.json` | Does candidate's experience fit the role? | Before generating documents |
| **Document Quality** | `doceval.json` | Do documents effectively present the fit? | After generating documents |

### Why Both Matter

**matcheval.json** answers:
- Should I apply at all?
- What strengths to highlight?
- What gaps need addressing?
- Realistic match score?

**doceval.json** answers:
- Does resume highlight the right achievements?
- Are JD keywords represented?
- Are gaps mitigated in cover letter?
- Is it ATS-compatible?
- Ready to submit?

## Generating a Document Evaluation

Use the **doceval schema** (`schema/doceval.go`) after generating resume/cover letter:

```
Evaluate document quality and output doceval.json:

1. Score each category (0-10 scale):
   - keyword_coverage, achievement_relevance, gap_mitigation
   - quantification, ats_compatibility, narrative_coherence
   - value_proposition (cover letter)

2. Record findings (strengths and gaps):
   - id: DS001 for doc strength, DG001 for doc gap
   - document: resume, cover_letter, both
   - Include currentText and suggestedText for fixes

3. Compute decision:
   - Status: excellent/good/needs_work/major_revision
   - readyToSubmit: boolean

4. List prioritized action items
```

Output follows `schema/doceval.schema.json`.

## Application Directory Structure

For each job application, create:

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

## Quick Reference Commands

```bash
# Analyze JD
sprofile jd analyze jobdescription.md

# Match profile to JD
sprofile match --profile {id} --jd jobdescription.md

# Generate resume
sprofile resume generate --profile {id} --jd jobdescription.md -o resume.md

# Generate cover letter
sprofile cover generate --profile {id} --jd jobdescription.md -o coverletter.md

# Export to PDF (via pandoc)
pandoc resume.md -o resume.pdf --pdf-engine=lualatex
```

## Match Category Reference

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

## Severity Levels

| Severity | Description | Impact |
|----------|-------------|--------|
| `critical` | Missing must-have | Blocks |
| `high` | Significant gap | Blocks |
| `medium` | Addressable with effort | Flag |
| `low` | Nice-to-have missing | Note |
| `info` | Informational (strengths) | Highlight |

## Pass Criteria Defaults

```go
PassCriteria{
    MaxCritical: 0,    // No critical gaps allowed
    MaxHigh:     0,    // No high gaps allowed
    MaxMedium:   -1,   // Unlimited medium (flag only)
    MinScore:    70.0, // 70% minimum match
}
```

For senior/executive roles, use stricter criteria:

```go
StrictPassCriteria{
    MaxCritical: 0,
    MaxHigh:     0,
    MaxMedium:   2,    // Max 2 medium gaps
    MinScore:    80.0, // 80% minimum match
}
```

## Document Quality Categories

Standard categories defined in `schema/doceval.go`:

| Category | Document | Description |
|----------|----------|-------------|
| `keyword_coverage` | Both | JD keywords represented in documents |
| `achievement_relevance` | Resume | Top achievements align with JD |
| `gap_mitigation` | Cover Letter | Known gaps addressed |
| `quantification` | Resume | Achievements include metrics |
| `ats_compatibility` | Resume | Format compatible with ATS |
| `narrative_coherence` | Both | Documents tell coherent story |
| `summary_impact` | Resume | Summary hooks the reader |
| `value_proposition` | Cover Letter | Clear unique value stated |
| `company_fit` | Cover Letter | Connection to company needs |

## Document Quality Criteria

```go
DocEvalCriteria{
    MaxCritical:       0,     // No critical issues
    MaxHigh:           0,     // No high issues
    MaxMedium:         2,     // Max 2 medium issues
    MinScore:          75.0,  // 75% minimum quality
    MinKeywordCoverage: 70.0, // 70% JD keywords
}
```

## Document Decision Status

| Status | Description | Action |
|--------|-------------|--------|
| `excellent` | Ready to submit | Submit as-is |
| `good` | Minor improvements optional | Submit or polish |
| `needs_work` | Should address issues | Revise before submitting |
| `major_revision` | Significant rewrite needed | Major edits required |
