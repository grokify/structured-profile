# Product Requirements Document: structured-profile

**Version:** 1.0.0
**Status:** Draft
**Last Updated:** 2026-02-08

## Executive Summary

**structured-profile** is a Go library and CLI for managing professional career data, generating tailored resumes and cover letters, evaluating resume quality against job descriptions using LLM-as-a-Judge, tracking job applications, and preparing for interviews with structured Q&A content.

## Problem Statement

### Current Challenges

1. **Fragmented Career Data**: Professional information (experiences, skills, achievements) is scattered across multiple documents, making it difficult to maintain and update.

2. **Manual Resume Tailoring**: Creating opportunity-specific resumes requires manually selecting and reorganizing content for each application, a time-consuming and error-prone process.

3. **No Quality Feedback**: Candidates lack objective feedback on how well their resume matches a specific job description before submitting.

4. **Lost Application Context**: Interview feedback, application outcomes, and lessons learned are not systematically captured for future improvement.

5. **Unstructured Interview Prep**: Interview preparation materials are ad-hoc, not linked to specific roles or past experiences.

### Target Users

- Job seekers managing multiple applications
- Career changers needing to highlight transferable skills
- Technical professionals with verifiable online credentials
- Recruiters and career coaches managing candidate profiles

## Product Vision

A comprehensive career management system that:

1. Maintains a single source of truth for all professional data
2. Automatically generates tailored resumes using domain-based filtering
3. Provides LLM-powered quality evaluation of resumes against job descriptions
4. Tracks the full application lifecycle with feedback integration
5. Generates structured interview preparation content

## Functional Requirements

### FR-1: Master Profile Management

**FR-1.1**: Store comprehensive professional data in a structured, versioned format.

- Personal information and contact details
- Work experiences in STAR format (Situation, Task, Action, Result)
- Education and certifications with verification URLs
- Skills taxonomy with proficiency levels
- Verifiable online credentials (GitHub, StackOverflow, LinkedIn, etc.)
- Publications, patents, speaking engagements
- Professional memberships and awards

**FR-1.2**: Support multiple data backends.

- JSON/YAML file storage (initial implementation)
- PostgreSQL database (future web application)

**FR-1.3**: Provide CRUD operations for all profile entities.

**FR-1.4**: Track entity history with timestamps and soft deletes.

### FR-2: Opportunity-Specific Resume Generation

**FR-2.1**: Parse and analyze job descriptions to extract:

- Required skills and experience levels
- Preferred qualifications
- Company information and culture indicators
- Keywords for ATS optimization

**FR-2.2**: Select relevant experiences from master profile based on:

- Skill match scoring
- Tag and keyword relevance
- Domain alignment
- Configurable filters

**FR-2.3**: Generate tailored resumes with:

- Domain-specific skill highlighting
- Reordered achievements based on relevance
- Customized summaries per opportunity
- Multiple output formats (Markdown, PDF, DOCX, HTML)

**FR-2.4**: Support resume templates and themes.

### FR-3: Cover Letter Generation

**FR-3.1**: Generate cover letters that reference:

- Selected STAR achievements from the tailored resume
- Company-specific research
- Position requirements alignment

**FR-3.2**: Support cover letter templates with variable substitution.

**FR-3.3**: Maintain cover letter versions per opportunity.

### FR-4: LLM-as-a-Judge Resume Evaluation

**FR-4.1**: Evaluate resume quality against job descriptions using structured-evaluation framework.

**FR-4.2**: Provide weighted scoring across categories:

| Category | Weight | Description |
|----------|--------|-------------|
| Skills Alignment | 25% | Match between resume skills and JD requirements |
| Experience Level | 25% | Seniority and domain experience match |
| Achievement Clarity | 20% | STAR format quality and quantification |
| Industry Fit | 15% | Domain and industry experience relevance |
| Communication | 15% | Writing quality and presentation |

**FR-4.3**: Generate actionable findings with severity levels:

- Critical: Missing must-have requirements
- High: Significant skill gaps
- Medium: Nice-to-have gaps
- Low: Minor improvements
- Info: Observations

**FR-4.4**: Provide pass/fail decision with rationale.

**FR-4.5**: Support multiple judge configurations for consensus evaluation.

**FR-4.6**: Track evaluation history for improvement over time.

### FR-5: Application Tracking

**FR-5.1**: Track application lifecycle stages:

- Draft → Submitted → Screening → Interview → Offer → Decision

**FR-5.2**: Store per-application data:

- Job description (raw and parsed)
- Company research
- Resume version used
- Cover letter version used
- Evaluation report
- Application date and status updates

**FR-5.3**: Record interview rounds with:

- Interview type (phone, technical, behavioral, onsite)
- Interviewer names (optional)
- Questions asked
- Self-assessment notes
- Interviewer feedback (when available)

**FR-5.4**: Capture outcome and lessons learned:

- Offer details or rejection reason (if known)
- What worked well
- What to improve
- Questions to add to interview prep

### FR-6: Interview Preparation

**FR-6.1**: Generate H5P-compatible question sets using h5p-go.

**FR-6.2**: Organize questions by:

- Section (Behavioral, Technical, Domain-specific)
- Topic (Leadership, System Design, etc.)
- Difficulty (Easy, Medium, Hard)

**FR-6.3**: Link questions to:

- Specific STAR experiences for behavioral questions
- Skills from master profile for technical questions
- Past interview feedback

**FR-6.4**: Include rich feedback:

- Tips before answering
- Feedback on correct/incorrect answers
- Score-based overall assessment

**FR-6.5**: Track practice history and scores.

### FR-7: Export and Integration

**FR-7.1**: Export formats:

- Markdown (for version control, Pandoc conversion)
- PDF (via Pandoc)
- DOCX (via Pandoc)
- JSON (machine-readable IR)
- H5P (for interview prep)

**FR-7.2**: Integration points:

- GitHub for portfolio verification
- LinkedIn for profile import/export
- StackOverflow for reputation data
- ATS systems (future)

## Non-Functional Requirements

### NFR-1: Performance

- Profile load time < 100ms for JSON backend
- Resume generation < 1s
- Evaluation report generation < 30s (depends on LLM)

### NFR-2: Scalability

- Support profiles with 100+ experiences
- Support 1000+ applications per profile
- Database backend for multi-user web deployment

### NFR-3: Data Privacy

- Sensitive data (contact info) encryption at rest
- No PII in logs
- GDPR-compliant data export/deletion

### NFR-4: Extensibility

- Plugin architecture for new output formats
- Custom evaluation rubrics
- Additional data backends

### NFR-5: Usability

- CLI with intuitive commands
- JSON Schema validation for data files
- Comprehensive error messages

## Success Metrics

| Metric | Target |
|--------|--------|
| Time to generate tailored resume | < 5 minutes |
| Resume-JD match score improvement | +15% after feedback iteration |
| Interview prep coverage | 80% of asked questions in prep set |
| Application tracking adoption | 90% of applications tracked |

## Dependencies

| Dependency | Purpose |
|------------|---------|
| github.com/agentplexus/multi-agent-spec | Agent team definitions and orchestration |
| github.com/agentplexus/structured-evaluation | LLM-as-a-Judge framework |
| github.com/grokify/h5p-go | Interview prep Q&A generation |
| github.com/grokify/gopandoc | Markdown to PDF/DOCX conversion |
| github.com/grokify/mogo | Go utility functions |

## Agent Team Architecture

Resume generation, cover letter creation, and quality evaluation are powered by a multi-agent system defined using the [multi-agent-spec](https://github.com/agentplexus/multi-agent-spec) framework. This enables coordinated LLM-powered workflows with explicit data flow between specialized agents.

### Agent Team: Resume Pipeline

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Resume Pipeline Agent Team                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────┐     ┌─────────────┐     ┌─────────────┐                   │
│  │  JD Analyst │────▶│   Resume    │────▶│   Resume    │                   │
│  │   Agent     │     │  Generator  │     │  Evaluator  │                   │
│  └─────────────┘     │   Agent     │     │   Agent     │                   │
│        │             └─────────────┘     └─────────────┘                   │
│        │                    │                   │                           │
│        │                    ▼                   ▼                           │
│        │             ┌─────────────┐     ┌─────────────┐                   │
│        └────────────▶│ Cover Letter│     │   Refiner   │◀──────────────────┤
│                      │  Generator  │     │   Agent     │                   │
│                      │   Agent     │     └─────────────┘                   │
│                      └─────────────┘           │                           │
│                                                ▼                           │
│                                         [Improved Resume]                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Agent Definitions

| Agent | Model | Purpose | Inputs | Outputs |
|-------|-------|---------|--------|---------|
| `jd-analyst` | sonnet | Parse JD, extract requirements, keywords, culture signals | Job description text | Parsed JD (skills, requirements, keywords) |
| `resume-generator` | sonnet | Select relevant experiences, generate tailored resume | Master profile, Parsed JD | Tailored resume (Markdown) |
| `cover-letter-generator` | sonnet | Generate personalized cover letter with STAR references | Master profile, Parsed JD, Resume | Cover letter (Markdown) |
| `resume-evaluator` | opus | LLM-as-a-Judge evaluation against JD | Resume, Parsed JD | Evaluation report (structured-evaluation) |
| `resume-refiner` | sonnet | Improve resume based on evaluation findings | Resume, Evaluation report | Refined resume (Markdown) |

### Workflow Type: DAG (Directed Acyclic Graph)

```yaml
workflow:
  type: dag
  steps:
    - name: analyze-jd
      agent: jd-analyst
      inputs:
        - name: job_description
          type: string
      outputs:
        - name: parsed_jd
          type: object

    - name: generate-resume
      agent: resume-generator
      depends_on: [analyze-jd]
      inputs:
        - name: master_profile
          type: object
        - name: parsed_jd
          type: object
          from: analyze-jd.parsed_jd
      outputs:
        - name: tailored_resume
          type: string

    - name: generate-cover-letter
      agent: cover-letter-generator
      depends_on: [analyze-jd, generate-resume]
      inputs:
        - name: master_profile
          type: object
        - name: parsed_jd
          type: object
          from: analyze-jd.parsed_jd
        - name: resume
          type: string
          from: generate-resume.tailored_resume
      outputs:
        - name: cover_letter
          type: string

    - name: evaluate-resume
      agent: resume-evaluator
      depends_on: [analyze-jd, generate-resume]
      inputs:
        - name: resume
          type: string
          from: generate-resume.tailored_resume
        - name: parsed_jd
          type: object
          from: analyze-jd.parsed_jd
      outputs:
        - name: evaluation_report
          type: object

    - name: refine-resume
      agent: resume-refiner
      depends_on: [generate-resume, evaluate-resume]
      inputs:
        - name: resume
          type: string
          from: generate-resume.tailored_resume
        - name: evaluation
          type: object
          from: evaluate-resume.evaluation_report
      outputs:
        - name: refined_resume
          type: string
```

### LLM-as-a-Judge Integration

The `resume-evaluator` agent uses the structured-evaluation framework to produce consistent, actionable feedback:

1. **Rubric-Based Scoring**: Each evaluation category has defined score anchors (0-2.9 Poor, 3-4.9 Weak, 5-6.9 Adequate, 7-8.9 Good, 9-10 Excellent)
2. **Findings Generation**: Specific gaps identified with severity, recommendations, and effort estimates
3. **Pass/Fail Decision**: Based on configurable criteria (max critical findings, minimum score)
4. **Judge Metadata**: Model, temperature, prompt version tracked for reproducibility

### Iterative Refinement Loop

The pipeline supports iterative improvement:

1. Generate initial resume
2. Evaluate against JD
3. If evaluation fails or is conditional:
   - Refiner agent addresses findings
   - Re-evaluate refined resume
4. Repeat until passing or max iterations reached

## Out of Scope (v1.0)

- Web-based UI (planned for v2.0)
- Multi-user collaboration
- Real-time job board integration
- Automated application submission
- Video interview analysis

## Appendix

### A. STAR Experience Format

```json
{
  "id": "uuid",
  "role": "Senior Product Manager",
  "company": "Acme Corp",
  "situation": "Legacy API platform had 40% customer churn due to poor DX",
  "task": "Redesign API platform to improve developer adoption",
  "action": "Led cross-functional team to implement OpenAPI-first design...",
  "result": "Reduced churn by 60%, increased API calls by 300%",
  "skills": ["API Design", "Product Strategy", "Developer Experience"],
  "tags": ["platform", "devx", "api"],
  "metrics": {
    "churn_reduction": "60%",
    "api_growth": "300%"
  }
}
```

### B. Evaluation Report Example

```json
{
  "review_type": "resume-jd-match",
  "weighted_score": 7.8,
  "categories": [
    {"category": "skills_alignment", "score": 8.5, "weight": 0.25},
    {"category": "experience_level", "score": 7.0, "weight": 0.25},
    {"category": "achievement_clarity", "score": 8.0, "weight": 0.20},
    {"category": "industry_fit", "score": 7.5, "weight": 0.15},
    {"category": "communication", "score": 8.0, "weight": 0.15}
  ],
  "decision": {
    "status": "conditional",
    "rationale": "Strong match with minor skill gaps in required Python experience"
  }
}
```
