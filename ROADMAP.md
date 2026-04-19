# Product Roadmap: structured-profile

**Version:** 1.0.0
**Status:** Draft
**Last Updated:** 2026-02-08

## Roadmap Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          structured-profile Roadmap                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  v0.1.0          v0.2.0          v0.3.0          v1.0.0          v2.0.0    │
│  Foundation      Resume Gen      Agent Team      Full CLI        Web App   │
│  ────────────    ──────────      ──────────      ────────        ───────   │
│  │              │               │               │               │          │
│  ▼              ▼               ▼               ▼               ▼          │
│  • Schema       • JD Parser     • LLM Providers • All Commands  • REST API │
│  • JSON Store   • Matcher       • Agent Defs    • Application   • Web UI   │
│  • Migration    • Resume Svc    • Orchestrator  • Interview     • Auth     │
│  • Basic CLI    • MD Export     • LLM-as-Judge  • H5P Export    • Multi-   │
│                 • PDF/DOCX      • Cover Letter  • Full Docs       user     │
│                                                                             │
│  ◄──────────── MVP ────────────►                                           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Version 0.1.0 - Foundation

**Theme:** Core data model and storage

**Target:** Establish the foundational schema and JSON storage backend, migrate existing gocvresume data.

### Features

| Feature | Description | Priority |
|---------|-------------|----------|
| Core Schema | Profile, Tenure, Position, Achievement (STAR), Skill, Education, Certification | P0 |
| JSON Store | File-based storage with full CRUD operations | P0 |
| Data Migration | Convert gocvresume/jwdata to structured-profile format | P0 |
| Basic CLI | `profile init`, `profile show`, experience commands | P1 |
| Schema Validation | JSON Schema generation and validation | P1 |

### Deliverables

- [ ] `schema/` package with all core types
- [ ] `store/json/` implementation
- [ ] `cmd/migrate/` tool for gocvresume migration
- [ ] `cmd/sprofile/` with profile and experience commands
- [ ] JSON Schema files in `schema/`
- [ ] Unit tests with >80% coverage
- [ ] README with quickstart guide

### Acceptance Criteria

1. All gocvresume data successfully migrated to JSON format
2. Profile CRUD operations work via CLI
3. Experience hierarchy (Tenure → Position → Achievement) maintained
4. STAR format preserved for all achievements
5. Domain-based filtering works for positions

### Release Checklist

- [ ] All tests passing
- [ ] golangci-lint clean
- [ ] Migration tested with real data
- [ ] README updated
- [ ] Tagged v0.1.0

---

## Version 0.2.0 - Resume Generation

**Theme:** Opportunity-specific resume generation

**Target:** Generate tailored resumes from master profile based on job descriptions.

### Features

| Feature | Description | Priority |
|---------|-------------|----------|
| Opportunity Schema | Job opportunity, parsed JD, company info | P0 |
| JD Parser | Extract skills, keywords, requirements from job descriptions | P0 |
| Profile-JD Matcher | Score relevance of experiences to job requirements | P0 |
| Resume Service | Generate tailored resumes with filtered content | P0 |
| Markdown Export | Export resume to Markdown format | P0 |
| PDF/DOCX Export | Export via Pandoc integration | P1 |
| Resume Templates | Support multiple resume formats/themes | P2 |

### Deliverables

- [ ] `schema/opportunity.go` and `schema/resume.go`
- [ ] `jdparser/` package
- [ ] `matcher/` package
- [ ] `service/resume.go`
- [ ] `export/markdown.go` and `export/pandoc.go`
- [ ] CLI commands: `resume generate`, `resume export`
- [ ] Sample job descriptions for testing

### Acceptance Criteria

1. JD parser extracts required/preferred skills with >90% accuracy
2. Matcher correctly ranks experiences by relevance
3. Generated resume includes only relevant content
4. PDF output matches expected formatting
5. Multiple resumes can be generated for same profile

### Release Checklist

- [ ] JD parser tested with 10+ real job descriptions
- [ ] Matcher accuracy validated
- [ ] PDF output visually verified
- [ ] Documentation updated
- [ ] Tagged v0.2.0

---

## Version 0.3.0 - LLM-Powered Agent Team

**Theme:** Multi-agent LLM system for resume generation and evaluation

**Target:** Implement an LLM-powered agent team using multi-agent-spec to parse job descriptions, generate tailored resumes/cover letters, and evaluate quality using LLM-as-a-Judge.

### Features

| Feature | Description | Priority |
|---------|-------------|----------|
| LLM Provider Abstraction | Support Anthropic, OpenAI, Bedrock | P0 |
| Agent Team Definitions | Define agents using multi-agent-spec | P0 |
| Agent Orchestrator | DAG workflow execution for agent team | P0 |
| JD Analyst Agent | LLM-powered job description parsing | P0 |
| Resume Generator Agent | LLM-based tailored resume generation | P0 |
| Resume Evaluator Agent | LLM-as-a-Judge resume evaluation | P0 |
| Resume Refiner Agent | Iterative improvement based on evaluation | P0 |
| Cover Letter Generator Agent | LLM-based cover letter generation | P1 |
| Evaluation Report | Structured report with findings | P0 |
| Cover Letter Export | MD/PDF/DOCX output | P1 |

### Agent Team: Resume Pipeline

```
JD Analyst ──┬──▶ Resume Generator ──┬──▶ Resume Evaluator ──▶ Refiner
             │                       │
             └──▶ Cover Letter Generator
```

| Agent | Model | Purpose |
|-------|-------|---------|
| `jd-analyst` | sonnet | Parse JD, extract requirements/keywords |
| `resume-generator` | sonnet | Select experiences, generate tailored resume |
| `cover-letter-generator` | sonnet | Generate personalized cover letter |
| `resume-evaluator` | opus | LLM-as-a-Judge evaluation |
| `resume-refiner` | sonnet | Improve based on evaluation findings |

### Evaluation Categories

| Category | Weight | Description |
|----------|--------|-------------|
| Skills Alignment | 25% | Match between resume skills and JD |
| Experience Level | 25% | Seniority and domain experience |
| Achievement Clarity | 20% | STAR format quality and metrics |
| Industry Fit | 15% | Domain and industry relevance |
| Communication | 15% | Writing quality and presentation |

### Deliverables

- [ ] `llm/` package with provider implementations
- [ ] `agents/resume-pipeline/` agent definitions (Markdown + YAML frontmatter)
- [ ] `agents/teams/resume-pipeline.json` team definition
- [ ] `orchestrator/` package for workflow execution
- [ ] `service/evaluation.go`
- [ ] Evaluation rubric configuration
- [ ] `schema/coverletter.go`
- [ ] `service/coverletter.go`
- [ ] CLI commands: `eval resume`, `eval report`, `cover generate`
- [ ] structured-evaluation integration
- [ ] multi-agent-spec integration

### Acceptance Criteria

1. Agent team executes full DAG workflow end-to-end
2. JD analyst extracts skills/keywords with >90% accuracy
3. Resume generator produces relevant, tailored content
4. Evaluation produces consistent scores for same input
5. Findings identify actionable improvements
6. Cover letters reference relevant STAR achievements
7. Pass/fail decision matches human judgment >80%
8. Iterative refinement improves evaluation score

### Release Checklist

- [ ] LLM providers tested with mocked responses
- [ ] Agent definitions validated with multi-agent-spec
- [ ] Orchestrator tested with full workflow
- [ ] Evaluation accuracy validated against human ratings
- [ ] Cover letter quality reviewed
- [ ] Cost tracking implemented
- [ ] Tagged v0.3.0

---

## Version 1.0.0 - Full CLI & Application Tracking

**Theme:** Complete career management system

**Target:** Full-featured CLI with application tracking, interview feedback, and interview preparation.

### Features

| Feature | Description | Priority |
|---------|-------------|----------|
| Application Schema | Full application lifecycle tracking | P0 |
| Application Service | Status tracking, outcome recording | P0 |
| Interview Schema | Interview rounds, questions, feedback | P0 |
| Feedback Integration | Learn from interview feedback | P0 |
| Interview Prep Schema | Question sets, sections, difficulty | P1 |
| Interview Prep Service | Generate prep from profile + feedback | P1 |
| H5P Export | Export prep sets to H5P format | P1 |
| Full CLI | All commands implemented | P0 |
| Documentation | Comprehensive user guide | P0 |

### Application Lifecycle

```
Draft → Submitted → Screening → Interview → Offer → Decision
                                   │
                                   ├── Round 1 (Phone)
                                   ├── Round 2 (Technical)
                                   ├── Round 3 (Behavioral)
                                   └── Final
```

### Deliverables

- [ ] `schema/application.go` and `schema/interview.go`
- [ ] `service/application.go`
- [ ] `schema/interviewprep.go`
- [ ] `service/interview.go`
- [ ] `export/h5p.go`
- [ ] Full CLI command set
- [ ] User documentation (MkDocs)
- [ ] Example workflows

### Acceptance Criteria

1. Full application lifecycle can be tracked via CLI
2. Interview questions captured and categorized
3. Feedback incorporated into future prep
4. H5P export produces valid quiz files
5. Documentation covers all features

### Release Checklist

- [ ] All CLI commands tested
- [ ] End-to-end workflow validated
- [ ] H5P output tested in H5P player
- [ ] Documentation complete
- [ ] CHANGELOG updated
- [ ] Tagged v1.0.0

---

## Version 2.0.0 - Web Application (Future)

**Theme:** Multi-user web application

**Target:** Web-based interface for managing profiles, generating resumes, and tracking applications.

### Features

| Feature | Description | Priority |
|---------|-------------|----------|
| PostgreSQL Backend | Database storage with migrations | P0 |
| REST API | JSON API for all operations | P0 |
| Authentication | User accounts, OAuth | P0 |
| Web UI | React/Vue frontend | P0 |
| Multi-user | Separate profiles per user | P0 |
| Real-time Updates | WebSocket for status updates | P2 |
| Job Board Integration | Import from LinkedIn, Indeed | P2 |
| Analytics Dashboard | Application success rates | P2 |

### Architecture

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   Web UI     │────▶│   REST API   │────▶│  PostgreSQL  │
│  (React)     │     │   (Go)       │     │              │
└──────────────┘     └──────────────┘     └──────────────┘
                            │
                            ▼
                     ┌──────────────┐
                     │  LLM APIs    │
                     │ (Anthropic)  │
                     └──────────────┘
```

### Deliverables

- [ ] `store/postgres/` implementation
- [ ] Database migrations
- [ ] REST API handlers
- [ ] Authentication middleware
- [ ] Web frontend
- [ ] Docker deployment
- [ ] API documentation (OpenAPI)

### Acceptance Criteria

1. Users can sign up and manage profiles
2. All CLI functionality available via web
3. Data isolated between users
4. Responsive design for mobile
5. Sub-second response times

### Release Checklist

- [ ] Security audit completed
- [ ] Load testing passed
- [ ] GDPR compliance verified
- [ ] Deployment documentation
- [ ] Tagged v2.0.0

---

## Feature Backlog (Future Versions)

### v2.1.0 - Enhanced Intelligence

- [ ] AI-powered resume suggestions
- [ ] Automatic skill extraction from experiences
- [ ] Interview question prediction
- [ ] Salary negotiation guidance

### v2.2.0 - Integration Hub

- [ ] LinkedIn profile import
- [ ] Indeed job search integration
- [ ] ATS compatibility checker
- [ ] Calendar integration for interviews

### v2.3.0 - Analytics & Insights

- [ ] Application success rate tracking
- [ ] Skill gap analysis
- [ ] Market demand insights
- [ ] Compensation benchmarking

### v2.4.0 - Collaboration

- [ ] Recruiter/coach view
- [ ] Resume review requests
- [ ] Peer feedback system
- [ ] Team management (for recruiters)

---

## Versioning Strategy

This project follows [Semantic Versioning](https://semver.org/):

- **MAJOR** (X.0.0): Breaking changes to CLI, API, or data format
- **MINOR** (0.X.0): New features, backward compatible
- **PATCH** (0.0.X): Bug fixes, backward compatible

### Data Migration Policy

- Breaking schema changes require migration scripts
- Migration scripts included in release
- Backward compatibility maintained within major version
- Data export always available before major upgrades

---

## Success Metrics by Version

| Version | Key Metric | Target |
|---------|------------|--------|
| v0.1.0 | Migration success rate | 100% |
| v0.2.0 | Resume generation time | < 5s |
| v0.3.0 | Evaluation accuracy | > 80% |
| v1.0.0 | CLI command coverage | 100% |
| v2.0.0 | API response time | < 200ms |

---

## Dependencies & External Systems

### Required Dependencies

| Dependency | Version | Purpose | Version |
|------------|---------|---------|---------|
| Go | 1.22+ | Runtime | Required |
| Pandoc | 2.0+ | PDF/DOCX export | v0.2.0+ |
| LLM API | Any | Evaluation | v0.3.0+ |
| PostgreSQL | 15+ | Database | v2.0.0+ |

### Integrated Libraries

| Library | Purpose | Integration |
|---------|---------|-------------|
| multi-agent-spec | Agent team definitions and orchestration | v0.3.0 |
| structured-evaluation | LLM-as-a-Judge | v0.3.0 |
| h5p-go | Interview prep export | v1.0.0 |
| gopandoc | Document conversion | v0.2.0 |

---

## Release Schedule (Tentative)

| Version | Target | Status |
|---------|--------|--------|
| v0.1.0 | Q1 2026 | Planning |
| v0.2.0 | Q1 2026 | Planning |
| v0.3.0 | Q2 2026 | Planning |
| v1.0.0 | Q2 2026 | Planning |
| v2.0.0 | Q4 2026 | Planning |

---

## Feedback & Contributions

This roadmap is a living document. Feedback and contributions welcome via:

- GitHub Issues for feature requests
- Pull Requests for contributions
- Discussions for roadmap feedback
