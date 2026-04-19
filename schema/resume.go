package schema

// Resume represents a generated resume tailored for a specific opportunity.
type Resume struct {
	BaseEntity
	ProfileID     string `json:"profile_id"`
	OpportunityID string `json:"opportunity_id,omitempty"`

	// Version tracking
	Version     string `json:"version,omitempty"`
	VersionNote string `json:"version_note,omitempty"`

	// Generation settings
	Domain  string         `json:"domain,omitempty"` // Domain filter used
	Options *ResumeOptions `json:"options,omitempty"`

	// Content (generated)
	Content *ResumeContent `json:"content,omitempty"`

	// Output
	MarkdownOutput string `json:"markdown_output,omitempty"`
}

// NewResume creates a new Resume for a profile and opportunity.
func NewResume(profileID, opportunityID string) *Resume {
	return &Resume{
		BaseEntity:    NewBaseEntity(),
		ProfileID:     profileID,
		OpportunityID: opportunityID,
		Version:       "1.0",
	}
}

// ResumeOptions contains configuration for resume generation.
type ResumeOptions struct {
	// Section inclusion
	IncludeContact         bool `json:"include_contact"`
	IncludeSummary         bool `json:"include_summary"`
	IncludeExperience      bool `json:"include_experience"`
	IncludeEducation       bool `json:"include_education"`
	IncludeCertifications  bool `json:"include_certifications"`
	IncludeSkills          bool `json:"include_skills"`
	IncludePublications    bool `json:"include_publications"`
	IncludeCommunity       bool `json:"include_community"`

	// Formatting
	CollapseTenurePositions  bool `json:"collapse_tenure_positions,omitempty"`
	DescriptionWithoutCounts bool `json:"description_without_counts,omitempty"`

	// Limits
	MaxExperiences   int `json:"max_experiences,omitempty"`
	MaxAchievements  int `json:"max_achievements,omitempty"`
	MaxSkills        int `json:"max_skills,omitempty"`

	// Output
	OutputBaseFilename string `json:"output_base_filename,omitempty"`
	MarginScalar       int    `json:"margin_scalar,omitempty"`
	MarginUnit         string `json:"margin_unit,omitempty"` // "cm", "in", "mm"
}

// DefaultResumeOptions returns sensible default options.
func DefaultResumeOptions() *ResumeOptions {
	return &ResumeOptions{
		IncludeContact:        true,
		IncludeSummary:        true,
		IncludeExperience:     true,
		IncludeEducation:      true,
		IncludeCertifications: true,
		IncludeSkills:         true,
		MarginScalar:          1,
		MarginUnit:            "cm",
	}
}

// ResumeContent contains the actual content of a generated resume.
type ResumeContent struct {
	// Header
	Name     string   `json:"name"`
	Email    string   `json:"email,omitempty"`
	Phone    string   `json:"phone,omitempty"`
	Location string   `json:"location,omitempty"`
	Links    []Link   `json:"links,omitempty"`

	// Summary
	Summary string `json:"summary,omitempty"`

	// Experience (selected and ordered)
	Experiences []ResumeExperience `json:"experiences,omitempty"`

	// Skills (selected and ordered)
	Skills []string `json:"skills,omitempty"`

	// Education
	Education []Education `json:"education,omitempty"`

	// Certifications
	Certifications []Certification `json:"certifications,omitempty"`

	// Additional Sections
	Publications []Publication `json:"publications,omitempty"`
}

// ResumeExperience represents a work experience entry in a resume.
type ResumeExperience struct {
	Company     string   `json:"company"`
	Title       string   `json:"title"`
	Location    string   `json:"location,omitempty"`
	StartDate   Date     `json:"start_date"`
	EndDate     *Date    `json:"end_date,omitempty"`
	Description string   `json:"description,omitempty"`
	Achievements []string `json:"achievements,omitempty"` // STAR strings
	Skills       []string `json:"skills,omitempty"`
}

// CoverLetter represents a generated cover letter.
type CoverLetter struct {
	BaseEntity
	ProfileID     string `json:"profile_id"`
	OpportunityID string `json:"opportunity_id"`

	// Version tracking
	Version     string `json:"version,omitempty"`
	VersionNote string `json:"version_note,omitempty"`

	// Target
	TargetCompany  string `json:"target_company"`
	TargetPosition string `json:"target_position"`
	HiringManager  string `json:"hiring_manager,omitempty"`

	// Content
	Opening    string   `json:"opening,omitempty"`    // Introduction paragraph
	Body       string   `json:"body,omitempty"`       // Main content
	Closing    string   `json:"closing,omitempty"`    // Conclusion
	STARRefs   []string `json:"star_refs,omitempty"`  // Achievement IDs referenced

	// Output
	MarkdownOutput string `json:"markdown_output,omitempty"`
}

// NewCoverLetter creates a new CoverLetter.
func NewCoverLetter(profileID, opportunityID string) *CoverLetter {
	return &CoverLetter{
		BaseEntity:    NewBaseEntity(),
		ProfileID:     profileID,
		OpportunityID: opportunityID,
		Version:       "1.0",
	}
}

// FullText returns the complete cover letter text.
func (cl *CoverLetter) FullText() string {
	text := ""
	if cl.Opening != "" {
		text += cl.Opening + "\n\n"
	}
	if cl.Body != "" {
		text += cl.Body + "\n\n"
	}
	if cl.Closing != "" {
		text += cl.Closing
	}
	return text
}

// CoverLetterTemplate represents a template for cover letter generation.
type CoverLetterTemplate struct {
	BaseEntity
	ProfileID   string `json:"profile_id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	// Template sections with variables
	OpeningTemplate string `json:"opening_template,omitempty"`
	BodyTemplate    string `json:"body_template,omitempty"`
	ClosingTemplate string `json:"closing_template,omitempty"`

	// Supported variables: {{.Company}}, {{.Position}}, {{.HiringManager}}, {{.STAR1}}, etc.
}

// NewCoverLetterTemplate creates a new CoverLetterTemplate.
func NewCoverLetterTemplate(name string) *CoverLetterTemplate {
	return &CoverLetterTemplate{
		BaseEntity: NewBaseEntity(),
		Name:       name,
	}
}
