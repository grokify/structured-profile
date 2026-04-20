# Document Categories

Standard categories for document quality evaluation defined in `schema/doceval.go`.

## Category Reference

| Category | Document | Weight | Description |
|----------|----------|--------|-------------|
| `keyword_coverage` | Both | 20% | JD keywords represented |
| `achievement_relevance` | Resume | 20% | Top achievements align with JD |
| `gap_mitigation` | Cover Letter | 15% | Known gaps addressed |
| `quantification` | Resume | 15% | Achievements include metrics |
| `ats_compatibility` | Resume | 10% | Format compatible with ATS |
| `narrative_coherence` | Both | 10% | Documents tell coherent story |
| `value_proposition` | Cover Letter | 10% | Clear unique value stated |

## Core Categories

### keyword_coverage

**Document:** Both
**Weight:** 20%

Are JD keywords and terminology represented in the documents?

| Score | Criteria |
|-------|----------|
| 9-10 | 90%+ JD keywords present, natural usage |
| 7-8.9 | 70-90% keywords, mostly natural |
| 5-6.9 | 50-70% keywords, some forced |
| 0-4.9 | <50% keywords missing |

**Common Issues:**

- Missing key technical terms
- JD uses different terminology than resume
- Over-stuffing keywords unnaturally

### achievement_relevance

**Document:** Resume
**Weight:** 20%

Do top achievements directly address JD requirements?

| Score | Criteria |
|-------|----------|
| 9-10 | Top 3 achievements directly address JD |
| 7-8.9 | Most achievements relevant, good ordering |
| 5-6.9 | Some relevant achievements buried |
| 0-4.9 | Achievements don't align with JD |

**Common Issues:**

- Best achievements not at top
- Irrelevant achievements prominent
- Missing achievements that match JD

### gap_mitigation

**Document:** Cover Letter
**Weight:** 15%

Are known gaps from matcheval addressed or mitigated?

| Score | Criteria |
|-------|----------|
| 9-10 | Gaps proactively addressed with transferable skills |
| 7-8.9 | Major gaps addressed, minor ones acceptable |
| 5-6.9 | Some gaps addressed, others ignored |
| 0-4.9 | Gaps not addressed, obvious omissions |

**Strategies:**

- Highlight transferable skills
- Show learning agility
- Connect adjacent experience
- Express enthusiasm to grow

### quantification

**Document:** Resume
**Weight:** 15%

Do achievements include metrics and measurable impact?

| Score | Criteria |
|-------|----------|
| 9-10 | 80%+ achievements have metrics (%, $, scale) |
| 7-8.9 | 60-80% quantified, meaningful metrics |
| 5-6.9 | 40-60% quantified, some vague |
| 0-4.9 | <40% quantified, mostly generic |

**Good Examples:**

- "Grew platform to 70K+ developers"
- "Reduced deployment time by 80%"
- "Generated $500M ARR"

**Poor Examples:**

- "Improved platform performance"
- "Led team successfully"
- "Delivered project on time"

### ats_compatibility

**Document:** Resume
**Weight:** 10%

Is the format compatible with Applicant Tracking Systems?

| Score | Criteria |
|-------|----------|
| 9-10 | Clean format, standard sections, no tables |
| 7-8.9 | Mostly ATS-friendly, minor issues |
| 5-6.9 | Some ATS issues (columns, headers) |
| 0-4.9 | Heavy formatting, likely ATS problems |

**ATS Best Practices:**

- Use standard section headers
- Avoid tables, columns, graphics
- Use standard fonts
- Include keywords in text, not just headers
- Save as PDF from plain source

### narrative_coherence

**Document:** Both
**Weight:** 10%

Do documents tell a coherent career story?

| Score | Criteria |
|-------|----------|
| 9-10 | Clear progression, consistent theme, compelling |
| 7-8.9 | Logical flow, minor gaps in narrative |
| 5-6.9 | Disjointed, unclear progression |
| 0-4.9 | No coherent story, random listing |

**Elements:**

- Career progression visible
- Consistent theme/focus
- Achievements build on each other
- Clear "why this role" narrative

### value_proposition

**Document:** Cover Letter
**Weight:** 10%

Is unique value clearly stated?

| Score | Criteria |
|-------|----------|
| 9-10 | Clear, specific value prop tied to JD |
| 7-8.9 | Value stated, mostly specific |
| 5-6.9 | Generic value statements |
| 0-4.9 | No clear value proposition |

**Good Example:**

> "My experience building a 70K-developer platform at RingCentral, combined with my direct payments authentication work at Arcot (selected by Visa/MasterCard), positions me uniquely for PayPal's developer platform challenges."

**Poor Example:**

> "I am a hard-working professional with many years of experience in technology."

## Resume-Specific Categories

### summary_impact

Does the summary immediately hook the reader?

### experience_ordering

Are experiences ordered by relevance, not just recency?

### skills_alignment

Does skills section mirror JD terminology?

### formatting

Is formatting clean and professional?

## Cover Letter-Specific Categories

### opening_hook

Does the opening grab attention?

### company_fit

Is connection to company needs clear?

### call_to_action

Is there a clear next step?
