// Package service provides business logic for resume and cover letter generation.
package service

import (
	"context"
	"fmt"

	"github.com/grokify/structured-profile/matcher"
	"github.com/grokify/structured-profile/schema"
	"github.com/grokify/structured-profile/store"
)

// ResumeService handles resume generation.
type ResumeService struct {
	store   store.Store
	matcher *matcher.Matcher
}

// NewResumeService creates a new ResumeService.
func NewResumeService(s store.Store) *ResumeService {
	return &ResumeService{
		store:   s,
		matcher: matcher.New(),
	}
}

// NewResumeServiceWithMatcher creates a ResumeService with a custom matcher.
func NewResumeServiceWithMatcher(s store.Store, m *matcher.Matcher) *ResumeService {
	return &ResumeService{
		store:   s,
		matcher: m,
	}
}

// GenerateInput contains input for resume generation.
type GenerateInput struct {
	ProfileID     string
	OpportunityID string              // Optional: if provided, uses JD matching
	Domain        string              // Optional: domain filter (e.g., "devx", "iam")
	Options       *schema.ResumeOptions
	JDOverride    *schema.JobDescParsed // Optional: override JD instead of loading from opportunity
}

// GenerateResult contains the result of resume generation.
type GenerateResult struct {
	Resume      *schema.Resume
	MatchResult *matcher.MatchResult
}

// Generate creates a tailored resume.
func (rs *ResumeService) Generate(ctx context.Context, input GenerateInput) (*GenerateResult, error) {
	// Load profile
	profile, err := rs.store.GetFullProfile(ctx, input.ProfileID)
	if err != nil {
		return nil, fmt.Errorf("failed to load profile: %w", err)
	}

	// Set default options
	options := input.Options
	if options == nil {
		options = schema.DefaultResumeOptions()
	}

	// Load opportunity JD if provided
	var jd *schema.JobDescParsed
	var opportunityID string

	if input.JDOverride != nil {
		jd = input.JDOverride
	} else if input.OpportunityID != "" {
		opp, err := rs.store.GetOpportunity(ctx, input.ProfileID, input.OpportunityID)
		if err != nil {
			return nil, fmt.Errorf("failed to load opportunity: %w", err)
		}
		jd = opp.JobDescParsed
		opportunityID = opp.ID
	}

	// Create resume
	resume := schema.NewResume(input.ProfileID, opportunityID)
	resume.Domain = input.Domain
	resume.Options = options

	// Generate content
	result := &GenerateResult{
		Resume: resume,
	}

	// If we have a JD, use matching
	if jd != nil {
		matchResult := rs.matcher.Match(profile, jd)
		result.MatchResult = matchResult
		resume.Content = rs.generateContentWithMatch(profile, matchResult, input.Domain, options)
	} else {
		resume.Content = rs.generateContentWithoutMatch(profile, input.Domain, options)
	}

	return result, nil
}

// generateContentWithMatch creates resume content using JD matching results.
func (rs *ResumeService) generateContentWithMatch(
	profile *schema.FullProfile,
	matchResult *matcher.MatchResult,
	domain string,
	options *schema.ResumeOptions,
) *schema.ResumeContent {
	content := &schema.ResumeContent{}

	// Header
	if options.IncludeContact {
		content.Name = profile.Profile.Name
		content.Email = profile.Profile.Email
		content.Phone = profile.Profile.Phone
		content.Location = profile.Profile.Location
		content.Links = profile.Profile.Links
	}

	// Summary
	if options.IncludeSummary {
		content.Summary = rs.selectSummary(profile, domain)
	}

	// Experience - use matched achievements
	if options.IncludeExperience {
		content.Experiences = rs.buildExperiences(profile, matchResult, domain, options)
	}

	// Skills - prioritize matched skills
	if options.IncludeSkills {
		content.Skills = rs.selectSkills(profile, matchResult, domain, options.MaxSkills)
	}

	// Education
	if options.IncludeEducation {
		content.Education = profile.Education
	}

	// Certifications
	if options.IncludeCertifications {
		content.Certifications = profile.Certifications
	}

	// Publications
	if options.IncludePublications {
		content.Publications = profile.Publications
	}

	return content
}

// generateContentWithoutMatch creates resume content without JD matching.
func (rs *ResumeService) generateContentWithoutMatch(
	profile *schema.FullProfile,
	domain string,
	options *schema.ResumeOptions,
) *schema.ResumeContent {
	content := &schema.ResumeContent{}

	// Header
	if options.IncludeContact {
		content.Name = profile.Profile.Name
		content.Email = profile.Profile.Email
		content.Phone = profile.Profile.Phone
		content.Location = profile.Profile.Location
		content.Links = profile.Profile.Links
	}

	// Summary
	if options.IncludeSummary {
		content.Summary = rs.selectSummary(profile, domain)
	}

	// Experience - use domain filtering only
	if options.IncludeExperience {
		content.Experiences = rs.buildExperiences(profile, nil, domain, options)
	}

	// Skills
	if options.IncludeSkills {
		content.Skills = rs.selectSkills(profile, nil, domain, options.MaxSkills)
	}

	// Education
	if options.IncludeEducation {
		content.Education = profile.Education
	}

	// Certifications
	if options.IncludeCertifications {
		content.Certifications = profile.Certifications
	}

	// Publications
	if options.IncludePublications {
		content.Publications = profile.Publications
	}

	return content
}

// selectSummary returns the appropriate summary for the domain.
func (rs *ResumeService) selectSummary(profile *schema.FullProfile, domain string) string {
	return profile.Profile.Summaries.ForDomain(domain)
}

// buildExperiences builds resume experiences from tenures.
func (rs *ResumeService) buildExperiences(
	profile *schema.FullProfile,
	matchResult *matcher.MatchResult,
	domain string,
	options *schema.ResumeOptions,
) []schema.ResumeExperience {
	var experiences []schema.ResumeExperience

	// Build a set of top achievements if matching was used
	topAchievementIDs := make(map[string]bool)
	if matchResult != nil {
		for i, ra := range matchResult.RankedAchievements {
			if options.MaxAchievements > 0 && i >= options.MaxAchievements*3 {
				break // Keep pool of top achievements
			}
			topAchievementIDs[ra.Achievement.ID] = true
		}
	}

	for i, tenure := range profile.Tenures {
		if options.MaxExperiences > 0 && i >= options.MaxExperiences {
			break
		}

		if options.CollapseTenurePositions && tenure.CollapsedInfo != nil {
			// Collapsed view
			exp := rs.buildCollapsedExperience(tenure, domain, topAchievementIDs, options)
			experiences = append(experiences, exp)
		} else {
			// Expanded view - one experience per position
			for _, position := range tenure.Positions {
				exp := rs.buildPositionExperience(tenure, position, domain, topAchievementIDs, options)
				experiences = append(experiences, exp)
			}
		}
	}

	return experiences
}

// buildCollapsedExperience builds a single experience from a collapsed tenure.
func (rs *ResumeService) buildCollapsedExperience(
	tenure schema.Tenure,
	domain string,
	topAchievementIDs map[string]bool,
	options *schema.ResumeOptions,
) schema.ResumeExperience {
	// Get title for domain
	title := tenure.CollapsedInfo.TitleForDomain(domain)
	if title == "" && len(tenure.Positions) > 0 {
		title = tenure.Positions[0].Title
	}

	exp := schema.ResumeExperience{
		Company:   tenure.Company,
		Title:     title,
		Location:  tenure.Location,
		StartDate: tenure.StartDate,
		EndDate:   tenure.EndDate,
	}

	// Add description from first position
	if len(tenure.Positions) > 0 {
		exp.Description = tenure.Positions[0].DescriptionForDomain(options.DescriptionWithoutCounts)
	}

	// Collect achievements from all positions
	exp.Achievements = rs.selectAchievements(tenure.Positions, domain, topAchievementIDs, options)

	// Collect skills
	exp.Skills = rs.selectPositionSkills(tenure.Positions, domain)

	return exp
}

// buildPositionExperience builds an experience from a single position.
func (rs *ResumeService) buildPositionExperience(
	tenure schema.Tenure,
	position schema.Position,
	domain string,
	topAchievementIDs map[string]bool,
	options *schema.ResumeOptions,
) schema.ResumeExperience {
	exp := schema.ResumeExperience{
		Company:     tenure.Company,
		Title:       position.Title,
		Location:    tenure.Location,
		StartDate:   position.StartDate,
		EndDate:     position.EndDate,
		Description: position.DescriptionForDomain(options.DescriptionWithoutCounts),
	}

	// Get achievements for domain
	positions := []schema.Position{position}
	exp.Achievements = rs.selectAchievements(positions, domain, topAchievementIDs, options)

	// Get skills for domain
	exp.Skills = position.SkillsForDomain(domain)

	return exp
}

// selectAchievements selects achievements based on domain and match results.
func (rs *ResumeService) selectAchievements(
	positions []schema.Position,
	domain string,
	topAchievementIDs map[string]bool,
	options *schema.ResumeOptions,
) []string {
	var achievements []string
	maxPerPosition := options.MaxAchievements
	if maxPerPosition == 0 {
		maxPerPosition = 5 // Default max per position
	}

	for _, position := range positions {
		// Get achievements in domain order
		domainAchievements := position.AchievementsForDomain(domain)

		count := 0
		for _, a := range domainAchievements {
			if count >= maxPerPosition {
				break
			}

			// If we have match results, prioritize matched achievements
			if len(topAchievementIDs) > 0 && !topAchievementIDs[a.ID] {
				continue
			}

			starStr := a.STARString()
			if starStr != "" {
				achievements = append(achievements, starStr)
				count++
			}
		}

		// If we didn't get enough from matched, fill with others
		if count < maxPerPosition && len(topAchievementIDs) > 0 {
			for _, a := range domainAchievements {
				if count >= maxPerPosition {
					break
				}
				if topAchievementIDs[a.ID] {
					continue // Already added
				}
				starStr := a.STARString()
				if starStr != "" {
					achievements = append(achievements, starStr)
					count++
				}
			}
		}
	}

	return achievements
}

// selectPositionSkills collects skills from positions for a domain.
func (rs *ResumeService) selectPositionSkills(positions []schema.Position, domain string) []string {
	skillSet := make(map[string]bool)
	var skills []string

	for _, position := range positions {
		for _, skill := range position.SkillsForDomain(domain) {
			if !skillSet[skill] {
				skillSet[skill] = true
				skills = append(skills, skill)
			}
		}
	}

	return skills
}

// selectSkills selects skills for the resume, prioritizing matched skills.
func (rs *ResumeService) selectSkills(
	profile *schema.FullProfile,
	matchResult *matcher.MatchResult,
	domain string,
	maxSkills int,
) []string {
	skillSet := make(map[string]bool)
	var skills []string

	// Add matched skills first (if available)
	if matchResult != nil {
		for _, skill := range matchResult.MatchedRequiredSkills {
			if !skillSet[skill] {
				skillSet[skill] = true
				skills = append(skills, skill)
			}
		}
		for _, skill := range matchResult.MatchedPreferredSkills {
			if !skillSet[skill] {
				skillSet[skill] = true
				skills = append(skills, skill)
			}
		}
	}

	// Add profile skills
	for _, skill := range profile.Skills {
		if !skillSet[skill.Name] {
			skillSet[skill.Name] = true
			skills = append(skills, skill.Name)
		}
	}

	// Limit if specified
	if maxSkills > 0 && len(skills) > maxSkills {
		skills = skills[:maxSkills]
	}

	return skills
}
