package schema

import "strings"

// Skill represents a professional skill with category and proficiency level.
type Skill struct {
	BaseEntity
	ProfileID string `json:"profile_id,omitempty"`

	Name        string `json:"name"`
	Category    string `json:"category,omitempty"`    // "technical", "domain", "soft"
	Subcategory string `json:"subcategory,omitempty"` // "cloud", "iam", "leadership"
	Level       string `json:"level,omitempty"`       // "expert", "proficient", "familiar"

	// Verification
	Verified  bool   `json:"verified,omitempty"`
	VerifyURL string `json:"verify_url,omitempty"`
}

// NewSkill creates a new Skill.
func NewSkill(name string) *Skill {
	return &Skill{
		BaseEntity: NewBaseEntity(),
		Name:       name,
	}
}

// NewSkillWithCategory creates a new Skill with category and level.
func NewSkillWithCategory(name, category, level string) *Skill {
	return &Skill{
		BaseEntity: NewBaseEntity(),
		Name:       name,
		Category:   category,
		Level:      level,
	}
}

// SkillLevel represents proficiency levels.
type SkillLevel string

const (
	SkillLevelExpert     SkillLevel = "expert"
	SkillLevelProficient SkillLevel = "proficient"
	SkillLevelFamiliar   SkillLevel = "familiar"
	SkillLevelBeginner   SkillLevel = "beginner"
)

// SkillCategory represents skill categories.
type SkillCategory string

const (
	SkillCategoryTechnical SkillCategory = "technical"
	SkillCategoryDomain    SkillCategory = "domain"
	SkillCategorySoft      SkillCategory = "soft"
	SkillCategoryTool      SkillCategory = "tool"
	SkillCategoryLanguage  SkillCategory = "language"
)

// SkillGroup represents a group of related skills.
type SkillGroup struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Skills      []string `json:"skills"`
}

// SkillSet represents a collection of skills grouped by category.
type SkillSet struct {
	Technical []string `json:"technical,omitempty"`
	Domain    []string `json:"domain,omitempty"`
	Soft      []string `json:"soft,omitempty"`
	Tools     []string `json:"tools,omitempty"`
	Languages []string `json:"languages,omitempty"`
}

// All returns all skills flattened into a single slice.
func (ss SkillSet) All() []string {
	var all []string
	all = append(all, ss.Technical...)
	all = append(all, ss.Domain...)
	all = append(all, ss.Soft...)
	all = append(all, ss.Tools...)
	all = append(all, ss.Languages...)
	return all
}

// Contains returns true if the skill set contains the given skill.
func (ss SkillSet) Contains(skill string) bool {
	skill = strings.ToLower(skill)
	for _, s := range ss.All() {
		if strings.ToLower(s) == skill {
			return true
		}
	}
	return false
}
