package schema

import "time"

// Application represents a job application and its lifecycle.
type Application struct {
	BaseEntity
	ProfileID     string `json:"profile_id,omitempty"`
	OpportunityID string `json:"opportunity_id"`

	// Status Tracking
	Status    ApplicationStatus `json:"status"`
	AppliedAt *time.Time        `json:"applied_at,omitempty"`

	// Documents Used (version identifiers)
	ResumeVersion      string `json:"resume_version,omitempty"`
	CoverLetterVersion string `json:"cover_letter_version,omitempty"`

	// Interview Rounds
	Interviews []Interview `json:"interviews,omitempty"`

	// Outcome
	Outcome *ApplicationOutcome `json:"outcome,omitempty"`

	// Notes
	Notes string `json:"notes,omitempty"`
}

// NewApplication creates a new Application for an opportunity.
func NewApplication(opportunityID string) *Application {
	return &Application{
		BaseEntity:    NewBaseEntity(),
		OpportunityID: opportunityID,
		Status:        ApplicationStatusDraft,
		Interviews:    []Interview{},
	}
}

// Submit marks the application as submitted.
func (a *Application) Submit() {
	now := time.Now().UTC()
	a.AppliedAt = &now
	a.Status = ApplicationStatusSubmitted
	a.Touch()
}

// AddInterview adds an interview round to the application.
func (a *Application) AddInterview(interview Interview) {
	interview.ApplicationID = a.ID
	a.Interviews = append(a.Interviews, interview)
	a.Status = ApplicationStatusInterview
	a.Touch()
}

// ApplicationStatus represents the status of an application.
type ApplicationStatus string

const (
	ApplicationStatusDraft      ApplicationStatus = "draft"
	ApplicationStatusSubmitted  ApplicationStatus = "submitted"
	ApplicationStatusScreening  ApplicationStatus = "screening"
	ApplicationStatusInterview  ApplicationStatus = "interview"
	ApplicationStatusOffer      ApplicationStatus = "offer"
	ApplicationStatusAccepted   ApplicationStatus = "accepted"
	ApplicationStatusRejected   ApplicationStatus = "rejected"
	ApplicationStatusWithdrawn  ApplicationStatus = "withdrawn"
	ApplicationStatusOnHold     ApplicationStatus = "on_hold"
)

// IsActive returns true if the application is still in progress.
func (s ApplicationStatus) IsActive() bool {
	switch s {
	case ApplicationStatusDraft, ApplicationStatusSubmitted,
		ApplicationStatusScreening, ApplicationStatusInterview,
		ApplicationStatusOffer, ApplicationStatusOnHold:
		return true
	}
	return false
}

// IsFinal returns true if the application has reached a final state.
func (s ApplicationStatus) IsFinal() bool {
	switch s {
	case ApplicationStatusAccepted, ApplicationStatusRejected, ApplicationStatusWithdrawn:
		return true
	}
	return false
}

// ApplicationOutcome contains the result of an application.
type ApplicationOutcome struct {
	Result          string   `json:"result"` // "offer", "rejected", "withdrawn"
	OfferDetails    string   `json:"offer_details,omitempty"`
	RejectionReason string   `json:"rejection_reason,omitempty"`
	LessonsLearned  []string `json:"lessons_learned,omitempty"`
	Strengths       []string `json:"strengths,omitempty"`
	Improvements    []string `json:"improvements,omitempty"`
	DecisionDate    Date     `json:"decision_date,omitempty"`
}

// Interview represents an interview round within an application.
type Interview struct {
	BaseEntity
	ApplicationID string `json:"application_id,omitempty"`

	// Interview Details
	Round        int           `json:"round"`
	Type         InterviewType `json:"type"`
	ScheduledAt  *time.Time    `json:"scheduled_at,omitempty"`
	DurationMins int           `json:"duration_mins,omitempty"`

	// Participants
	Interviewers []string `json:"interviewers,omitempty"`

	// Content
	Questions []InterviewQuestion `json:"questions,omitempty"`

	// Assessment
	SelfAssessment *SelfAssessment    `json:"self_assessment,omitempty"`
	Feedback       *InterviewFeedback `json:"feedback,omitempty"`

	// Notes
	Notes string `json:"notes,omitempty"`
}

// NewInterview creates a new Interview.
func NewInterview(round int, interviewType InterviewType) *Interview {
	return &Interview{
		BaseEntity: NewBaseEntity(),
		Round:      round,
		Type:       interviewType,
		Questions:  []InterviewQuestion{},
	}
}

// AddQuestion adds a question to the interview record.
func (i *Interview) AddQuestion(q InterviewQuestion) {
	i.Questions = append(i.Questions, q)
	i.Touch()
}

// InterviewType represents types of interviews.
type InterviewType string

const (
	InterviewTypePhone      InterviewType = "phone"
	InterviewTypeTechnical  InterviewType = "technical"
	InterviewTypeBehavioral InterviewType = "behavioral"
	InterviewTypeSystem     InterviewType = "system_design"
	InterviewTypeCoding     InterviewType = "coding"
	InterviewTypeOnsite     InterviewType = "onsite"
	InterviewTypePanel      InterviewType = "panel"
	InterviewTypeHiring     InterviewType = "hiring_manager"
	InterviewTypeFinal      InterviewType = "final"
	InterviewTypeOther      InterviewType = "other"
)

// InterviewQuestion represents a question asked during an interview.
type InterviewQuestion struct {
	Question       string   `json:"question"`
	Category       string   `json:"category,omitempty"` // "behavioral", "technical", "situational"
	MyAnswer       string   `json:"my_answer,omitempty"`
	IdealAnswer    string   `json:"ideal_answer,omitempty"`
	RelatedSTARID  string   `json:"related_star_id,omitempty"` // Achievement ID
	Tags           []string `json:"tags,omitempty"`
	Difficulty     string   `json:"difficulty,omitempty"` // "easy", "medium", "hard"
	AddToPrep      bool     `json:"add_to_prep,omitempty"` // Flag to add to prep bank
	WentWell       bool     `json:"went_well,omitempty"`
	Notes          string   `json:"notes,omitempty"`
}

// SelfAssessment contains self-evaluation after an interview.
type SelfAssessment struct {
	OverallRating  int      `json:"overall_rating"` // 1-5
	WentWell       []string `json:"went_well,omitempty"`
	CouldImprove   []string `json:"could_improve,omitempty"`
	QuestionsToAdd []string `json:"questions_to_add,omitempty"`
	KeyTakeaways   []string `json:"key_takeaways,omitempty"`
	Notes          string   `json:"notes,omitempty"`
}

// InterviewFeedback contains feedback received from interviewers.
type InterviewFeedback struct {
	Source         string   `json:"source"`          // "recruiter", "interviewer", "self"
	Rating         string   `json:"rating,omitempty"` // "strong_yes", "yes", "no", "strong_no"
	Strengths      []string `json:"strengths,omitempty"`
	Concerns       []string `json:"concerns,omitempty"`
	Recommendation string   `json:"recommendation,omitempty"`
	Notes          string   `json:"notes,omitempty"`
}

// FeedbackRating represents interview feedback ratings.
type FeedbackRating string

const (
	FeedbackRatingStrongYes FeedbackRating = "strong_yes"
	FeedbackRatingYes       FeedbackRating = "yes"
	FeedbackRatingMaybe     FeedbackRating = "maybe"
	FeedbackRatingNo        FeedbackRating = "no"
	FeedbackRatingStrongNo  FeedbackRating = "strong_no"
)
