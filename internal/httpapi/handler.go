package httpapi

import (
 "context"
 "encoding/json"
 "net/http"
 "os"
 "strings"
 "time"
 "github.com/jeevanragula/career-os/internal/ai"
 "github.com/jeevanragula/career-os/internal/career"
 "github.com/jeevanragula/career-os/internal/discovery"
 "github.com/jeevanragula/career-os/internal/jobs"
 "github.com/jeevanragula/career-os/internal/jobs/providers"
 "github.com/jeevanragula/career-os/internal/store"
)

type Handler struct { Store *store.Store; AI ai.Client; Career *career.HTTP; AuthUser string; AuthPassword string }

func NewHandler(stores ...*store.Store) http.Handler {
 var s *store.Store
 if len(stores)>0 { s=stores[0] }
 aiClient:=ai.Client{BaseURL:os.Getenv("AI_BASE_URL"),APIKey:os.Getenv("AI_API_KEY"),Model:os.Getenv("AI_MODEL")}
 h:=&Handler{Store:s,AI:aiClient,AuthUser:os.Getenv("CAREEROS_AUTH_USER"),AuthPassword:os.Getenv("CAREEROS_AUTH_PASSWORD")}
 if s!=nil { h.Career=&career.HTTP{Store:&career.Store{DB:s.DB},AI:aiClient} }
 mux:=http.NewServeMux()
 mux.HandleFunc("GET /healthz",health)
 mux.HandleFunc("GET /readyz",func(w http.ResponseWriter,r *http.Request){if h.Store==nil{writeJSON(w,503,map[string]string{"status":"not_ready"});return};writeJSON(w,200,map[string]string{"status":"ready"})})
 mux.HandleFunc("GET /api/status",h.status)
 mux.HandleFunc("GET /api/jobs",h.listJobs)
 mux.HandleFunc("POST /api/discover",h.discover)
 mux.HandleFunc("POST /api/discover/automatic",h.automaticDiscover)
 mux.HandleFunc("GET /api/discovery/candidates",h.discoveryCandidates)
 mux.HandleFunc("GET /api/discovery/runs/{id}",h.discoveryRun)
 mux.HandleFunc("POST /api/analyze",h.analyze)
 mux.HandleFunc("POST /api/resume",h.resume)
 mux.HandleFunc("POST /api/applications",h.application)
 mux.HandleFunc("GET /api/dashboard",func(w http.ResponseWriter,r *http.Request){if h.Career==nil{writeJSON(w,503,map[string]string{"error":"database not configured"});return};h.Career.Dashboard(w,r)})
 mux.HandleFunc("GET /api/agents",func(w http.ResponseWriter,r *http.Request){h.Career.Agents(w,r)})
 mux.HandleFunc("POST /api/agents/run",func(w http.ResponseWriter,r *http.Request){if h.Career==nil{writeJSON(w,503,map[string]string{"error":"database not configured"});return};h.Career.RunAgent(w,r)})
 mux.HandleFunc("POST /api/recommendations/",func(w http.ResponseWriter,r *http.Request){if h.Career==nil{writeJSON(w,503,map[string]string{"error":"database not configured"});return};h.Career.CompleteRecommendation(w,r)})
 mux.HandleFunc("POST /api/tasks/",func(w http.ResponseWriter,r *http.Request){if h.Career==nil{writeJSON(w,503,map[string]string{"error":"database not configured"});return};h.Career.CompleteTask(w,r)})
 mux.HandleFunc("GET /dashboard.css",func(w http.ResponseWriter,r *http.Request){http.ServeFile(w,r,"web/dashboard.css")})
 mux.HandleFunc("GET /dashboard.js",func(w http.ResponseWriter,r *http.Request){http.ServeFile(w,r,"web/dashboard.js")})
 mux.HandleFunc("GET /",h.index)
 return withBasicAuth(mux,h.AuthUser,h.AuthPassword)
}
func withBasicAuth(next http.Handler,user,password string) http.Handler {
 if user=="" || password=="" { return next }
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
  if r.URL.Path=="/healthz" || r.URL.Path=="/readyz" { next.ServeHTTP(w,r); return }
  u,p,ok:=r.BasicAuth()
  if !ok || u!=user || p!=password { w.Header().Set("WWW-Authenticate", "Basic realm=CareerOS"); w.WriteHeader(http.StatusUnauthorized); return }
  next.ServeHTTP(w,r)
 })
}

func health(w http.ResponseWriter,r *http.Request){writeJSON(w,200,map[string]string{"status":"ok"})}
func (h *Handler) status(w http.ResponseWriter,r *http.Request) {
 aiConfigured := h.AI.BaseURL != "" && h.AI.Model != ""
 tavilyConfigured := os.Getenv("TAVILY_API_KEY") != ""
 writeJSON(w,200,map[string]any{
  "database": h.Store != nil,
  "tavily": tavilyConfigured,
  "ai": aiConfigured,
  "ai_provider": h.AI.BaseURL,
 })
}

type automaticDiscoverRequest struct {
 Roles []string
 Domains []string
 Locations []string
}

func (h *Handler) automaticDiscover(w http.ResponseWriter,r *http.Request) {
 if h.Store==nil { writeJSON(w,503,map[string]string{"error":"database not configured"}); return }
 var in automaticDiscoverRequest
 _ = json.NewDecoder(r.Body).Decode(&in)
 ds:=discovery.Store{DB:h.Store.DB}
 runID,err:=ds.CreateRun(r.Context(),map[string]any{"roles":in.Roles,"domains":in.Domains,"locations":in.Locations})
 if err!=nil { writeJSON(w,500,map[string]string{"error":err.Error()}); return }
 go h.runAutomaticDiscovery(runID,in)
 writeJSON(w,202,map[string]any{"status":"started","run_id":runID})
}

func (h *Handler) discoveryRun(w http.ResponseWriter,r *http.Request) {
 if h.Store==nil { writeJSON(w,503,map[string]string{"error":"database not configured"}); return }
 run,err:=(discovery.Store{DB:h.Store.DB}).GetRun(r.Context(),r.PathValue("id"))
 if err!=nil { writeJSON(w,404,map[string]string{"error":"discovery run not found"}); return }
 writeJSON(w,200,run)
}

func (h *Handler) runAutomaticDiscovery(runID string,in automaticDiscoverRequest) {
 ctx,cancel:=context.WithTimeout(context.Background(),2*time.Minute)
 defer cancel()
 ds:=discovery.Store{DB:h.Store.DB}
 fail:=func(err error){ _=ds.UpdateRun(context.Background(),runID,"failed",0,err.Error()) }

 started:=time.Now().UTC()
 engine:=discovery.Engine{Search:discovery.NewTavily()}
 candidates,err:=engine.Discover(ctx,in.Roles,in.Domains,in.Locations,started.Year())
 if err!=nil { fail(err); return }
 if err=ds.SaveCandidates(ctx,candidates);err!=nil { fail(err); return }

 // First-pass discovery intentionally stops here. Do not crawl or harvest yet.
 result,_:=json.Marshal(map[string]any{
  "search_results":len(candidates),
  "next_step":"review search results before crawling",
 })
 _,_=h.Store.DB.ExecContext(context.Background(),
  "UPDATE discovery_runs SET status='completed',discovered_count=$2,error='',finished_at=now(),query=query || $3::jsonb WHERE id=$1",
  runID,len(candidates),string(result))
}

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
func (h *Handler) analyze(w http.ResponseWriter,r *http.Request){var in analyzeRequest;if err:=json.NewDecoder(r.Body).Decode(&in);err!=nil{writeJSON(w,400,map[string]string{"error":err.Error()});return};if h.Store==nil{writeJSON(w,503,map[string]string{"error":"database not configured"});return};if h.AI.BaseURL==""||h.AI.Model==""{writeJSON(w,503,map[string]string{"error":"AI_BASE_URL and AI_MODEL are required"});return};err:=(&ai.Analyzer{Client:h.AI,Store:h.Store}).Analyze(r.Context(),in.JobID);if err!=nil{writeJSON(w,502,map[string]string{"error":err.Error()});return};writeJSON(w,200,map[string]string{"status":"analyzed"})}
func (h *Handler) index(w http.ResponseWriter,r *http.Request){http.ServeFile(w,r,"web/index.html")}
func writeJSON(w http.ResponseWriter,status int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(status);_ = json.NewEncoder(w).Encode(v)}
type resumeRequest struct { JobID string `json:"job_id"` }
func (h *Handler) resume(w http.ResponseWriter,r *http.Request){var in resumeRequest;if err:=json.NewDecoder(r.Body).Decode(&in);err!=nil{writeJSON(w,400,map[string]string{"error":err.Error()});return};if h.Store==nil||h.AI.BaseURL==""||h.AI.Model==""{writeJSON(w,503,map[string]string{"error":"database and AI configuration are required"});return};content,err:=(&ai.Analyzer{Client:h.AI,Store:h.Store}).GenerateResume(r.Context(),in.JobID);if err!=nil{writeJSON(w,502,map[string]string{"error":err.Error()});return};id,err:=h.Store.SaveResume(r.Context(),in.JobID,"Tailored Resume","Targeted resume","job-tailored",content,h.AI.Model);if err!=nil{writeJSON(w,500,map[string]string{"error":err.Error()});return};writeJSON(w,200,map[string]any{"variant_id":id,"content":content})}
type applicationRequest struct { JobID string `json:"job_id"`; ResumeVariantID string `json:"resume_variant_id"` }
func (h *Handler) application(w http.ResponseWriter,r *http.Request){var in applicationRequest;if err:=json.NewDecoder(r.Body).Decode(&in);err!=nil{writeJSON(w,400,map[string]string{"error":err.Error()});return};if h.Store==nil{writeJSON(w,503,map[string]string{"error":"database not configured"});return};id,err:=h.Store.CreateApplication(r.Context(),in.JobID,in.ResumeVariantID);if err!=nil{writeJSON(w,500,map[string]string{"error":err.Error()});return};writeJSON(w,200,map[string]any{"application_id":id,"status":"preparing","message":"Application package created. Submission still requires explicit human approval."})}

func(h *Handler) discoveryCandidates(w http.ResponseWriter,r *http.Request){if h.Store==nil{writeJSON(w,503,map[string]string{"error":"database not configured"});return};items,e:=(discovery.Store{DB:h.Store.DB}).ListCandidates(r.Context(),50);if e!=nil{writeJSON(w,500,map[string]string{"error":e.Error()});return};writeJSON(w,200,map[string]any{"candidates":items})}
