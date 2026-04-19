// Package store provides data persistence interfaces and implementations.
package store

import (
	"context"
	"errors"

	"github.com/grokify/structured-profile/schema"
)

// Common errors
var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidID     = errors.New("invalid ID")
)

// Store defines the interface for profile data persistence.
// Implementations can use different backends (JSON files, PostgreSQL, etc.)
type Store interface {
	// Profile operations - single file approach stores entire profile
	GetFullProfile(ctx context.Context, id string) (*schema.FullProfile, error)
	ListProfiles(ctx context.Context) ([]schema.Profile, error)
	SaveFullProfile(ctx context.Context, fp *schema.FullProfile) error
	DeleteProfile(ctx context.Context, id string) error

	// Profile-level operations (convenience methods that work on FullProfile)
	GetProfile(ctx context.Context, id string) (*schema.Profile, error)
	SaveProfile(ctx context.Context, p *schema.Profile) error

	// Search operations
	SearchAchievementsByTags(ctx context.Context, profileID string, tags []string) ([]schema.Achievement, error)
	SearchAchievementsBySkills(ctx context.Context, profileID string, skills []string) ([]schema.Achievement, error)

	// Opportunity operations
	GetOpportunity(ctx context.Context, profileID, opportunityID string) (*schema.Opportunity, error)
	ListOpportunities(ctx context.Context, profileID string) ([]schema.Opportunity, error)
	SaveOpportunity(ctx context.Context, profileID string, opp *schema.Opportunity) error
	DeleteOpportunity(ctx context.Context, profileID, opportunityID string) error

	// Application operations
	GetApplication(ctx context.Context, profileID, applicationID string) (*schema.Application, error)
	ListApplications(ctx context.Context, profileID string) ([]schema.Application, error)
	SaveApplication(ctx context.Context, profileID string, app *schema.Application) error
	DeleteApplication(ctx context.Context, profileID, applicationID string) error

	// CoverLetter operations
	GetCoverLetter(ctx context.Context, profileID, coverLetterID string) (*schema.CoverLetter, error)
	ListCoverLetters(ctx context.Context, profileID string) ([]schema.CoverLetter, error)
	SaveCoverLetter(ctx context.Context, profileID string, cl *schema.CoverLetter) error
	DeleteCoverLetter(ctx context.Context, profileID, coverLetterID string) error

	// CoverLetterTemplate operations
	GetCoverLetterTemplate(ctx context.Context, profileID, templateID string) (*schema.CoverLetterTemplate, error)
	ListCoverLetterTemplates(ctx context.Context, profileID string) ([]schema.CoverLetterTemplate, error)
	SaveCoverLetterTemplate(ctx context.Context, profileID string, tmpl *schema.CoverLetterTemplate) error
	DeleteCoverLetterTemplate(ctx context.Context, profileID, templateID string) error

	// File path operations (for JSON store)
	ProfilePath(id string) string

	// Close releases any resources held by the store
	Close() error
}

// ProfileStore is an alias for Store for backward compatibility.
type ProfileStore = Store

// ReadOnlyStore provides read-only access to profile data.
type ReadOnlyStore interface {
	GetFullProfile(ctx context.Context, id string) (*schema.FullProfile, error)
	ListProfiles(ctx context.Context) ([]schema.Profile, error)
	GetProfile(ctx context.Context, id string) (*schema.Profile, error)
	SearchAchievementsByTags(ctx context.Context, profileID string, tags []string) ([]schema.Achievement, error)
	SearchAchievementsBySkills(ctx context.Context, profileID string, skills []string) ([]schema.Achievement, error)
}
