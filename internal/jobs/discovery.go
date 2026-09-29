package jobs

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
)

type DiscoveryRun struct {
	ID string
	SourceID string
	QueryID string
	StartedAt time.Time
	FinishedAt *time.Time
	Status string
	ObservedCount int
	InsertedCount int
	UpdatedCount int
	FailedCount int
	ErrorMessage string
	Cursor string
}

func QueryFingerprint(q DiscoveryQuery) string {
	parts := []string{joinNormalized(q.Keywords), joinNormalized(q.Companies), joinNormalized(q.Locations), joinNormalized(q.RemoteModes)}
	h := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return hex.EncodeToString(h[:])
}

func joinNormalized(values []string) string {
	out := make([]string, 0, len(values))
	for _, v := range values { out = append(out, NormalizeText(v)) }
	return strings.Join(out, ",")
}
