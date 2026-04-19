// Package jdparser provides job description parsing functionality.
// It extracts structured information from raw job description text.
package jdparser

import (
	"regexp"
	"strings"

	"github.com/grokify/structured-profile/schema"
)

// Parser extracts structured data from job descriptions.
type Parser struct {
	// TechSkills is the dictionary of known technical skills to match.
	TechSkills []string

	// SoftSkills is the dictionary of known soft skills to match.
	SoftSkills []string

	// SeniorityKeywords maps seniority levels to their keywords.
	SeniorityKeywords map[schema.SeniorityLevel][]string
}

// New creates a new Parser with default skill dictionaries.
func New() *Parser {
	return &Parser{
		TechSkills:        DefaultTechSkills(),
		SoftSkills:        DefaultSoftSkills(),
		SeniorityKeywords: DefaultSeniorityKeywords(),
	}
}

// Parse extracts structured data from a job description.
func (p *Parser) Parse(rawJD string) *schema.JobDescParsed {
	if rawJD == "" {
		return nil
	}

	normalized := normalizeText(rawJD)
	lines := strings.Split(normalized, "\n")

	result := &schema.JobDescParsed{}

	// Extract keywords
	result.Keywords = p.extractKeywords(normalized)

	// Extract skills (categorized as required/preferred)
	result.RequiredSkills, result.PreferredSkills = p.extractSkills(normalized, lines)

	// Extract experience years
	result.ExperienceYears = p.extractExperienceYears(normalized)

	// Extract seniority level
	result.SeniorityLevel = string(p.extractSeniorityLevel(normalized))

	// Extract responsibilities and qualifications
	result.Responsibilities = p.extractSection(lines, responsibilitySectionHeaders)
	result.Qualifications = p.extractSection(lines, qualificationSectionHeaders)

	// Extract team info
	result.TeamSize = p.extractTeamSize(normalized)
	result.ReportingTo = p.extractReportingTo(normalized)

	return result
}

// extractKeywords extracts important keywords from the job description.
func (p *Parser) extractKeywords(text string) []string {
	keywords := make(map[string]bool)

	// Extract tech skills as keywords
	for _, skill := range p.TechSkills {
		if containsWord(text, skill) {
			keywords[skill] = true
		}
	}

	// Extract domain-specific keywords
	for _, kw := range domainKeywords {
		if containsWord(text, kw) {
			keywords[kw] = true
		}
	}

	return mapKeys(keywords)
}

// extractSkills extracts required and preferred skills from the job description.
func (p *Parser) extractSkills(text string, lines []string) (required, preferred []string) {
	requiredSet := make(map[string]bool)
	preferredSet := make(map[string]bool)

	// Find required/preferred sections
	inPreferredSection := false

	for _, line := range lines {
		lower := strings.ToLower(line)

		// Check for section headers
		if containsAny(lower, []string{"required", "must have", "requirements", "qualifications"}) {
			inPreferredSection = false
		} else if containsAny(lower, []string{"preferred", "nice to have", "bonus", "plus", "desired"}) {
			inPreferredSection = true
		}

		// Extract skills from line
		for _, skill := range p.TechSkills {
			if containsWord(lower, strings.ToLower(skill)) {
				if inPreferredSection {
					preferredSet[skill] = true
				} else {
					requiredSet[skill] = true
				}
			}
		}

		for _, skill := range p.SoftSkills {
			if containsWord(lower, strings.ToLower(skill)) {
				if inPreferredSection {
					preferredSet[skill] = true
				} else {
					requiredSet[skill] = true
				}
			}
		}
	}

	// If no section-based extraction worked, scan the whole text
	if len(requiredSet) == 0 && len(preferredSet) == 0 {
		for _, skill := range p.TechSkills {
			if containsWord(text, strings.ToLower(skill)) {
				requiredSet[skill] = true
			}
		}
		for _, skill := range p.SoftSkills {
			if containsWord(text, strings.ToLower(skill)) {
				requiredSet[skill] = true
			}
		}
	}

	return mapKeys(requiredSet), mapKeys(preferredSet)
}

// extractExperienceYears extracts years of experience required.
func (p *Parser) extractExperienceYears(text string) int {
	// Pattern: "X+ years", "X years", "X-Y years"
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(\d+)\+?\s*(?:years?|yrs?)(?:\s+of)?\s+(?:experience|exp)`),
		regexp.MustCompile(`(?:experience|exp)(?:\s+of)?\s*:?\s*(\d+)\+?\s*(?:years?|yrs?)`),
		regexp.MustCompile(`(\d+)\s*-\s*\d+\s*(?:years?|yrs?)`),
		regexp.MustCompile(`minimum\s+(?:of\s+)?(\d+)\s*(?:years?|yrs?)`),
		regexp.MustCompile(`at\s+least\s+(\d+)\s*(?:years?|yrs?)`),
	}

	lower := strings.ToLower(text)
	for _, pattern := range patterns {
		matches := pattern.FindStringSubmatch(lower)
		if len(matches) >= 2 {
			years := 0
			for _, c := range matches[1] {
				if c >= '0' && c <= '9' {
					years = years*10 + int(c-'0')
				}
			}
			if years > 0 && years < 50 { // Sanity check
				return years
			}
		}
	}

	return 0
}

// extractSeniorityLevel determines the seniority level from the job description.
func (p *Parser) extractSeniorityLevel(text string) schema.SeniorityLevel {
	lower := strings.ToLower(text)

	// Check from highest to lowest seniority
	seniorityOrder := []schema.SeniorityLevel{
		schema.SeniorityExecutive,
		schema.SeniorityVP,
		schema.SeniorityDirector,
		schema.SeniorityManager,
		schema.SeniorityLead,
		schema.SeniorityPrincipal,
		schema.SeniorityStaff,
		schema.SenioritySenior,
		schema.SeniorityMid,
		schema.SeniorityEntry,
	}

	for _, level := range seniorityOrder {
		keywords := p.SeniorityKeywords[level]
		for _, kw := range keywords {
			if containsWord(lower, kw) {
				return level
			}
		}
	}

	// Infer from experience years
	years := p.extractExperienceYears(text)
	switch {
	case years >= 10:
		return schema.SenioritySenior
	case years >= 5:
		return schema.SeniorityMid
	case years >= 2:
		return schema.SeniorityMid
	case years > 0:
		return schema.SeniorityEntry
	}

	return ""
}

// extractSection extracts bullet points from a named section.
func (p *Parser) extractSection(lines []string, headers []string) []string {
	var results []string
	inSection := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)

		// Check if this is a section header
		isHeader := false
		for _, h := range headers {
			if strings.Contains(lower, h) && len(trimmed) < 80 {
				isHeader = true
				inSection = true
				break
			}
		}

		if isHeader {
			continue
		}

		// Check if we've hit a new section
		if inSection && isSectionHeader(trimmed) {
			break
		}

		// Collect bullet points
		if inSection && isBulletPoint(trimmed) {
			item := cleanBulletPoint(trimmed)
			if len(item) > 10 { // Skip very short items
				results = append(results, item)
			}
		}
	}

	return results
}

// extractTeamSize extracts team size information.
func (p *Parser) extractTeamSize(text string) string {
	lower := strings.ToLower(text)

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`team\s+of\s+(\d+)`),
		regexp.MustCompile(`(\d+)\s*(?:person|people|member)\s+team`),
		regexp.MustCompile(`manage\s+(?:a\s+)?team\s+of\s+(\d+)`),
	}

	for _, pattern := range patterns {
		matches := pattern.FindStringSubmatch(lower)
		if len(matches) >= 2 {
			return matches[1]
		}
	}

	return ""
}

// extractReportingTo extracts reporting structure.
func (p *Parser) extractReportingTo(text string) string {
	lower := strings.ToLower(text)

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`report(?:s|ing)?\s+(?:to|into)\s+(?:the\s+)?([a-z\s]+?)(?:\.|,|$)`),
	}

	for _, pattern := range patterns {
		matches := pattern.FindStringSubmatch(lower)
		if len(matches) >= 2 {
			title := strings.TrimSpace(matches[1])
			if len(title) > 3 && len(title) < 50 {
				return title
			}
		}
	}

	return ""
}

// Section header lists
var responsibilitySectionHeaders = []string{
	"responsibilities",
	"what you'll do",
	"what you will do",
	"your role",
	"job duties",
	"key responsibilities",
	"core responsibilities",
	"day to day",
	"in this role",
}

var qualificationSectionHeaders = []string{
	"qualifications",
	"requirements",
	"what you'll need",
	"what you will need",
	"who you are",
	"about you",
	"skills and experience",
	"minimum qualifications",
	"basic qualifications",
}

// Domain keywords for extraction
var domainKeywords = []string{
	"API", "REST", "GraphQL", "microservices", "distributed systems",
	"cloud", "AWS", "GCP", "Azure", "Kubernetes", "Docker",
	"CI/CD", "DevOps", "SRE", "observability", "monitoring",
	"security", "authentication", "authorization", "IAM", "SSO", "OAuth",
	"data pipeline", "ETL", "data warehouse", "analytics",
	"machine learning", "ML", "AI", "NLP",
	"frontend", "backend", "full-stack", "fullstack",
	"mobile", "iOS", "Android", "React Native", "Flutter",
	"agile", "scrum", "kanban",
	"startup", "enterprise", "scale", "growth",
}
