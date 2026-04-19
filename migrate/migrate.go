// Package migrate provides tools to migrate data from gocvresume to structured-profile format.
package migrate

import (
	"github.com/grokify/structured-profile/schema"
)

// MigrateJohnWangProfile creates a FullProfile from the hardcoded John Wang resume data.
// This demonstrates the migration from gocvresume's Go code data to structured-profile JSON.
func MigrateJohnWangProfile(profileID string) (*schema.FullProfile, error) {
	fp := schema.NewFullProfile("John Wang")
	if profileID != "" {
		fp.Profile.ID = profileID
	}

	// Basic info
	fp.Profile.Email = "johncwang@gmail.com"
	fp.Profile.Links = []schema.Link{
		schema.NewLink("linkedin", "https://www.linkedin.com/in/johncwang/"),
		schema.NewLink("website", "https://grokify.github.io"),
		schema.NewLink("github", "https://github.com/grokify"),
		schema.NewLink("stackoverflow", "http://stackoverflow.com/users/1908967/grokify"),
	}

	// Summaries by domain
	fp.Profile.Summaries = schema.Summaries{
		Default: "John is a technical, hands-on product management executive with 20+ years of experience working with C-suite leaders, engineering, and sales teams to deliver innovative customer solutions. He has a MBA and Bachelors in Computer Science, and is an ex-CISSP. He specializes in platforms including IAM, APIs, infra, and open source.",
		ByDomain: map[string]string{
			"devx":     "John is an award-winning, technical, hands-on product management executive with 20+ years of experience working with C-suite leaders, engineering, sales, business development, and customers to deliver innovative customer solutions. He has a MBA and Bachelors in Computer Science, and is an ex-CISSP. He uses and specializes in developer platforms including AI (AI agents/assistants, MCP, multi-agent systems), developer experience, container and serverless platforms, APIs, SDKs, CI/CD pipelines, marketplaces/ecosystems, developer advocacy, developer communities, and open source. His developer programs have won 7 industry awards in 4 years including Best Public API, Best Communications API, and Best Developer Dashboard.",
			"iam":      "John is a technical, hands-on product management executive with 20+ years of experience working with C-suite leaders, engineering, and sales teams to deliver innovative customer solutions. He has a MBA and Bachelors in Computer Science, and is an ex-CISSP. He specializes in platforms including IAM, APIs, infra, and open source.",
			"platform": "John is a technical, hands-on product management executive with 20+ years of experience working with C-suite leaders, engineering, marketing, and sales teams to deliver innovative customer solutions. He has a MBA and Bachelors in Computer Science, and is an ex-CISSP. He specializes in platforms including APIs, IAM, analytics, infra, and open source.",
		},
	}

	// Add tenures (work experience)
	fp.Tenures = migrateWorkExperience()

	// Add education
	fp.Education = migrateEducation()

	// Add certifications
	fp.Certifications = migrateCertifications()

	// Add publications (articles, patents, papers)
	fp.Publications = migratePublications()

	// Add credentials (GitHub, StackOverflow, etc.)
	fp.Credentials = migrateCredentials()

	// Add skills
	fp.Skills = migrateSkills()

	return fp, nil
}

func migrateWorkExperience() []schema.Tenure {
	return []schema.Tenure{
		migrateSaviynt(),
		migrateSabbatical(),
		migrateFastly(),
		migrateRingCentral(),
		migrateProofpoint(),
		migrateGrokbase(),
		migrateZLTechnologies(),
		migrateSelfEmployed(),
		migrateArcotSystems(),
	}
}

func migrateSaviynt() schema.Tenure {
	tenure := schema.NewTenure("Saviynt", schema.NewDate(2023, 7))

	// VP Platform
	vpPlatform := schema.NewPosition("VP Platform", schema.NewDate(2024, 10))
	vpPlatform.Description = "Lead platform product team"
	vpPlatform.DomainConfigs = []schema.PositionDomainConfig{
		{
			Domain: "devx",
			Skills: []string{"Product Strategy", "API Design", "OAuth 2.0", "OIDC", "SCIM"},
			AchievementOrder: []string{"api-mvp", "token-exchange"},
		},
		{
			Domain: "iam",
			Skills: []string{"Product Strategy", "IAM", "IGA", "PAM", "OAuth 2.0", "OIDC", "SCIM"},
			AchievementOrder: []string{"api-mvp", "token-exchange", "marketplace"},
		},
	}

	vpPlatformAchievements := []schema.Achievement{
		*createSTARAchievement("api-mvp",
			[]string{"api", "devx", "platform"},
			[]string{"API Design", "Product Management", "Go"},
			"Own platform APIs",
			"to drive adoption",
			"by designing and launching RESTful APIs with OpenAPI specs",
			"resulting in improved developer experience and adoption"),
		*createSTARAchievement("token-exchange",
			[]string{"api", "iam", "oauth"},
			[]string{"OAuth 2.0", "OIDC", "Security"},
			"Launch OAuth 2.0 Token Exchange with OIDC support",
			"to enable secure token federation",
			"by implementing RFC 8693 with OIDC integration",
			"enabling customers to integrate with their IdPs"),
		*createSTARAchievement("marketplace",
			[]string{"marketplace", "integrations"},
			[]string{"Product Management", "Integrations"},
			"Build integration marketplace strategy",
			"to expand ecosystem",
			"by defining marketplace architecture and partner program",
			"to drive platform adoption"),
	}
	for _, a := range vpPlatformAchievements {
		vpPlatform.AddAchievement(a)
	}

	// Sr. Director Platform
	srDir := schema.NewPosition("Sr. Director Platform", schema.NewDate(2023, 7))
	srDir.EndDate = ptrDate(schema.NewDate(2024, 9))
	srDir.Description = "Platform product management leadership"

	srDirAchievements := []schema.Achievement{
		*createSTARAchievement("api-strategy",
			[]string{"api", "strategy"},
			[]string{"Product Strategy", "API Design"},
			"Define API strategy",
			"to modernize platform",
			"by creating comprehensive API roadmap",
			"enabling developer self-service"),
	}
	for _, a := range srDirAchievements {
		srDir.AddAchievement(a)
	}

	tenure.AddPosition(*vpPlatform)
	tenure.AddPosition(*srDir)

	return *tenure
}

func migrateSabbatical() schema.Tenure {
	tenure := schema.NewTenure("Sabbatical", schema.NewDate(2023, 1))
	tenure.EndDate = ptrDate(schema.NewDate(2023, 7))

	pos := schema.NewPosition("R&D", schema.NewDate(2023, 1))
	pos.EndDate = ptrDate(schema.NewDate(2023, 7))
	pos.DomainConfigs = []schema.PositionDomainConfig{
		{
			Domain: "devx",
			Skills: []string{"Generative AI", "ChatGPT", "GitHub Copilot", "Cloudflare Pages", "CI/CD", "DevOps"},
			AchievementOrder: []string{"ai", "sso", "gitops-site"},
		},
	}

	achievements := []schema.Achievement{
		*createSTARAchievement("ai",
			[]string{"ai", "generative-ai"},
			[]string{"Generative AI", "ChatGPT", "GitHub Copilot"},
			"Use and evaluate AI",
			"to improve productivity",
			"tools including OpenAI ChatGPT, GitHub Copilot, and Midjourney",
			"to improve productivity"),
		*createSTARAchievement("sso",
			[]string{"iam", "sso"},
			[]string{"1Password", "LastPass", "SSO"},
			"Evaluate SSO solutions including 1Password and LastPass",
			"for enterprise and personal use",
			"evaluating features and integrations",
			"selecting 1Password due to DevOps CI/CD integrations with verified GitHub commits"),
		*createSTARAchievement("gitops-site",
			[]string{"devops", "cicd"},
			[]string{"Cloudflare Pages", "GitHub", "CI/CD"},
			"Launch content website",
			"to share knowledge",
			"using CI/CD with Cloudflare Pages and GitHub",
			"to positive response and automatic deployments"),
	}
	for _, a := range achievements {
		pos.AddAchievement(a)
	}

	tenure.AddPosition(*pos)
	return *tenure
}

func migrateFastly() schema.Tenure {
	tenure := schema.NewTenure("Fastly", schema.NewDate(2022, 6))
	tenure.EndDate = ptrDate(schema.NewDate(2022, 12))

	pos := schema.NewPosition("VP, Product Management, Platform", schema.NewDate(2022, 6))
	pos.EndDate = ptrDate(schema.NewDate(2022, 12))
	pos.Description = "Led team of 5 product managers covering platform, APIs, access, billing, and marketplace"
	pos.DomainConfigs = []schema.PositionDomainConfig{
		{
			Domain: "devx",
			Skills: []string{"Product Strategy", "Acquisition Integration", "OpenAPI/Swagger", "RAML", "Social Login", "RBAC", "Terraform"},
			AchievementOrder: []string{"platform-overview-tech", "unblock", "devx", "iamdevops"},
		},
		{
			Domain: "iam",
			Skills: []string{"Product Strategy", "IAM", "SSO", "Service Accounts", "Social Login", "Passwordless Auth", "RBAC", "Okta Auth0", "AWS Cognito", "Terraform", "Argo CD"},
			AchievementOrder: []string{"platform-overview-iam", "unblock-iam", "iamdevops"},
		},
	}

	achievements := []schema.Achievement{
		*createSTARAchievement("platform-overview-tech",
			[]string{"platform", "strategy"},
			[]string{"Product Strategy", "Team Leadership"},
			"Led product effort to become a multi-product company with PX/UX, user console, IAM, APIs, billing, and marketplace",
			"requiring underlying, scalable, and secure platform",
			"with a team of 7 product managers",
			"resulting in shared vision and strategy at the exec, leadership and team levels"),
		*createSTARAchievement("platform-overview-iam",
			[]string{"platform", "iam"},
			[]string{"Product Strategy", "IAM", "Auth0"},
			"Led product effort to become a multi-product company with IAM, Billing, Marketplace, and Integrations",
			"requiring underlying, scalable, and secure IAM platform",
			"evaluating and selecting identity solutions",
			"resulting in shared vision and strategy. Evaluated and selected Okta Auth0 as the IAM platform vs. AWS Cognito"),
		*createSTARAchievement("unblock",
			[]string{"leadership", "technical-debt"},
			[]string{"Problem Solving", "Cross-functional Leadership"},
			"Unblock important, unschedulable projects including multi-product web console and scaling challenges",
			"requiring technical understanding",
			"by understanding capabilities/tradeoffs and partnering with UX and engineering",
			"to eliminate technical debt and improve team confidence"),
		*createSTARAchievement("unblock-iam",
			[]string{"leadership", "iam"},
			[]string{"Problem Solving", "IAM"},
			"Unblock important, unschedulable projects including multi-product web console and scaling challenges",
			"requiring IAM platform decisions",
			"by understanding capabilities/tradeoffs and partnering with UX and engineering",
			"to eliminate technical debt and improve team confidence, prioritizing IAM as a key strategic, multi-product initiative"),
		*createSTARAchievement("devx",
			[]string{"api", "devx"},
			[]string{"API Design", "OpenAPI", "RAML"},
			"Drive consistent, best-practices API design",
			"to improve developer experience",
			"by evaluating and discussing API design patterns including REST vs. JSON-API and OpenAPI vs RAML",
			"to drive to easier adoption, standardization, and SDK maintenance"),
		*createSTARAchievement("iamdevops",
			[]string{"iam", "devops"},
			[]string{"IAM", "Terraform", "Argo CD"},
			"Own and deliver key IAM and DevOps features",
			"to support customer requirements",
			"including service accounts, Terraform updates and plans for Argo CD integration",
			"to drive DevSecOps efficiencies and hard customer requirements"),
	}
	for _, a := range achievements {
		pos.AddAchievement(a)
	}

	tenure.AddPosition(*pos)
	return *tenure
}

func migrateRingCentral() schema.Tenure {
	tenure := schema.NewTenure("RingCentral", schema.NewDate(2015, 2))
	tenure.EndDate = ptrDate(schema.NewDate(2022, 6))

	// AVP Position
	avp := schema.NewPosition("AVP (Associate VP), Product Management, Platform", schema.NewDate(2020, 10))
	avp.EndDate = ptrDate(schema.NewDate(2022, 6))
	avp.DomainConfigs = []schema.PositionDomainConfig{
		{
			Domain: "devx",
			Skills: []string{"Product Strategy", "Websockets", "Android SDK", "iOS SDK", "Acquisition Integration"},
			AchievementOrder: []string{"platform-overview", "newbusiness", "videotranscription"},
		},
		{
			Domain: "iam",
			Skills: []string{"Product Strategy", "OIDC", "PKCE", "Microservices", "Service Mesh", "V2MOM", "Acquisition Integration"},
			AchievementOrder: []string{"platform-overview", "iam", "newbusiness-iam", "internalplatform"},
		},
	}

	avpAchievements := []schema.Achievement{
		*createSTARAchievement("platform-overview",
			[]string{"platform", "leadership"},
			[]string{"Product Strategy", "Team Leadership", "V2MOM"},
			"Own overall platform strategy including API/SDK products, developer ecosystem/marketplace, developer experience, developer advocacy, and Labs integrations, with a team of 10 product managers, 5 developer advocates/SDK engineers, and 2 tech writers",
			"to manage geographically distributed, cross-functional teams using North Star metrics, MBOs and V2MOM",
			"building and scaling platform organization",
			"resulting in over $500 million attached ARR, 70,000+ developers, 200+ marketplace apps"),
		*createSTARAchievement("iam",
			[]string{"iam", "security"},
			[]string{"OIDC", "PKCE", "JWT"},
			"Advance IAM and security leadership",
			"to improve security posture",
			"by releasing OIDC, PKCE, and JWT, along with deprecating Implicit Grant",
			"resulting in easier integration with no-code OIDC and transparent SDK enhancements"),
		*createSTARAchievement("newbusiness",
			[]string{"business", "product"},
			[]string{"SMS API", "AI APIs", "Video SDK"},
			"Deliver new lines of business",
			"to expand revenue",
			"including high-volume SMS API and app, AI APIs, real-time websocket APIs, and video SDKs (Android & iOS)",
			"with finance-approved $50m revenue target"),
		*createSTARAchievement("newbusiness-iam",
			[]string{"business", "iam"},
			[]string{"SMS API", "AI APIs", "2FA"},
			"Deliver new lines of business",
			"to expand revenue",
			"including high-volume SMS API and app, AI APIs and video SDK, and analysis of 2FA SMS & authenticator app product line",
			"with finance-approved $50m revenue target"),
		*createSTARAchievement("internalplatform",
			[]string{"platform", "devops"},
			[]string{"API Gateway", "Kubernetes", "Docker"},
			"Drive engineering efficiency across 1000 engineers and 100s of microservices",
			"to improve deployment velocity",
			"by driving API gateway, service mesh, and Kubernetes(K8s)/Docker",
			"reducing deployment times from once a quarter to once a week"),
		*createSTARAchievement("videotranscription",
			[]string{"ai", "video"},
			[]string{"Video API", "Transcription", "ML Evaluation"},
			"Launch video transcription API and evaluate services from Google, Amazon, IBM, and Microsoft",
			"to select best provider",
			"by creating test recording set with 'golden' transcriptions to compare against using metrics such as Diarization Error Rate (DER) and Jaccard Error Rate (JER)",
			"to select best combination of quality and price"),
	}
	for _, a := range avpAchievements {
		avp.AddAchievement(a)
	}

	// Sr. Director Position
	srDir := schema.NewPosition("Sr. Director, Product Management, Platform", schema.NewDate(2017, 4))
	srDir.EndDate = ptrDate(schema.NewDate(2020, 9))
	srDir.DomainConfigs = []schema.PositionDomainConfig{
		{
			Domain: "devx",
			Skills: []string{"Acquisition Integration", "GCP Dialogflow", "Amazon Alexa", "GDPR/CCPA", "OpenAPI Generator", "AI Agent Assist", "SEO", "Mixpanel"},
			AchievementOrder: []string{"product-overview", "acquistion", "chatbot", "ipaas", "marketing"},
		},
	}

	srDirAchievements := []schema.Achievement{
		*createSTARAchievement("product-overview",
			[]string{"platform", "awards"},
			[]string{"JavaScript", "React", "WebRTC", "Amazon Alexa"},
			"Launched multiple products including JS/React-based Embeddable widget UI SDK and React component library, Amazon Alexa skills with voice-based biometric authentication",
			"to drive adoption",
			"building and releasing products",
			"resulting in platform adoption of over 70% of enterprise customers. 5 industry awards, 2 personal awards, and a call out from Goldman Sachs investment analyst"),
		*createSTARAchievement("iam",
			[]string{"iam", "privacy"},
			[]string{"RBAC", "PBAC", "GDPR", "CCPA", "SCIM 2.0", "HIPAA"},
			"Launched IAM and privacy features including RBAC/PBAC, Authorized app management, GDPR/CCPA, SCIM 2.0 API, HIPAA compliance, and Amazon Alexa voice-based biometric authentication",
			"to meet compliance requirements",
			"designing and launching security features",
			"resulting in platform adoption of over 70% of enterprise customers"),
		*createSTARAchievement("chatbot",
			[]string{"ai", "chatbot"},
			[]string{"Amazon Alexa", "Chatbot API"},
			"Enable real-time chat applications via Amazon Alexa and interactive chatbot engagement in Slack-like chat communications",
			"to expand use cases",
			"by launching Alexa integration and chatbot API",
			"resulting in the launch of two Alexa skills and over 20 chatbots"),
		*createSTARAchievement("acquistion",
			[]string{"m&a", "integration"},
			[]string{"Acquisition Integration", "IAM", "API"},
			"Successfully integrated multiple acquisitions, consumer chat, voice and digital contact center SaaS solutions",
			"to create unified product",
			"releasing integrated solutions for IAM, API, marketplace, and documentation",
			"resulting in a seamless multi-product experience"),
		*createSTARAchievement("ipaas",
			[]string{"integrations", "ipaas"},
			[]string{"Zapier", "Workato", "Tray.io", "Mulesoft", "Dell Boomi"},
			"Drove integration capabilities",
			"to expand ecosystem",
			"by partnering with iPaaS workflow solutions like Zapier, Workato, Tray.io, Mulesoft, and Dell Boomi",
			"resulting in integrations with over 1000 apps listed on the marketplace"),
		*createSTARAchievement("marketing",
			[]string{"marketing", "seo"},
			[]string{"SEO", "Content Marketing"},
			"Increase customer awareness of products",
			"to drive adoption",
			"by partnering with marketing on SEO, release notes, and blog posts",
			"resulting in customer adoption of key products and first page SERP placement for key search terms"),
	}
	for _, a := range srDirAchievements {
		srDir.AddAchievement(a)
	}

	// Director Position
	dir := schema.NewPosition("Director, Product Management, Platform", schema.NewDate(2015, 2))
	dir.EndDate = ptrDate(schema.NewDate(2017, 3))
	dir.DomainConfigs = []schema.PositionDomainConfig{
		{
			Domain: "devx",
			Skills: []string{"AWS Lambda", "AWS API Gateway", "Heroku", "Sumo Logic", "OAuth 2.0", "OpenAPI/Swagger", "Webhooks", "SQL", "Tableau"},
			AchievementOrder: []string{"integrations", "productlaunches", "community"},
		},
	}

	dirAchievements := []schema.Achievement{
		*createSTARAchievement("integrations",
			[]string{"integrations", "sales"},
			[]string{"Integrations", "Sales Enablement"},
			"Handle sales objections regarding lack of integrations",
			"to unblock sales",
			"by proposing, receiving funding releasing over 20 integrations in 3 months",
			"to eliminate sales blocker and attain new accounts"),
		*createSTARAchievement("productlaunches",
			[]string{"api", "platform"},
			[]string{"OAuth 2.0", "WebRTC", "SDKs"},
			"Launch API program including OAuth 2.0 Authorization Code Grant & Implicit Grant, WebRTC, voice streaming API, open source SDKs, etc.",
			"to build developer platform",
			"designing and launching APIs",
			"growing program to 1 million API calls per day at 50% QoQ, and contributing to moving from Gartner Visionary to Leader"),
		*createSTARAchievement("community",
			[]string{"devrel", "community"},
			[]string{"Developer Relations", "Open Source", "Stack Overflow"},
			"Drive developer experience and grow community engagement",
			"to build developer community",
			"by launching online community, maintaining open source SDKs, gaining Stack Overflow tag, publishing tutorials, publishing Heroku applications, running Medium blog and launching Slack-like real-time support channel",
			"resulting in 10+ official SDKs, 250 blog articles, 230 page developer guide, and 246 Stack Overflow questions"),
	}
	for _, a := range dirAchievements {
		dir.AddAchievement(a)
	}

	tenure.AddPosition(*avp)
	tenure.AddPosition(*srDir)
	tenure.AddPosition(*dir)

	return *tenure
}

func migrateProofpoint() schema.Tenure {
	tenure := schema.NewTenure("Proofpoint", schema.NewDate(2013, 1))
	tenure.EndDate = ptrDate(schema.NewDate(2015, 2))

	gpm := schema.NewPosition("Group Product Manager, Platform", schema.NewDate(2013, 8))
	gpm.EndDate = ptrDate(schema.NewDate(2015, 2))
	gpm.DomainConfigs = []schema.PositionDomainConfig{
		{
			Domain: "devx",
			Skills: []string{"search", "continuity", "backup", "IAM", "GRC Finance", "APIs", "REST", "RAML"},
			AchievementOrder: []string{"platform"},
		},
	}

	gpmAchievements := []schema.Achievement{
		*createSTARAchievement("platform",
			[]string{"platform", "saas"},
			[]string{"IAM", "Elasticsearch", "SAML"},
			"Owned multi-tenant SaaS solutions for IAM (Shibboleth), search (Elasticsearch), and backup",
			"to move from hosted/on-premises solutions",
			"building and deploying cloud solutions",
			"driving cloud migration and earning Proofpoint's \"Critical Impact Award\""),
	}
	for _, a := range gpmAchievements {
		gpm.AddAchievement(a)
	}

	pmm := schema.NewPosition("Sr. Product Marketing Manager", schema.NewDate(2013, 1))
	pmm.EndDate = ptrDate(schema.NewDate(2013, 8))
	pmm.DomainConfigs = []schema.PositionDomainConfig{
		{
			Domain: "default",
			Skills: []string{"competitive analysis", "product positioning", "product messaging", "product launch", "marketing collateral"},
			AchievementOrder: []string{"pmm"},
		},
	}

	pmmAchievements := []schema.Achievement{
		*createSTARAchievement("pmm",
			[]string{"marketing", "positioning"},
			[]string{"Product Marketing", "Competitive Analysis"},
			"Responsible for product positioning, messaging, and launch for Proofpoint Archive",
			"to drive sales",
			"including creation of white papers, data sheets, and competitive battle cards",
			"resulting in strong sales and Gartner leadership position"),
	}
	for _, a := range pmmAchievements {
		pmm.AddAchievement(a)
	}

	tenure.AddPosition(*gpm)
	tenure.AddPosition(*pmm)

	return *tenure
}

func migrateGrokbase() schema.Tenure {
	tenure := schema.NewTenure("Grokbase", schema.NewDate(2011, 4))
	tenure.EndDate = ptrDate(schema.NewDate(2012, 12))

	pos := schema.NewPosition("Founder, Product Manager, Full Stack Developer", schema.NewDate(2011, 4))
	pos.EndDate = ptrDate(schema.NewDate(2012, 12))
	pos.DomainConfigs = []schema.PositionDomainConfig{
		{
			Domain: "default",
			Skills: []string{"AI/Clustering", "Product-Led Growth", "operations", "sysadmin", "LAMP", "MongoDB", "RHEL"},
			AchievementOrder: []string{"grokbaseoverview"},
		},
	}

	achievements := []schema.Achievement{
		*createSTARAchievement("grokbaseoverview",
			[]string{"startup", "seo"},
			[]string{"Full Stack", "SEO", "Product-Led Growth"},
			"Improve the usability of open source mailing lists",
			"to help developers",
			"by designing, launching, and promoting Stack Overflow-like website",
			"attaining a 12,500 Alexa rank and recognition by developers, recruiters, and venture capital"),
	}
	for _, a := range achievements {
		pos.AddAchievement(a)
	}

	tenure.AddPosition(*pos)
	return *tenure
}

func migrateZLTechnologies() schema.Tenure {
	tenure := schema.NewTenure("ZL Technologies", schema.NewDate(2007, 8))
	tenure.EndDate = ptrDate(schema.NewDate(2011, 3))

	pos := schema.NewPosition("Lead Product Manager", schema.NewDate(2007, 8))
	pos.EndDate = ptrDate(schema.NewDate(2011, 3))
	pos.Description = "Lead team of 4 product managers and 2 tech writers"
	pos.DomainConfigs = []schema.PositionDomainConfig{
		{
			Domain: "default",
			Skills: []string{"GRC Finance", "FRCP", "AI Concept Search", "AWS EC2", "AWS S3"},
			AchievementOrder: []string{"pm"},
		},
	}

	achievements := []schema.Achievement{
		*createSTARAchievement("pm",
			[]string{"grc", "ediscovery"},
			[]string{"Information Governance", "eDiscovery", "Compliance"},
			"Owned flagship ZL Unified Archive information governance solution covering FINRA/SEC 17a-4 regulatory compliance, eDiscovery, and records management",
			"to serve enterprise customers",
			"building and managing product",
			"resulting multiple Fortune 50 accounts and product awards"),
	}
	for _, a := range achievements {
		pos.AddAchievement(a)
	}

	tenure.AddPosition(*pos)
	return *tenure
}

func migrateSelfEmployed() schema.Tenure {
	tenure := schema.NewTenure("Self-employed", schema.NewDate(2004, 7))
	tenure.EndDate = ptrDate(schema.NewDate(2007, 7))

	pos := schema.NewPosition("Full Stack Developer", schema.NewDate(2004, 7))
	pos.EndDate = ptrDate(schema.NewDate(2007, 7))
	pos.DomainConfigs = []schema.PositionDomainConfig{
		{
			Domain: "default",
			Skills: []string{"LAMP", "Oracle RDBMS", "RHEL", "SEO", "AI/Clustering", "Unsupervised Learning", "PLG", "Affiliate Marketing", "Full Stack"},
			AchievementOrder: []string{"overview"},
		},
	}

	achievements := []schema.Achievement{
		*createSTARAchievement("overview",
			[]string{"entrepreneurship", "ecommerce"},
			[]string{"Full Stack", "LAMP", "SEO"},
			"Build e-commerce, affiliate-marketing, and content sites",
			"to generate revenue",
			"using LAMP stack and unsupervised machine learning",
			"resulting in high search engine rankings and revenue"),
	}
	for _, a := range achievements {
		pos.AddAchievement(a)
	}

	tenure.AddPosition(*pos)
	return *tenure
}

func migrateArcotSystems() schema.Tenure {
	tenure := schema.NewTenure("Arcot Systems (Acquired for $200M)", schema.NewDate(1999, 1))
	tenure.EndDate = ptrDate(schema.NewDate(2004, 6))

	// Director
	dir := schema.NewPosition("Director, Product Marketing", schema.NewDate(2003, 7))
	dir.EndDate = ptrDate(schema.NewDate(2004, 6))
	dir.DomainConfigs = []schema.PositionDomainConfig{
		{
			Domain: "default",
			Skills: []string{"Corporate Ventures", "Anti-Phishing"},
			AchievementOrder: []string{"overview"},
		},
	}

	dirAchievements := []schema.Achievement{
		*createSTARAchievement("overview",
			[]string{"product-marketing", "anti-phishing"},
			[]string{"Product Marketing", "Anti-Phishing"},
			"Managed new product lines",
			"to expand business",
			"such as patent-pending anti-phishing web browser extension by working with researchers and the Anti-Phishing Working Group (APWG)",
			"resulting in potential new lines of business"),
	}
	for _, a := range dirAchievements {
		dir.AddAchievement(a)
	}

	// Product Manager
	pm := schema.NewPosition("Product Manager", schema.NewDate(2001, 7))
	pm.EndDate = ptrDate(schema.NewDate(2003, 6))
	pm.DomainConfigs = []schema.PositionDomainConfig{
		{
			Domain: "default",
			Skills: []string{"MFA", "2FA", "Cryptography", "IAM", "Authorization", "LDAP"},
			AchievementOrder: []string{"pm", "patent"},
		},
	}

	pmAchievements := []schema.Achievement{
		*createSTARAchievement("pm",
			[]string{"iam", "mfa"},
			[]string{"MFA", "2FA", "IAM", "SSO"},
			"Owned product management for software MFA/2FA authenticator product line, and PBAC/SSO access control solution",
			"to drive enterprise sales",
			"launching multiple product releases including strong-authentication client, LDAP integration, and IAM web-SSO access control solution",
			"resulting in selection by Visa and MasterCard for Verified by Visa and SecureCode web-scale authentication"),
		*createSTARAchievement("patent",
			[]string{"invention", "iam"},
			[]string{"X.509", "mTLS", "Cryptography"},
			"Proposed and launched new patented feature to overcome key, consistent sales objection from customers planted by competitors",
			"to eliminate sales blocker",
			"by providing industry standard X.509 mTLS interface to Arcot's proprietary \"software smart card\" solution",
			"resulting in over $1m in pre-release direct sales and eliminating the sales objection"),
	}
	for _, a := range pmAchievements {
		pm.AddAchievement(a)
	}

	// Sales Engineer
	se := schema.NewPosition("Sales Engineer", schema.NewDate(1999, 1))
	se.EndDate = ptrDate(schema.NewDate(2001, 6))
	se.DomainConfigs = []schema.PositionDomainConfig{
		{
			Domain: "default",
			Skills: []string{"MFA", "2FA", "MS IIS", "MS SQL Server"},
			AchievementOrder: []string{"se"},
		},
	}

	seAchievements := []schema.Achievement{
		*createSTARAchievement("se",
			[]string{"sales", "presales"},
			[]string{"Sales Engineering", "MFA", "Presentations"},
			"Customer-facing sales-engineering ownership for strong, 2-factor authentication solution, covering Western U.S., APAC, and later EMEA",
			"to close deals",
			"including sales presentations and custom sales demos",
			"exceeding sales quota every quarter"),
	}
	for _, a := range seAchievements {
		se.AddAchievement(a)
	}

	tenure.AddPosition(*dir)
	tenure.AddPosition(*pm)
	tenure.AddPosition(*se)

	return *tenure
}

func migrateEducation() []schema.Education {
	return []schema.Education{
		{
			BaseEntity:  schema.NewBaseEntity(),
			Institution: "Stanford Graduate School of Business",
			Degree:      "SEP (Stanford Executive Program)",
			Display:     true,
		},
		{
			BaseEntity:  schema.NewBaseEntity(),
			Institution: "Rensselaer Polytechnic Institute",
			Degree:      "MBA",
			Display:     true,
		},
		{
			BaseEntity:  schema.NewBaseEntity(),
			Institution: "University of Pennsylvania",
			Degree:      "BAS Computer Science",
			Honors:      "University Scholar",
			Display:     true,
		},
		{
			BaseEntity:  schema.NewBaseEntity(),
			Institution: "University of Pennsylvania Wharton School of Business",
			Degree:      "BS Economics, Finance",
			Honors:      "University Scholar",
			Display:     true,
		},
	}
}

func migrateCertifications() []schema.Certification {
	return []schema.Certification{
		{
			BaseEntity: schema.NewBaseEntity(),
			Name:       "API Security Architect",
			Issuer:     "API Academy",
			Display:    true,
		},
		{
			BaseEntity: schema.NewBaseEntity(),
			Name:       "REST API",
			Issuer:     "HackerRank",
			Display:    true,
		},
		{
			BaseEntity: schema.NewBaseEntity(),
			Name:       "Fundamentals of RESTful API Design",
			Issuer:     "Apigee",
			Display:    true,
		},
		{
			BaseEntity: schema.NewBaseEntity(),
			Name:       "Structuring Machine Learning Projects",
			Issuer:     "DeepLearning.AI",
			Display:    true,
		},
		{
			BaseEntity: schema.NewBaseEntity(),
			Name:       "Machine Learning",
			Issuer:     "Stanford University",
			Display:    true,
		},
		{
			BaseEntity: schema.NewBaseEntity(),
			Name:       "Machine Learning Foundations for Product Managers",
			Issuer:     "Duke University",
			Display:    true,
		},
		{
			BaseEntity: schema.NewBaseEntity(),
			Name:       "Databricks Accredited Generative AI Fundamentals",
			Issuer:     "Databricks",
			Display:    true,
		},
		{
			BaseEntity: schema.NewBaseEntity(),
			Name:       "CSPO (Certified Scrum Product Owner)",
			Issuer:     "Scrum Alliance",
			Display:    true,
		},
		{
			BaseEntity: schema.NewBaseEntity(),
			Name:       "CSM (Certified Scrum Master)",
			Issuer:     "Scrum Alliance",
			Display:    true,
		},
		{
			BaseEntity: schema.NewBaseEntity(),
			Name:       "PSPO I, PSM III, PSD I, SPS, PAL, PAL-EBM",
			Issuer:     "Scrum.org",
			Display:    true,
		},
		{
			BaseEntity: schema.NewBaseEntity(),
			Name:       "SQL (Advanced)",
			Issuer:     "HackerRank",
			Display:    true,
		},
		{
			BaseEntity: schema.NewBaseEntity(),
			Name:       "Go (Basic)",
			Issuer:     "HackerRank",
			Display:    true,
		},
		{
			BaseEntity: schema.NewBaseEntity(),
			Name:       "Python (Basic)",
			Issuer:     "HackerRank",
			Display:    true,
		},
	}
}

func migratePublications() []schema.Publication {
	return []schema.Publication{
		// Patents
		{
			BaseEntity:  schema.NewBaseEntity(),
			Type:        "patent",
			Title:       "Methods and systems for associating a team with a meeting",
			Description: "Automatically associate team channels with meetings",
			URL:         "https://patents.google.com/patent/US20220207488",
			Display:     true,
		},
		{
			BaseEntity:  schema.NewBaseEntity(),
			Type:        "patent",
			Title:       "Collaborative communications environment and automatic account creation thereof",
			Description: "Automatic account creation and team assignment to drive viral growth",
			URL:         "https://patents.google.com/patent/US10805101",
			Display:     true,
		},
		{
			BaseEntity:  schema.NewBaseEntity(),
			Type:        "patent",
			Title:       "Method and apparatus for cryptographic key storage wherein key servers are authenticated by possession and secure distribution of stored keys",
			Description: "Strong software-based multi-factor authentication method to use standard X.509 mTLS to authenticate clients and servers",
			URL:         "https://patents.google.com/patent/US7711122",
			Display:     true,
		},
		// Papers
		{
			BaseEntity:  schema.NewBaseEntity(),
			Type:        "paper",
			Title:       "Comparing Exclusionary and Investigative Approaches for Electronic Discovery using the TREC Enron Corpus",
			Description: "Evaluating the use of machine learning and concept search to improve eDiscovery",
			URL:         "http://trec.nist.gov/pubs/trec18/papers/zlti.legal.pdf",
			Display:     true,
		},
		// Articles
		{
			BaseEntity:  schema.NewBaseEntity(),
			Type:        "article",
			Title:       "Introducing OAuth 2.0 Token Exchange and OpenID Connect (OIDC) Support",
			URL:         "https://developers.saviynt.com/blog/introducing-oauth2-token-exchange-with-oidc",
			Display:     true,
		},
		{
			BaseEntity:  schema.NewBaseEntity(),
			Type:        "article",
			Title:       "SCIM 2.0 API is Now Available on RingCentral",
			URL:         "https://medium.com/ringcentral-developers/scim-2-0-api-is-now-available-on-ringcentral-47cbb9a4c7fb",
			Display:     true,
		},
		{
			BaseEntity:  schema.NewBaseEntity(),
			Type:        "article",
			Title:       "Answer to 'What is PKCE actually protecting?'",
			URL:         "https://security.stackexchange.com/a/186997/105337",
			Display:     true,
		},
	}
}

func migrateCredentials() []schema.VerifiableCredential {
	creds := []schema.VerifiableCredential{
		{
			BaseEntity: schema.NewBaseEntity(),
			Type:       "github",
			Username:   "grokify",
			ProfileURL: "https://github.com/grokify",
			Data: schema.CredentialData{
				Repositories: 300,
			},
		},
		{
			BaseEntity: schema.NewBaseEntity(),
			Type:       "stackoverflow",
			Username:   "grokify",
			ProfileURL: "http://stackoverflow.com/users/1908967/grokify",
			Data: schema.CredentialData{
				Reputation:   16000,
				GoldBadges:   5,
				SilverBadges: 30,
				BronzeBadges: 50,
			},
		},
		{
			BaseEntity: schema.NewBaseEntity(),
			Type:       "linkedin",
			Username:   "johncwang",
			ProfileURL: "https://www.linkedin.com/in/johncwang/",
		},
		{
			BaseEntity: schema.NewBaseEntity(),
			Type:       "medium",
			Username:   "grokify",
			ProfileURL: "https://medium.com/@grokify",
			Data: schema.CredentialData{
				Articles: 40,
			},
		},
	}
	// Mark all as verified
	for i := range creds {
		creds[i].Verify()
	}
	return creds
}

func migrateSkills() []schema.Skill {
	skills := []schema.Skill{
		// Technical - Languages
		{BaseEntity: schema.NewBaseEntity(), Name: "Go", Category: "technical", Level: "expert"},
		{BaseEntity: schema.NewBaseEntity(), Name: "Python", Category: "technical", Level: "intermediate"},
		{BaseEntity: schema.NewBaseEntity(), Name: "JavaScript", Category: "technical", Level: "intermediate"},
		{BaseEntity: schema.NewBaseEntity(), Name: "SQL", Category: "technical", Level: "advanced"},

		// Technical - APIs
		{BaseEntity: schema.NewBaseEntity(), Name: "REST APIs", Category: "technical", Level: "expert"},
		{BaseEntity: schema.NewBaseEntity(), Name: "OpenAPI/Swagger", Category: "technical", Level: "expert"},
		{BaseEntity: schema.NewBaseEntity(), Name: "OAuth 2.0", Category: "technical", Level: "expert"},
		{BaseEntity: schema.NewBaseEntity(), Name: "OIDC", Category: "technical", Level: "expert"},
		{BaseEntity: schema.NewBaseEntity(), Name: "SCIM", Category: "technical", Level: "advanced"},
		{BaseEntity: schema.NewBaseEntity(), Name: "SAML", Category: "technical", Level: "advanced"},

		// Technical - IAM
		{BaseEntity: schema.NewBaseEntity(), Name: "IAM", Category: "technical", Level: "expert"},
		{BaseEntity: schema.NewBaseEntity(), Name: "SSO", Category: "technical", Level: "expert"},
		{BaseEntity: schema.NewBaseEntity(), Name: "MFA/2FA", Category: "technical", Level: "expert"},
		{BaseEntity: schema.NewBaseEntity(), Name: "RBAC", Category: "technical", Level: "expert"},

		// Technical - Cloud & DevOps
		{BaseEntity: schema.NewBaseEntity(), Name: "AWS", Category: "technical", Level: "advanced"},
		{BaseEntity: schema.NewBaseEntity(), Name: "Kubernetes", Category: "technical", Level: "intermediate"},
		{BaseEntity: schema.NewBaseEntity(), Name: "Docker", Category: "technical", Level: "intermediate"},
		{BaseEntity: schema.NewBaseEntity(), Name: "Terraform", Category: "technical", Level: "intermediate"},
		{BaseEntity: schema.NewBaseEntity(), Name: "CI/CD", Category: "technical", Level: "advanced"},

		// Technical - AI
		{BaseEntity: schema.NewBaseEntity(), Name: "Generative AI", Category: "technical", Level: "advanced"},
		{BaseEntity: schema.NewBaseEntity(), Name: "Machine Learning", Category: "technical", Level: "intermediate"},
		{BaseEntity: schema.NewBaseEntity(), Name: "LLM APIs", Category: "technical", Level: "advanced"},

		// Product Management
		{BaseEntity: schema.NewBaseEntity(), Name: "Product Strategy", Category: "product", Level: "expert"},
		{BaseEntity: schema.NewBaseEntity(), Name: "Product Roadmapping", Category: "product", Level: "expert"},
		{BaseEntity: schema.NewBaseEntity(), Name: "Product-Led Growth", Category: "product", Level: "expert"},
		{BaseEntity: schema.NewBaseEntity(), Name: "Developer Experience", Category: "product", Level: "expert"},
		{BaseEntity: schema.NewBaseEntity(), Name: "API Design", Category: "product", Level: "expert"},

		// Leadership
		{BaseEntity: schema.NewBaseEntity(), Name: "Team Leadership", Category: "leadership", Level: "expert"},
		{BaseEntity: schema.NewBaseEntity(), Name: "Cross-functional Leadership", Category: "leadership", Level: "expert"},
		{BaseEntity: schema.NewBaseEntity(), Name: "Executive Communication", Category: "leadership", Level: "expert"},
		{BaseEntity: schema.NewBaseEntity(), Name: "Strategic Planning", Category: "leadership", Level: "expert"},
	}
	return skills
}

// Helper functions

func createSTARAchievement(name string, tags, skills []string, situation, task, action, result string) *schema.Achievement {
	a := schema.NewSTARAchievement(name, situation, task, action, result)
	a.Tags = tags
	a.Skills = skills
	return a
}

func ptrDate(d schema.Date) *schema.Date {
	return &d
}

