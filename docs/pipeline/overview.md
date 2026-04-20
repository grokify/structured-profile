# Pipeline Overview

The structured-profile pipeline transforms a profile and job description into tailored application materials.

## Pipeline Flow

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

## Pipeline Steps

### 1. JD Analysis

Extract structured requirements from the job description:

- Role title and level (IC, Manager, Director, VP, etc.)
- Required and preferred skills
- Experience requirements
- Semantic mappings (e.g., "technical enablement" → "developer advocacy")
- Company context and culture signals

**Output:** `jdanalysis.json`

### 2. Profile Match

Score the profile against JD requirements:

- Evaluate each category (technical skills, leadership, domain experience, etc.)
- Identify strengths and gaps
- Calculate overall match score
- Determine pass/fail status

**Output:** `matcheval.json`

### 3. Generate Resume

Create a tailored resume:

- Select relevant achievements based on match results
- Use domain-specific summary if available
- Mirror JD terminology where authentic
- Order experiences by relevance

**Output:** `resume.md`

### 4. Generate Cover Letter

Create a compelling cover letter:

- Express interest in specific role
- Highlight top 2-3 relevant achievements
- Address known gaps if appropriate
- Connect experience to company needs

**Output:** `coverletter.md`

### 5. Document Quality Check

Evaluate the generated documents:

- Check JD keyword coverage
- Verify achievement relevance
- Assess gap mitigation
- Check ATS compatibility
- Evaluate narrative coherence

**Output:** `doceval.json`

### 6. Refine if Needed

If document quality score is below threshold:

- Address high-priority issues
- Improve keyword coverage
- Strengthen quantification
- Enhance narrative flow

### 7. Export PDFs

Generate final documents:

```bash
pandoc resume.md -o resume.pdf --pdf-engine=lualatex
pandoc coverletter.md -o coverletter.pdf --pdf-engine=lualatex
```

## Decision Points

### Should I Apply?

Check `matcheval.json`:

| Score | Decision |
|-------|----------|
| 80%+ | Strong fit, apply |
| 70-80% | Good fit, address gaps in cover letter |
| 60-70% | Borderline, consider if role is stretch |
| <60% | Weak fit, may not be worth applying |

### Are Documents Ready?

Check `doceval.json`:

| Status | Action |
|--------|--------|
| `excellent` | Submit as-is |
| `good` | Optional polish |
| `needs_work` | Revise before submitting |
| `major_revision` | Significant edits required |

## Output Files

```
applications/app_{date}_{company}_{role}/
├── jobdescription.md    # Original JD
├── jdanalysis.json      # Step 1 output
├── matcheval.json       # Step 2 output
├── resume.md            # Step 3 output
├── coverletter.md       # Step 4 output
├── doceval.json         # Step 5 output
├── resume.pdf           # Step 7 output
└── coverletter.pdf      # Step 7 output
```
