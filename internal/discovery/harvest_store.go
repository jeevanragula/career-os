package discovery

import (
 "context"
 "database/sql"
 "time"
)

type HarvestRun struct { ID string; StartedAt time.Time }

func StartHarvestRun(ctx context.Context,db *sql.DB,key string,candidates int)(HarvestRun,error){
 var r HarvestRun
 err:=db.QueryRowContext(ctx,`INSERT INTO opportunity_harvest_runs(run_key,candidates_seen,status) VALUES($1,$2,'running')
 ON CONFLICT(run_key) DO UPDATE SET started_at=now(),candidates_seen=EXCLUDED.candidates_seen,status='running',error=NULL
 RETURNING id,started_at`,key,candidates).Scan(&r.ID,&r.StartedAt)
 return r,err
}
func FinishHarvestRun(ctx context.Context,db *sql.DB,id string,pages,observed,ingested int64,status string,errText string) error {
 _,err:=db.ExecContext(ctx,`UPDATE opportunity_harvest_runs SET completed_at=now(),pages_fetched=$2,jobs_observed=$3,jobs_ingested=$4,status=$5,error=NULLIF($6,'') WHERE id=$1`,id,pages,observed,ingested,status,errText);return err
}
func MarkCareerURL(ctx context.Context,db *sql.DB,name,url string) error {
 _,err:=db.ExecContext(ctx,`UPDATE discovered_companies SET career_url=$2,status='verified',last_seen_at=now() WHERE name=$1`,name,url);return err
}
