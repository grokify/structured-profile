package schema

import "time"

// VerifiableCredential represents a verifiable online credential.
// Examples: GitHub profile, StackOverflow reputation, LinkedIn profile.
type VerifiableCredential struct {
	BaseEntity
	ProfileID string `json:"profile_id,omitempty"`

	Type       string `json:"type"` // "github", "stackoverflow", "linkedin", "portfolio"
	Username   string `json:"username,omitempty"`
	ProfileURL string `json:"profile_url"`

	// Platform-specific data
	Data CredentialData `json:"data,omitempty"`

	// Verification
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
}

// NewVerifiableCredential creates a new VerifiableCredential.
func NewVerifiableCredential(credType, profileURL string) *VerifiableCredential {
	return &VerifiableCredential{
		BaseEntity: NewBaseEntity(),
		Type:       credType,
		ProfileURL: profileURL,
	}
}

// IsVerified returns true if the credential has been verified.
func (vc *VerifiableCredential) IsVerified() bool {
	return vc.VerifiedAt != nil
}

// Verify marks the credential as verified at the current time.
func (vc *VerifiableCredential) Verify() {
	now := time.Now().UTC()
	vc.VerifiedAt = &now
	vc.Touch()
}

// CredentialType represents known credential types.
type CredentialType string

const (
	CredentialTypeGitHub        CredentialType = "github"
	CredentialTypeStackOverflow CredentialType = "stackoverflow"
	CredentialTypeLinkedIn      CredentialType = "linkedin"
	CredentialTypePortfolio     CredentialType = "portfolio"
	CredentialTypeWebsite       CredentialType = "website"
	CredentialTypeMedium        CredentialType = "medium"
	CredentialTypeDevTo         CredentialType = "devto"
)

// CredentialData contains platform-specific data for credentials.
type CredentialData struct {
	// GitHub
	Repositories  int      `json:"repositories,omitempty"`
	Stars         int      `json:"stars,omitempty"`
	Contributions int      `json:"contributions,omitempty"`
	Languages     []string `json:"languages,omitempty"`
	Followers     int      `json:"followers,omitempty"`

	// StackOverflow
	Reputation   int      `json:"reputation,omitempty"`
	GoldBadges   int      `json:"gold_badges,omitempty"`
	SilverBadges int      `json:"silver_badges,omitempty"`
	BronzeBadges int      `json:"bronze_badges,omitempty"`
	TopTags      []string `json:"top_tags,omitempty"`
	Answers      int      `json:"answers,omitempty"`
	Questions    int      `json:"questions,omitempty"`

	// LinkedIn
	Connections  int `json:"connections,omitempty"`
	Endorsements int `json:"endorsements,omitempty"`

	// General
	Articles   int `json:"articles,omitempty"`
	Projects   int `json:"projects,omitempty"`
	LastActive Date `json:"last_active,omitempty"`
}

// Publication represents a published work (article, paper, patent, book).
type Publication struct {
	BaseEntity
	ProfileID string `json:"profile_id,omitempty"`

	Type        string `json:"type"` // "article", "paper", "patent", "book", "talk"
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`

	// For patents
	PatentNumber string `json:"patent_number,omitempty"`

	// Publication details
	Publisher string `json:"publisher,omitempty"`
	Date      Date   `json:"date,omitempty"`

	// Display control
	Display bool `json:"display"`
}

// NewPublication creates a new Publication.
func NewPublication(pubType, title string) *Publication {
	return &Publication{
		BaseEntity: NewBaseEntity(),
		Type:       pubType,
		Title:      title,
		Display:    true,
	}
}

// PublicationType represents publication types.
type PublicationType string

const (
	PublicationTypeArticle PublicationType = "article"
	PublicationTypePaper   PublicationType = "paper"
	PublicationTypePatent  PublicationType = "patent"
	PublicationTypeBook    PublicationType = "book"
	PublicationTypeTalk    PublicationType = "talk"
	PublicationTypeProject PublicationType = "project"
)
