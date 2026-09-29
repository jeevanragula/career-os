package store

import "context"

func (s *Store) SaveResume(ctx context.Context, jobID, title, summary, strategy, content, model string) (string,error) {
    tx,err:=s.DB.BeginTx(ctx,nil);if err!=nil{return "",err};defer tx.Rollback()
    var id string
    if err=tx.QueryRowContext(ctx,"INSERT INTO resume_variants(job_id,strategy,title,summary,model_provider,model_name,prompt_version) VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING id",jobID,strategy,title,summary,"openai-compatible",model,"v1").Scan(&id);err!=nil{return "",err}
    if _,err=tx.ExecContext(ctx,"INSERT INTO resume_variant_content(variant_id,format,content) VALUES($1,'markdown',$2)",id,content);err!=nil{return "",err}
    return id,tx.Commit()
}
