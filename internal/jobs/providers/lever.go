package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jeevanragula/career-os/internal/jobs"
)

type LeverAdapter struct{ Client *http.Client }

func (LeverAdapter) Name() string { return "lever-public-postings" }

func (a LeverAdapter) Discover(ctx context.Context, source jobs.SourceConfig, q jobs.DiscoveryQuery) (jobs.SourceResult, error) {
	client := a.Client
	if client == nil { client = http.DefaultClient }
	base := strings.TrimRight(source.BaseURL, "/")
	if base == "" { return jobs.SourceResult{}, fmt.Errorf("lever base URL is required") }
	u, err := url.Parse(base + "/" + url.PathEscape(source.Name))
	if err != nil { return jobs.SourceResult{}, err }
	params := u.Query()
	params.Set("mode", "json")
	if q.PageSize > 0 { params.Set("limit", strconv.Itoa(q.PageSize)) }
	if q.Cursor != "" { params.Set("skip", q.Cursor) }
	u.RawQuery = params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil { return jobs.SourceResult{}, err }
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil { return jobs.SourceResult{}, err }
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 { return jobs.SourceResult{}, fmt.Errorf("lever returned HTTP %d", resp.StatusCode) }
	var payload []struct {
		ID string `json:"id"`
		Text string `json:"text"`
		Categories struct { Location string `json:"location"`; Commitment string `json:"commitment"` } `json:"categories"`
		Description string `json:"description"`
		DescriptionPlain string `json:"descriptionPlain"`
		HostedURL string `json:"hostedUrl"`
		WorkplaceType string `json:"workplaceType"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil { return jobs.SourceResult{}, err }
	now := time.Now().UTC()
	out := make([]jobs.JobObservation, 0, len(payload))
	for _, p := range payload {
		desc := p.DescriptionPlain
		if desc == "" { desc = p.Description }
		out = append(out, jobs.JobObservation{Source: source.Name, SourceJobID: p.ID, SourceURL: p.HostedURL, ObservedAt: now, RawTitle: p.Text, RawCompany: source.Name, RawLocation: p.Categories.Location, RawDescription: desc})
	}
	return jobs.SourceResult{Observations: out}, nil
}
