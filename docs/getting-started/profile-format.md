# Profile Format

The structured-profile format uses JSON to represent career data with STAR-format achievements.

## Full Profile Structure

```json
{
  "profile": {
    "id": "unique-id",
    "name": "Full Name",
    "email": "email@example.com",
    "phone": "+1-555-123-4567",
    "location": "San Francisco, CA",
    "links": [
      {"type": "linkedin", "url": "https://linkedin.com/in/..."},
      {"type": "github", "url": "https://github.com/..."},
      {"type": "website", "url": "https://..."}
    ],
    "summaries": {
      "default": "Default professional summary...",
      "by_domain": {
        "devx": "Developer experience focused summary...",
        "iam": "Identity and security focused summary...",
        "platform": "Platform engineering focused summary..."
      }
    }
  },
  "tenures": [...],
  "skills": [...],
  "education": [...],
  "certifications": [...],
  "publications": [...]
}
```

## Tenures and Positions

A tenure represents time at a company, containing one or more positions:

```json
{
  "tenures": [
    {
      "id": "uuid",
      "company": "Acme Corp",
      "start_date": "2020-01",
      "end_date": "2024-06",
      "positions": [
        {
          "id": "uuid",
          "title": "VP Platform",
          "description": "Lead platform product team",
          "start_date": "2023-01",
          "end_date": "2024-06",
          "achievements": [...]
        },
        {
          "id": "uuid",
          "title": "Sr. Director Platform",
          "start_date": "2020-01",
          "end_date": "2022-12",
          "achievements": [...]
        }
      ]
    }
  ]
}
```

## Achievements (STAR Format)

Achievements use the STAR format for structured impact:

```json
{
  "achievements": [
    {
      "id": "uuid",
      "name": "api-platform-launch",
      "situation": "Company needed to expand developer ecosystem",
      "task": "to drive platform adoption and revenue",
      "action": "by designing and launching REST APIs with OAuth 2.0, building 10+ SDKs, and growing developer advocacy team",
      "result": "resulting in 70K+ developers, $500M ARR, and 7 industry awards",
      "skills": ["api-design", "oauth2", "sdk-development", "leadership"],
      "tags": ["platform", "devx", "growth"],
      "metrics": {
        "developers": 70000,
        "arr": 500000000,
        "awards": 7
      }
    }
  ]
}
```

### STAR Fields

| Field | Description | Example |
|-------|-------------|---------|
| `situation` | Context/challenge | "Company needed developer ecosystem" |
| `task` | Goal/objective | "to drive platform adoption" |
| `action` | What you did | "by launching REST APIs with OAuth 2.0" |
| `result` | Measurable outcome | "resulting in 70K+ developers and $500M ARR" |

## Skills

```json
{
  "skills": [
    {
      "name": "Go",
      "level": "expert",
      "years": 8,
      "tags": ["programming", "backend"]
    },
    {
      "name": "OAuth 2.0",
      "level": "expert",
      "years": 10,
      "tags": ["security", "iam"]
    }
  ]
}
```

### Skill Levels

| Level | Description |
|-------|-------------|
| `beginner` | Learning, limited practical experience |
| `intermediate` | Comfortable, can work independently |
| `advanced` | Strong proficiency, can mentor others |
| `expert` | Deep expertise, industry recognition |

## Education

```json
{
  "education": [
    {
      "institution": "Stanford University",
      "degree": "MS",
      "field": "Computer Science",
      "year": 2010,
      "honors": "Dean's List"
    }
  ]
}
```

## Certifications

```json
{
  "certifications": [
    {
      "name": "API Security Architect",
      "issuer": "API Academy",
      "date": "2023-06",
      "url": "https://..."
    }
  ]
}
```

## Domain-Specific Summaries

Use `summaries.by_domain` to store role-specific summaries:

```json
{
  "summaries": {
    "default": "Technical product leader...",
    "by_domain": {
      "devx": "Award-winning developer platform executive with 20+ years building developer ecosystems...",
      "iam": "Identity and authentication platform executive with deep expertise in OAuth 2.0, OIDC...",
      "platform": "Platform engineering leader with experience scaling APIs to 1M+ calls/day..."
    }
  }
}
```

The resume generator selects the appropriate summary based on job requirements.
