# Match Categories

Standard categories for profile-to-JD matching defined in `schema/matcheval.go`.

## Category Reference

| Category | Description | Match Type | Weight |
|----------|-------------|------------|--------|
| `technical_skills` | Technical expertise | semantic | 20% |
| `domain_experience` | Industry knowledge | semantic | 15% |
| `leadership` | Team/people management | threshold | 20% |
| `years_experience` | Total relevant years | threshold | 15% |
| `platform_scale` | Scale of systems managed | threshold | 10% |
| `apis_sdks` | API/SDK experience | exact | 10% |
| `identity_security` | IAM/auth expertise | semantic | varies |
| `compliance` | Regulatory experience | semantic | 10% |
| `thought_leadership` | Industry presence | boolean | 10% |
| `developer_relations` | DevRel/advocacy | semantic | varies |
| `team_management` | Direct reports | threshold | varies |
| `cross_functional` | Cross-org collaboration | semantic | varies |

!!! note "Weights are examples"
    Actual weights vary based on job requirements. They should sum to 1.0.

## Category Definitions

### technical_skills

Technical expertise and proficiency.

**Match Type:** `semantic` - Equivalent technologies accepted

**Example Matches:**

| JD Requirement | Profile Experience | Match? |
|----------------|-------------------|--------|
| "Go experience" | "Go, Python" | ✅ Exact |
| "REST API design" | "API design, OpenAPI" | ✅ Semantic |
| "Kubernetes" | "Docker only" | ⚠️ Partial |

### domain_experience

Industry and domain knowledge.

**Match Type:** `semantic` - Related industries count

**Example Matches:**

| JD Requirement | Profile Experience | Match? |
|----------------|-------------------|--------|
| "Fintech experience" | "Payments at Visa" | ✅ Direct |
| "Enterprise SaaS" | "B2B software" | ✅ Semantic |
| "Healthcare" | "Fintech only" | ❌ Gap |

### leadership

Team and people management experience.

**Match Type:** `threshold` - Meets minimum scope

**Scoring:**

| Team Size | Score |
|-----------|-------|
| 20+ direct reports | 10 |
| 10-19 direct reports | 9 |
| 5-9 direct reports | 8 |
| 1-4 direct reports | 6 |
| No direct reports | 4 |

### years_experience

Total years of relevant experience.

**Match Type:** `threshold` - Meets or exceeds requirement

**Scoring:**

| Requirement | Actual | Score |
|-------------|--------|-------|
| 10+ years | 15 years | 10 (exceeds) |
| 10+ years | 10 years | 8 (meets) |
| 10+ years | 8 years | 6 (close) |
| 10+ years | 5 years | 3 (below) |

### platform_scale

Scale of systems managed.

**Match Type:** `threshold`

**Indicators:**

- Users/developers served
- ARR managed
- Engineers supported
- API calls/day

### apis_sdks

API and SDK development experience.

**Match Type:** `exact` - Specific technologies required

**Indicators:**

- REST API design
- OpenAPI/Swagger
- SDK development
- Developer portals

### identity_security

Identity and access management expertise.

**Match Type:** `semantic`

**Indicators:**

- OAuth 2.0, OIDC, SAML
- MFA/2FA implementation
- RBAC/PBAC
- SSO integration

### compliance

Regulatory and compliance experience.

**Match Type:** `semantic`

**Frameworks:**

- HIPAA (healthcare)
- GDPR/CCPA (privacy)
- SOC 2 (security)
- PCI-DSS (payments)
- FINRA/SEC (finance)

### thought_leadership

Industry presence and recognition.

**Match Type:** `boolean` - Present or absent

**Indicators:**

- Patents
- Publications
- Conference speaking
- Standards body participation
- Industry awards

### developer_relations

Developer advocacy and community building.

**Match Type:** `semantic`

**Indicators:**

- Developer advocacy team
- Technical content creation
- Community management
- Stack Overflow presence
- Open source contributions

## Match Types

| Type | Description | Example |
|------|-------------|---------|
| `exact` | Specific match required | "Python required" |
| `semantic` | Equivalent accepted | "DevRel" ≈ "Technical Enablement" |
| `partial` | Related but not identical | "Go" for "Python" role |
| `threshold` | Meets minimum | "10+ years" ≥ 10 |
| `boolean` | Present/absent | "Has patents" |
| `weighted` | Combined factors | Leadership from multiple inputs |

## Scoring Guide

| Score | Status | Meaning |
|-------|--------|---------|
| 9-10 | `pass` | Exceeds requirements |
| 7-8.9 | `pass` | Meets requirements |
| 5-6.9 | `warn` | Partial match, addressable |
| 0-4.9 | `fail` | Significant gap |
