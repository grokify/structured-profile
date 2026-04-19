package schema

// Profile represents a professional profile with contact information,
// summaries, and links to online presence.
type Profile struct {
	BaseEntity

	// Personal Information
	Name     string `json:"name"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Location string `json:"location,omitempty"`

	// Online Presence
	Links []Link `json:"links,omitempty"`

	// Summaries by domain/filter
	Summaries Summaries `json:"summaries,omitempty"`
}

// NewProfile creates a new Profile with the given name.
func NewProfile(name string) *Profile {
	return &Profile{
		BaseEntity: NewBaseEntity(),
		Name:       name,
		Links:      []Link{},
		Summaries:  Summaries{},
	}
}

// Link represents a URL link with type and optional display text.
type Link struct {
	Type string `json:"type"`           // "linkedin", "github", "website", "portfolio", "stackoverflow"
	URL  string `json:"url"`            // Full URL
	Text string `json:"text,omitempty"` // Optional display text
}

// NewLink creates a new Link.
func NewLink(linkType, url string) Link {
	return Link{Type: linkType, URL: url}
}

// Summaries contains professional summaries keyed by domain.
type Summaries struct {
	Default  string            `json:"default,omitempty"`
	ByDomain map[string]string `json:"by_domain,omitempty"`
}

// ForDomain returns the summary for the specified domain, or default if not found.
func (s Summaries) ForDomain(domain string) string {
	if domain != "" && s.ByDomain != nil {
		if summary, ok := s.ByDomain[domain]; ok {
			return summary
		}
	}
	return s.Default
}

// SetDomain sets a domain-specific summary.
func (s *Summaries) SetDomain(domain, summary string) {
	if s.ByDomain == nil {
		s.ByDomain = make(map[string]string)
	}
	s.ByDomain[domain] = summary
}

// FullProfile contains a Profile with all related entities loaded.
// This is the structure used for single-file JSON storage.
type FullProfile struct {
	Profile              Profile                `json:"profile"`
	Tenures              []Tenure               `json:"tenures,omitempty"`
	Skills               []Skill                `json:"skills,omitempty"`
	Education            []Education            `json:"education,omitempty"`
	Certifications       []Certification        `json:"certifications,omitempty"`
	Credentials          []VerifiableCredential `json:"credentials,omitempty"`
	Publications         []Publication          `json:"publications,omitempty"`
	Opportunities        []Opportunity          `json:"opportunities,omitempty"`
	Applications         []Application          `json:"applications,omitempty"`
	CoverLetters         []CoverLetter          `json:"cover_letters,omitempty"`
	CoverLetterTemplates []CoverLetterTemplate  `json:"cover_letter_templates,omitempty"`
	PrepSets             []InterviewPrepSet     `json:"prep_sets,omitempty"`
}

// NewFullProfile creates a new FullProfile with the given name.
func NewFullProfile(name string) *FullProfile {
	return &FullProfile{
		Profile: *NewProfile(name),
	}
}

// FindTenure finds a tenure by ID.
func (fp *FullProfile) FindTenure(id string) *Tenure {
	for i := range fp.Tenures {
		if fp.Tenures[i].ID == id {
			return &fp.Tenures[i]
		}
	}
	return nil
}

// FindPosition finds a position by ID across all tenures.
func (fp *FullProfile) FindPosition(id string) *Position {
	for i := range fp.Tenures {
		for j := range fp.Tenures[i].Positions {
			if fp.Tenures[i].Positions[j].ID == id {
				return &fp.Tenures[i].Positions[j]
			}
		}
	}
	return nil
}

// FindAchievement finds an achievement by ID across all positions.
func (fp *FullProfile) FindAchievement(id string) *Achievement {
	for i := range fp.Tenures {
		for j := range fp.Tenures[i].Positions {
			for k := range fp.Tenures[i].Positions[j].Achievements {
				if fp.Tenures[i].Positions[j].Achievements[k].ID == id {
					return &fp.Tenures[i].Positions[j].Achievements[k]
				}
			}
		}
	}
	return nil
}

// AllAchievements returns all achievements across all tenures and positions.
func (fp *FullProfile) AllAchievements() []Achievement {
	var achievements []Achievement
	for i := range fp.Tenures {
		for j := range fp.Tenures[i].Positions {
			achievements = append(achievements, fp.Tenures[i].Positions[j].Achievements...)
		}
	}
	return achievements
}

// AllPositions returns all positions across all tenures.
func (fp *FullProfile) AllPositions() []Position {
	var positions []Position
	for i := range fp.Tenures {
		positions = append(positions, fp.Tenures[i].Positions...)
	}
	return positions
}
