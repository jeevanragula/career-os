package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jeevanragula/career-os/internal/discovery"
	"github.com/jeevanragula/career-os/internal/jobs"
	"github.com/jeevanragula/career-os/internal/store"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" { log.Fatal("DATABASE_URL is required") }
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()
	s, err := store.Open(ctx, dsn)
	if err != nil { log.Fatal(err) }
	defer s.Close()

	engine := discovery.Engine{Search: discovery.NewBrave()}
	candidates, err := engine.Discover(ctx, split(os.Getenv("CAREEROS_DISCOVERY_ROLES")), split(os.Getenv("CAREEROS_DISCOVERY_DOMAINS")), split(os.Getenv("CAREEROS_DISCOVERY_LOCATIONS")), time.Now().Year())
	if err != nil { log.Fatal(err) }
	resolved := engine.ResolveCareerPages(ctx, candidates)
	result := engine.HarvestJobs(ctx, resolved, 75)

	ingested := 0
	for _, o := range result.Observations {
		in := jobs.NormalizeInput{Source:o.Source, SourceJobID:o.SourceJobID, URL:o.SourceURL, Title:o.RawTitle, Company:o.RawCompany, Location:o.RawLocation, Description:o.RawDescription, ObservedAt:o.ObservedAt}
		id, err := s.UpsertJob(ctx, in)
		if err != nil { log.Printf("upsert %q: %v", o.RawTitle, err); continue }
		if err := s.InsertObservation(ctx, id, o); err != nil { log.Printf("observation %q: %v", o.RawTitle, err); continue }
		if err := s.UpsertVersion(ctx, id, in); err != nil { log.Printf("version %q: %v", o.RawTitle, err); continue }
		ingested++
	}
	log.Printf("CareerOS harvest: candidates=%d resolved=%d pages=%d jobs=%d ingested=%d", len(candidates), len(resolved), result.Pages, len(result.Observations), ingested)
}

func split(v string) []string {
	var out []string
	for _, x := range strings.Split(v, ",") {
		if x = strings.TrimSpace(x); x != "" { out = append(out, x) }
	}
	return out
}
