package export

import (
	"strings"
	"testing"

	"github.com/grokify/structured-profile/schema"
)

func TestNewMarkdownExporter(t *testing.T) {
	e := NewMarkdownExporter()
	if e == nil {
		t.Fatal("expected non-nil exporter")
	}
}

func TestExportNil(t *testing.T) {
	e := NewMarkdownExporter()

	// Nil resume
	_, err := e.Export(nil)
	if err == nil {
		t.Error("expected error for nil resume")
	}

	// Nil content
	_, err = e.Export(&schema.Resume{})
	if err == nil {
		t.Error("expected error for nil content")
	}
}

func TestExportHeader(t *testing.T) {
	e := NewMarkdownExporter()

	content := &schema.ResumeContent{
		Name:     "John Doe",
		Email:    "john@example.com",
		Phone:    "555-1234",
		Location: "San Francisco, CA",
		Links: []schema.Link{
			{Type: "LinkedIn", URL: "https://linkedin.com/in/johndoe", Text: "LinkedIn"},
			{Type: "GitHub", URL: "https://github.com/johndoe"},
		},
	}

	md, err := e.ExportContent(content)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	// Check header
	if !strings.Contains(md, "# John Doe") {
		t.Error("expected name in header")
	}

	// Check contact info
	if !strings.Contains(md, "john@example.com") {
		t.Error("expected email")
	}
	if !strings.Contains(md, "555-1234") {
		t.Error("expected phone")
	}
	if !strings.Contains(md, "San Francisco, CA") {
		t.Error("expected location")
	}

	// Check links
	if !strings.Contains(md, "[LinkedIn](https://linkedin.com/in/johndoe)") {
		t.Error("expected LinkedIn link")
	}
	if !strings.Contains(md, "[GitHub](https://github.com/johndoe)") {
		t.Error("expected GitHub link")
	}
}

func TestExportSummary(t *testing.T) {
	e := NewMarkdownExporter()

	content := &schema.ResumeContent{
		Name:    "John Doe",
		Summary: "Experienced software engineer with expertise in backend systems and cloud infrastructure.",
	}

	md, err := e.ExportContent(content)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	if !strings.Contains(md, "## Summary") {
		t.Error("expected Summary section")
	}
	if !strings.Contains(md, "Experienced software engineer") {
		t.Error("expected summary content")
	}
}

func TestExportExperience(t *testing.T) {
	e := NewMarkdownExporter()

	endDate := schema.Date{Year: 2023, Month: 12}
	content := &schema.ResumeContent{
		Name: "John Doe",
		Experiences: []schema.ResumeExperience{
			{
				Company:     "Acme Corp",
				Title:       "Senior Software Engineer",
				Location:    "San Francisco, CA",
				StartDate:   schema.Date{Year: 2020, Month: 1},
				EndDate:     &endDate,
				Description: "Led backend development team",
				Achievements: []string{
					"Built microservices architecture reducing latency by 50%",
					"Mentored 3 junior engineers",
				},
			},
		},
	}

	md, err := e.ExportContent(content)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	// Check experience section
	if !strings.Contains(md, "## Experience") {
		t.Error("expected Experience section")
	}
	if !strings.Contains(md, "### Senior Software Engineer | Acme Corp") {
		t.Error("expected title and company")
	}
	if !strings.Contains(md, "Jan 2020 - Dec 2023") {
		t.Error("expected date range")
	}
	if !strings.Contains(md, "Led backend development team") {
		t.Error("expected description")
	}
	if !strings.Contains(md, "- Built microservices") {
		t.Error("expected achievement bullet")
	}
}

func TestExportCurrentPosition(t *testing.T) {
	e := NewMarkdownExporter()

	content := &schema.ResumeContent{
		Name: "John Doe",
		Experiences: []schema.ResumeExperience{
			{
				Company:   "Current Inc",
				Title:     "Staff Engineer",
				StartDate: schema.Date{Year: 2022, Month: 6},
				EndDate:   nil, // Current position
			},
		},
	}

	md, err := e.ExportContent(content)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	if !strings.Contains(md, "Jun 2022 - Present") {
		t.Error("expected 'Present' for current position")
	}
}

func TestExportSkills(t *testing.T) {
	e := NewMarkdownExporter()

	content := &schema.ResumeContent{
		Name:   "John Doe",
		Skills: []string{"Go", "Python", "Kubernetes", "AWS", "PostgreSQL"},
	}

	md, err := e.ExportContent(content)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	if !strings.Contains(md, "## Skills") {
		t.Error("expected Skills section")
	}
	if !strings.Contains(md, "Go, Python, Kubernetes, AWS, PostgreSQL") {
		t.Error("expected skills list")
	}
}

func TestExportEducation(t *testing.T) {
	e := NewMarkdownExporter()

	endDate := schema.Date{Year: 2018, Month: 5}
	content := &schema.ResumeContent{
		Name: "John Doe",
		Education: []schema.Education{
			{
				Institution: "MIT",
				Degree:      "BS Computer Science",
				Field:       "Computer Science",
				StartDate:   schema.Date{Year: 2014, Month: 9},
				EndDate:     &endDate,
				Honors:      "Magna Cum Laude",
			},
		},
	}

	md, err := e.ExportContent(content)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	if !strings.Contains(md, "## Education") {
		t.Error("expected Education section")
	}
	if !strings.Contains(md, "### MIT") {
		t.Error("expected institution")
	}
	if !strings.Contains(md, "BS Computer Science") {
		t.Error("expected degree")
	}
	if !strings.Contains(md, "*Magna Cum Laude*") {
		t.Error("expected honors in italics")
	}
}

func TestExportCertifications(t *testing.T) {
	e := NewMarkdownExporter()

	content := &schema.ResumeContent{
		Name: "John Doe",
		Certifications: []schema.Certification{
			{
				Name:      "AWS Solutions Architect",
				Issuer:    "Amazon Web Services",
				IssueDate: schema.Date{Year: 2023, Month: 3},
			},
			{
				Name:   "CKA",
				Issuer: "CNCF",
			},
		},
	}

	md, err := e.ExportContent(content)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	if !strings.Contains(md, "## Certifications") {
		t.Error("expected Certifications section")
	}
	if !strings.Contains(md, "**AWS Solutions Architect**") {
		t.Error("expected cert name in bold")
	}
	if !strings.Contains(md, "Amazon Web Services") {
		t.Error("expected issuer")
	}
	if !strings.Contains(md, "Mar 2023") {
		t.Error("expected issue date")
	}
}

func TestExportFullResume(t *testing.T) {
	e := NewMarkdownExporter()

	endDate := schema.Date{Year: 2023, Month: 12}
	eduEndDate := schema.Date{Year: 2018, Month: 5}

	resume := &schema.Resume{
		Content: &schema.ResumeContent{
			Name:     "Jane Developer",
			Email:    "jane@example.com",
			Phone:    "555-5678",
			Location: "Seattle, WA",
			Links: []schema.Link{
				{Type: "GitHub", URL: "https://github.com/jane"},
			},
			Summary: "Senior software engineer with 10+ years of experience.",
			Experiences: []schema.ResumeExperience{
				{
					Company:     "Tech Giant",
					Title:       "Staff Engineer",
					StartDate:   schema.Date{Year: 2020, Month: 1},
					EndDate:     nil, // Current
					Description: "Leading platform team",
					Achievements: []string{
						"Designed new API gateway",
						"Improved system reliability to 99.99%",
					},
				},
				{
					Company:   "Startup Inc",
					Title:     "Senior Engineer",
					StartDate: schema.Date{Year: 2018, Month: 3},
					EndDate:   &endDate,
					Achievements: []string{
						"Built microservices from scratch",
					},
				},
			},
			Skills: []string{"Go", "Kubernetes", "AWS", "PostgreSQL"},
			Education: []schema.Education{
				{
					Institution: "Stanford",
					Degree:      "MS Computer Science",
					EndDate:     &eduEndDate,
				},
			},
			Certifications: []schema.Certification{
				{Name: "AWS Certified", Issuer: "AWS"},
			},
		},
	}

	md, err := e.Export(resume)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	// Check all sections are present
	sections := []string{
		"# Jane Developer",
		"## Summary",
		"## Experience",
		"## Skills",
		"## Education",
		"## Certifications",
	}

	for _, section := range sections {
		if !strings.Contains(md, section) {
			t.Errorf("expected section %q", section)
		}
	}

	// Log the output for inspection
	t.Logf("Generated Markdown:\n%s", md)
}

func TestExportWithCustomTemplate(t *testing.T) {
	tmpl := DefaultTemplate()
	e := NewMarkdownExporterWithTemplate(tmpl)

	resume := &schema.Resume{
		Content: &schema.ResumeContent{
			Name:    "Template Test",
			Email:   "test@example.com",
			Summary: "Testing custom template.",
			Skills:  []string{"Go", "Testing"},
		},
	}

	md, err := e.Export(resume)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	if !strings.Contains(md, "# Template Test") {
		t.Error("expected name from template")
	}
}

func TestExportCoverLetter(t *testing.T) {
	cl := &schema.CoverLetter{
		HiringManager:  "John Smith",
		TargetCompany:  "Tech Corp",
		TargetPosition: "Staff Engineer",
		Opening:        "I am writing to express my interest in the Staff Engineer position.",
		Body:           "With my experience in distributed systems and team leadership, I believe I would be a great fit.",
		Closing:        "Thank you for your consideration. I look forward to hearing from you.",
	}

	md, err := ExportCoverLetter(cl)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	if !strings.Contains(md, "Dear John Smith,") {
		t.Error("expected hiring manager in greeting")
	}
	if !strings.Contains(md, "I am writing to express") {
		t.Error("expected opening")
	}
	if !strings.Contains(md, "distributed systems") {
		t.Error("expected body")
	}
	if !strings.Contains(md, "Thank you for your consideration") {
		t.Error("expected closing")
	}
}

func TestExportCoverLetterNoManager(t *testing.T) {
	cl := &schema.CoverLetter{
		Opening: "I am interested in joining your team.",
	}

	md, err := ExportCoverLetter(cl)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	if !strings.Contains(md, "Dear Hiring Manager,") {
		t.Error("expected default greeting")
	}
}

func TestExportCoverLetterNil(t *testing.T) {
	_, err := ExportCoverLetter(nil)
	if err == nil {
		t.Error("expected error for nil cover letter")
	}
}

func TestFormatDateRange(t *testing.T) {
	tests := []struct {
		name     string
		start    schema.Date
		end      *schema.Date
		expected string
	}{
		{
			name:     "with end date",
			start:    schema.Date{Year: 2020, Month: 1},
			end:      &schema.Date{Year: 2023, Month: 12},
			expected: "Jan 2020 - Dec 2023",
		},
		{
			name:     "current (nil end)",
			start:    schema.Date{Year: 2022, Month: 6},
			end:      nil,
			expected: "Jun 2022 - Present",
		},
		{
			name:     "current (zero end)",
			start:    schema.Date{Year: 2022, Month: 6},
			end:      &schema.Date{},
			expected: "Jun 2022 - Present",
		},
		{
			name:     "zero start",
			start:    schema.Date{},
			end:      &schema.Date{Year: 2020, Month: 5},
			expected: "May 2020",
		},
		{
			name:     "both zero",
			start:    schema.Date{},
			end:      nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDateRange(tt.start, tt.end)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestEmptySections(t *testing.T) {
	e := NewMarkdownExporter()

	// Minimal content
	content := &schema.ResumeContent{
		Name: "Minimal User",
	}

	md, err := e.ExportContent(content)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	// Should only have name, no section headers
	if strings.Contains(md, "## Summary") {
		t.Error("should not have empty Summary section")
	}
	if strings.Contains(md, "## Experience") {
		t.Error("should not have empty Experience section")
	}
	if strings.Contains(md, "## Skills") {
		t.Error("should not have empty Skills section")
	}
	if strings.Contains(md, "## Education") {
		t.Error("should not have empty Education section")
	}
	if strings.Contains(md, "## Certifications") {
		t.Error("should not have empty Certifications section")
	}
}
