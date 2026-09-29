package jobs

// Requirement represents one source-grounded job requirement.
type Requirement struct {
	ID string
	JobVersionID string
	Category string
	Text string
	Priority string
	SourceQuote string
}

// RequirementMatch links a requirement to canonical career evidence.
type RequirementMatch struct {
	RequirementID string
	ClaimID string
	Status string
	Confidence float64
	Rationale string
}

type JobAnalysis struct {
	ID string
	JobVersionID string
	Summary string
	Requirements []Requirement
	Matches []RequirementMatch
	OverallConfidence float64
}
