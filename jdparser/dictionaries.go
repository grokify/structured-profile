package jdparser

import "github.com/grokify/structured-profile/schema"

// DefaultTechSkills returns a list of common technical skills to match.
func DefaultTechSkills() []string {
	return []string{
		// Programming Languages
		"Go", "Golang", "Python", "Java", "JavaScript", "TypeScript",
		"C", "C++", "C#", "Rust", "Ruby", "PHP", "Swift", "Kotlin",
		"Scala", "Perl", "R", "MATLAB", "Julia", "Elixir", "Erlang",
		"Haskell", "Clojure", "F#", "Objective-C", "Dart", "Lua",

		// Frontend
		"React", "Angular", "Vue", "Vue.js", "Svelte", "Next.js", "Nuxt",
		"HTML", "CSS", "SASS", "LESS", "Tailwind", "Bootstrap",
		"Redux", "MobX", "Webpack", "Vite", "Babel", "jQuery",

		// Backend
		"Node.js", "Express", "Django", "Flask", "FastAPI",
		"Spring", "Spring Boot", "Rails", "Ruby on Rails", "Laravel",
		"ASP.NET", ".NET", ".NET Core", "Gin", "Echo", "Fiber",

		// Databases
		"PostgreSQL", "MySQL", "MongoDB", "Redis", "Elasticsearch",
		"DynamoDB", "Cassandra", "CouchDB", "Neo4j", "SQLite",
		"Oracle", "SQL Server", "MariaDB", "InfluxDB", "TimescaleDB",
		"Snowflake", "BigQuery", "Redshift",

		// Cloud & Infrastructure
		"AWS", "Amazon Web Services", "GCP", "Google Cloud", "Azure",
		"Kubernetes", "K8s", "Docker", "Terraform", "CloudFormation",
		"Ansible", "Puppet", "Chef", "Vagrant", "Packer",
		"EC2", "S3", "Lambda", "ECS", "EKS", "Fargate",
		"GKE", "Cloud Run", "Cloud Functions", "App Engine",
		"AKS", "Azure Functions", "Blob Storage",

		// DevOps & CI/CD
		"Jenkins", "GitLab CI", "GitHub Actions", "CircleCI", "Travis CI",
		"ArgoCD", "Flux", "Spinnaker", "Harness", "Tekton",
		"Prometheus", "Grafana", "Datadog", "New Relic", "Splunk",
		"ELK", "Elasticsearch", "Logstash", "Kibana", "Jaeger", "Zipkin",

		// APIs & Protocols
		"REST", "RESTful", "GraphQL", "gRPC", "WebSocket",
		"OpenAPI", "Swagger", "Protobuf", "JSON", "XML",
		"SOAP", "HTTP", "HTTPS", "TCP", "UDP",

		// Security
		"OAuth", "OAuth2", "OIDC", "OpenID Connect", "SAML",
		"JWT", "SSL", "TLS", "PKI", "mTLS",
		"OWASP", "Penetration Testing", "Vulnerability Assessment",
		"IAM", "RBAC", "ABAC", "Zero Trust",
		"Encryption", "Cryptography", "HashiCorp Vault", "Secrets Management",

		// Data & ML
		"Spark", "Hadoop", "Kafka", "Flink", "Airflow",
		"TensorFlow", "PyTorch", "scikit-learn", "Keras",
		"Pandas", "NumPy", "Jupyter", "MLflow", "Kubeflow",
		"Databricks", "dbt", "Looker", "Tableau", "Power BI",

		// Mobile
		"iOS", "Android", "React Native", "Flutter", "Xamarin",
		"Swift", "SwiftUI", "Kotlin", "Jetpack Compose",

		// Testing
		"Unit Testing", "Integration Testing", "E2E Testing",
		"Jest", "Mocha", "Chai", "Cypress", "Selenium", "Playwright",
		"JUnit", "TestNG", "pytest", "RSpec", "Cucumber",

		// Version Control
		"Git", "GitHub", "GitLab", "Bitbucket", "SVN",

		// Architecture
		"Microservices", "Monolith", "Serverless", "Event-Driven",
		"DDD", "Domain-Driven Design", "CQRS", "Event Sourcing",
		"Service Mesh", "Istio", "Linkerd", "Envoy",

		// Message Queues
		"Kafka", "RabbitMQ", "SQS", "SNS", "NATS",
		"ActiveMQ", "ZeroMQ", "Pulsar",

		// Observability
		"Monitoring", "Logging", "Tracing", "Alerting",
		"OpenTelemetry", "APM", "SLOs", "SLIs",
	}
}

// DefaultSoftSkills returns a list of common soft skills to match.
func DefaultSoftSkills() []string {
	return []string{
		"Leadership", "Communication", "Collaboration",
		"Problem Solving", "Critical Thinking", "Analytical",
		"Team Player", "Teamwork", "Cross-functional",
		"Mentoring", "Coaching", "Training",
		"Project Management", "Time Management", "Prioritization",
		"Adaptability", "Flexibility", "Self-motivated",
		"Attention to Detail", "Detail-oriented",
		"Creative", "Innovative", "Strategic Thinking",
		"Customer Focus", "Client-facing", "Stakeholder Management",
		"Written Communication", "Verbal Communication", "Presentation",
		"Conflict Resolution", "Negotiation", "Influence",
		"Decision Making", "Judgment", "Initiative",
		"Empathy", "Emotional Intelligence", "Interpersonal",
		"Ownership", "Accountability", "Reliability",
		"Fast-paced", "Startup Mindset", "Entrepreneurial",
		"Agile", "Scrum", "Kanban",
	}
}

// DefaultSeniorityKeywords returns keywords for each seniority level.
func DefaultSeniorityKeywords() map[schema.SeniorityLevel][]string {
	return map[schema.SeniorityLevel][]string{
		schema.SeniorityEntry: {
			"entry level", "entry-level", "junior", "associate",
			"new grad", "new graduate", "recent graduate",
			"intern", "internship", "apprentice",
			"0-2 years", "1-2 years", "0-1 years",
		},
		schema.SeniorityMid: {
			"mid level", "mid-level", "intermediate",
			"2-5 years", "3-5 years", "2-4 years", "3-4 years",
		},
		schema.SenioritySenior: {
			"senior", "sr.", "sr ", "experienced",
			"5+ years", "5-7 years", "5-8 years", "6+ years", "7+ years",
		},
		schema.SeniorityStaff: {
			"staff", "senior staff",
			"8+ years", "10+ years",
		},
		schema.SeniorityPrincipal: {
			"principal", "distinguished",
			"10+ years", "12+ years", "15+ years",
		},
		schema.SeniorityLead: {
			"lead", "tech lead", "technical lead", "team lead",
			"engineering lead",
		},
		schema.SeniorityManager: {
			"manager", "engineering manager", "software manager",
			"development manager", "people manager",
		},
		schema.SeniorityDirector: {
			"director", "engineering director", "senior director",
			"associate director",
		},
		schema.SeniorityVP: {
			"vp", "vice president", "vp of engineering",
			"vp engineering", "svp",
		},
		schema.SeniorityExecutive: {
			"cto", "chief technology", "chief technical",
			"cio", "chief information", "chief architect",
			"c-level", "c-suite", "executive",
		},
	}
}
