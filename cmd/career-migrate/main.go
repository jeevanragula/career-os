package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "sort"
    "strings"
    "time"

    "github.com/jeevanragula/career-os/internal/store"
)

func main() {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        log.Fatal("DATABASE_URL is required")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

    s, err := store.Open(ctx, dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer s.Close()

    if _, err = s.DB.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations(version TEXT PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT now())`); err != nil {
        log.Fatal(err)
    }

    dir := os.Getenv("MIGRATIONS_DIR")
    if dir == "" {
        dir = "database/migrations"
    }

    entries, err := os.ReadDir(dir)
    if err != nil {
        log.Fatal(err)
    }

    var files []string
    for _, e := range entries {
        if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
            files = append(files, e.Name())
        }
    }
    sort.Strings(files)

    // The original Docker setup applied migrations 001+ without recording
    // migration history. For an existing legacy DB, baseline only the core
    // schema migrations through 014. Discovery/harvest migrations 015+ must
    // actually execute because those tables may be absent.
    var migrationCount int
    if err := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM schema_migrations`).Scan(&migrationCount); err != nil {
        log.Fatal(err)
    }
    if migrationCount == 0 {
        var legacySkillsExists bool
        if err := s.DB.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema='public' AND table_name='skills')`).Scan(&legacySkillsExists); err != nil {
            log.Fatal(err)
        }
        if legacySkillsExists {
            for _, name := range files {
                if name >= "001_" && name < "015_" {
                    if _, err := s.DB.ExecContext(ctx, `INSERT INTO schema_migrations(version) VALUES($1) ON CONFLICT DO NOTHING`, name); err != nil {
                        log.Fatalf("baseline migration %s: %v", name, err)
                    }
                }
            }
            fmt.Println("baselined legacy migrations 001-014")
        }
    }

    for _, name := range files {
        var exists bool
        if err = s.DB.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version=$1)`, name).Scan(&exists); err != nil {
            log.Fatal(err)
        }
        if exists {
            continue
        }

        raw, err := os.ReadFile(filepath.Join(dir, name))
        if err != nil {
            log.Fatal(err)
        }
        if _, err = s.DB.ExecContext(ctx, string(raw)); err != nil {
            log.Fatalf("migration %s: %v", name, err)
        }
        if _, err = s.DB.ExecContext(ctx, `INSERT INTO schema_migrations(version) VALUES($1)`, name); err != nil {
            log.Fatal(err)
        }
        fmt.Println("applied", name)
    }
}
