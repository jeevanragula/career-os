package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jeevanragula/career-os/internal/jobs"
)

type AshbyAdapter struct{ Client *http.Client }

func (AshbyAdapter) Name() string { return "ashby-public-job-postings" }

func (a AshbyAdapter) Discover(ctx context.Context, source jobs.SourceConfig, q jobs.DiscoveryQuery) (jobs.SourceResult, error) {
	client := a.Client
	if client == nil { client = http.DefaultClient }
	u, err := url.Parse(source.BaseURL)
	if err != nil { return jobs.SourceResult{}, err }
	if u.Scheme == "" || u.Host == "" { return jobs.SourceResult{}, fmt.Errorf("ashby base URL is required") }
	params := u.Query()
	params.Set("includeCompensation", "false")
	u.RawQuery = params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil { return jobs.SourceResult{}, err }
	resp, err := client.Do(req)
	if err != nil { return jobs.SourceResult{}, err }
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 { return jobs.SourceResult{}, fmt.Errorf("ashby returned HTTP %d", resp.StatusCode) }
	var payload struct { Jobs []struct {
		ID string `json:"id"`; Title string `json:"title"`; Location string `json:"location"`; JobURL string `json:"jobUrl"`; ApplyURL string `json:"applyUrl"`; Description string `json:"descriptionPlain"`; EmploymentType string `json:"employmentType"`
	} `json:"jobs"` }
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil { return jobs.SourceResult{}, err }
	now := time.Now().UTC()
	out := make([]jobs.JobObservation, 0, len(payload.Jobs))
	for _, p := range payload.Jobs {
		sourceURL := p.JobURL
		if sourceURL == "" { sourceURL = p.ApplyURL }
		out = append(out, jobs.JobObservation{Source: source.Name, SourceJobID: p.ID, SourceURL: sourceURL, ObservedAt: now, RawTitle: p.Title, RawCompany: source.Name, RawLocation: strings.TrimSpace(p.Location), RawDescription: p.Description})
	}
	return jobs.SourceResult{Observations: out}, nil
}
