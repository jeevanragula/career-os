package jobs

import "time"

// Job is the canonical identity used by downstream analysis.
type Job struct {
	ID             string
	Company        string
	Title          string
	CanonicalURL   string
	Location       string
	RemoteMode     string
	EmploymentType string
	FirstSeenAt    time.Time
	LastSeenAt     time.Time
	Status         string
}

// JobObservation preserves one source's representation of a job.
type JobObservation struct {
	ID             string
	JobID          string
	Source         string
	SourceJobID    string
	SourceURL      string
	ObservedAt     time.Time
	RawTitle       string
	RawCompany     string
	RawLocation    string
	RawDescription string
}

// JobVersion represents a materially different normalized posting.
type JobVersion struct {
	ID             string
	JobID          string
	ContentHash    string
	ObservedAt     time.Time
	Title          string
	Company        string
	Location       string
	RemoteMode     string
	EmploymentType string
	Description    string
}

// NormalizeInput is the source-independent input to deterministic normalization.
type NormalizeInput struct {
	Source      string
	SourceJobID string
	URL         string
	Title       string
	Company     string
	Location    string
	RemoteMode  string
	Employment  string
	Description string
	ObservedAt  time.Time
}
