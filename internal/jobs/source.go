package jobs

import "context"

type SourceConfig struct {
	ID string
	Name string
	Provider string
	BaseURL string
	Enabled bool
	Authorization string
	RateLimitPerMin int
}

type DiscoveryQuery struct {
	ID string
	Keywords []string
	Companies []string
	Locations []string
	RemoteModes []string
	PageSize int
	Cursor string
}

type SourceResult struct {
	Observations []JobObservation
	NextCursor string
}

type SourceAdapter interface {
	Name() string
	Discover(ctx context.Context, source SourceConfig, query DiscoveryQuery) (SourceResult, error)
}
