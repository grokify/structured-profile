---
name: parse-jd
description: Parse a job description with semantic understanding
model: haiku
tools: [Read]
invocation: /parse-jd <jd-file>
---

# Parse Job Description

Analyze a job description file and extract structured requirements with semantic understanding.

## Usage

```
/parse-jd path/to/job-description.md
```

## Process

1. Read the job description file
2. Extract role information (title, level, type)
3. Identify required and preferred skills
4. Build semantic mappings (e.g., "technical enablement" → related terms)
5. Extract context (company, location, compensation)

## Output

Return a structured analysis:

```json
{
  "role": {
    "title": "Technical Enablement Lead",
    "level": "lead",
    "type": "ic"
  },
  "requirements": {
    "required_skills": ["programming", "technical training", "AI coding tools"],
    "preferred_skills": ["demo coaching", "enterprise workflows"],
    "experience_years": 7
  },
  "semantic_mappings": {
    "technical enablement": ["developer advocacy", "technical training", "demos"],
    "AI coding tools": ["Claude Code", "GitHub Copilot", "Cursor"]
  },
  "context": {
    "company": "Anthropic",
    "remote": true,
    "compensation": {"min": 270000, "max": 310000}
  }
}
```

This skill is useful for quickly understanding what a job is looking for before generating a tailored resume.
