---
name: generate-resume
description: Generate a tailored resume using the full AI pipeline
model: sonnet
tools: [Read, Write, Bash]
invocation: /generate-resume --profile <id> --jd <file> [--domain <domain>] [--output <file>]
---

# Generate Tailored Resume

Generate a tailored resume by running the full AI-powered pipeline:
1. Analyze job description semantically
2. Match profile achievements to requirements
3. Generate tailored resume
4. Evaluate quality
5. Refine if needed

## Usage

```
/generate-resume --profile john-wang --jd job.md --output resume.pdf
/generate-resume --profile john-wang --jd job.md --domain devx
```

## Parameters

- `--profile`: Profile ID (required)
- `--jd`: Path to job description markdown file (required)
- `--domain`: Domain filter - devx, iam, platform (optional)
- `--output`: Output file path with extension for format (optional, defaults to stdout markdown)

## Process

### Step 1: Analyze JD
Use jd-analyst agent to extract requirements with semantic understanding.

### Step 2: Match Profile
Use profile-matcher agent to:
- Map profile skills to JD requirements (exact + semantic)
- Rank achievements by relevance
- Identify gaps and positioning strategy

### Step 3: Generate Resume
Use resume-generator agent to create tailored resume:
- Craft targeted summary
- Select and order relevant achievements
- Mirror JD terminology where authentic

### Step 4: Evaluate
Use resume-evaluator agent to score:
- Job fit (40 points)
- Achievement quality (25 points)
- Summary effectiveness (15 points)
- Structure (10 points)
- Language (10 points)

### Step 5: Refine (if score < 85)
Use resume-refiner agent to address feedback.

## Output

- **Markdown**: To stdout or file
- **PDF/DOCX**: Via Pandoc export

## Example

```bash
# Generate for Anthropic role
/generate-resume --profile john-wang \
  --jd ~/jobs/anthropic-technical-enablement.md \
  --domain devx \
  --output resume-anthropic.pdf
```

This will produce a resume optimized for the Technical Enablement Lead role, highlighting developer advocacy, AI tools experience, and training capabilities.
