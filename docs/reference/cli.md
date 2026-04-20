# CLI Reference

The `sprofile` CLI provides commands for profile management and document generation.

## Installation

```bash
go install github.com/grokify/structured-profile/cmd/sprofile@latest
```

## Commands

### sprofile jd analyze

Analyze a job description and extract requirements.

```bash
sprofile jd analyze <jd-file> [flags]
```

**Arguments:**

| Argument | Description |
|----------|-------------|
| `jd-file` | Path to job description markdown file |

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `-o, --output` | Output file path | stdout |
| `--format` | Output format (json, yaml) | json |

**Example:**

```bash
sprofile jd analyze jobdescription.md -o jdanalysis.json
```

### sprofile match

Match a profile against a job description.

```bash
sprofile match [flags]
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--profile` | Profile ID | (required) |
| `--profile-dir` | Profile directory | `~/.sprofile/profiles/` |
| `--jd` | Job description file | (required) |
| `-o, --output` | Output file path | stdout |
| `--strict` | Use strict pass criteria | false |

**Example:**

```bash
sprofile match --profile john-doe --jd jobdescription.md -o matcheval.json
```

### sprofile resume generate

Generate a tailored resume.

```bash
sprofile resume generate [flags]
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--profile` | Profile ID | (required) |
| `--profile-dir` | Profile directory | `~/.sprofile/profiles/` |
| `--jd` | Job description file | (required) |
| `-o, --output` | Output file path | stdout |
| `--domain` | Domain filter (devx, iam, platform) | auto |

**Example:**

```bash
sprofile resume generate \
  --profile john-doe \
  --jd jobdescription.md \
  --domain devx \
  -o resume.md
```

**Output Formats:**

| Extension | Format |
|-----------|--------|
| `.md` | Markdown |
| `.pdf` | PDF (via pandoc) |
| `.docx` | Word (via pandoc) |

### sprofile cover generate

Generate a cover letter.

```bash
sprofile cover generate [flags]
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--profile` | Profile ID | (required) |
| `--profile-dir` | Profile directory | `~/.sprofile/profiles/` |
| `--jd` | Job description file | (required) |
| `-o, --output` | Output file path | stdout |
| `--template` | Cover letter template ID | default |

**Example:**

```bash
sprofile cover generate \
  --profile john-doe \
  --jd jobdescription.md \
  -o coverletter.md
```

### sprofile doceval

Evaluate document quality.

```bash
sprofile doceval [flags]
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--resume` | Resume file path | (required) |
| `--cover` | Cover letter file path | optional |
| `--jd` | Job description file | (required) |
| `--matcheval` | Reference matcheval.json | optional |
| `-o, --output` | Output file path | stdout |
| `--strict` | Use strict criteria | false |

**Example:**

```bash
sprofile doceval \
  --resume resume.md \
  --cover coverletter.md \
  --jd jobdescription.md \
  -o doceval.json
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--help` | Show help |
| `--version` | Show version |
| `-v, --verbose` | Verbose output |

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SPROFILE_DIR` | Profile directory | `~/.sprofile/profiles/` |
| `SPROFILE_MODEL` | LLM model for AI agents | sonnet |
| `ANTHROPIC_API_KEY` | API key for Claude | (required for AI) |

## PDF Export

For PDF export, use pandoc:

```bash
# Basic PDF
pandoc resume.md -o resume.pdf

# With custom font (requires LuaLaTeX)
pandoc resume.md -o resume.pdf --pdf-engine=lualatex

# With YAML frontmatter for styling
cat > resume.md << 'EOF'
---
geometry: margin=2cm
mainfont: Lato
fontsize: 11pt
---

# John Doe
...
EOF

pandoc resume.md -o resume.pdf --pdf-engine=lualatex
```

## Workflow Example

Complete workflow for a job application:

```bash
# 1. Create application directory
mkdir -p applications/app_2026-04-19_acme_vp-platform
cd applications/app_2026-04-19_acme_vp-platform

# 2. Save job description
# (paste JD into jobdescription.md)

# 3. Analyze JD
sprofile jd analyze jobdescription.md -o jdanalysis.json

# 4. Match profile
sprofile match --profile john-doe --jd jobdescription.md -o matcheval.json

# 5. Generate resume
sprofile resume generate --profile john-doe --jd jobdescription.md -o resume.md

# 6. Generate cover letter
sprofile cover generate --profile john-doe --jd jobdescription.md -o coverletter.md

# 7. Evaluate documents
sprofile doceval --resume resume.md --cover coverletter.md --jd jobdescription.md -o doceval.json

# 8. Export to PDF
pandoc resume.md -o resume.pdf --pdf-engine=lualatex
pandoc coverletter.md -o coverletter.pdf --pdf-engine=lualatex
```
