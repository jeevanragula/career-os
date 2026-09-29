package ai

import (
    "context"
    "encoding/json"
    "fmt"
    "strconv"
    "strings"
    "github.com/jeevanragula/career-os/internal/store"
)

type Analyzer struct { Client Client; Store *store.Store }

func AnalyzerOutput(raw string) (string,float64,[]store.AnalysisRequirement,error) {
    raw=strings.TrimSpace(strings.TrimPrefix(raw,"\\x60\\x60\\x60json"))
    raw=strings.TrimSuffix(raw,"\\x60\\x60\\x60")
    var out map[string]any
    if err:=json.Unmarshal([]byte(strings.TrimSpace(raw)),&out);err!=nil{return "",0,nil,fmt.Errorf("AI returned invalid JSON: %w",err)}
    summary,_:=out["summary"].(string)
    confidence,_:=out["overall_confidence"].(float64)
    var reqs []store.AnalysisRequirement
    list,_:=out["requirements"].([]any)
    for _,item:=range list {
        m,_:=item.(map[string]any)
        r:=store.AnalysisRequirement{}
        r.Category,_=m["category"].(string);r.Text,_=m["text"].(string);r.Priority,_=m["priority"].(string);r.SourceQuote,_=m["source_quote"].(string)
        matches,_:=m["matches"].([]any)
        for _,mi:=range matches { mm,_:=mi.(map[string]any); am:=store.AnalysisMatch{};am.ClaimID,_=mm["claim_id"].(string);am.Status,_=mm["status"].(string);am.Rationale,_=mm["rationale"].(string);if v,ok:=mm["confidence"].(float64);ok{am.Confidence=v};r.Matches=append(r.Matches,am) }
        reqs=append(reqs,r)
    }
    return summary,confidence,reqs,nil
}

func (a Analyzer) Analyze(ctx context.Context, jobID string) error {
    j,err:=a.Store.GetJobForAnalysis(ctx,jobID);if err!=nil{return err}
    claims,err:=a.Store.ApplicationSafeClaims(ctx);if err!=nil{return err}
    var b strings.Builder
    for _,c:=range claims{fmt.Fprintf(&b,"- %s: %s\n",c.ID,c.Text)}
    system:="You are CareerOS Job Analysis. Return ONLY valid JSON. Never invent career evidence. Extract requirements from the job description and match only against supplied claims. Use status matched, partial, unmatched, or needs_verification. Every requirement must include an exact short source_quote from the job description."
    user:=fmt.Sprintf("JOB: %s at %s\nLOCATION: %s\nDESCRIPTION:\n%s\n\nAPPLICATION-SAFE CAREER CLAIMS:\n%s\n\nReturn JSON with summary, overall_confidence, and requirements. Each requirement needs category, text, priority, source_quote and matches. Each match needs claim_id, status, confidence and rationale.",j.Title,j.Company,j.Location,j.Description,b.String())
    raw,err:=a.Client.Complete(ctx,system,user);if err!=nil{return err}
    summary,confidence,reqs,err:=AnalyzerOutput(raw);if err!=nil{return err}
    if confidence<0||confidence>1{confidence=0}
    return a.Store.SaveAnalysis(ctx,j.VersionID,summary,a.Client.Model,"v1",confidence,reqs)
}

var _=strconv.Itoa
