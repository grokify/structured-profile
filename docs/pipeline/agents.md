# AI Agents

The AI agent pipeline provides semantic matching and generation capabilities beyond keyword-based approaches.

## Agent Pipeline

```
┌─────────────┐    ┌─────────────────┐    ┌──────────────────┐
│ jd-analyst  │───▶│ profile-matcher │───▶│ resume-generator │
└─────────────┘    └─────────────────┘    └──────────────────┘
                                                   │
                                                   ▼
┌─────────────────┐    ┌──────────────────┐    ┌──────────────────┐
│ resume-refiner  │◀───│ resume-evaluator │◀───│ (resume draft)   │
└─────────────────┘    └──────────────────┘    └──────────────────┘
```

## Agent Specifications

Agent specs are defined in `agents/specs/agents/`:

### jd-analyst

**Purpose:** Semantic job description parsing

**Capabilities:**

- Extract role title, level, and type (IC/management)
- Identify required and preferred skills
- Create semantic mappings (e.g., "GTM" → "go-to-market")
- Extract company context and culture signals

**Model:** haiku (fast, cost-effective)

### profile-matcher

**Purpose:** Profile-to-JD matching with scoring

**Capabilities:**

- Score achievements against requirements
- Identify exact and semantic skill matches
- Perform gap analysis
- Recommend positioning strategy

**Scoring Criteria:**

| Factor | Points |
|--------|--------|
| Direct skill match | 40 |
| Semantic skill match | 25 |
| Impact alignment | 20 |
| Recency bonus | 10 |
| Quantification bonus | 5 |

**Model:** sonnet (balanced quality/speed)

### resume-generator

**Purpose:** Tailored resume creation

**Capabilities:**

- Select and order relevant achievements
- Craft role-specific summary
- Mirror JD terminology authentically
- Format for ATS compatibility

**Model:** sonnet

### resume-evaluator

**Purpose:** LLM-as-Judge quality scoring

**Capabilities:**

- Score keyword coverage
- Evaluate achievement relevance
- Check quantification
- Assess narrative coherence
- Identify improvement areas

**Model:** sonnet

### resume-refiner

**Purpose:** Feedback-based improvement

**Capabilities:**

- Address evaluator feedback
- Improve weak sections
- Enhance quantification
- Strengthen narrative

**Model:** sonnet

### resume-coordinator

**Purpose:** Pipeline orchestration

**Capabilities:**

- Load profile and JD
- Spawn and coordinate agents
- Pass context between stages
- Handle errors and reporting

**Model:** sonnet

## Using Agents with Claude Code

```markdown
# Example prompt to resume-coordinator

Generate a tailored resume for the PayPal Sr. Director role.

Profile: ~/.sprofile/profiles/john-doe.json
JD: applications/app_2026-04-19_paypal_srdir/jobdescription.md

Run the full pipeline:
1. Analyze JD with jd-analyst
2. Match profile with profile-matcher
3. Generate resume with resume-generator
4. Evaluate with resume-evaluator
5. Refine if score < 85
6. Output final resume.md
```

## Agent Configuration

Agents are configured in `agents/specs/deployments/local.json`:

```json
{
  "deployment": "local",
  "agents": {
    "jd-analyst": {
      "model": "haiku",
      "temperature": 0.3
    },
    "profile-matcher": {
      "model": "sonnet",
      "temperature": 0.3
    },
    "resume-generator": {
      "model": "sonnet",
      "temperature": 0.5
    }
  }
}
```

## Semantic Matching Examples

The agent pipeline excels at semantic understanding:

| JD Says | Profile Has | Match Type |
|---------|-------------|------------|
| "technical enablement" | "developer advocacy" | Semantic ✓ |
| "AI coding tools" | "GitHub Copilot evaluation" | Semantic ✓ |
| "GTM partnership" | "sales engineering" | Semantic ✓ |
| "10+ years experience" | "15 years experience" | Threshold ✓ |
| "Python required" | "Go expertise" | Partial (flag) |

## Keyword vs Semantic Matching

| Approach | Match Rate | False Negatives |
|----------|------------|-----------------|
| Keyword only | ~17% | High (misses synonyms) |
| Semantic | ~90% | Low |

The semantic approach recognizes equivalent skills and experiences that keyword matching would miss.
