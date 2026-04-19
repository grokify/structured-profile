package export

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/grokify/structured-profile/schema"
)

// PandocExporter exports resumes to PDF and DOCX via Pandoc.
type PandocExporter struct {
	// PandocPath is the path to the pandoc binary.
	// If empty, uses "pandoc" from PATH.
	PandocPath string

	// Options for document formatting.
	Options PandocOptions
}

// PandocOptions contains formatting options for Pandoc export.
type PandocOptions struct {
	// Margin settings
	MarginTop    string // e.g., "1cm", "0.5in"
	MarginBottom string
	MarginLeft   string
	MarginRight  string

	// PDF settings
	PDFEngine string // e.g., "pdflatex", "xelatex", "lualatex"

	// Font settings
	MainFont string // e.g., "Times New Roman"
	FontSize string // e.g., "11pt", "12pt"

	// Template (optional)
	TemplatePath string

	// Additional raw Pandoc arguments
	ExtraArgs []string
}

// DefaultPandocOptions returns sensible default options.
func DefaultPandocOptions() PandocOptions {
	return PandocOptions{
		MarginTop:    "1cm",
		MarginBottom: "1cm",
		MarginLeft:   "1.5cm",
		MarginRight:  "1.5cm",
		PDFEngine:    "pdflatex",
		FontSize:     "11pt",
	}
}

// NewPandocExporter creates a new PandocExporter with default options.
func NewPandocExporter() *PandocExporter {
	return &PandocExporter{
		Options: DefaultPandocOptions(),
	}
}

// NewPandocExporterWithOptions creates a PandocExporter with custom options.
func NewPandocExporterWithOptions(opts PandocOptions) *PandocExporter {
	return &PandocExporter{
		Options: opts,
	}
}

// ExportFormat represents the output format for Pandoc export.
type ExportFormat string

const (
	FormatPDF  ExportFormat = "pdf"
	FormatDOCX ExportFormat = "docx"
	FormatHTML ExportFormat = "html"
)

// ExportResult contains the result of a Pandoc export.
type ExportResult struct {
	// OutputPath is the path to the generated file.
	OutputPath string

	// Markdown is the intermediate Markdown content.
	Markdown string

	// Stderr contains any warnings from Pandoc.
	Stderr string
}

// Export converts a Resume to the specified format.
func (e *PandocExporter) Export(resume *schema.Resume, format ExportFormat, outputPath string) (*ExportResult, error) {
	// First, convert to Markdown
	mdExporter := NewMarkdownExporter()
	markdown, err := mdExporter.Export(resume)
	if err != nil {
		return nil, fmt.Errorf("markdown export failed: %w", err)
	}

	return e.ExportMarkdown(markdown, format, outputPath)
}

// ExportMarkdown converts Markdown content to the specified format.
func (e *PandocExporter) ExportMarkdown(markdown string, format ExportFormat, outputPath string) (*ExportResult, error) {
	// Check if pandoc is available
	pandocPath := e.PandocPath
	if pandocPath == "" {
		pandocPath = "pandoc"
	}

	// Build Pandoc arguments
	args := e.buildArgs(format, outputPath)

	// Run Pandoc
	cmd := exec.Command(pandocPath, args...)
	cmd.Stdin = strings.NewReader(markdown)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("pandoc failed: %w\nstderr: %s", err, stderr.String())
	}

	return &ExportResult{
		OutputPath: outputPath,
		Markdown:   markdown,
		Stderr:     stderr.String(),
	}, nil
}

// buildArgs builds Pandoc command-line arguments.
func (e *PandocExporter) buildArgs(format ExportFormat, outputPath string) []string {
	args := []string{
		"-f", "markdown",
		"-t", string(format),
		"-o", outputPath,
	}

	// Add margins
	if e.Options.MarginTop != "" {
		args = append(args, "-V", fmt.Sprintf("margin-top=%s", e.Options.MarginTop))
	}
	if e.Options.MarginBottom != "" {
		args = append(args, "-V", fmt.Sprintf("margin-bottom=%s", e.Options.MarginBottom))
	}
	if e.Options.MarginLeft != "" {
		args = append(args, "-V", fmt.Sprintf("margin-left=%s", e.Options.MarginLeft))
	}
	if e.Options.MarginRight != "" {
		args = append(args, "-V", fmt.Sprintf("margin-right=%s", e.Options.MarginRight))
	}

	// PDF-specific options
	if format == FormatPDF {
		if e.Options.PDFEngine != "" {
			args = append(args, "--pdf-engine="+e.Options.PDFEngine)
		}
		if e.Options.MainFont != "" {
			args = append(args, "-V", fmt.Sprintf("mainfont=%s", e.Options.MainFont))
		}
		if e.Options.FontSize != "" {
			args = append(args, "-V", fmt.Sprintf("fontsize=%s", e.Options.FontSize))
		}
	}

	// Template
	if e.Options.TemplatePath != "" {
		args = append(args, "--template="+e.Options.TemplatePath)
	}

	// Extra args
	args = append(args, e.Options.ExtraArgs...)

	return args
}

// ExportToPDF is a convenience method for PDF export.
func (e *PandocExporter) ExportToPDF(resume *schema.Resume, outputPath string) (*ExportResult, error) {
	return e.Export(resume, FormatPDF, outputPath)
}

// ExportToDOCX is a convenience method for DOCX export.
func (e *PandocExporter) ExportToDOCX(resume *schema.Resume, outputPath string) (*ExportResult, error) {
	return e.Export(resume, FormatDOCX, outputPath)
}

// ExportToHTML is a convenience method for HTML export.
func (e *PandocExporter) ExportToHTML(resume *schema.Resume, outputPath string) (*ExportResult, error) {
	return e.Export(resume, FormatHTML, outputPath)
}

// IsPandocAvailable checks if Pandoc is installed and available.
func IsPandocAvailable() bool {
	_, err := exec.LookPath("pandoc")
	return err == nil
}

// GetPandocVersion returns the Pandoc version string.
func GetPandocVersion() (string, error) {
	cmd := exec.Command("pandoc", "--version")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	// First line contains version
	lines := strings.Split(string(out), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0]), nil
	}
	return "", fmt.Errorf("no version output")
}

// ExportResumeToFile is a high-level function to export a resume to file.
// It automatically determines format from the output filename extension.
func ExportResumeToFile(resume *schema.Resume, outputPath string) (*ExportResult, error) {
	ext := strings.ToLower(filepath.Ext(outputPath))

	var format ExportFormat
	switch ext {
	case ".pdf":
		format = FormatPDF
	case ".docx":
		format = FormatDOCX
	case ".html":
		format = FormatHTML
	case ".md", ".markdown":
		// Direct Markdown export
		mdExporter := NewMarkdownExporter()
		md, err := mdExporter.Export(resume)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(outputPath, []byte(md), 0644); err != nil {
			return nil, err
		}
		return &ExportResult{OutputPath: outputPath, Markdown: md}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", ext)
	}

	exporter := NewPandocExporter()
	return exporter.Export(resume, format, outputPath)
}

// ExportCoverLetterToFile exports a cover letter to file.
func ExportCoverLetterToFile(cl *schema.CoverLetter, outputPath string) (*ExportResult, error) {
	markdown, err := ExportCoverLetter(cl)
	if err != nil {
		return nil, err
	}

	ext := strings.ToLower(filepath.Ext(outputPath))

	switch ext {
	case ".md", ".markdown":
		if err := os.WriteFile(outputPath, []byte(markdown), 0644); err != nil {
			return nil, err
		}
		return &ExportResult{OutputPath: outputPath, Markdown: markdown}, nil
	case ".pdf":
		exporter := NewPandocExporter()
		return exporter.ExportMarkdown(markdown, FormatPDF, outputPath)
	case ".docx":
		exporter := NewPandocExporter()
		return exporter.ExportMarkdown(markdown, FormatDOCX, outputPath)
	default:
		return nil, fmt.Errorf("unsupported format: %s", ext)
	}
}
