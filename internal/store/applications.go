package store

import "context"

func (s *Store) CreateApplication(ctx context.Context,jobID,resumeVariantID string)(string,error){
    var id string
    err:=s.DB.QueryRowContext(ctx,"INSERT INTO applications(job_id,resume_variant_id,status) VALUES($1,NULLIF($2,'')::uuid,'preparing') RETURNING id",jobID,resumeVariantID).Scan(&id)
    return id,err
}
