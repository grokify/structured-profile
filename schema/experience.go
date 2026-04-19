package schema

import (
	"strings"
)

// Tenure represents a period of employment at a company.
// A tenure can contain multiple positions (promotions/role changes).
type Tenure struct {
	BaseEntity
	ProfileID string `json:"profile_id,omitempty"`

	// Company Information
	Company     string `json:"company"`
	Description string `json:"description,omitempty"`
	Location    string `json:"location,omitempty"`

	// Date Range
	StartDate Date  `json:"start_date"`
	EndDate   *Date `json:"end_date,omitempty"` // nil = current

	// Positions within this tenure
	Positions []Position `json:"positions,omitempty"`

	// Collapsed display options (for multi-position tenures)
	CollapsedInfo *CollapsedInfo `json:"collapsed_info,omitempty"`
}

// NewTenure creates a new Tenure at the given company.
func NewTenure(company string, startDate Date) *Tenure {
	return &Tenure{
		BaseEntity: NewBaseEntity(),
		Company:    company,
		StartDate:  startDate,
		Positions:  []Position{},
	}
}

// IsCurrent returns true if the tenure is ongoing.
func (t *Tenure) IsCurrent() bool {
	return t.EndDate == nil || t.EndDate.IsZero()
}

// DateRange returns the tenure's date range.
func (t *Tenure) DateRange() DateRange {
	return DateRange{Start: t.StartDate, End: t.EndDate}
}

// AddPosition adds a position to the tenure.
func (t *Tenure) AddPosition(p Position) {
	p.TenureID = t.ID
	t.Positions = append(t.Positions, p)
}

// CollapsedInfo provides display options for collapsing multiple positions.
type CollapsedInfo struct {
	TitleCollapsedDefault string            `json:"title_collapsed_default"`
	TitleCollapsedAlts    map[string]string `json:"title_collapsed_alts,omitempty"` // domain -> title
	FormerTitles          string            `json:"former_titles,omitempty"`
}

// TitleForDomain returns the collapsed title for a domain.
func (c *CollapsedInfo) TitleForDomain(domain string) string {
	if c == nil {
		return ""
	}
	if domain != "" && c.TitleCollapsedAlts != nil {
		if title, ok := c.TitleCollapsedAlts[domain]; ok {
			return title
		}
	}
	return c.TitleCollapsedDefault
}

// Position represents a specific role/title within a tenure.
type Position struct {
	BaseEntity
	TenureID string `json:"tenure_id,omitempty"`

	// Position Details
	Title                    string `json:"title"`
	Description              string `json:"description,omitempty"`
	DescriptionWithoutCounts string `json:"description_without_counts,omitempty"`

	// Date Range (within the tenure)
	StartDate Date  `json:"start_date"`
	EndDate   *Date `json:"end_date,omitempty"`

	// Contact (for reference checks)
	Email        string `json:"email,omitempty"`
	ManagerEmail string `json:"manager_email,omitempty"`
	LocationCode string `json:"location_code,omitempty"` // UN/LOCODE

	// Achievements (STAR format)
	Achievements []Achievement `json:"achievements,omitempty"`

	// Domain-specific configurations
	DomainConfigs []PositionDomainConfig `json:"domain_configs,omitempty"`

	// Default skills (when no domain specified)
	SkillsDefault []string `json:"skills_default,omitempty"`

	// Notes per domain
	Notes Summaries `json:"notes,omitempty"`
}

// NewPosition creates a new Position with the given title.
func NewPosition(title string, startDate Date) *Position {
	return &Position{
		BaseEntity:    NewBaseEntity(),
		Title:         title,
		StartDate:     startDate,
		Achievements:  []Achievement{},
		DomainConfigs: []PositionDomainConfig{},
	}
}

// IsCurrent returns true if the position is ongoing.
func (p *Position) IsCurrent() bool {
	return p.EndDate == nil || p.EndDate.IsZero()
}

// DateRange returns the position's date range.
func (p *Position) DateRange() DateRange {
	return DateRange{Start: p.StartDate, End: p.EndDate}
}

// AddAchievement adds an achievement to the position.
func (p *Position) AddAchievement(a Achievement) {
	a.PositionID = p.ID
	p.Achievements = append(p.Achievements, a)
}

// GetDomainConfig returns the domain config for the specified domain.
func (p *Position) GetDomainConfig(domain string) *PositionDomainConfig {
	for i := range p.DomainConfigs {
		if p.DomainConfigs[i].Domain == domain {
			return &p.DomainConfigs[i]
		}
	}
	return nil
}

// SkillsForDomain returns skills for a specific domain, or default skills.
func (p *Position) SkillsForDomain(domain string) []string {
	if cfg := p.GetDomainConfig(domain); cfg != nil {
		return cfg.Skills
	}
	return p.SkillsDefault
}

// AchievementsForDomain returns achievements ordered for a specific domain.
func (p *Position) AchievementsForDomain(domain string) []Achievement {
	cfg := p.GetDomainConfig(domain)
	if cfg == nil || len(cfg.AchievementOrder) == 0 {
		// Return all non-skipped achievements in original order
		var result []Achievement
		for _, a := range p.Achievements {
			if !a.SkipDisplay {
				result = append(result, a)
			}
		}
		return result
	}

	// Build map of achievements by name
	byName := make(map[string]Achievement)
	for _, a := range p.Achievements {
		byName[a.Name] = a
	}

	// Return achievements in specified order
	var result []Achievement
	for _, name := range cfg.AchievementOrder {
		if a, ok := byName[name]; ok && !a.SkipDisplay {
			result = append(result, a)
		}
	}
	return result
}

// DescriptionForDomain returns the appropriate description based on options.
func (p *Position) DescriptionForDomain(withoutCounts bool) string {
	if withoutCounts && p.DescriptionWithoutCounts != "" {
		return p.DescriptionWithoutCounts
	}
	return p.Description
}

// PositionDomainConfig provides domain-specific configuration for a position.
type PositionDomainConfig struct {
	BaseEntity
	PositionID string `json:"position_id,omitempty"`

	Domain           string   `json:"domain"`
	Alias            string   `json:"alias,omitempty"` // Points to another domain config
	Skills           []string `json:"skills,omitempty"`
	AchievementOrder []string `json:"achievement_order,omitempty"` // Achievement names in display order
}

// NewPositionDomainConfig creates a new domain config.
func NewPositionDomainConfig(domain string) PositionDomainConfig {
	return PositionDomainConfig{
		BaseEntity: NewBaseEntity(),
		Domain:     domain,
	}
}

// Achievement represents a professional accomplishment in STAR format.
// STAR = Situation, Task, Action, Result
type Achievement struct {
	BaseEntity
	PositionID string `json:"position_id,omitempty"`

	// Identifier for ordering and referencing
	Name string `json:"name"`

	// STAR Format
	Situation string `json:"situation,omitempty"`
	Task      string `json:"task,omitempty"`
	Action    string `json:"action,omitempty"`
	Result    string `json:"result,omitempty"`

	// Metadata
	Skills  []string `json:"skills,omitempty"`
	Tags    []string `json:"tags,omitempty"`
	Metrics Metrics  `json:"metrics,omitempty"`

	// Display control
	SkipDisplay bool `json:"skip_display,omitempty"`
}

// NewAchievement creates a new Achievement with the given name.
func NewAchievement(name string) *Achievement {
	return &Achievement{
		BaseEntity: NewBaseEntity(),
		Name:       name,
	}
}

// NewSTARAchievement creates a new Achievement with full STAR content.
func NewSTARAchievement(name, situation, task, action, result string) *Achievement {
	return &Achievement{
		BaseEntity: NewBaseEntity(),
		Name:       name,
		Situation:  situation,
		Task:       task,
		Action:     action,
		Result:     result,
	}
}

// STARString returns the achievement as a single STAR-formatted string.
func (a *Achievement) STARString() string {
	parts := []string{}
	if a.Situation != "" {
		parts = append(parts, a.Situation)
	}
	if a.Task != "" {
		parts = append(parts, a.Task)
	}
	if a.Action != "" {
		parts = append(parts, a.Action)
	}
	if a.Result != "" {
		parts = append(parts, a.Result)
	}
	return strings.TrimSpace(strings.Join(parts, " "))
}

// HasTag returns true if the achievement has the specified tag.
func (a *Achievement) HasTag(tag string) bool {
	for _, t := range a.Tags {
		if strings.EqualFold(t, tag) {
			return true
		}
	}
	return false
}

// HasSkill returns true if the achievement has the specified skill.
func (a *Achievement) HasSkill(skill string) bool {
	for _, s := range a.Skills {
		if strings.EqualFold(s, skill) {
			return true
		}
	}
	return false
}

// Metrics contains quantifiable results from an achievement.
type Metrics struct {
	Values map[string]string `json:"values,omitempty"`
}

// Get returns a metric value by key.
func (m Metrics) Get(key string) string {
	if m.Values == nil {
		return ""
	}
	return m.Values[key]
}

// Set sets a metric value.
func (m *Metrics) Set(key, value string) {
	if m.Values == nil {
		m.Values = make(map[string]string)
	}
	m.Values[key] = value
}
