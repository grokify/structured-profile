# Dual Evaluation

Run **two evaluations** for each job application to ensure both fit and presentation quality.

## Why Two Evaluations?

| Evaluation | Purpose | Question Answered |
|------------|---------|-------------------|
| **Profile Match** (`matcheval.json`) | Assess raw fit | Should I apply? |
| **Document Quality** (`doceval.json`) | Assess presentation | Are my documents effective? |

A candidate might be a 90% fit, but their resume might only show 70% of that fit. Or the resume might overclaim and not be supported by the actual profile.

## Profile Match Evaluation

### What It Measures

| Category | Description |
|----------|-------------|
| `technical_skills` | Technical expertise match |
| `domain_experience` | Industry knowledge |
| `leadership` | Team/people management |
| `years_experience` | Total relevant years |
| `platform_scale` | Scale of systems managed |
| `apis_sdks` | API/SDK experience |
| `identity_security` | IAM/auth expertise |
| `compliance` | Regulatory experience |
| `thought_leadership` | Industry presence |

### Match Types

| Type | Description | Example |
|------|-------------|---------|
| `exact` | Specific skill required | "Go experience" |
| `semantic` | Equivalent skills accepted | "developer advocacy" ≈ "technical enablement" |
| `partial` | Related but not identical | "Python" for "JavaScript" role |
| `threshold` | Meets minimum level | "10+ years experience" |
| `boolean` | Present or absent | "Has patents" |
| `weighted` | Combined sub-factors | Leadership score from multiple inputs |

### Pass Criteria

```go
PassCriteria{
    MaxCritical: 0,    // No critical gaps allowed
    MaxHigh:     0,    // No high gaps allowed
    MaxMedium:   -1,   // Unlimited medium (flag only)
    MinScore:    70.0, // 70% minimum match
}
```

### Decision Status

| Status | Meaning | Action |
|--------|---------|--------|
| `pass` | Strong match | Proceed with application |
| `conditional` | Good match with caveats | Address gaps in cover letter |
| `fail` | Has blocking gaps | Do not apply |
| `human_review` | Borderline | Manual review needed |

## Document Quality Evaluation

### What It Measures

| Category | Document | Weight |
|----------|----------|--------|
| `keyword_coverage` | Both | 20% |
| `achievement_relevance` | Resume | 20% |
| `gap_mitigation` | Cover Letter | 15% |
| `quantification` | Resume | 15% |
| `ats_compatibility` | Resume | 10% |
| `narrative_coherence` | Both | 10% |
| `value_proposition` | Cover Letter | 10% |

### Evaluation Criteria

```go
DocEvalCriteria{
    MaxCritical:       0,     // No critical issues
    MaxHigh:           0,     // No high issues
    MaxMedium:         2,     // Max 2 medium issues
    MinScore:          75.0,  // 75% minimum quality
    MinKeywordCoverage: 70.0, // 70% JD keywords
}
```

### Decision Status

| Status | Meaning | Action |
|--------|---------|--------|
| `excellent` | Ready to submit | Submit as-is |
| `good` | Minor improvements optional | Submit or polish |
| `needs_work` | Should address issues | Revise before submitting |
| `major_revision` | Significant rewrite needed | Major edits required |

## Example Workflow

```bash
# Step 1: Evaluate profile fit
sprofile match --profile john-doe --jd jobdescription.md
# Output: matcheval.json (92% match, PASS)

# Step 2: Generate documents
sprofile resume generate --profile john-doe --jd jobdescription.md -o resume.md
sprofile cover generate --profile john-doe --jd jobdescription.md -o coverletter.md

# Step 3: Evaluate document quality
sprofile doceval --resume resume.md --cover coverletter.md --jd jobdescription.md
# Output: doceval.json (88% quality, GOOD)

# Step 4: Export (if quality is acceptable)
pandoc resume.md -o resume.pdf --pdf-engine=lualatex
```

## Severity Levels

Both evaluations use the same severity scale:

| Severity | Description | Impact |
|----------|-------------|--------|
| `critical` | Must-have missing | Blocks |
| `high` | Significant issue | Blocks |
| `medium` | Addressable | Flag for attention |
| `low` | Nice-to-have | Note only |
| `info` | Strength/FYI | Highlight |

## Finding ID Format

| Evaluation | Strengths | Gaps |
|------------|-----------|------|
| matcheval | S001, S002, ... | G001, G002, ... |
| doceval | DS001, DS002, ... | DG001, DG002, ... |
