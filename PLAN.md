# Implementation Plan: structured-profile

**Version:** 1.0.0
**Status:** Draft
**Last Updated:** 2026-02-08

## Overview

This document outlines the implementation plan for `structured-profile`, breaking down the work into phases, epics, and tasks with dependencies.

## Phase 1: Foundation (Core Schema & JSON Store)

### Epic 1.1: Project Setup

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 1.1.1 | Initialize Go module with `go mod init` | None | S |
| 1.1.2 | Set up directory structure per TRD | 1.1.1 | S |
| 1.1.3 | Add initial dependencies to go.mod | 1.1.1 | S |
| 1.1.4 | Create Makefile with build, test, lint targets | 1.1.2 | S |
| 1.1.5 | Set up golangci-lint configuration | 1.1.2 | S |
| 1.1.6 | Create initial README.md | 1.1.2 | S |

### Epic 1.2: Core Schema Types

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 1.2.1 | Implement `BaseEntity` with ID, timestamps | 1.1.2 | S |
| 1.2.2 | Implement `Profile` type | 1.2.1 | S |
| 1.2.3 | Implement `Tenure` type | 1.2.1 | S |
| 1.2.4 | Implement `Position` and `PositionDomainConfig` | 1.2.1 | M |
| 1.2.5 | Implement `Achievement` with STAR fields | 1.2.1 | M |
| 1.2.6 | Implement `Skill` and `SkillCategory` | 1.2.1 | S |
| 1.2.7 | Implement `Education` and `Certification` | 1.2.1 | S |
| 1.2.8 | Implement `VerifiableCredential` | 1.2.1 | S |
| 1.2.9 | Implement `Date` helper type (YYYYMM format) | 1.2.1 | S |
| 1.2.10 | Add JSON Schema generation for all types | 1.2.2-1.2.9 | M |
| 1.2.11 | Write unit tests for all schema types | 1.2.2-1.2.9 | M |

### Epic 1.3: Store Interface & JSON Backend

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 1.3.1 | Define `Store` interface | 1.2.11 | M |
| 1.3.2 | Implement JSON file store structure | 1.3.1 | M |
| 1.3.3 | Implement Profile CRUD in JSON store | 1.3.2 | M |
| 1.3.4 | Implement Tenure CRUD in JSON store | 1.3.2 | M |
| 1.3.5 | Implement Position CRUD in JSON store | 1.3.2 | M |
| 1.3.6 | Implement Achievement CRUD in JSON store | 1.3.2 | M |
| 1.3.7 | Implement Skills CRUD in JSON store | 1.3.2 | S |
| 1.3.8 | Implement Education/Certification CRUD | 1.3.2 | S |
| 1.3.9 | Implement Credential CRUD | 1.3.2 | S |
| 1.3.10 | Implement search by tags/skills | 1.3.6 | M |
| 1.3.11 | Write integration tests for JSON store | 1.3.3-1.3.10 | M |

### Epic 1.4: Data Migration from gocvresume

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 1.4.1 | Create migration tool `cmd/migrate` | 1.3.11 | M |
| 1.4.2 | Map gocvresume types to structured-profile | 1.4.1 | M |
| 1.4.3 | Convert jwdata to JSON master data | 1.4.2 | L |
| 1.4.4 | Validate migrated data | 1.4.3 | M |
| 1.4.5 | Document migration process | 1.4.4 | S |

## Phase 2: Resume Generation

### Epic 2.1: Opportunity Schema

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 2.1.1 | Implement `Opportunity` type | 1.2.1 | S |
| 2.1.2 | Implement `JobDescParsed` type | 2.1.1 | S |
| 2.1.3 | Implement `CompanyInfo` type | 2.1.1 | S |
| 2.1.4 | Add Opportunity CRUD to Store interface | 2.1.1 | S |
| 2.1.5 | Implement Opportunity in JSON store | 2.1.4 | M |

### Epic 2.2: Job Description Parser

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 2.2.1 | Create `jdparser` package structure | 2.1.5 | S |
| 2.2.2 | Implement keyword extraction | 2.2.1 | M |
| 2.2.3 | Implement skills extraction | 2.2.1 | M |
| 2.2.4 | Implement experience level detection | 2.2.1 | M |
| 2.2.5 | Add LLM-powered JD analysis (optional) | 2.2.2-2.2.4 | L |
| 2.2.6 | Write tests with sample JDs | 2.2.2-2.2.4 | M |

### Epic 2.3: Profile-JD Matcher

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 2.3.1 | Create `matcher` package structure | 2.2.6 | S |
| 2.3.2 | Implement skill matching algorithm | 2.3.1 | M |
| 2.3.3 | Implement tag/keyword matching | 2.3.1 | M |
| 2.3.4 | Implement experience relevance scoring | 2.3.1 | M |
| 2.3.5 | Implement achievement ordering by relevance | 2.3.4 | M |
| 2.3.6 | Write matcher unit tests | 2.3.2-2.3.5 | M |

### Epic 2.4: Resume Service

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 2.4.1 | Implement `Resume` type | 1.2.1 | S |
| 2.4.2 | Create `service/resume.go` | 2.3.6 | M |
| 2.4.3 | Implement `GenerateTailoredResume()` | 2.4.2 | L |
| 2.4.4 | Implement domain-based filtering | 2.4.3 | M |
| 2.4.5 | Implement summary selection by domain | 2.4.3 | S |
| 2.4.6 | Write resume generation tests | 2.4.3-2.4.5 | M |

### Epic 2.5: Markdown Export

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 2.5.1 | Create `export/markdown.go` | 2.4.6 | M |
| 2.5.2 | Implement header/contact rendering | 2.5.1 | S |
| 2.5.3 | Implement summary rendering | 2.5.1 | S |
| 2.5.4 | Implement experience/STAR rendering | 2.5.1 | M |
| 2.5.5 | Implement education/certs rendering | 2.5.1 | S |
| 2.5.6 | Implement skills rendering | 2.5.1 | S |
| 2.5.7 | Add template support | 2.5.2-2.5.6 | M |
| 2.5.8 | Write golden tests for Markdown output | 2.5.7 | M |

### Epic 2.6: PDF/DOCX Export

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 2.6.1 | Create `export/pandoc.go` | 2.5.8 | M |
| 2.6.2 | Implement PDF export via gopandoc | 2.6.1 | M |
| 2.6.3 | Implement DOCX export via gopandoc | 2.6.1 | M |
| 2.6.4 | Add margin and formatting options | 2.6.2 | S |
| 2.6.5 | Write export integration tests | 2.6.2-2.6.4 | M |

## Phase 3: Cover Letter Generation

### Epic 3.1: Cover Letter Schema

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 3.1.1 | Implement `CoverLetter` type | 1.2.1 | S |
| 3.1.2 | Implement `CoverLetterTemplate` type | 3.1.1 | S |
| 3.1.3 | Add CoverLetter to Store interface | 3.1.1 | S |
| 3.1.4 | Implement CoverLetter in JSON store | 3.1.3 | M |

### Epic 3.2: Cover Letter Service

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 3.2.1 | Create `service/coverletter.go` | 3.1.4 | M |
| 3.2.2 | Implement template variable substitution | 3.2.1 | M |
| 3.2.3 | Implement STAR reference integration | 3.2.1, 2.4.6 | M |
| 3.2.4 | Add company research integration | 3.2.1, 2.1.5 | M |
| 3.2.5 | Write cover letter generation tests | 3.2.2-3.2.4 | M |

### Epic 3.3: Cover Letter Export

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 3.3.1 | Add cover letter to Markdown exporter | 3.2.5, 2.5.8 | M |
| 3.3.2 | Add cover letter to PDF/DOCX exporter | 3.3.1, 2.6.5 | M |
| 3.3.3 | Write export tests | 3.3.2 | M |

## Phase 4: LLM-as-a-Judge Evaluation

### Epic 4.1: LLM Provider Abstraction

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 4.1.1 | Create `llm` package structure | 1.1.2 | S |
| 4.1.2 | Define `Provider` interface | 4.1.1 | S |
| 4.1.3 | Implement Anthropic provider | 4.1.2 | M |
| 4.1.4 | Implement OpenAI provider | 4.1.2 | M |
| 4.1.5 | Implement Bedrock provider | 4.1.2 | M |
| 4.1.6 | Add provider configuration | 4.1.3-4.1.5 | M |
| 4.1.7 | Write provider tests (mocked) | 4.1.3-4.1.5 | M |

### Epic 4.2: Agent Team Definitions (multi-agent-spec)

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 4.2.1 | Create `agents/resume-pipeline/` directory structure | 1.1.2 | S |
| 4.2.2 | Define `jd-analyst` agent (JD parsing via LLM) | 4.2.1 | M |
| 4.2.3 | Define `resume-generator` agent (LLM-powered resume tailoring) | 4.2.1 | M |
| 4.2.4 | Define `cover-letter-generator` agent | 4.2.1 | M |
| 4.2.5 | Define `resume-evaluator` agent (LLM-as-a-Judge) | 4.2.1 | M |
| 4.2.6 | Define `resume-refiner` agent (iterative improvement) | 4.2.1 | M |
| 4.2.7 | Create `teams/resume-pipeline.json` team definition | 4.2.2-4.2.6 | M |
| 4.2.8 | Validate agent definitions with multi-agent-spec | 4.2.7 | S |

### Epic 4.3: Agent Team Orchestrator

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 4.3.1 | Create `orchestrator` package structure | 4.1.7, 4.2.8 | M |
| 4.3.2 | Implement agent loader (from multi-agent-spec) | 4.3.1 | M |
| 4.3.3 | Implement DAG workflow runner | 4.3.2 | L |
| 4.3.4 | Implement step input/output resolution | 4.3.3 | M |
| 4.3.5 | Implement agent execution (LLM calls) | 4.3.4 | L |
| 4.3.6 | Implement shared context management | 4.3.5 | M |
| 4.3.7 | Write orchestrator unit tests | 4.3.5-4.3.6 | M |
| 4.3.8 | Write orchestrator integration tests | 4.3.7 | L |

### Epic 4.4: Evaluation Service

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 4.4.1 | Create `service/evaluation.go` | 4.3.8 | M |
| 4.4.2 | Define resume evaluation rubric | 4.4.1 | M |
| 4.4.3 | Implement category evaluation prompts | 4.4.2 | L |
| 4.4.4 | Implement findings generation | 4.4.3 | L |
| 4.4.5 | Integrate structured-evaluation report | 4.4.4 | M |
| 4.4.6 | Implement evaluation storage | 4.4.5 | M |
| 4.4.7 | Write evaluation service tests | 4.4.5-4.4.6 | M |

### Epic 4.5: Evaluation Reporting

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 4.5.1 | Implement terminal report rendering | 4.4.7 | M |
| 4.5.2 | Implement JSON report export | 4.4.7 | S |
| 4.5.3 | Implement Markdown report export | 4.4.7 | M |
| 4.5.4 | Add evaluation history tracking | 4.4.6 | M |
| 4.5.5 | Write report generation tests | 4.5.1-4.5.4 | M |

## Phase 5: Application Tracking

### Epic 5.1: Application Schema

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 5.1.1 | Implement `Application` type | 1.2.1 | S |
| 5.1.2 | Implement `ApplicationStatus` enum | 5.1.1 | S |
| 5.1.3 | Implement `ApplicationOutcome` type | 5.1.1 | S |
| 5.1.4 | Add Application to Store interface | 5.1.1 | S |
| 5.1.5 | Implement Application in JSON store | 5.1.4 | M |

### Epic 5.2: Interview Schema

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 5.2.1 | Implement `Interview` type | 1.2.1 | M |
| 5.2.2 | Implement `InterviewType` enum | 5.2.1 | S |
| 5.2.3 | Implement `InterviewQuestion` type | 5.2.1 | M |
| 5.2.4 | Implement `SelfAssessment` type | 5.2.1 | S |
| 5.2.5 | Implement `InterviewFeedback` type | 5.2.1 | S |
| 5.2.6 | Add Interview to Store interface | 5.2.1 | S |
| 5.2.7 | Implement Interview in JSON store | 5.2.6 | M |

### Epic 5.3: Application Service

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 5.3.1 | Create `service/application.go` | 5.1.5, 5.2.7 | M |
| 5.3.2 | Implement application lifecycle | 5.3.1 | M |
| 5.3.3 | Implement status transitions | 5.3.2 | M |
| 5.3.4 | Implement interview recording | 5.3.1 | M |
| 5.3.5 | Implement feedback capture | 5.3.4 | M |
| 5.3.6 | Implement outcome recording | 5.3.1 | M |
| 5.3.7 | Write application service tests | 5.3.2-5.3.6 | M |

### Epic 5.4: Feedback Integration

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 5.4.1 | Implement lessons learned extraction | 5.3.7 | M |
| 5.4.2 | Implement strength/weakness aggregation | 5.4.1 | M |
| 5.4.3 | Link feedback to master profile updates | 5.4.2 | L |
| 5.4.4 | Add feedback to interview prep | 5.4.3 | M |
| 5.4.5 | Write feedback integration tests | 5.4.1-5.4.4 | M |

## Phase 6: Interview Preparation

### Epic 6.1: Interview Prep Schema

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 6.1.1 | Implement `InterviewPrepSet` type | 1.2.1 | M |
| 6.1.2 | Implement `PrepSection` type | 6.1.1 | S |
| 6.1.3 | Implement `PrepQuestion` type | 6.1.1 | M |
| 6.1.4 | Implement `PrepAnswer` type | 6.1.3 | S |
| 6.1.5 | Add PrepSet to Store interface | 6.1.1 | S |
| 6.1.6 | Implement PrepSet in JSON store | 6.1.5 | M |

### Epic 6.2: Interview Service

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 6.2.1 | Create `service/interview.go` | 6.1.6 | M |
| 6.2.2 | Implement question generation from profile | 6.2.1 | L |
| 6.2.3 | Link questions to STAR achievements | 6.2.2 | M |
| 6.2.4 | Implement difficulty progression | 6.2.2 | M |
| 6.2.5 | Import questions from interview feedback | 5.4.5 | M |
| 6.2.6 | Write interview prep service tests | 6.2.2-6.2.5 | M |

### Epic 6.3: H5P Export

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 6.3.1 | Create `export/h5p.go` | 6.2.6 | M |
| 6.3.2 | Implement QuestionSet generation | 6.3.1 | M |
| 6.3.3 | Implement section/topic mapping | 6.3.2 | M |
| 6.3.4 | Implement h5p-go extensions | 6.3.2 | M |
| 6.3.5 | Implement H5P package export | 6.3.3-6.3.4 | M |
| 6.3.6 | Write H5P export tests | 6.3.5 | M |

## Phase 7: CLI Application

### Epic 7.1: CLI Framework

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 7.1.1 | Set up Cobra CLI framework | 1.1.2 | M |
| 7.1.2 | Implement root command | 7.1.1 | S |
| 7.1.3 | Add configuration loading | 7.1.2 | M |
| 7.1.4 | Add output formatting (table, JSON) | 7.1.2 | M |
| 7.1.5 | Implement error handling | 7.1.2 | M |

### Epic 7.2: Profile Commands

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 7.2.1 | Implement `profile init` | 7.1.5, 1.3.11 | M |
| 7.2.2 | Implement `profile show` | 7.2.1 | M |
| 7.2.3 | Implement `profile export` | 7.2.1 | M |
| 7.2.4 | Implement `tenure` subcommands | 7.2.1 | M |
| 7.2.5 | Implement `position` subcommands | 7.2.4 | M |
| 7.2.6 | Implement `achievement` subcommands | 7.2.5 | M |

### Epic 7.3: Resume Commands

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 7.3.1 | Implement `resume generate` | 7.1.5, 2.4.6 | M |
| 7.3.2 | Implement `resume list` | 7.3.1 | S |
| 7.3.3 | Implement `resume export` | 7.3.1, 2.6.5 | M |

### Epic 7.4: Evaluation Commands

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 7.4.1 | Implement `eval resume` | 7.1.5, 4.2.7 | M |
| 7.4.2 | Implement `eval report` | 7.4.1 | M |
| 7.4.3 | Implement `eval history` | 7.4.1 | M |

### Epic 7.5: Application Commands

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 7.5.1 | Implement `apply create` | 7.1.5, 5.3.7 | M |
| 7.5.2 | Implement `apply list` | 7.5.1 | S |
| 7.5.3 | Implement `apply status` | 7.5.1 | M |
| 7.5.4 | Implement `apply outcome` | 7.5.1 | M |
| 7.5.5 | Implement `interview` subcommands | 7.5.1 | M |

### Epic 7.6: Interview Prep Commands

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 7.6.1 | Implement `prep generate` | 7.1.5, 6.2.6 | M |
| 7.6.2 | Implement `prep export` | 7.6.1, 6.3.6 | M |
| 7.6.3 | Implement `prep practice` | 7.6.1 | L |

## Phase 8: PostgreSQL Backend (Future)

### Epic 8.1: Database Implementation

| Task | Description | Dependencies | Estimate |
|------|-------------|--------------|----------|
| 8.1.1 | Create migration scripts | 7.6.3 | M |
| 8.1.2 | Implement PostgreSQL store | 8.1.1 | L |
| 8.1.3 | Implement transactions | 8.1.2 | M |
| 8.1.4 | Write integration tests | 8.1.2-8.1.3 | L |

## Dependency Graph

```
Phase 1 ──┬── Phase 2 ──┬── Phase 3
          │             │
          │             └── Phase 4
          │
          └── Phase 5 ───── Phase 6
                              │
All Phases ─────────────── Phase 7 ───── Phase 8
```

## Estimates Key

- **S** (Small): < 2 hours
- **M** (Medium): 2-8 hours
- **L** (Large): 1-3 days
- **XL** (Extra Large): > 3 days

## Risk Mitigation

| Risk | Mitigation |
|------|------------|
| LLM API rate limits | Implement caching, retry logic |
| gocvresume migration complexity | Start with subset, iterate |
| Pandoc dependency | Document installation, provide fallback |
| H5P format changes | Pin h5p-go version, write compatibility tests |

## Definition of Done

Each epic is considered complete when:

1. All tasks are implemented
2. Unit tests pass with >80% coverage
3. Integration tests pass
4. Documentation updated
5. Code reviewed and merged
6. golangci-lint passes with no errors
