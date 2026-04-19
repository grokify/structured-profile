package service

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/grokify/structured-profile/matcher"
	"github.com/grokify/structured-profile/schema"
	"github.com/grokify/structured-profile/store"
)

// CoverLetterService handles cover letter generation.
type CoverLetterService struct {
	store   store.Store
	matcher *matcher.Matcher
}

// NewCoverLetterService creates a new CoverLetterService.
func NewCoverLetterService(s store.Store) *CoverLetterService {
	return &CoverLetterService{
		store:   s,
		matcher: matcher.New(),
	}
}

// NewCoverLetterServiceWithMatcher creates a CoverLetterService with a custom matcher.
func NewCoverLetterServiceWithMatcher(s store.Store, m *matcher.Matcher) *CoverLetterService {
	return &CoverLetterService{
		store:   s,
		matcher: m,
	}
}

// GenerateCoverLetterInput contains input for cover letter generation.
type GenerateCoverLetterInput struct {
	ProfileID     string
	OpportunityID string                      // Optional: opportunity with JD
	TemplateID    string                      // Optional: use a saved template
	Template      *schema.CoverLetterTemplate // Optional: provide template directly
	NumSTAR       int                         // Number of STAR achievements to include (default: 3)
	Domain        string                      // Optional: domain filter for achievements

	// JD override fields (used when OpportunityID is not provided)
	JDOverride     *schema.JobDescParsed // Optional: provide JD directly
	TargetCompany  string                // Company name (used with JDOverride)
	TargetPosition string                // Position title (used with JDOverride)
	HiringManager  string                // Hiring manager name (used with JDOverride)
}

// GenerateCoverLetterResult contains the result of cover letter generation.
type GenerateCoverLetterResult struct {
	CoverLetter *schema.CoverLetter
	MatchResult *matcher.MatchResult
}

// TemplateData contains data available for template substitution.
type TemplateData struct {
	// Company and position info
	Company        string
	Position       string
	HiringManager  string

	// Profile info
	Name           string
	Email          string
	Phone          string

	// STAR achievements (indexed)
	STAR           []string  // All selected STAR strings
	STAR1          string    // First STAR
	STAR2          string    // Second STAR
	STAR3          string    // Third STAR
	STAR4          string    // Fourth STAR
	STAR5          string    // Fifth STAR

	// Skills
	MatchedSkills  []string  // Skills that matched the JD
	TopSkills      string    // Comma-separated top skills

	// Additional context
	YearsExperience int
	Domain          string
}

// Generate creates a tailored cover letter.
func (cls *CoverLetterService) Generate(ctx context.Context, input GenerateCoverLetterInput) (*GenerateCoverLetterResult, error) {
	// Load profile
	profile, err := cls.store.GetFullProfile(ctx, input.ProfileID)
	if err != nil {
		return nil, fmt.Errorf("failed to load profile: %w", err)
	}

	// Build opportunity from input (either load or create from override)
	var opp *schema.Opportunity
	var jd *schema.JobDescParsed

	if input.OpportunityID != "" {
		// Load opportunity from store
		opp, err = cls.store.GetOpportunity(ctx, input.ProfileID, input.OpportunityID)
		if err != nil {
			return nil, fmt.Errorf("failed to load opportunity: %w", err)
		}
		jd = opp.JobDescParsed
	} else if input.JDOverride != nil {
		// Create virtual opportunity from override
		opp = &schema.Opportunity{
			Company:       input.TargetCompany,
			Position:      input.TargetPosition,
			HiringManager: input.HiringManager,
			JobDescParsed: input.JDOverride,
		}
		jd = input.JDOverride
	} else {
		return nil, fmt.Errorf("either OpportunityID or JDOverride is required")
	}

	// Get template
	tmpl, err := cls.resolveTemplate(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve template: %w", err)
	}

	// Run matching to rank achievements
	var matchResult *matcher.MatchResult
	if jd != nil {
		matchResult = cls.matcher.Match(profile, jd)
	}

	// Build template data
	numSTAR := input.NumSTAR
	if numSTAR == 0 {
		numSTAR = 3
	}
	data := cls.buildTemplateData(profile, opp, matchResult, input.Domain, numSTAR)

	// Generate cover letter
	cl := schema.NewCoverLetter(input.ProfileID, input.OpportunityID)
	cl.TargetCompany = opp.Company
	cl.TargetPosition = opp.Position
	cl.HiringManager = opp.HiringManager

	// Render template sections
	if tmpl != nil {
		cl.Opening, err = cls.renderSection(tmpl.OpeningTemplate, data)
		if err != nil {
			return nil, fmt.Errorf("failed to render opening: %w", err)
		}

		cl.Body, err = cls.renderSection(tmpl.BodyTemplate, data)
		if err != nil {
			return nil, fmt.Errorf("failed to render body: %w", err)
		}

		cl.Closing, err = cls.renderSection(tmpl.ClosingTemplate, data)
		if err != nil {
			return nil, fmt.Errorf("failed to render closing: %w", err)
		}
	} else {
		// Use default content
		cl.Opening, cl.Body, cl.Closing = cls.generateDefaultContent(data)
	}

	// Track which STAR achievements were used
	cl.STARRefs = cls.getSTARRefs(profile, matchResult, input.Domain, numSTAR)

	return &GenerateCoverLetterResult{
		CoverLetter: cl,
		MatchResult: matchResult,
	}, nil
}

// resolveTemplate gets the template to use for generation.
func (cls *CoverLetterService) resolveTemplate(ctx context.Context, input GenerateCoverLetterInput) (*schema.CoverLetterTemplate, error) {
	// Direct template takes priority
	if input.Template != nil {
		return input.Template, nil
	}

	// Load template by ID
	if input.TemplateID != "" {
		tmpl, err := cls.store.GetCoverLetterTemplate(ctx, input.ProfileID, input.TemplateID)
		if err != nil {
			return nil, err
		}
		return tmpl, nil
	}

	// No template - will use defaults
	return nil, nil
}

// buildTemplateData creates the data structure for template rendering.
func (cls *CoverLetterService) buildTemplateData(
	profile *schema.FullProfile,
	opp *schema.Opportunity,
	matchResult *matcher.MatchResult,
	domain string,
	numSTAR int,
) *TemplateData {
	data := &TemplateData{
		Company:       opp.Company,
		Position:      opp.Position,
		HiringManager: opp.HiringManager,
		Name:          profile.Profile.Name,
		Email:         profile.Profile.Email,
		Phone:         profile.Profile.Phone,
		Domain:        domain,
		STAR:          []string{},
	}

	// Get top STAR achievements
	starStrings := cls.getTopSTARStrings(profile, matchResult, domain, numSTAR)
	data.STAR = starStrings

	// Set indexed STAR values
	if len(starStrings) > 0 {
		data.STAR1 = starStrings[0]
	}
	if len(starStrings) > 1 {
		data.STAR2 = starStrings[1]
	}
	if len(starStrings) > 2 {
		data.STAR3 = starStrings[2]
	}
	if len(starStrings) > 3 {
		data.STAR4 = starStrings[3]
	}
	if len(starStrings) > 4 {
		data.STAR5 = starStrings[4]
	}

	// Get matched skills
	if matchResult != nil {
		data.MatchedSkills = append(matchResult.MatchedRequiredSkills, matchResult.MatchedPreferredSkills...)
		if len(data.MatchedSkills) > 5 {
			data.TopSkills = strings.Join(data.MatchedSkills[:5], ", ")
		} else {
			data.TopSkills = strings.Join(data.MatchedSkills, ", ")
		}
	}

	// Calculate years of experience
	data.YearsExperience = cls.calculateYearsExperience(profile)

	return data
}

// getTopSTARStrings returns the top N STAR achievement strings.
func (cls *CoverLetterService) getTopSTARStrings(
	profile *schema.FullProfile,
	matchResult *matcher.MatchResult,
	domain string,
	n int,
) []string {
	var starStrings []string

	if matchResult != nil {
		// Use ranked achievements from matcher
		for _, ra := range matchResult.RankedAchievements {
			if len(starStrings) >= n {
				break
			}
			star := ra.Achievement.STARString()
			if star != "" {
				starStrings = append(starStrings, star)
			}
		}
	} else {
		// Fall back to collecting achievements from positions
		for _, tenure := range profile.Tenures {
			if len(starStrings) >= n {
				break
			}
			for _, position := range tenure.Positions {
				if len(starStrings) >= n {
					break
				}
				// Use domain-aware achievement selection if domain specified
				var achievements []schema.Achievement
				if domain != "" {
					achievements = position.AchievementsForDomain(domain)
				} else {
					achievements = position.Achievements
				}
				for _, a := range achievements {
					if len(starStrings) >= n {
						break
					}
					star := a.STARString()
					if star != "" {
						starStrings = append(starStrings, star)
					}
				}
			}
		}
	}

	return starStrings
}

// getSTARRefs returns achievement IDs for the top STAR achievements.
func (cls *CoverLetterService) getSTARRefs(
	profile *schema.FullProfile,
	matchResult *matcher.MatchResult,
	domain string,
	n int,
) []string {
	var refs []string

	if matchResult != nil {
		for _, ra := range matchResult.RankedAchievements {
			if len(refs) >= n {
				break
			}
			refs = append(refs, ra.Achievement.ID)
		}
	} else {
		// Collect achievement IDs from positions
		for _, tenure := range profile.Tenures {
			if len(refs) >= n {
				break
			}
			for _, position := range tenure.Positions {
				if len(refs) >= n {
					break
				}
				var achievements []schema.Achievement
				if domain != "" {
					achievements = position.AchievementsForDomain(domain)
				} else {
					achievements = position.Achievements
				}
				for _, a := range achievements {
					if len(refs) >= n {
						break
					}
					refs = append(refs, a.ID)
				}
			}
		}
	}

	return refs
}

// calculateYearsExperience estimates total years of professional experience.
func (cls *CoverLetterService) calculateYearsExperience(profile *schema.FullProfile) int {
	if len(profile.Tenures) == 0 {
		return 0
	}

	// Find earliest start date
	var earliest schema.Date
	for _, tenure := range profile.Tenures {
		if earliest.Year == 0 || tenure.StartDate.Year < earliest.Year {
			earliest = tenure.StartDate
		}
	}

	if earliest.Year == 0 {
		return 0
	}

	// Approximate years from earliest start to now (assuming 2024)
	return 2024 - earliest.Year
}

// renderSection renders a template section with the given data.
func (cls *CoverLetterService) renderSection(templateText string, data *TemplateData) (string, error) {
	if templateText == "" {
		return "", nil
	}

	tmpl, err := template.New("section").Parse(templateText)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return strings.TrimSpace(buf.String()), nil
}

// generateDefaultContent creates default cover letter content.
func (cls *CoverLetterService) generateDefaultContent(data *TemplateData) (opening, body, closing string) {
	// Default opening
	greeting := "Hiring Manager"
	if data.HiringManager != "" {
		greeting = data.HiringManager
	}
	opening = fmt.Sprintf(
		"I am writing to express my strong interest in the %s position at %s. "+
			"With %d years of experience in software development, I am confident I can make valuable contributions to your team.",
		data.Position, data.Company, data.YearsExperience,
	)

	// Default body with STAR achievements
	var bodyParts []string
	bodyParts = append(bodyParts, "My background includes significant accomplishments that align with your needs:")

	for i, star := range data.STAR {
		if i >= 3 {
			break
		}
		bodyParts = append(bodyParts, fmt.Sprintf("• %s", star))
	}

	if len(data.MatchedSkills) > 0 {
		bodyParts = append(bodyParts, fmt.Sprintf(
			"\nMy technical expertise includes %s, which directly matches your requirements.",
			data.TopSkills,
		))
	}

	body = strings.Join(bodyParts, "\n")

	// Default closing
	closing = fmt.Sprintf(
		"I am excited about the opportunity to bring my skills and experience to %s. "+
			"I would welcome the chance to discuss how I can contribute to your team. "+
			"Thank you for considering my application, %s.",
		data.Company, greeting,
	)

	return opening, body, closing
}

// SaveCoverLetter saves a generated cover letter to the store.
func (cls *CoverLetterService) SaveCoverLetter(ctx context.Context, cl *schema.CoverLetter) error {
	return cls.store.SaveCoverLetter(ctx, cl.ProfileID, cl)
}

// GetCoverLetter retrieves a cover letter by ID.
func (cls *CoverLetterService) GetCoverLetter(ctx context.Context, profileID, coverLetterID string) (*schema.CoverLetter, error) {
	return cls.store.GetCoverLetter(ctx, profileID, coverLetterID)
}

// ListCoverLetters lists all cover letters for a profile.
func (cls *CoverLetterService) ListCoverLetters(ctx context.Context, profileID string) ([]schema.CoverLetter, error) {
	return cls.store.ListCoverLetters(ctx, profileID)
}

// DeleteCoverLetter deletes a cover letter.
func (cls *CoverLetterService) DeleteCoverLetter(ctx context.Context, profileID, coverLetterID string) error {
	return cls.store.DeleteCoverLetter(ctx, profileID, coverLetterID)
}

// DefaultCoverLetterTemplate returns a default template for cover letters.
func DefaultCoverLetterTemplate() *schema.CoverLetterTemplate {
	return &schema.CoverLetterTemplate{
		BaseEntity:  schema.NewBaseEntity(),
		Name:        "Default",
		Description: "Standard professional cover letter template",
		OpeningTemplate: `I am writing to express my strong interest in the {{.Position}} position at {{.Company}}. ` +
			`With {{.YearsExperience}} years of experience in software development, I am confident I can make valuable contributions to your team.`,
		BodyTemplate: `My background includes significant accomplishments that align with your needs:

• {{.STAR1}}
{{if .STAR2}}• {{.STAR2}}{{end}}
{{if .STAR3}}• {{.STAR3}}{{end}}
{{if .TopSkills}}
My technical expertise includes {{.TopSkills}}, which directly matches your requirements.{{end}}`,
		ClosingTemplate: `I am excited about the opportunity to bring my skills and experience to {{.Company}}. ` +
			`I would welcome the chance to discuss how I can contribute to your team. Thank you for considering my application.`,
	}
}
