package schema

// InterviewPrepSet represents a set of interview preparation questions.
type InterviewPrepSet struct {
	BaseEntity
	ProfileID     string `json:"profile_id,omitempty"`
	OpportunityID string `json:"opportunity_id,omitempty"` // Optional: role-specific

	Title       string `json:"title"`
	Description string `json:"description,omitempty"`

	Sections []PrepSection `json:"sections,omitempty"`

	// H5P Export Settings
	PassPercentage int             `json:"pass_percentage,omitempty"`
	FeedbackRanges []FeedbackRange `json:"feedback_ranges,omitempty"`
}

// NewInterviewPrepSet creates a new InterviewPrepSet.
func NewInterviewPrepSet(title string) *InterviewPrepSet {
	return &InterviewPrepSet{
		BaseEntity:     NewBaseEntity(),
		Title:          title,
		Sections:       []PrepSection{},
		PassPercentage: 70,
	}
}

// AddSection adds a section to the prep set.
func (ps *InterviewPrepSet) AddSection(section PrepSection) {
	ps.Sections = append(ps.Sections, section)
	ps.Touch()
}

// TotalQuestions returns the total number of questions across all sections.
func (ps *InterviewPrepSet) TotalQuestions() int {
	count := 0
	for _, s := range ps.Sections {
		count += len(s.Questions)
	}
	return count
}

// PrepSection represents a section of interview prep questions.
type PrepSection struct {
	Name      string         `json:"name"`  // "Behavioral", "Technical"
	Topic     string         `json:"topic"` // "Leadership", "System Design"
	Questions []PrepQuestion `json:"questions,omitempty"`
}

// NewPrepSection creates a new PrepSection.
func NewPrepSection(name, topic string) PrepSection {
	return PrepSection{
		Name:      name,
		Topic:     topic,
		Questions: []PrepQuestion{},
	}
}

// AddQuestion adds a question to the section.
func (ps *PrepSection) AddQuestion(q PrepQuestion) {
	ps.Questions = append(ps.Questions, q)
}

// PrepQuestion represents an interview prep question.
type PrepQuestion struct {
	BaseEntity
	Question   string `json:"question"`
	Difficulty string `json:"difficulty,omitempty"` // "easy", "medium", "hard"
	Category   string `json:"category,omitempty"`   // "behavioral", "technical", etc.

	// Answer Options (for multiple choice)
	Answers []PrepAnswer `json:"answers,omitempty"`

	// For open-ended questions
	SampleAnswer string `json:"sample_answer,omitempty"`
	KeyPoints    []string `json:"key_points,omitempty"`

	// Linked Content
	RelatedSTARID string   `json:"related_star_id,omitempty"`
	RelatedSkills []string `json:"related_skills,omitempty"`

	// Learning
	LearningObjective string `json:"learning_objective,omitempty"`
	Explanation       string `json:"explanation,omitempty"`

	// Source Tracking
	Source          string `json:"source,omitempty"`           // "interview:app-123", "book:xyz"
	SourceInterview string `json:"source_interview,omitempty"` // Interview ID if from real interview
}

// NewPrepQuestion creates a new PrepQuestion.
func NewPrepQuestion(question string) *PrepQuestion {
	return &PrepQuestion{
		BaseEntity: NewBaseEntity(),
		Question:   question,
		Answers:    []PrepAnswer{},
	}
}

// AddAnswer adds an answer option to the question.
func (pq *PrepQuestion) AddAnswer(answer PrepAnswer) {
	pq.Answers = append(pq.Answers, answer)
}

// IsMultipleChoice returns true if the question has answer options.
func (pq *PrepQuestion) IsMultipleChoice() bool {
	return len(pq.Answers) > 0
}

// PrepAnswer represents an answer option for a prep question.
type PrepAnswer struct {
	Text     string `json:"text"`
	Correct  bool   `json:"correct"`
	Feedback string `json:"feedback,omitempty"`
	Tip      string `json:"tip,omitempty"`
}

// NewPrepAnswer creates a new PrepAnswer.
func NewPrepAnswer(text string, correct bool) PrepAnswer {
	return PrepAnswer{
		Text:    text,
		Correct: correct,
	}
}

// FeedbackRange represents score-based feedback for H5P export.
type FeedbackRange struct {
	From int    `json:"from"`
	To   int    `json:"to"`
	Text string `json:"text"`
}

// NewFeedbackRange creates a new FeedbackRange.
func NewFeedbackRange(from, to int, text string) FeedbackRange {
	return FeedbackRange{
		From: from,
		To:   to,
		Text: text,
	}
}

// DefaultFeedbackRanges returns sensible default feedback ranges.
func DefaultFeedbackRanges() []FeedbackRange {
	return []FeedbackRange{
		{From: 0, To: 50, Text: "Needs more practice. Review the material and try again."},
		{From: 51, To: 70, Text: "Good effort! Some areas need improvement."},
		{From: 71, To: 85, Text: "Well done! You're on the right track."},
		{From: 86, To: 100, Text: "Excellent! You're well prepared."},
	}
}

// QuestionDifficulty represents question difficulty levels.
type QuestionDifficulty string

const (
	DifficultyEasy   QuestionDifficulty = "easy"
	DifficultyMedium QuestionDifficulty = "medium"
	DifficultyHard   QuestionDifficulty = "hard"
)

// QuestionCategory represents question categories.
type QuestionCategory string

const (
	QuestionCategoryBehavioral  QuestionCategory = "behavioral"
	QuestionCategoryTechnical   QuestionCategory = "technical"
	QuestionCategorySituational QuestionCategory = "situational"
	QuestionCategoryCase        QuestionCategory = "case"
	QuestionCategoryBrainteaser QuestionCategory = "brainteaser"
)
