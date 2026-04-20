# Quick Start

## 1. Create Your Profile

Create a profile JSON file at `~/.sprofile/profiles/your-name.json`:

```json
{
  "profile": {
    "id": "john-doe",
    "name": "John Doe",
    "email": "john@example.com",
    "links": [
      {"type": "linkedin", "url": "https://linkedin.com/in/johndoe"},
      {"type": "github", "url": "https://github.com/johndoe"}
    ],
    "summaries": {
      "default": "Platform leader with 15+ years experience building developer ecosystems..."
    }
  },
  "tenures": [
    {
      "company": "Acme Corp",
      "start_date": "2020-01",
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
  ],
  "skills": [
    {"name": "Go", "level": "expert", "years": 8},
    {"name": "API Design", "level": "expert", "years": 15}
  ],
  "education": [
    {
      "institution": "Stanford University",
      "degree": "MS Computer Science",
      "year": 2010
    }
  ]
}
```

## 2. Prepare a Job Description

Save the job description as `jobdescription.md`:

```markdown
# VP of Platform Engineering

## About the Role

We are seeking a VP of Platform Engineering to lead our platform team...

## Requirements

- 10+ years of software engineering experience
- 5+ years in engineering leadership roles
- Deep expertise in API design and developer platforms
- Experience with Kubernetes and microservices
```

## 3. Run the Pipeline

```bash
# Match profile to JD (generates matcheval.json)
sprofile match --profile john-doe --jd jobdescription.md

# Generate tailored resume
sprofile resume generate --profile john-doe --jd jobdescription.md -o resume.md

# Generate cover letter
sprofile cover generate --profile john-doe --jd jobdescription.md -o coverletter.md

# Export to PDF
pandoc resume.md -o resume.pdf --pdf-engine=lualatex
pandoc coverletter.md -o coverletter.pdf --pdf-engine=lualatex
```

## 4. Review Results

Check your generated files:

- `matcheval.json` - Profile-to-JD match score and findings
- `resume.md` / `resume.pdf` - Tailored resume
- `coverletter.md` / `coverletter.pdf` - Cover letter
- `doceval.json` - Document quality assessment

## Next Steps

- [Profile Format](profile-format.md) - Learn the full profile schema
- [Dual Evaluation](../pipeline/dual-evaluation.md) - Understand the evaluation system
- [AI Agents](../pipeline/agents.md) - Use semantic matching
