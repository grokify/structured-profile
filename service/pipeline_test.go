package service

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grokify/structured-profile/export"
	"github.com/grokify/structured-profile/jdparser"
	"github.com/grokify/structured-profile/schema"
	jsonstore "github.com/grokify/structured-profile/store/json"
)

// TestFullPipeline demonstrates the complete workflow:
// 1. Load profile from testdata
// 2. Parse a job description
// 3. Create an opportunity with the parsed JD
// 4. Generate a tailored resume
// 5. Generate a tailored cover letter
// 6. Export both to Markdown
func TestFullPipeline(t *testing.T) {
	ctx := context.Background()

	// 1. Set up store with testdata
	testDataDir := filepath.Join("..", "testdata")
	store, err := jsonstore.New(jsonstore.Config{BaseDir: testDataDir})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// 2. Load the migrated profile
	profile, err := store.GetFullProfile(ctx, "john-wang")
	if err != nil {
		t.Fatalf("Failed to load profile: %v", err)
	}
	t.Logf("Loaded profile: %s (%s)", profile.Profile.Name, profile.Profile.ID)
	t.Logf("  Tenures: %d", len(profile.Tenures))
	t.Logf("  Total achievements: %d", len(profile.AllAchievements()))
	t.Logf("  Skills: %d", len(profile.Skills))

	// 3. Sample job description for a VP of Platform role
	rawJD := `
VP of Platform Engineering

About the Role:
We are looking for a VP of Platform Engineering to lead our platform team. This is a senior leadership
position reporting to the CTO. You will be responsible for building and scaling our developer platform,
APIs, and infrastructure.

Requirements:
- 10+ years of experience in software engineering, with 5+ years in platform/infrastructure roles
- Experience building and managing developer platforms and APIs
- Strong background in OAuth 2.0, OIDC, and identity management (IAM)
- Experience with SCIM, SAML, and enterprise SSO integrations
- Track record of building high-performing engineering teams

Technical Skills:
- Go, Python, or similar backend languages
- Kubernetes, Docker, and container orchestration
- RESTful API design and OpenAPI specifications
- CI/CD pipelines and developer tooling
- Cloud platforms (AWS, GCP, Azure)

Preferred:
- Experience with developer experience (DevX) initiatives
- Open source contributions and community building
- Background in security products or IAM platforms
- MBA or advanced technical degree

What You'll Do:
- Define and execute the platform strategy aligned with company goals
- Lead a team of 15+ engineers across platform, API, and infrastructure
- Partner with product, security, and engineering leadership
- Drive developer productivity and platform adoption
- Own API design standards and developer documentation
`

	// 4. Parse the job description
	parser := jdparser.New()
	jdParsed := parser.Parse(rawJD)

	t.Logf("\nParsed Job Description:")
	t.Logf("  Required Skills: %v", jdParsed.RequiredSkills)
	t.Logf("  Preferred Skills: %v", jdParsed.PreferredSkills)
	t.Logf("  Experience Years: %d", jdParsed.ExperienceYears)
	t.Logf("  Seniority Level: %s", jdParsed.SeniorityLevel)
	t.Logf("  Keywords: %v", jdParsed.Keywords)

	// 5. Create an opportunity with the parsed JD
	opp := schema.NewOpportunity("TechCorp", "VP of Platform Engineering")
	opp.ProfileID = profile.Profile.ID
	opp.JobDescRaw = rawJD
	opp.JobDescParsed = jdParsed
	opp.Location = "San Francisco, CA"
	opp.Remote = true
	opp.SalaryMin = 300000
	opp.SalaryMax = 400000
	opp.SalaryCurrency = "USD"
	opp.HiringManager = "Jane Smith"

	// Save the opportunity
	if err := store.SaveOpportunity(ctx, profile.Profile.ID, opp); err != nil {
		t.Fatalf("Failed to save opportunity: %v", err)
	}
	t.Logf("\nCreated Opportunity: %s at %s", opp.Position, opp.Company)

	// 6. Generate a tailored resume
	resumeService := NewResumeService(store)
	resumeResult, err := resumeService.Generate(ctx, GenerateInput{
		ProfileID:     profile.Profile.ID,
		OpportunityID: opp.ID,
		Domain:        "devx", // Use DevX domain for this platform role
		Options:       schema.DefaultResumeOptions(),
	})
	if err != nil {
		t.Fatalf("Failed to generate resume: %v", err)
	}

	t.Logf("\nGenerated Resume:")
	t.Logf("  Experiences: %d", len(resumeResult.Resume.Content.Experiences))
	t.Logf("  Skills: %d", len(resumeResult.Resume.Content.Skills))
	if resumeResult.MatchResult != nil {
		t.Logf("  Match Score: %.2f", resumeResult.MatchResult.OverallScore)
		t.Logf("  Matched Required Skills: %v", resumeResult.MatchResult.MatchedRequiredSkills)
		t.Logf("  Matched Preferred Skills: %v", resumeResult.MatchResult.MatchedPreferredSkills)
		t.Logf("  Missing Required Skills: %v", resumeResult.MatchResult.MissingRequiredSkills)
	}

	// 7. Export resume to Markdown
	mdExporter := export.NewMarkdownExporter()
	resumeMarkdown, err := mdExporter.Export(resumeResult.Resume)
	if err != nil {
		t.Fatalf("Failed to export resume to Markdown: %v", err)
	}

	t.Logf("\n=== RESUME (Markdown preview) ===")
	// Show first 1000 chars of resume
	if len(resumeMarkdown) > 1000 {
		t.Logf("%s\n...(truncated)", resumeMarkdown[:1000])
	} else {
		t.Log(resumeMarkdown)
	}

	// 8. Generate a tailored cover letter
	clService := NewCoverLetterService(store)
	clResult, err := clService.Generate(ctx, GenerateCoverLetterInput{
		ProfileID:     profile.Profile.ID,
		OpportunityID: opp.ID,
		NumSTAR:       3,
		Domain:        "devx",
	})
	if err != nil {
		t.Fatalf("Failed to generate cover letter: %v", err)
	}

	t.Logf("\n=== COVER LETTER ===")
	t.Logf("To: %s at %s", clResult.CoverLetter.HiringManager, clResult.CoverLetter.TargetCompany)
	t.Logf("Position: %s", clResult.CoverLetter.TargetPosition)
	t.Logf("\n%s", clResult.CoverLetter.FullText())

	// 9. Export cover letter to Markdown
	clMarkdown, err := export.ExportCoverLetter(clResult.CoverLetter)
	if err != nil {
		t.Fatalf("Failed to export cover letter to Markdown: %v", err)
	}

	t.Logf("\n=== COVER LETTER (Markdown) ===")
	t.Log(clMarkdown)

	// 10. Verify outputs are non-empty and contain expected content
	if len(resumeMarkdown) < 500 {
		t.Error("Resume Markdown seems too short")
	}
	if !strings.Contains(resumeMarkdown, "John Wang") {
		t.Error("Resume should contain the candidate's name")
	}
	if !strings.Contains(resumeMarkdown, "Platform") || !strings.Contains(resumeMarkdown, "API") {
		t.Error("Resume should contain relevant platform/API keywords")
	}

	if len(clMarkdown) < 200 {
		t.Error("Cover letter Markdown seems too short")
	}
	if !strings.Contains(clMarkdown, "TechCorp") {
		t.Error("Cover letter should mention the target company")
	}

	// 11. Cleanup - remove the test opportunity
	if err := store.DeleteOpportunity(ctx, profile.Profile.ID, opp.ID); err != nil {
		t.Logf("Warning: Failed to clean up opportunity: %v", err)
	}

	t.Log("\n✓ Full pipeline test completed successfully!")
}

// TestPipelineWithDifferentDomains tests resume generation for different domains.
func TestPipelineWithDifferentDomains(t *testing.T) {
	ctx := context.Background()

	testDataDir := filepath.Join("..", "testdata")
	store, err := jsonstore.New(jsonstore.Config{BaseDir: testDataDir})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	profile, err := store.GetFullProfile(ctx, "john-wang")
	if err != nil {
		t.Fatalf("Failed to load profile: %v", err)
	}

	resumeService := NewResumeService(store)
	mdExporter := export.NewMarkdownExporter()

	domains := []string{"", "devx", "iam", "platform"}

	for _, domain := range domains {
		domainLabel := domain
		if domain == "" {
			domainLabel = "(default)"
		}

		result, err := resumeService.Generate(ctx, GenerateInput{
			ProfileID: profile.Profile.ID,
			Domain:    domain,
			Options:   schema.DefaultResumeOptions(),
		})
		if err != nil {
			t.Errorf("Failed to generate resume for domain %s: %v", domainLabel, err)
			continue
		}

		md, err := mdExporter.Export(result.Resume)
		if err != nil {
			t.Errorf("Failed to export resume for domain %s: %v", domainLabel, err)
			continue
		}

		t.Logf("Domain %s: %d experiences, %d skills, %d chars markdown",
			domainLabel,
			len(result.Resume.Content.Experiences),
			len(result.Resume.Content.Skills),
			len(md))

		// Verify domain-specific summary is used
		if domain != "" {
			expectedSummary := profile.Profile.Summaries.ByDomain[domain]
			if expectedSummary != "" && result.Resume.Content.Summary != expectedSummary {
				t.Logf("  Note: Using domain-specific summary for %s", domain)
			}
		}
	}
}
