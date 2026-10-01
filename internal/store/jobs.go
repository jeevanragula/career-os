package store

import (
 "context"
 "database/sql"
 "encoding/json"
 "fmt"
 "time"
 "github.com/jackc/pgx/v5"
 "github.com/jackc/pgx/v5/pgconn"
 "github.com/jackc/pgx/v5/stdlib"
 "github.com/jeevanragula/career-os/internal/jobs"
)

func (s *Store) UpsertJob(ctx context.Context,in jobs.NormalizeInput)(string,error){
 u:=jobs.CanonicalizeURL(in.URL);var id string
 err:=s.DB.QueryRowContext(ctx,`INSERT INTO jobs(company,title,canonical_url,location,remote_mode,employment_type,first_seen_at,last_seen_at,updated_at)
 VALUES($1,$2,NULLIF($3,''),NULLIF($4,''),NULLIF($5,''),NULLIF($6,''),$7,$7,$7)
 ON CONFLICT(canonical_url) WHERE canonical_url IS NOT NULL
 DO UPDATE SET company=EXCLUDED.company,title=EXCLUDED.title,location=EXCLUDED.location,remote_mode=EXCLUDED.remote_mode,employment_type=EXCLUDED.employment_type,last_seen_at=EXCLUDED.last_seen_at,updated_at=EXCLUDED.updated_at,status='open'
 RETURNING id`,in.Company,in.Title,u,in.Location,in.RemoteMode,in.Employment,in.ObservedAt).Scan(&id)
 if err!=nil{return "",fmt.Errorf("upsert job: %w",err)};return id,nil
}
func(s *Store)InsertObservation(ctx context.Context,jobID string,o jobs.JobObservation)error{_,err:=s.DB.ExecContext(ctx,`INSERT INTO job_observations(job_id,source,source_job_id,source_url,observed_at,raw_title,raw_company,raw_location,raw_description) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)`,jobID,o.Source,o.SourceJobID,o.SourceURL,o.ObservedAt,o.RawTitle,o.RawCompany,o.RawLocation,o.RawDescription);return err}
func(s *Store)UpsertVersion(ctx context.Context,jobID string,in jobs.NormalizeInput)error{_,err:=s.DB.ExecContext(ctx,`INSERT INTO job_versions(job_id,content_hash,observed_at,title,company,location,remote_mode,employment_type,description) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9) ON CONFLICT(job_id,content_hash) DO NOTHING`,jobID,jobs.ContentHash(in),in.ObservedAt,in.Title,in.Company,in.Location,in.RemoteMode,in.Employment,in.Description);return err}
func(s *Store)SaveJobEvidence(ctx context.Context,jobID,url,kind string,details map[string]any)error{raw,_:=json.Marshal(details);_,err:=s.DB.ExecContext(ctx,`INSERT INTO job_source_evidence(job_id,source_url,evidence_type,observed_at,details) VALUES($1,$2,$3,now(),$4) ON CONFLICT(job_id,source_url,evidence_type) DO UPDATE SET observed_at=now(),details=EXCLUDED.details`,jobID,url,kind,raw);return err}
func(s *Store)MarkStaleJobs(ctx context.Context,cutoff time.Time)(int64,error){r,err:=s.DB.ExecContext(ctx,`UPDATE jobs SET status='expired',updated_at=now() WHERE status='open' AND last_seen_at < $1`,cutoff);if err!=nil{return 0,err};return r.RowsAffected()}
func(s *Store)ListJobs(ctx context.Context,limit int)([]jobs.Job,error){if limit<=0||limit>200{limit=50};rows,err:=s.DB.QueryContext(ctx,`SELECT id,company,title,COALESCE(canonical_url,''),COALESCE(location,''),COALESCE(remote_mode,''),COALESCE(employment_type,''),first_seen_at,last_seen_at,status FROM jobs ORDER BY last_seen_at DESC LIMIT $1`,limit);if err!=nil{return nil,err};defer rows.Close();var out []jobs.Job;for rows.Next(){var j jobs.Job;if err:=rows.Scan(&j.ID,&j.Company,&j.Title,&j.CanonicalURL,&j.Location,&j.RemoteMode,&j.EmploymentType,&j.FirstSeenAt,&j.LastSeenAt,&j.Status);err!=nil{return nil,err};out=append(out,j)};return out,rows.Err()}
var _=sql.ErrNoRows;var _=pgx.ErrNoRows;var _=pgconn.PgError{};var _=stdlib.GetDefaultDriver
