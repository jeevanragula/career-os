package store

import "context"

type JobForAnalysis struct { ID, VersionID, Title, Company, Location, Description string }
type CareerClaim struct { ID, Text, Confidentiality string }
type AnalysisRequirement struct { Category, Text, Priority, SourceQuote string; Matches []AnalysisMatch }
type AnalysisMatch struct { ClaimID, Status, Rationale string; Confidence float64 }

func (s *Store) GetJobForAnalysis(ctx context.Context, jobID string) (JobForAnalysis,error) {
    var x JobForAnalysis
    err:=s.DB.QueryRowContext(ctx, "SELECT j.id,v.id,v.title,v.company,COALESCE(v.location,''),v.description FROM jobs j JOIN job_versions v ON v.job_id=j.id WHERE j.id=$1 ORDER BY v.observed_at DESC LIMIT 1",jobID).Scan(&x.ID,&x.VersionID,&x.Title,&x.Company,&x.Location,&x.Description)
    return x,err
}
func (s *Store) ApplicationSafeClaims(ctx context.Context)([]CareerClaim,error){
    rows,err:=s.DB.QueryContext(ctx, "SELECT id,text,confidentiality FROM claims WHERE verification_status IN ('verified','user_asserted') AND confidentiality IN ('public','application_safe') ORDER BY claim_key");if err!=nil{return nil,err};defer rows.Close()
    var out []CareerClaim
    for rows.Next(){var c CareerClaim;if err:=rows.Scan(&c.ID,&c.Text,&c.Confidentiality);err!=nil{return nil,err};out=append(out,c)}
    return out,rows.Err()
}
func (s *Store) SaveAnalysis(ctx context.Context,versionID,summary,model,prompt string,confidence float64,reqs []AnalysisRequirement) error {
    tx,err:=s.DB.BeginTx(ctx,nil);if err!=nil{return err};defer tx.Rollback();var aid string
    if err=tx.QueryRowContext(ctx, "INSERT INTO job_analyses(job_version_id,summary,overall_confidence,model_provider,model_name,prompt_version) VALUES($1,$2,$3,$4,$5,$6) RETURNING id",versionID,summary,confidence,"openai-compatible",model,prompt).Scan(&aid);err!=nil{return err}
    for _,r:=range reqs{
        var rid string
        if err=tx.QueryRowContext(ctx, "INSERT INTO job_requirements(analysis_id,category,requirement_text,priority,source_quote) VALUES($1,$2,$3,$4,$5) RETURNING id",aid,r.Category,r.Text,r.Priority,r.SourceQuote).Scan(&rid);err!=nil{return err}
        for _,m:=range r.Matches{if _,err=tx.ExecContext(ctx,"INSERT INTO requirement_matches(requirement_id,claim_id,status,confidence,rationale) VALUES($1,$2,$3,$4,$5)",rid,m.ClaimID,m.Status,m.Confidence,m.Rationale);err!=nil{return err}}
    }
    return tx.Commit()
}
func (s *Store) CountJobs(ctx context.Context)(int,error){var n int;err:=s.DB.QueryRowContext(ctx,"SELECT count(*) FROM jobs").Scan(&n);return n,err}
