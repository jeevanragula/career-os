package discovery

import (
 "bytes"
 "context"
 "encoding/json"
 "fmt"
 "net/http"
 "os"
)

type SearchResult struct { Title string; URL string; Description string; Source string }
type SearchProvider interface { Search(context.Context,string,int)([]SearchResult,error) }

type Tavily struct { APIKey string; Client *http.Client }

func NewTavily() *Tavily {
 return &Tavily{APIKey: os.Getenv("TAVILY_API_KEY"), Client: http.DefaultClient}
}

func (t *Tavily) Search(ctx context.Context, q string, count int) ([]SearchResult, error) {
 if t.APIKey == "" { return nil, fmt.Errorf("TAVILY_API_KEY is not configured") }
 if count <= 0 { count = 10 }
 payload := map[string]any{
  "api_key": t.APIKey, "query": q, "search_depth": "basic",
  "max_results": count, "include_answer": false, "include_raw_content": false,
 }
 body, err := json.Marshal(payload); if err != nil { return nil, err }
 req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.tavily.com/search", bytes.NewReader(body))
 if err != nil { return nil, err }
 req.Header.Set("Content-Type", "application/json")
 req.Header.Set("Accept", "application/json")
 resp, err := t.Client.Do(req); if err != nil { return nil, err }
 defer resp.Body.Close()
 if resp.StatusCode < 200 || resp.StatusCode >= 300 { return nil, fmt.Errorf("tavily search returned HTTP %d", resp.StatusCode) }
 var p struct { Results []struct { Title string; URL string; Content string } }
 if err := json.NewDecoder(resp.Body).Decode(&p); err != nil { return nil, err }
 out := make([]SearchResult, 0, len(p.Results))
 for _, x := range p.Results { out = append(out, SearchResult{Title:x.Title, URL:x.URL, Description:x.Content, Source:"tavily"}) }
 return out, nil
}
