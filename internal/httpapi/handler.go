package httpapi

import (
    "context"
    "encoding/json"
    "net/http"
    "os"
    "strings"
    "github.com/jeevanragula/career-os/internal/ai"
    "github.com/jeevanragula/career-os/internal/jobs"
    "github.com/jeevanragula/career-os/internal/jobs/providers"
    "github.com/jeevanragula/career-os/internal/store"
)

type Handler struct { Store *store.Store; AI ai.Client }

func NewHandler(s *store.Store) http.Handler {
    h:=&Handler{Store:s,AI:ai.Client{BaseURL:os.Getenv("AI_BASE_URL"),APIKey:os.Getenv("AI_API_KEY"),Model:os.Getenv("AI_MODEL")}}
    mux:=http.NewServeMux()
    mux.HandleFunc("GET /healthz",health)
    mux.HandleFunc("GET /readyz",func(w http.ResponseWriter,r *http.Request){if h.Store==nil{writeJSON(w,503,map[string]string{"status":"not_ready"});return};writeJSON(w,200,map[string]string{"status":"ready"})})
    mux.HandleFunc("GET /api/jobs",h.listJobs)
    mux.HandleFunc("POST /api/discover",h.discover)
    mux.HandleFunc("POST /api/analyze",h.analyze)
    mux.HandleFunc("GET /",h.index)
    return mux
}
func health(w http.ResponseWriter,r *http.Request){writeJSON(w,200,map[string]string{"status":"ok"})}
type discoverRequest struct { Provider string `json:"provider"`; Name string `json:"name"`; BaseURL string `json:"base_url"`; Keywords []string `json:"keywords"` }
func (h *Handler) discover(w http.ResponseWriter,r *http.Request){
    var in discoverRequest;if err:=json.NewDecoder(r.Body).Decode(&in);err!=nil{writeJSON(w,400,map[string]string{"error":err.Error()});return}
    if h.Store==nil{writeJSON(w,503,map[string]string{"error":"database not configured"});return}
    var adapter jobs.SourceAdapter
    switch strings.ToLower(in.Provider){case "lever":adapter=providers.LeverAdapter{};case "ashby":adapter=providers.AshbyAdapter{};default:writeJSON(w,400,map[string]string{"error":"provider must be lever or ashby"});return}
    result,err:=adapter.Discover(r.Context(),jobs.SourceConfig{Name:in.Name,BaseURL:in.BaseURL},jobs.DiscoveryQuery{Keywords:in.Keywords});if err!=nil{writeJSON(w,502,map[string]string{"error":err.Error()});return}
    ingested:=0
    for _,o:=range result.Observations{if !matches(o,in.Keywords){continue};id,err:=h.Store.UpsertJob(r.Context(),jobs.NormalizeInput{Source:o.Source,SourceJobID:o.SourceJobID,URL:o.SourceURL,Title:o.RawTitle,Company:o.RawCompany,Location:o.RawLocation,Description:o.RawDescription,ObservedAt:o.ObservedAt});if err!=nil{writeJSON(w,500,map[string]string{"error":err.Error()});return};if err=h.Store.InsertObservation(r.Context(),id,o);err!=nil{writeJSON(w,500,map[string]string{"error":err.Error()});return};if err=h.Store.UpsertVersion(r.Context(),id,jobs.NormalizeInput{Source:o.Source,SourceJobID:o.SourceJobID,URL:o.SourceURL,Title:o.RawTitle,Company:o.RawCompany,Location:o.RawLocation,Description:o.RawDescription,ObservedAt:o.ObservedAt});err!=nil{writeJSON(w,500,map[string]string{"error":err.Error()});return};ingested++}
    writeJSON(w,200,map[string]any{"observed":len(result.Observations),"ingested":ingested})
}
func matches(o jobs.JobObservation,keywords []string) bool {if len(keywords)==0{return true};hay:=strings.ToLower(o.RawTitle+" "+o.RawDescription+" "+o.RawLocation);for _,k:=range keywords{if strings.Contains(hay,strings.ToLower(strings.TrimSpace(k))){return true}};return false}
func (h *Handler) listJobs(w http.ResponseWriter,r *http.Request){if h.Store==nil{writeJSON(w,503,map[string]string{"error":"database not configured"});return};items,err:=h.Store.ListJobs(r.Context(),50);if err!=nil{writeJSON(w,500,map[string]string{"error":err.Error()});return};writeJSON(w,200,items)}
type analyzeRequest struct { JobID string `json:"job_id"` }
func (h *Handler) analyze(w http.ResponseWriter,r *http.Request){var in analyzeRequest;if err:=json.NewDecoder(r.Body).Decode(&in);err!=nil{writeJSON(w,400,map[string]string{"error":err.Error()});return};if h.Store==nil{writeJSON(w,503,map[string]string{"error":"database not configured"});return};if h.AI.BaseURL==""||h.AI.Model==""{writeJSON(w,503,map[string]string{"error":"AI_BASE_URL and AI_MODEL are required"});return};err:=(&ai.Analyzer{Client:h.AI,Store:h.Store}).Analyze(context.Background(),in.JobID);if err!=nil{writeJSON(w,502,map[string]string{"error":err.Error()});return};writeJSON(w,200,map[string]string{"status":"analyzed"})}
func (h *Handler) index(w http.ResponseWriter,r *http.Request){http.ServeFile(w,r,"web/index.html")}
func writeJSON(w http.ResponseWriter,status int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(status);_ = json.NewEncoder(w).Encode(v)}
