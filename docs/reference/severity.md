# Severity Levels

Severity levels used in both matcheval and doceval findings.

## Severity Reference

| Severity | Description | Impact | Action |
|----------|-------------|--------|--------|
| `critical` | Must-have missing | Blocks | Do not apply |
| `high` | Significant gap | Blocks | Address before applying |
| `medium` | Addressable with effort | Flag | Address in cover letter |
| `low` | Nice-to-have missing | Note | Optional to address |
| `info` | Informational (often strengths) | Highlight | Emphasize in interview |

## Severity by Context

### Profile Match (matcheval)

| Severity | Example Gap |
|----------|-------------|
| `critical` | Missing required certification (e.g., medical license) |
| `high` | 5 years experience when 10+ required |
| `medium` | No direct industry experience but transferable skills |
| `low` | Missing "nice to have" skill |
| `info` | Strength: exceeds requirement |

### Document Quality (doceval)

| Severity | Example Issue |
|----------|---------------|
| `critical` | Wrong company name in cover letter |
| `high` | Missing entire required section |
| `medium` | Weak quantification in achievements |
| `low` | Formatting inconsistency |
| `info` | Strength: excellent keyword coverage |

## Decision Logic

### matcheval Decision

```go
if counts.Critical > criteria.MaxCritical {
    status = "fail"
} else if counts.High > criteria.MaxHigh {
    status = "fail"
} else if score < criteria.MinScore {
    status = "human_review"
} else if counts.Medium > criteria.MaxMedium {
    status = "conditional"
} else {
    status = "pass"
}
```

### doceval Decision

```go
if counts.Critical > criteria.MaxCritical {
    status = "major_revision"
} else if counts.High > criteria.MaxHigh {
    status = "needs_work"
} else if score < criteria.MinScore {
    status = "needs_work"
} else if counts.Medium > criteria.MaxMedium {
    status = "good"
} else {
    status = "excellent"
}
```

## Effort Levels

For addressable gaps, effort level indicates remediation difficulty:

| Effort | Description | Example |
|--------|-------------|---------|
| `low` | Quick fix | Reword a bullet point |
| `medium` | Moderate work | Research and address a skill gap |
| `high` | Significant investment | Obtain certification, gain experience |

## Finding Counts

Both evaluations track finding counts:

```go
type FindingCounts struct {
    Critical int `json:"critical"`
    High     int `json:"high"`
    Medium   int `json:"medium"`
    Low      int `json:"low"`
    Info     int `json:"info"`
    Total    int `json:"total"`
}
```

## Pass Criteria Defaults

### matcheval

```go
PassCriteria{
    MaxCritical: 0,    // None allowed
    MaxHigh:     0,    // None allowed
    MaxMedium:   -1,   // Unlimited (flag only)
    MinScore:    70.0, // 70% minimum
}
```

### doceval

```go
DocEvalCriteria{
    MaxCritical:        0,    // None allowed
    MaxHigh:            0,    // None allowed
    MaxMedium:          2,    // Max 2
    MinScore:           75.0, // 75% minimum
    MinKeywordCoverage: 70.0, // 70% keywords
}
```

## Strict Criteria (Senior Roles)

For senior/executive roles, use stricter thresholds:

```go
StrictPassCriteria{
    MaxCritical: 0,
    MaxHigh:     0,
    MaxMedium:   2,
    MinScore:    80.0,
}

StrictDocEvalCriteria{
    MaxCritical:        0,
    MaxHigh:            0,
    MaxMedium:          1,
    MinScore:           85.0,
    MinKeywordCoverage: 80.0,
}
```
