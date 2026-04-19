// Package export provides resume and cover letter export functionality.
package export

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/grokify/structured-profile/schema"
)

// MarkdownExporter exports resumes to Markdown format.
type MarkdownExporter struct {
	// Template allows customization of the Markdown output.
	// If nil, uses the default template.
	Template *template.Template
}

// NewMarkdownExporter creates a new MarkdownExporter with the default template.
func NewMarkdownExporter() *MarkdownExporter {
	return &MarkdownExporter{}
}

// NewMarkdownExporterWithTemplate creates a MarkdownExporter with a custom template.
func NewMarkdownExporterWithTemplate(tmpl *template.Template) *MarkdownExporter {
	return &MarkdownExporter{Template: tmpl}
}

// Export converts a Resume to Markdown.
func (e *MarkdownExporter) Export(resume *schema.Resume) (string, error) {
	if resume == nil || resume.Content == nil {
		return "", fmt.Errorf("resume or content is nil")
	}

	if e.Template != nil {
		return e.exportWithTemplate(resume)
	}

	return e.exportDefault(resume)
}

// ExportContent converts ResumeContent directly to Markdown.
func (e *MarkdownExporter) ExportContent(content *schema.ResumeContent) (string, error) {
	if content == nil {
		return "", fmt.Errorf("content is nil")
	}

	resume := &schema.Resume{Content: content}
	return e.Export(resume)
}

// exportWithTemplate renders using a custom template.
func (e *MarkdownExporter) exportWithTemplate(resume *schema.Resume) (string, error) {
	var buf bytes.Buffer
	if err := e.Template.Execute(&buf, resume); err != nil {
		return "", fmt.Errorf("template execution failed: %w", err)
	}
	return buf.String(), nil
}

// exportDefault renders using the default format.
func (e *MarkdownExporter) exportDefault(resume *schema.Resume) (string, error) {
	var sections []string

	content := resume.Content

	// Header
	if header := e.renderHeader(content); header != "" {
		sections = append(sections, header)
	}

	// Summary
	if summary := e.renderSummary(content); summary != "" {
		sections = append(sections, summary)
	}

	// Experience
	if experience := e.renderExperience(content); experience != "" {
		sections = append(sections, experience)
	}

	// Skills
	if skills := e.renderSkills(content); skills != "" {
		sections = append(sections, skills)
	}

	// Education
	if education := e.renderEducation(content); education != "" {
		sections = append(sections, education)
	}

	// Certifications
	if certs := e.renderCertifications(content); certs != "" {
		sections = append(sections, certs)
	}

	return strings.Join(sections, "\n\n"), nil
}

// renderHeader renders the contact/header section.
func (e *MarkdownExporter) renderHeader(content *schema.ResumeContent) string {
	if content.Name == "" {
		return ""
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("# %s", content.Name))

	// Contact info line
	var contactParts []string
	if content.Email != "" {
		contactParts = append(contactParts, content.Email)
	}
	if content.Phone != "" {
		contactParts = append(contactParts, content.Phone)
	}
	if content.Location != "" {
		contactParts = append(contactParts, content.Location)
	}

	if len(contactParts) > 0 {
		lines = append(lines, strings.Join(contactParts, " | "))
	}

	// Links
	if len(content.Links) > 0 {
		var linkParts []string
		for _, link := range content.Links {
			text := link.Text
			if text == "" {
				text = link.Type
			}
			linkParts = append(linkParts, fmt.Sprintf("[%s](%s)", text, link.URL))
		}
		lines = append(lines, strings.Join(linkParts, " | "))
	}

	return strings.Join(lines, "\n")
}

// renderSummary renders the professional summary section.
func (e *MarkdownExporter) renderSummary(content *schema.ResumeContent) string {
	if content.Summary == "" {
		return ""
	}

	return fmt.Sprintf("## Summary\n\n%s", content.Summary)
}

// renderExperience renders the work experience section.
func (e *MarkdownExporter) renderExperience(content *schema.ResumeContent) string {
	if len(content.Experiences) == 0 {
		return ""
	}

	var lines []string
	lines = append(lines, "## Experience")

	for _, exp := range content.Experiences {
		// Title and company line
		dateStr := formatDateRange(exp.StartDate, exp.EndDate)
		lines = append(lines, fmt.Sprintf("\n### %s | %s", exp.Title, exp.Company))

		// Date and location
		var metaParts []string
		metaParts = append(metaParts, dateStr)
		if exp.Location != "" {
			metaParts = append(metaParts, exp.Location)
		}
		lines = append(lines, strings.Join(metaParts, " | "))

		// Description
		if exp.Description != "" {
			lines = append(lines, fmt.Sprintf("\n%s", exp.Description))
		}

		// Achievements
		if len(exp.Achievements) > 0 {
			lines = append(lines, "")
			for _, achievement := range exp.Achievements {
				lines = append(lines, fmt.Sprintf("- %s", achievement))
			}
		}
	}

	return strings.Join(lines, "\n")
}

// renderSkills renders the skills section.
func (e *MarkdownExporter) renderSkills(content *schema.ResumeContent) string {
	if len(content.Skills) == 0 {
		return ""
	}

	skillList := strings.Join(content.Skills, ", ")
	return fmt.Sprintf("## Skills\n\n%s", skillList)
}

// renderEducation renders the education section.
func (e *MarkdownExporter) renderEducation(content *schema.ResumeContent) string {
	if len(content.Education) == 0 {
		return ""
	}

	var lines []string
	lines = append(lines, "## Education")

	for _, edu := range content.Education {
		dateStr := formatDateRange(edu.StartDate, edu.EndDate)

		lines = append(lines, fmt.Sprintf("\n### %s", edu.Institution))

		var parts []string
		parts = append(parts, edu.Degree)
		if edu.Field != "" {
			parts = append(parts, edu.Field)
		}
		parts = append(parts, dateStr)

		lines = append(lines, strings.Join(parts, " | "))

		if edu.Honors != "" {
			lines = append(lines, fmt.Sprintf("*%s*", edu.Honors))
		}
	}

	return strings.Join(lines, "\n")
}

// renderCertifications renders the certifications section.
func (e *MarkdownExporter) renderCertifications(content *schema.ResumeContent) string {
	if len(content.Certifications) == 0 {
		return ""
	}

	var lines []string
	lines = append(lines, "## Certifications")
	lines = append(lines, "")

	for _, cert := range content.Certifications {
		line := fmt.Sprintf("- **%s**", cert.Name)
		if cert.Issuer != "" {
			line += fmt.Sprintf(" - %s", cert.Issuer)
		}
		if !cert.IssueDate.IsZero() {
			line += fmt.Sprintf(" (%s)", cert.IssueDate.DisplayString())
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// formatDateRange formats a date range for display.
func formatDateRange(start schema.Date, end *schema.Date) string {
	startStr := start.DisplayString()
	if startStr == "" || startStr == "Present" {
		startStr = ""
	}

	if end == nil || end.IsZero() {
		if startStr == "" {
			return ""
		}
		return fmt.Sprintf("%s - Present", startStr)
	}

	endStr := end.DisplayString()
	if startStr == "" {
		return endStr
	}

	return fmt.Sprintf("%s - %s", startStr, endStr)
}

// DefaultTemplate returns the default Markdown template.
func DefaultTemplate() *template.Template {
	tmpl := template.New("resume").Funcs(template.FuncMap{
		"join": strings.Join,
		"formatDateRange": func(start schema.Date, end *schema.Date) string {
			return formatDateRange(start, end)
		},
	})

	tmpl, _ = tmpl.Parse(defaultTemplateText)
	return tmpl
}

const defaultTemplateText = `# {{.Content.Name}}
{{if .Content.Email}}{{.Content.Email}}{{end}}{{if .Content.Phone}} | {{.Content.Phone}}{{end}}{{if .Content.Location}} | {{.Content.Location}}{{end}}

{{if .Content.Summary}}## Summary

{{.Content.Summary}}
{{end}}
{{if .Content.Experiences}}## Experience
{{range .Content.Experiences}}
### {{.Title}} | {{.Company}}
{{formatDateRange .StartDate .EndDate}}{{if .Location}} | {{.Location}}{{end}}
{{if .Description}}
{{.Description}}
{{end}}{{if .Achievements}}
{{range .Achievements}}- {{.}}
{{end}}{{end}}{{end}}{{end}}
{{if .Content.Skills}}## Skills

{{join .Content.Skills ", "}}
{{end}}
{{if .Content.Education}}## Education
{{range .Content.Education}}
### {{.Institution}}
{{.Degree}}{{if .Field}} | {{.Field}}{{end}} | {{formatDateRange .StartDate .EndDate}}
{{if .Honors}}*{{.Honors}}*{{end}}
{{end}}{{end}}
{{if .Content.Certifications}}## Certifications
{{range .Content.Certifications}}
- **{{.Name}}**{{if .Issuer}} - {{.Issuer}}{{end}}{{if not .IssueDate.IsZero}} ({{.IssueDate.DisplayString}}){{end}}
{{end}}{{end}}`

// ExportCoverLetter exports a cover letter to Markdown.
func ExportCoverLetter(cl *schema.CoverLetter) (string, error) {
	if cl == nil {
		return "", fmt.Errorf("cover letter is nil")
	}

	var lines []string

	// Header with date and recipient
	if cl.HiringManager != "" {
		lines = append(lines, fmt.Sprintf("Dear %s,", cl.HiringManager))
	} else {
		lines = append(lines, "Dear Hiring Manager,")
	}
	lines = append(lines, "")

	// Body
	if cl.Opening != "" {
		lines = append(lines, cl.Opening)
		lines = append(lines, "")
	}

	if cl.Body != "" {
		lines = append(lines, cl.Body)
		lines = append(lines, "")
	}

	if cl.Closing != "" {
		lines = append(lines, cl.Closing)
	}

	return strings.Join(lines, "\n"), nil
}
