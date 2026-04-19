package schema

import "time"

// timeNow is a variable to allow testing with mock time.
var timeNow = time.Now

// Education represents an educational background entry.
type Education struct {
	BaseEntity
	ProfileID string `json:"profile_id,omitempty"`

	Institution string `json:"institution"`
	Degree      string `json:"degree"`
	Field       string `json:"field,omitempty"`

	StartDate Date  `json:"start_date,omitempty"`
	EndDate   *Date `json:"end_date,omitempty"`

	Honors string `json:"honors,omitempty"` // e.g., "University Scholar", "Magna Cum Laude"
	GPA    string `json:"gpa,omitempty"`

	// Display control
	Display bool `json:"display"`
}

// NewEducation creates a new Education entry.
func NewEducation(institution, degree string) *Education {
	return &Education{
		BaseEntity:  NewBaseEntity(),
		Institution: institution,
		Degree:      degree,
		Display:     true,
	}
}

// DateRange returns the education's date range.
func (e *Education) DateRange() DateRange {
	return DateRange{Start: e.StartDate, End: e.EndDate}
}

// Certification represents a professional certification.
type Certification struct {
	BaseEntity
	ProfileID string `json:"profile_id,omitempty"`

	Name           string `json:"name"`
	Issuer         string `json:"issuer,omitempty"`
	IssueDate      Date   `json:"issue_date,omitempty"`
	ExpirationDate *Date  `json:"expiration_date,omitempty"`

	// Verification
	CredentialID  string `json:"credential_id,omitempty"`
	CredentialURL string `json:"credential_url,omitempty"`

	// Display control
	Display bool `json:"display"`
}

// NewCertification creates a new Certification.
func NewCertification(name string) *Certification {
	return &Certification{
		BaseEntity: NewBaseEntity(),
		Name:       name,
		Display:    true,
	}
}

// NewCertificationWithIssuer creates a new Certification with issuer details.
func NewCertificationWithIssuer(name, issuer string, issueDate Date) *Certification {
	return &Certification{
		BaseEntity: NewBaseEntity(),
		Name:       name,
		Issuer:     issuer,
		IssueDate:  issueDate,
		Display:    true,
	}
}

// IsExpired returns true if the certification has expired.
func (c *Certification) IsExpired() bool {
	if c.ExpirationDate == nil || c.ExpirationDate.IsZero() {
		return false
	}
	now := NewDateFromTime(timeNow())
	return c.ExpirationDate.Before(now)
}

// IsValid returns true if the certification is not expired.
func (c *Certification) IsValid() bool {
	return !c.IsExpired()
}
