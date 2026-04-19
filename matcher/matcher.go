// Package matcher provides profile-to-job-description matching functionality.
// It scores and ranks profile achievements based on relevance to job requirements.
package matcher

import (
	"sort"
	"strings"

	"github.com/grokify/structured-profile/schema"
)

// Matcher scores and ranks profile content against job descriptions.
type Matcher struct {
	// Weights for different matching criteria
	Weights MatchWeights
}

// MatchWeights configures the relative importance of different matching criteria.
type MatchWeights struct {
	RequiredSkill   float64 // Weight for matching required skills
	PreferredSkill  float64 // Weight for matching preferred skills
	Keyword         float64 // Weight for matching keywords
	Tag             float64 // Weight for matching tags
	ExperienceBonus float64 // Bonus for matching experience years
}

// DefaultWeights returns sensible default weights for matching.
func DefaultWeights() MatchWeights {
	return MatchWeights{
		RequiredSkill:   10.0,
		PreferredSkill:  5.0,
		Keyword:         3.0,
		Tag:             2.0,
		ExperienceBonus: 5.0,
	}
}

// New creates a new Matcher with default weights.
func New() *Matcher {
	return &Matcher{
		Weights: DefaultWeights(),
	}
}

// NewWithWeights creates a new Matcher with custom weights.
func NewWithWeights(weights MatchWeights) *Matcher {
	return &Matcher{
		Weights: weights,
	}
}

// MatchResult contains the result of matching a profile against a job description.
type MatchResult struct {
	// Ranked achievements by relevance
	RankedAchievements []RankedAchievement

	// Skills from the profile that match the JD
	MatchedRequiredSkills  []string
	MatchedPreferredSkills []string

	// Keywords from the JD that appear in the profile
	MatchedKeywords []string

	// Overall match score (0-100)
	OverallScore float64

	// Gap analysis
	MissingRequiredSkills  []string
	MissingPreferredSkills []string
}

// RankedAchievement contains an achievement with its relevance score.
type RankedAchievement struct {
	Achievement schema.Achievement
	Score       float64
	MatchedSkills   []string
	MatchedKeywords []string
	MatchedTags     []string
}

// Match scores a profile against a job description.
func (m *Matcher) Match(profile *schema.FullProfile, jd *schema.JobDescParsed) *MatchResult {
	if profile == nil || jd == nil {
		return &MatchResult{}
	}

	result := &MatchResult{}

	// Build skill and keyword sets from the profile
	profileSkills := m.extractProfileSkills(profile)
	profileKeywords := m.extractProfileKeywords(profile)

	// Find matched and missing skills
	result.MatchedRequiredSkills, result.MissingRequiredSkills = m.matchStrings(profileSkills, jd.RequiredSkills)
	result.MatchedPreferredSkills, result.MissingPreferredSkills = m.matchStrings(profileSkills, jd.PreferredSkills)
	result.MatchedKeywords, _ = m.matchStrings(profileKeywords, jd.Keywords)

	// Score and rank achievements
	achievements := profile.AllAchievements()
	result.RankedAchievements = m.rankAchievements(achievements, jd)

	// Calculate overall score
	result.OverallScore = m.calculateOverallScore(result, jd)

	return result
}

// extractProfileSkills extracts all skills from a profile.
func (m *Matcher) extractProfileSkills(profile *schema.FullProfile) map[string]bool {
	skills := make(map[string]bool)

	// From explicit skills
	for _, skill := range profile.Skills {
		skills[strings.ToLower(skill.Name)] = true
	}

	// From achievement skills
	for _, achievement := range profile.AllAchievements() {
		for _, skill := range achievement.Skills {
			skills[strings.ToLower(skill)] = true
		}
	}

	// From position skills
	for _, tenure := range profile.Tenures {
		for _, position := range tenure.Positions {
			for _, skill := range position.SkillsDefault {
				skills[strings.ToLower(skill)] = true
			}
			// Include domain-specific skills
			for _, cfg := range position.DomainConfigs {
				for _, skill := range cfg.Skills {
					skills[strings.ToLower(skill)] = true
				}
			}
		}
	}

	return skills
}

// extractProfileKeywords extracts keywords from a profile (skills + tags + titles).
func (m *Matcher) extractProfileKeywords(profile *schema.FullProfile) map[string]bool {
	keywords := make(map[string]bool)

	// Include all skills
	for k := range m.extractProfileSkills(profile) {
		keywords[k] = true
	}

	// Include tags from achievements
	for _, achievement := range profile.AllAchievements() {
		for _, tag := range achievement.Tags {
			keywords[strings.ToLower(tag)] = true
		}
	}

	// Include position titles and domains
	for _, tenure := range profile.Tenures {
		keywords[strings.ToLower(tenure.Company)] = true
		for _, position := range tenure.Positions {
			keywords[strings.ToLower(position.Title)] = true
			for _, cfg := range position.DomainConfigs {
				keywords[strings.ToLower(cfg.Domain)] = true
			}
		}
	}

	return keywords
}

// matchStrings finds matching and missing strings between profile set and JD list.
func (m *Matcher) matchStrings(profileSet map[string]bool, jdList []string) (matched, missing []string) {
	for _, item := range jdList {
		lower := strings.ToLower(item)
		if profileSet[lower] {
			matched = append(matched, item)
		} else {
			missing = append(missing, item)
		}
	}
	return matched, missing
}

// rankAchievements scores and ranks achievements by relevance to the JD.
func (m *Matcher) rankAchievements(achievements []schema.Achievement, jd *schema.JobDescParsed) []RankedAchievement {
	// Build JD sets for fast lookup
	requiredSkillSet := toSet(jd.RequiredSkills)
	preferredSkillSet := toSet(jd.PreferredSkills)
	keywordSet := toSet(jd.Keywords)

	var ranked []RankedAchievement

	for _, achievement := range achievements {
		ra := RankedAchievement{
			Achievement: achievement,
		}

		// Score based on skills
		for _, skill := range achievement.Skills {
			lower := strings.ToLower(skill)
			if requiredSkillSet[lower] {
				ra.Score += m.Weights.RequiredSkill
				ra.MatchedSkills = append(ra.MatchedSkills, skill)
			} else if preferredSkillSet[lower] {
				ra.Score += m.Weights.PreferredSkill
				ra.MatchedSkills = append(ra.MatchedSkills, skill)
			}
		}

		// Score based on tags
		for _, tag := range achievement.Tags {
			lower := strings.ToLower(tag)
			if keywordSet[lower] || requiredSkillSet[lower] || preferredSkillSet[lower] {
				ra.Score += m.Weights.Tag
				ra.MatchedTags = append(ra.MatchedTags, tag)
			}
		}

		// Score based on keywords in text
		achievementText := strings.ToLower(achievement.Name + " " + achievement.Situation + " " +
			achievement.Task + " " + achievement.Action + " " + achievement.Result)

		for _, kw := range jd.Keywords {
			if strings.Contains(achievementText, strings.ToLower(kw)) {
				ra.Score += m.Weights.Keyword
				ra.MatchedKeywords = append(ra.MatchedKeywords, kw)
			}
		}

		// Add quantified metrics bonus
		if len(achievement.Metrics.Values) > 0 {
			ra.Score += 2.0 // Bonus for quantified achievements
		}

		ranked = append(ranked, ra)
	}

	// Sort by score descending
	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].Score > ranked[j].Score
	})

	return ranked
}

// calculateOverallScore calculates a 0-100 overall match score.
func (m *Matcher) calculateOverallScore(result *MatchResult, jd *schema.JobDescParsed) float64 {
	if jd == nil {
		return 0
	}

	var score float64
	var maxScore float64

	// Required skills coverage (weighted heavily)
	if len(jd.RequiredSkills) > 0 {
		coverage := float64(len(result.MatchedRequiredSkills)) / float64(len(jd.RequiredSkills))
		score += coverage * 60 // Up to 60 points for required skills
		maxScore += 60
	}

	// Preferred skills coverage
	if len(jd.PreferredSkills) > 0 {
		coverage := float64(len(result.MatchedPreferredSkills)) / float64(len(jd.PreferredSkills))
		score += coverage * 20 // Up to 20 points for preferred skills
		maxScore += 20
	}

	// Keyword coverage
	if len(jd.Keywords) > 0 {
		coverage := float64(len(result.MatchedKeywords)) / float64(len(jd.Keywords))
		score += coverage * 10 // Up to 10 points for keywords
		maxScore += 10
	}

	// Achievement relevance bonus
	if len(result.RankedAchievements) > 0 {
		// Top achievement has a high score
		topScore := result.RankedAchievements[0].Score
		if topScore > 20 {
			score += 10 // Strong relevant achievements
		} else if topScore > 10 {
			score += 5 // Moderate relevant achievements
		}
		maxScore += 10
	}

	if maxScore == 0 {
		return 0
	}

	return (score / maxScore) * 100
}

// TopAchievements returns the top N achievements from a match result.
func (r *MatchResult) TopAchievements(n int) []schema.Achievement {
	if n > len(r.RankedAchievements) {
		n = len(r.RankedAchievements)
	}

	achievements := make([]schema.Achievement, n)
	for i := 0; i < n; i++ {
		achievements[i] = r.RankedAchievements[i].Achievement
	}
	return achievements
}

// FilterByMinScore returns achievements with score >= minScore.
func (r *MatchResult) FilterByMinScore(minScore float64) []RankedAchievement {
	var filtered []RankedAchievement
	for _, ra := range r.RankedAchievements {
		if ra.Score >= minScore {
			filtered = append(filtered, ra)
		}
	}
	return filtered
}

// toSet converts a string slice to a lowercase set.
func toSet(items []string) map[string]bool {
	set := make(map[string]bool)
	for _, item := range items {
		set[strings.ToLower(item)] = true
	}
	return set
}
