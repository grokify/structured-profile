// Package json provides a JSON file-based implementation of the Store interface.
// Each profile is stored as a single JSON file containing all related data.
package json

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/grokify/structured-profile/schema"
	"github.com/grokify/structured-profile/store"
)

// Store implements the store.Store interface using JSON files.
// Each profile is stored as a single JSON file: {baseDir}/{profileID}.json
type Store struct {
	baseDir string
	mu      sync.RWMutex
	cache   map[string]*schema.FullProfile // Optional in-memory cache
	useCache bool
}

// Config contains configuration options for the JSON store.
type Config struct {
	BaseDir  string // Directory to store profile JSON files
	UseCache bool   // Whether to cache profiles in memory
}

// New creates a new JSON store with the given configuration.
func New(cfg Config) (*Store, error) {
	if cfg.BaseDir == "" {
		cfg.BaseDir = "."
	}

	// Ensure base directory exists
	if err := os.MkdirAll(cfg.BaseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}

	return &Store{
		baseDir:  cfg.BaseDir,
		cache:    make(map[string]*schema.FullProfile),
		useCache: cfg.UseCache,
	}, nil
}

// NewWithDir creates a new JSON store with the given directory.
func NewWithDir(dir string) (*Store, error) {
	return New(Config{BaseDir: dir})
}

// ProfilePath returns the file path for a profile.
func (s *Store) ProfilePath(id string) string {
	return filepath.Join(s.baseDir, id+".json")
}

// GetFullProfile retrieves a full profile by ID.
func (s *Store) GetFullProfile(ctx context.Context, id string) (*schema.FullProfile, error) {
	if id == "" {
		return nil, store.ErrInvalidID
	}

	s.mu.RLock()
	if s.useCache {
		if fp, ok := s.cache[id]; ok {
			s.mu.RUnlock()
			return fp, nil
		}
	}
	s.mu.RUnlock()

	path := s.ProfilePath(id)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, store.ErrNotFound
		}
		return nil, fmt.Errorf("failed to read profile: %w", err)
	}

	var fp schema.FullProfile
	if err := json.Unmarshal(data, &fp); err != nil {
		return nil, fmt.Errorf("failed to parse profile: %w", err)
	}

	if s.useCache {
		s.mu.Lock()
		s.cache[id] = &fp
		s.mu.Unlock()
	}

	return &fp, nil
}

// ListProfiles returns a list of all profiles (metadata only).
func (s *Store) ListProfiles(ctx context.Context) ([]schema.Profile, error) {
	entries, err := os.ReadDir(s.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to list profiles: %w", err)
	}

	var profiles []schema.Profile
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		id := strings.TrimSuffix(entry.Name(), ".json")
		fp, err := s.GetFullProfile(ctx, id)
		if err != nil {
			continue // Skip invalid files
		}
		profiles = append(profiles, fp.Profile)
	}

	return profiles, nil
}

// SaveFullProfile saves a full profile to disk.
func (s *Store) SaveFullProfile(ctx context.Context, fp *schema.FullProfile) error {
	if fp == nil {
		return errors.New("profile cannot be nil")
	}
	if fp.Profile.ID == "" {
		return store.ErrInvalidID
	}

	fp.Profile.Touch()

	data, err := json.MarshalIndent(fp, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal profile: %w", err)
	}

	path := s.ProfilePath(fp.Profile.ID)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write profile: %w", err)
	}

	if s.useCache {
		s.mu.Lock()
		s.cache[fp.Profile.ID] = fp
		s.mu.Unlock()
	}

	return nil
}

// DeleteProfile deletes a profile by ID.
func (s *Store) DeleteProfile(ctx context.Context, id string) error {
	if id == "" {
		return store.ErrInvalidID
	}

	path := s.ProfilePath(id)
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return store.ErrNotFound
		}
		return fmt.Errorf("failed to delete profile: %w", err)
	}

	if s.useCache {
		s.mu.Lock()
		delete(s.cache, id)
		s.mu.Unlock()
	}

	return nil
}

// GetProfile retrieves profile metadata by ID.
func (s *Store) GetProfile(ctx context.Context, id string) (*schema.Profile, error) {
	fp, err := s.GetFullProfile(ctx, id)
	if err != nil {
		return nil, err
	}
	return &fp.Profile, nil
}

// SaveProfile saves profile metadata (loads full profile, updates, saves).
func (s *Store) SaveProfile(ctx context.Context, p *schema.Profile) error {
	if p == nil {
		return errors.New("profile cannot be nil")
	}

	fp, err := s.GetFullProfile(ctx, p.ID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			// New profile
			fp = &schema.FullProfile{Profile: *p}
		} else {
			return err
		}
	} else {
		fp.Profile = *p
	}

	return s.SaveFullProfile(ctx, fp)
}

// SearchAchievementsByTags searches for achievements with matching tags.
func (s *Store) SearchAchievementsByTags(ctx context.Context, profileID string, tags []string) ([]schema.Achievement, error) {
	fp, err := s.GetFullProfile(ctx, profileID)
	if err != nil {
		return nil, err
	}

	tagSet := make(map[string]bool)
	for _, t := range tags {
		tagSet[strings.ToLower(t)] = true
	}

	var matches []schema.Achievement
	for _, achievement := range fp.AllAchievements() {
		for _, tag := range achievement.Tags {
			if tagSet[strings.ToLower(tag)] {
				matches = append(matches, achievement)
				break
			}
		}
	}

	return matches, nil
}

// SearchAchievementsBySkills searches for achievements with matching skills.
func (s *Store) SearchAchievementsBySkills(ctx context.Context, profileID string, skills []string) ([]schema.Achievement, error) {
	fp, err := s.GetFullProfile(ctx, profileID)
	if err != nil {
		return nil, err
	}

	skillSet := make(map[string]bool)
	for _, sk := range skills {
		skillSet[strings.ToLower(sk)] = true
	}

	var matches []schema.Achievement
	for _, achievement := range fp.AllAchievements() {
		for _, skill := range achievement.Skills {
			if skillSet[strings.ToLower(skill)] {
				matches = append(matches, achievement)
				break
			}
		}
	}

	return matches, nil
}

// GetOpportunity retrieves an opportunity by ID from a profile.
func (s *Store) GetOpportunity(ctx context.Context, profileID, opportunityID string) (*schema.Opportunity, error) {
	fp, err := s.GetFullProfile(ctx, profileID)
	if err != nil {
		return nil, err
	}

	for i := range fp.Opportunities {
		if fp.Opportunities[i].ID == opportunityID {
			return &fp.Opportunities[i], nil
		}
	}

	return nil, store.ErrNotFound
}

// GetApplication retrieves an application by ID from a profile.
func (s *Store) GetApplication(ctx context.Context, profileID, applicationID string) (*schema.Application, error) {
	fp, err := s.GetFullProfile(ctx, profileID)
	if err != nil {
		return nil, err
	}

	for i := range fp.Applications {
		if fp.Applications[i].ID == applicationID {
			return &fp.Applications[i], nil
		}
	}

	return nil, store.ErrNotFound
}

// ListOpportunities returns all opportunities for a profile.
func (s *Store) ListOpportunities(ctx context.Context, profileID string) ([]schema.Opportunity, error) {
	fp, err := s.GetFullProfile(ctx, profileID)
	if err != nil {
		return nil, err
	}
	return fp.Opportunities, nil
}

// SaveOpportunity saves an opportunity to a profile.
func (s *Store) SaveOpportunity(ctx context.Context, profileID string, opp *schema.Opportunity) error {
	if opp == nil {
		return errors.New("opportunity cannot be nil")
	}

	fp, err := s.GetFullProfile(ctx, profileID)
	if err != nil {
		return err
	}

	opp.ProfileID = profileID
	opp.Touch()

	// Update existing or append new
	found := false
	for i := range fp.Opportunities {
		if fp.Opportunities[i].ID == opp.ID {
			fp.Opportunities[i] = *opp
			found = true
			break
		}
	}
	if !found {
		fp.Opportunities = append(fp.Opportunities, *opp)
	}

	return s.SaveFullProfile(ctx, fp)
}

// DeleteOpportunity deletes an opportunity from a profile.
func (s *Store) DeleteOpportunity(ctx context.Context, profileID, opportunityID string) error {
	fp, err := s.GetFullProfile(ctx, profileID)
	if err != nil {
		return err
	}

	found := false
	for i := range fp.Opportunities {
		if fp.Opportunities[i].ID == opportunityID {
			fp.Opportunities = append(fp.Opportunities[:i], fp.Opportunities[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return store.ErrNotFound
	}

	return s.SaveFullProfile(ctx, fp)
}

// ListApplications returns all applications for a profile.
func (s *Store) ListApplications(ctx context.Context, profileID string) ([]schema.Application, error) {
	fp, err := s.GetFullProfile(ctx, profileID)
	if err != nil {
		return nil, err
	}
	return fp.Applications, nil
}

// SaveApplication saves an application to a profile.
func (s *Store) SaveApplication(ctx context.Context, profileID string, app *schema.Application) error {
	if app == nil {
		return errors.New("application cannot be nil")
	}

	fp, err := s.GetFullProfile(ctx, profileID)
	if err != nil {
		return err
	}

	app.ProfileID = profileID
	app.Touch()

	// Update existing or append new
	found := false
	for i := range fp.Applications {
		if fp.Applications[i].ID == app.ID {
			fp.Applications[i] = *app
			found = true
			break
		}
	}
	if !found {
		fp.Applications = append(fp.Applications, *app)
	}

	return s.SaveFullProfile(ctx, fp)
}

// DeleteApplication deletes an application from a profile.
func (s *Store) DeleteApplication(ctx context.Context, profileID, applicationID string) error {
	fp, err := s.GetFullProfile(ctx, profileID)
	if err != nil {
		return err
	}

	found := false
	for i := range fp.Applications {
		if fp.Applications[i].ID == applicationID {
			fp.Applications = append(fp.Applications[:i], fp.Applications[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return store.ErrNotFound
	}

	return s.SaveFullProfile(ctx, fp)
}

// GetCoverLetter retrieves a cover letter by ID from a profile.
func (s *Store) GetCoverLetter(ctx context.Context, profileID, coverLetterID string) (*schema.CoverLetter, error) {
	fp, err := s.GetFullProfile(ctx, profileID)
	if err != nil {
		return nil, err
	}

	for i := range fp.CoverLetters {
		if fp.CoverLetters[i].ID == coverLetterID {
			return &fp.CoverLetters[i], nil
		}
	}

	return nil, store.ErrNotFound
}

// ListCoverLetters returns all cover letters for a profile.
func (s *Store) ListCoverLetters(ctx context.Context, profileID string) ([]schema.CoverLetter, error) {
	fp, err := s.GetFullProfile(ctx, profileID)
	if err != nil {
		return nil, err
	}
	return fp.CoverLetters, nil
}

// SaveCoverLetter saves a cover letter to a profile.
func (s *Store) SaveCoverLetter(ctx context.Context, profileID string, cl *schema.CoverLetter) error {
	if cl == nil {
		return errors.New("cover letter cannot be nil")
	}

	fp, err := s.GetFullProfile(ctx, profileID)
	if err != nil {
		return err
	}

	cl.ProfileID = profileID
	cl.Touch()

	// Update existing or append new
	found := false
	for i := range fp.CoverLetters {
		if fp.CoverLetters[i].ID == cl.ID {
			fp.CoverLetters[i] = *cl
			found = true
			break
		}
	}
	if !found {
		fp.CoverLetters = append(fp.CoverLetters, *cl)
	}

	return s.SaveFullProfile(ctx, fp)
}

// DeleteCoverLetter deletes a cover letter from a profile.
func (s *Store) DeleteCoverLetter(ctx context.Context, profileID, coverLetterID string) error {
	fp, err := s.GetFullProfile(ctx, profileID)
	if err != nil {
		return err
	}

	found := false
	for i := range fp.CoverLetters {
		if fp.CoverLetters[i].ID == coverLetterID {
			fp.CoverLetters = append(fp.CoverLetters[:i], fp.CoverLetters[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return store.ErrNotFound
	}

	return s.SaveFullProfile(ctx, fp)
}

// GetCoverLetterTemplate retrieves a cover letter template by ID from a profile.
func (s *Store) GetCoverLetterTemplate(ctx context.Context, profileID, templateID string) (*schema.CoverLetterTemplate, error) {
	fp, err := s.GetFullProfile(ctx, profileID)
	if err != nil {
		return nil, err
	}

	for i := range fp.CoverLetterTemplates {
		if fp.CoverLetterTemplates[i].ID == templateID {
			return &fp.CoverLetterTemplates[i], nil
		}
	}

	return nil, store.ErrNotFound
}

// ListCoverLetterTemplates returns all cover letter templates for a profile.
func (s *Store) ListCoverLetterTemplates(ctx context.Context, profileID string) ([]schema.CoverLetterTemplate, error) {
	fp, err := s.GetFullProfile(ctx, profileID)
	if err != nil {
		return nil, err
	}
	return fp.CoverLetterTemplates, nil
}

// SaveCoverLetterTemplate saves a cover letter template to a profile.
func (s *Store) SaveCoverLetterTemplate(ctx context.Context, profileID string, tmpl *schema.CoverLetterTemplate) error {
	if tmpl == nil {
		return errors.New("cover letter template cannot be nil")
	}

	fp, err := s.GetFullProfile(ctx, profileID)
	if err != nil {
		return err
	}

	tmpl.ProfileID = profileID
	tmpl.Touch()

	// Update existing or append new
	found := false
	for i := range fp.CoverLetterTemplates {
		if fp.CoverLetterTemplates[i].ID == tmpl.ID {
			fp.CoverLetterTemplates[i] = *tmpl
			found = true
			break
		}
	}
	if !found {
		fp.CoverLetterTemplates = append(fp.CoverLetterTemplates, *tmpl)
	}

	return s.SaveFullProfile(ctx, fp)
}

// DeleteCoverLetterTemplate deletes a cover letter template from a profile.
func (s *Store) DeleteCoverLetterTemplate(ctx context.Context, profileID, templateID string) error {
	fp, err := s.GetFullProfile(ctx, profileID)
	if err != nil {
		return err
	}

	found := false
	for i := range fp.CoverLetterTemplates {
		if fp.CoverLetterTemplates[i].ID == templateID {
			fp.CoverLetterTemplates = append(fp.CoverLetterTemplates[:i], fp.CoverLetterTemplates[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return store.ErrNotFound
	}

	return s.SaveFullProfile(ctx, fp)
}

// Close releases any resources held by the store.
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache = make(map[string]*schema.FullProfile)
	return nil
}

// ClearCache clears the in-memory cache.
func (s *Store) ClearCache() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache = make(map[string]*schema.FullProfile)
}

// Ensure Store implements the Store interface.
var _ store.Store = (*Store)(nil)
