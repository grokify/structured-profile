package schema

// Opportunity represents a job opportunity being pursued.
type Opportunity struct {
	BaseEntity
	ProfileID string `json:"profile_id,omitempty"`

	// Job Details
	Company       string `json:"company"`
	Position      string `json:"position"`
	Location      string `json:"location,omitempty"`
	Remote        bool   `json:"remote,omitempty"`
	URL           string `json:"url,omitempty"`
	HiringManager string `json:"hiring_manager,omitempty"`

	// Salary Information
	SalaryMin    int    `json:"salary_min,omitempty"`
	SalaryMax    int    `json:"salary_max,omitempty"`
	SalaryCurrency string `json:"salary_currency,omitempty"`

	// Job Description
	JobDescRaw    string         `json:"job_desc_raw,omitempty"`
	JobDescParsed *JobDescParsed `json:"job_desc_parsed,omitempty"`

	// Company Research
	CompanyInfo *CompanyInfo `json:"company_info,omitempty"`

	// Generated Documents (IDs reference stored documents)
	ResumeID      string `json:"resume_id,omitempty"`
	CoverLetterID string `json:"cover_letter_id,omitempty"`
	EvaluationID  string `json:"evaluation_id,omitempty"`

	// Notes
	Notes string `json:"notes,omitempty"`
}

// NewOpportunity creates a new Opportunity.
func NewOpportunity(company, position string) *Opportunity {
	return &Opportunity{
		BaseEntity: NewBaseEntity(),
		Company:    company,
		Position:   position,
	}
}

// JobDescParsed contains structured data extracted from a job description.
type JobDescParsed struct {
	RequiredSkills   []string `json:"required_skills,omitempty"`
	PreferredSkills  []string `json:"preferred_skills,omitempty"`
	ExperienceYears  int      `json:"experience_years,omitempty"`
	Keywords         []string `json:"keywords,omitempty"`
	Responsibilities []string `json:"responsibilities,omitempty"`
	Qualifications   []string `json:"qualifications,omitempty"`

	// Extracted metadata
	SeniorityLevel string `json:"seniority_level,omitempty"` // "entry", "mid", "senior", "lead", "executive"
	TeamSize       string `json:"team_size,omitempty"`
	ReportingTo    string `json:"reporting_to,omitempty"`
}

// AllSkills returns both required and preferred skills.
func (jp *JobDescParsed) AllSkills() []string {
	if jp == nil {
		return nil
	}
	var all []string
	all = append(all, jp.RequiredSkills...)
	all = append(all, jp.PreferredSkills...)
	return all
}

// CompanyInfo contains research about a company.
type CompanyInfo struct {
	Name        string   `json:"name"`
	Industry    string   `json:"industry,omitempty"`
	Size        string   `json:"size,omitempty"` // "startup", "small", "mid", "large", "enterprise"
	Founded     int      `json:"founded,omitempty"`
	Headquarters string  `json:"headquarters,omitempty"`
	Values      []string `json:"values,omitempty"`
	Culture     []string `json:"culture,omitempty"`
	TechStack   []string `json:"tech_stack,omitempty"`
	Products    []string `json:"products,omitempty"`
	Competitors []string `json:"competitors,omitempty"`
	Links       []Link   `json:"links,omitempty"`
	Notes       string   `json:"notes,omitempty"`
}

// CompanySize represents company size categories.
type CompanySize string

const (
	CompanySizeStartup    CompanySize = "startup"    // < 50 employees
	CompanySizeSmall      CompanySize = "small"      // 50-200 employees
	CompanySizeMid        CompanySize = "mid"        // 200-1000 employees
	CompanySizeLarge      CompanySize = "large"      // 1000-10000 employees
	CompanySizeEnterprise CompanySize = "enterprise" // > 10000 employees
)

// SeniorityLevel represents job seniority levels.
type SeniorityLevel string

const (
	SeniorityEntry     SeniorityLevel = "entry"
	SeniorityMid       SeniorityLevel = "mid"
	SenioritySenior    SeniorityLevel = "senior"
	SeniorityStaff     SeniorityLevel = "staff"
	SeniorityPrincipal SeniorityLevel = "principal"
	SeniorityLead      SeniorityLevel = "lead"
	SeniorityManager   SeniorityLevel = "manager"
	SeniorityDirector  SeniorityLevel = "director"
	SeniorityVP        SeniorityLevel = "vp"
	SeniorityExecutive SeniorityLevel = "executive"
)
