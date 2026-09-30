package career
import("encoding/json";"net/http";"strings";"github.com/jeevanragula/career-os/internal/ai")
type HTTP struct{Store *Store;AI ai.Client}
func(h HTTP)Dashboard(w http.ResponseWriter,r *http.Request){if h.Store==nil{write(w,503,map[string]string{"error":"database not configured"});return};d,e:=h.Store.Dashboard(r.Context());if e!=nil{write(w,500,map[string]string{"error":e.Error()});return};write(w,200,d)}
func(h HTTP)Agents(w http.ResponseWriter,r *http.Request){write(w,200,map[string]any{"agents":Registry})}
func(h HTTP)RunAgent(w http.ResponseWriter,r *http.Request){var in struct{Name string `json:"name"`;Trigger string `json:"trigger"`};if e:=json.NewDecoder(r.Body).Decode(&in);e!=nil{write(w,400,map[string]string{"error":e.Error()});return};if in.Trigger==""{in.Trigger="manual"};out,e:=RunAgent(r.Context(),h.Store,h.AI,in.Name,in.Trigger);if e!=nil{write(w,502,map[string]string{"error":e.Error()});return};write(w,200,out)}
func(h HTTP)CompleteRecommendation(w http.ResponseWriter,r *http.Request){id:=strings.TrimPrefix(r.URL.Path,"/api/recommendations/");status:=r.URL.Query().Get("status");if status==""{status="accepted"};if e:=h.Store.CompleteRecommendation(r.Context(),id,status);e!=nil{write(w,500,map[string]string{"error":e.Error()});return};write(w,200,map[string]string{"status":status})}
func(h HTTP)CompleteTask(w http.ResponseWriter,r *http.Request){id:=strings.TrimPrefix(r.URL.Path,"/api/tasks/");if e:=h.Store.CompleteTask(r.Context(),id);e!=nil{write(w,500,map[string]string{"error":e.Error()});return};write(w,200,map[string]string{"status":"done"})}
func write(w http.ResponseWriter,status int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(status);_=json.NewEncoder(w).Encode(v)}
