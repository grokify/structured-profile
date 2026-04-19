package export

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grokify/structured-profile/schema"
)

func TestNewPandocExporter(t *testing.T) {
	e := NewPandocExporter()
	if e == nil {
		t.Fatal("expected non-nil exporter")
	}

	// Check default options
	if e.Options.MarginTop != "1cm" {
		t.Errorf("expected default margin-top 1cm, got %s", e.Options.MarginTop)
	}
}

func TestDefaultPandocOptions(t *testing.T) {
	opts := DefaultPandocOptions()

	if opts.MarginTop != "1cm" {
		t.Error("expected default margin-top")
	}
	if opts.PDFEngine != "pdflatex" {
		t.Error("expected default PDF engine")
	}
	if opts.FontSize != "11pt" {
		t.Error("expected default font size")
	}
}

func TestBuildArgs(t *testing.T) {
	e := NewPandocExporter()
	e.Options = PandocOptions{
		MarginTop:    "2cm",
		MarginBottom: "2cm",
		MarginLeft:   "1in",
		MarginRight:  "1in",
		PDFEngine:    "xelatex",
		MainFont:     "Arial",
		FontSize:     "12pt",
		ExtraArgs:    []string{"--standalone"},
	}

	args := e.buildArgs(FormatPDF, "output.pdf")

	// Check format args
	if !containsArg(args, "-f", "markdown") {
		t.Error("expected markdown input format")
	}
	if !containsArg(args, "-t", "pdf") {
		t.Error("expected pdf output format")
	}
	if !containsArg(args, "-o", "output.pdf") {
		t.Error("expected output path")
	}

	// Check margins
	if !containsVar(args, "margin-top=2cm") {
		t.Error("expected margin-top variable")
	}
	if !containsVar(args, "margin-left=1in") {
		t.Error("expected margin-left variable")
	}

	// Check PDF options
	if !containsString(args, "--pdf-engine=xelatex") {
		t.Error("expected pdf-engine")
	}
	if !containsVar(args, "mainfont=Arial") {
		t.Error("expected mainfont variable")
	}

	// Check extra args
	if !containsString(args, "--standalone") {
		t.Error("expected extra arg --standalone")
	}
}

func TestIsPandocAvailable(t *testing.T) {
	available := IsPandocAvailable()
	if available {
		t.Log("Pandoc is available")
		version, err := GetPandocVersion()
		if err != nil {
			t.Logf("Could not get version: %v", err)
		} else {
			t.Logf("Pandoc version: %s", version)
		}
	} else {
		t.Log("Pandoc is not available - PDF/DOCX tests will be skipped")
	}
}

func TestExportToPDF(t *testing.T) {
	if !IsPandocAvailable() {
		t.Skip("Pandoc not available")
	}

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "pandoc-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test resume
	resume := createTestResume()
	outputPath := filepath.Join(tmpDir, "resume.pdf")

	e := NewPandocExporter()
	result, err := e.ExportToPDF(resume, outputPath)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	// Check file exists
	if _, err := os.Stat(result.OutputPath); os.IsNotExist(err) {
		t.Error("output file was not created")
	}

	// Check file has content
	info, err := os.Stat(result.OutputPath)
	if err != nil {
		t.Fatalf("failed to stat output: %v", err)
	}
	if info.Size() == 0 {
		t.Error("output file is empty")
	}

	t.Logf("Created PDF: %s (%d bytes)", result.OutputPath, info.Size())
}

func TestExportToDOCX(t *testing.T) {
	if !IsPandocAvailable() {
		t.Skip("Pandoc not available")
	}

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "pandoc-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test resume
	resume := createTestResume()
	outputPath := filepath.Join(tmpDir, "resume.docx")

	e := NewPandocExporter()
	result, err := e.ExportToDOCX(resume, outputPath)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	// Check file exists
	if _, err := os.Stat(result.OutputPath); os.IsNotExist(err) {
		t.Error("output file was not created")
	}

	info, err := os.Stat(result.OutputPath)
	if err != nil {
		t.Fatalf("failed to stat output: %v", err)
	}
	if info.Size() == 0 {
		t.Error("output file is empty")
	}

	t.Logf("Created DOCX: %s (%d bytes)", result.OutputPath, info.Size())
}

func TestExportToHTML(t *testing.T) {
	if !IsPandocAvailable() {
		t.Skip("Pandoc not available")
	}

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "pandoc-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test resume
	resume := createTestResume()
	outputPath := filepath.Join(tmpDir, "resume.html")

	e := NewPandocExporter()
	result, err := e.ExportToHTML(resume, outputPath)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	// Check file exists and has HTML content
	content, err := os.ReadFile(result.OutputPath)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	if !strings.Contains(string(content), "Jane Developer") {
		t.Error("HTML should contain name")
	}

	t.Logf("Created HTML: %s (%d bytes)", result.OutputPath, len(content))
}

func TestExportResumeToFile(t *testing.T) {
	if !IsPandocAvailable() {
		t.Skip("Pandoc not available")
	}

	tmpDir, err := os.MkdirTemp("", "pandoc-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	resume := createTestResume()

	// Test PDF by extension
	pdfPath := filepath.Join(tmpDir, "resume.pdf")
	result, err := ExportResumeToFile(resume, pdfPath)
	if err != nil {
		t.Fatalf("PDF export failed: %v", err)
	}
	if result.OutputPath != pdfPath {
		t.Error("wrong output path")
	}

	// Test Markdown export (no Pandoc needed)
	mdPath := filepath.Join(tmpDir, "resume.md")
	result, err = ExportResumeToFile(resume, mdPath)
	if err != nil {
		t.Fatalf("Markdown export failed: %v", err)
	}
	content, _ := os.ReadFile(mdPath)
	if !strings.Contains(string(content), "# Jane Developer") {
		t.Error("Markdown should contain name")
	}
}

func TestExportResumeToMarkdown(t *testing.T) {
	// This test doesn't need Pandoc
	tmpDir, err := os.MkdirTemp("", "export-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	resume := createTestResume()
	mdPath := filepath.Join(tmpDir, "resume.md")

	result, err := ExportResumeToFile(resume, mdPath)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	// Check Markdown in result
	if result.Markdown == "" {
		t.Error("expected markdown in result")
	}
	if !strings.Contains(result.Markdown, "Jane Developer") {
		t.Error("markdown should contain name")
	}
}

func TestExportCoverLetterToFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "export-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cl := &schema.CoverLetter{
		HiringManager: "John Smith",
		Opening:       "I am writing to express my interest.",
		Body:          "My experience includes...",
		Closing:       "Thank you for your consideration.",
	}

	// Test Markdown export
	mdPath := filepath.Join(tmpDir, "cover.md")
	result, err := ExportCoverLetterToFile(cl, mdPath)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	content, _ := os.ReadFile(mdPath)
	if !strings.Contains(string(content), "Dear John Smith") {
		t.Error("should contain greeting")
	}

	// Test PDF if Pandoc available
	if IsPandocAvailable() {
		pdfPath := filepath.Join(tmpDir, "cover.pdf")
		result, err = ExportCoverLetterToFile(cl, pdfPath)
		if err != nil {
			t.Fatalf("PDF export failed: %v", err)
		}
		info, _ := os.Stat(result.OutputPath)
		t.Logf("Created cover letter PDF: %d bytes", info.Size())
	}
}

func TestExportUnsupportedFormat(t *testing.T) {
	resume := createTestResume()
	_, err := ExportResumeToFile(resume, "output.xyz")
	if err == nil {
		t.Error("expected error for unsupported format")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("wrong error message: %v", err)
	}
}

func TestExportWithCustomOptions(t *testing.T) {
	if !IsPandocAvailable() {
		t.Skip("Pandoc not available")
	}

	tmpDir, err := os.MkdirTemp("", "pandoc-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	opts := PandocOptions{
		MarginTop:    "0.5in",
		MarginBottom: "0.5in",
		MarginLeft:   "0.75in",
		MarginRight:  "0.75in",
		FontSize:     "10pt",
	}

	e := NewPandocExporterWithOptions(opts)
	resume := createTestResume()
	outputPath := filepath.Join(tmpDir, "resume-custom.pdf")

	result, err := e.ExportToPDF(resume, outputPath)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	info, _ := os.Stat(result.OutputPath)
	t.Logf("Created PDF with custom options: %d bytes", info.Size())
}

// Helper functions

func createTestResume() *schema.Resume {
	endDate := schema.Date{Year: 2023, Month: 12}
	eduEndDate := schema.Date{Year: 2018, Month: 5}

	return &schema.Resume{
		Content: &schema.ResumeContent{
			Name:     "Jane Developer",
			Email:    "jane@example.com",
			Phone:    "555-5678",
			Location: "Seattle, WA",
			Links: []schema.Link{
				{Type: "LinkedIn", URL: "https://linkedin.com/in/jane"},
				{Type: "GitHub", URL: "https://github.com/jane"},
			},
			Summary: "Senior software engineer with 10+ years of experience in building scalable systems.",
			Experiences: []schema.ResumeExperience{
				{
					Company:   "Tech Giant",
					Title:     "Staff Engineer",
					StartDate: schema.Date{Year: 2020, Month: 1},
					EndDate:   nil,
					Achievements: []string{
						"Led redesign of core API platform serving 1M+ requests/day",
						"Mentored team of 5 engineers",
					},
				},
				{
					Company:   "Startup Inc",
					Title:     "Senior Engineer",
					StartDate: schema.Date{Year: 2018, Month: 3},
					EndDate:   &endDate,
					Achievements: []string{
						"Built microservices architecture from scratch",
						"Reduced deployment time by 90%",
					},
				},
			},
			Skills: []string{
				"Go", "Python", "TypeScript",
				"Kubernetes", "AWS", "GCP",
				"PostgreSQL", "Redis", "Kafka",
			},
			Education: []schema.Education{
				{
					Institution: "Stanford University",
					Degree:      "MS Computer Science",
					EndDate:     &eduEndDate,
				},
			},
			Certifications: []schema.Certification{
				{Name: "AWS Solutions Architect", Issuer: "AWS"},
				{Name: "CKA", Issuer: "CNCF"},
			},
		},
	}
}

func containsArg(args []string, flag, value string) bool {
	for i := 0; i < len(args)-1; i++ {
		if args[i] == flag && args[i+1] == value {
			return true
		}
	}
	return false
}

func containsVar(args []string, varValue string) bool {
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "-V" && args[i+1] == varValue {
			return true
		}
	}
	return false
}

func containsString(args []string, s string) bool {
	for _, arg := range args {
		if arg == s {
			return true
		}
	}
	return false
}
