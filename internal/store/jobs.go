package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/jeevanragula/career-os/internal/jobs"
)

func (s *Store) UpsertJob(ctx context.Context, in jobs.NormalizeInput) (string, error) {
	var id string
	err := s.DB.QueryRowContext(ctx, `
		INSERT INTO jobs(company,title,canonical_url,location,remote_mode,employment_type,first_seen_at,last_seen_at)
		VALUES(NULLIF($1,''),NULLIF($2,''),NULLIF($3,''),NULLIF($4,''),NULLIF($5,''),NULLIF($6,''),$7,$7)
		RETURNING id
	`, in.Company,in.Title,jobs.CanonicalizeURL(in.URL),in.Location,in.RemoteMode,in.Employment,in.ObservedAt).Scan(&id)
	if err == nil { return id,nil }
	if !isUniqueViolation(err) { return "",err }
	err = s.DB.QueryRowContext(ctx, `SELECT id FROM jobs WHERE canonical_url=$1 LIMIT 1`, jobs.CanonicalizeURL(in.URL)).Scan(&id)
	return id,err
}

func (s *Store) InsertObservation(ctx context.Context, jobID string, o jobs.JobObservation) error {
	_, err := s.DB.ExecContext(ctx,`
		INSERT INTO job_observations(job_id,source,source_job_id,source_url,observed_at,raw_title,raw_company,raw_location,raw_description)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)`,jobID,o.Source,o.SourceJobID,o.SourceURL,o.ObservedAt,o.RawTitle,o.RawCompany,o.RawLocation,o.RawDescription)
	return err
}

func (s *Store) UpsertVersion(ctx context.Context, jobID string, in jobs.NormalizeInput) error {
	_, err := s.DB.ExecContext(ctx,`
		INSERT INTO job_versions(job_id,content_hash,observed_at,title,company,location,remote_mode,employment_type,description)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT(job_id,content_hash) DO NOTHING`,jobID,jobs.ContentHash(in),in.ObservedAt,in.Title,in.Company,in.Location,in.RemoteMode,in.Employment,in.Description)
	return err
}

func isUniqueViolation(err error) bool {
	// The store treats any failed canonical insert as a lookup opportunity.
	// A future repository layer can inspect pgconn.PgError.Code for 23505.
	return err != nil
}

var _ = sql.ErrNoRows
var _ = time.Time{}
