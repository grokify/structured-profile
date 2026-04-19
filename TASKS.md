# Tasks

Open tasks for structured-profile development.

## v0.3.0 - Agent Team (Resume Pipeline)

### Epic: Agent Specs

- [ ] Create `specs/` directory structure
- [ ] Create `specs/plugin.json` manifest
- [ ] Create `specs/agents/jd-analyst.md` - Parse job descriptions, extract requirements
- [ ] Create `specs/agents/resume-generator.md` - Generate tailored resumes from profile
- [ ] Create `specs/agents/resume-evaluator.md` - LLM-as-a-Judge evaluation (outputs structured-evaluation JSON)
- [ ] Create `specs/agents/resume-refiner.md` - Improve resume based on evaluation findings
- [ ] Create `specs/agents/cover-letter-generator.md` - Generate personalized cover letters
- [ ] Create `specs/agents/resume-coordinator.md` - Orchestrates the pipeline
- [ ] Create `specs/teams/resume-pipeline.json` - DAG workflow definition
- [ ] Create `specs/deployments/claude-code.json` - Claude Code deployment config

### Epic: Claude Code Integration

- [ ] Generate `.claude/agents/` from specs using assistantkit
- [ ] Test agents as Claude Code subagents via Task tool
- [ ] Document agent invocation in README

### Epic: Evaluation Integration

- [ ] Define evaluation rubric for resume quality (categories, weights)
- [ ] Evaluator agent outputs `structured-evaluation` JSON format
- [ ] Write evaluation results to `evaluation.json` file
- [ ] Parent agent reads and processes evaluation results

### Epic: Skills/Commands

- [ ] Create `specs/skills/parse-jd.md` - Skill for parsing job descriptions
- [ ] Create `specs/skills/generate-resume.md` - Skill for resume generation
- [ ] Create `specs/commands/resume.md` - `/resume generate` command

---

## v1.0.0 - Full CLI & Application Tracking

### Epic: Application Lifecycle

- [ ] Create `service/application.go` - Application tracking service
- [ ] Implement application status transitions (Draft → Submitted → Screening → Interview → Offer → Decision)
- [ ] Add outcome recording and statistics

### Epic: Interview Tracking

- [ ] Create `schema/interview.go` - Interview rounds, questions, feedback
- [ ] Create `service/interview.go` - Interview management service
- [ ] Implement feedback integration for future prep

### Epic: Interview Prep

- [ ] Create `schema/interviewprep.go` - Question sets, sections, difficulty levels
- [ ] Create `service/interviewprep.go` - Generate prep from profile + feedback
- [ ] Create `export/h5p.go` - Export prep sets to H5P format

### Epic: CLI Commands

- [ ] Create `cmd/sprofile/` CLI application
- [ ] Implement `profile init`, `profile show` commands
- [ ] Implement `experience add`, `experience list` commands
- [ ] Implement `resume generate`, `resume export` commands
- [ ] Implement `cover generate`, `cover export` commands
- [ ] Implement `opportunity add`, `opportunity list` commands
- [ ] Implement `application submit`, `application status` commands
- [ ] Implement `interview add`, `interview feedback` commands
- [ ] Implement `prep generate`, `prep export` commands

### Epic: Documentation

- [ ] Create MkDocs site structure
- [ ] Write quickstart guide
- [ ] Document all CLI commands
- [ ] Add example workflows
- [ ] Generate API documentation

---

## v2.0.0 - Web Application (Future)

### Epic: Database Backend

- [ ] Create `store/postgres/` implementation
- [ ] Create database migrations
- [ ] Add connection pooling and error handling

### Epic: REST API

- [ ] Design OpenAPI specification
- [ ] Implement REST handlers
- [ ] Add authentication middleware (OAuth)
- [ ] Add rate limiting

### Epic: Web Frontend

- [ ] Set up React/Vue project
- [ ] Implement profile management UI
- [ ] Implement resume generation UI
- [ ] Implement application tracking UI
- [ ] Add responsive design

---

## Infrastructure & Tooling

### Epic: Claude Code Evaluation Enhancement

Context: Currently Claude Code Task tool returns unstructured messages. For better LLM-as-a-Judge integration, consider:

- [ ] Document current workaround (JSON in message + file-based output)
- [ ] File feature request: `PostAgentComplete` hook for capturing subagent output
- [ ] File feature request: Structured output type for Task tool
- [ ] Prototype hooks-based evaluation capture

### Epic: assistantkit Integration

- [ ] Verify `genagents` works with specs/ directory
- [ ] Test generation for Claude Code target
- [ ] Add to CI/CD pipeline for automatic regeneration

### Epic: Testing

- [ ] Add integration tests for agent pipeline (mocked)
- [ ] Add end-to-end test with real LLM (optional, CI skip)
- [ ] Validate evaluation output against structured-evaluation schema

---

## Backlog

### Ideas for Future Versions

- [ ] AI-powered resume suggestions
- [ ] Automatic skill extraction from experiences
- [ ] Interview question prediction
- [ ] Salary negotiation guidance
- [ ] LinkedIn profile import
- [ ] Indeed job search integration
- [ ] ATS compatibility checker
- [ ] Calendar integration for interviews
- [ ] Application success rate tracking
- [ ] Skill gap analysis
- [ ] Market demand insights
- [ ] Compensation benchmarking

---

## Completed

### v0.1.0 - Foundation

- [x] Core schema (Profile, Tenure, Position, Achievement, Skill, Education, Certification)
- [x] JSON store with full CRUD operations
- [x] Data migration from gocvresume format
- [x] Unit tests with good coverage

### v0.2.0 - Resume Generation

- [x] Opportunity schema and store operations
- [x] JD parser (`jdparser/`) - extracts skills, keywords, requirements
- [x] Profile-JD matcher (`matcher/`) - scores and ranks achievements
- [x] Resume service (`service/resume.go`) - generates tailored resumes
- [x] Markdown export (`export/markdown.go`)
- [x] PDF/DOCX export via Pandoc (`export/pandoc.go`)

### v0.3.0 - Cover Letter (Non-LLM Parts)

- [x] CoverLetter and CoverLetterTemplate store operations
- [x] CoverLetter service (`service/coverletter.go`) - template-based generation
- [x] CoverLetter export (Markdown, PDF, DOCX)
- [x] Tests for all cover letter functionality
