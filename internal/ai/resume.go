package ai

import (
    "context"
    "fmt"
    "strings"
    "github.com/jeevanragula/career-os/internal/store"
)

func (a Analyzer) GenerateResume(ctx context.Context, jobID string) (string,error) {
    j,err:=a.Store.GetJobForAnalysis(ctx,jobID);if err!=nil{return "",err}
    claims,err:=a.Store.ApplicationSafeClaims(ctx);if err!=nil{return "",err}
    var b strings.Builder
    for _,c:=range claims{fmt.Fprintf(&b,"- %s: %s\n",c.ID,c.Text)}
    system:="You are CareerOS Resume Agent. Generate a concise ATS-friendly Markdown resume using ONLY the supplied career claims. Never invent metrics, employers, dates, technologies, ownership, or achievements. Tailor emphasis to the target job. Do not mention the internal claim IDs in the visible resume."
    user:=fmt.Sprintf("TARGET JOB: %s at %s\nLOCATION: %s\nJOB DESCRIPTION:\n%s\n\nVERIFIED/APPLICATION-SAFE CLAIMS:\n%s\n\nProduce Markdown with name placeholder, headline, summary, skills, experience, selected projects, and patents where supported. Keep chronology truthful.",j.Title,j.Company,j.Location,j.Description,b.String())
    return a.Client.Complete(ctx,system,user)
}
