package main

import (
 "context"
 "embed"
 "fmt"
 "log"
 "os"
 "path/filepath"
 "sort"
 "strings"
 "time"
 "github.com/jeevanragula/career-os/internal/store"
)

//go:embed ../../database/migrations/*.sql
var migrationFiles embed.FS

func main(){
 dsn:=os.Getenv("DATABASE_URL");if dsn==""{log.Fatal("DATABASE_URL is required")}
 ctx,cancel:=context.WithTimeout(context.Background(),5*time.Minute);defer cancel()
 s,err:=store.Open(ctx,dsn);if err!=nil{log.Fatal(err)};defer s.Close()
 if _,err=s.DB.ExecContext(ctx,`CREATE TABLE IF NOT EXISTS schema_migrations(version TEXT PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT now())`);err!=nil{log.Fatal(err)}
 entries,err:=migrationFiles.ReadDir("../../database/migrations");if err!=nil{log.Fatal(err)}
 var files []string;for _,e:=range entries{if !e.IsDir()&&strings.HasSuffix(e.Name(),".sql"){files=append(files,e.Name())}};sort.Strings(files)
 for _,name:=range files{
  var exists bool;if err=s.DB.QueryRowContext(ctx,`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version=$1)`,name).Scan(&exists);err!=nil{log.Fatal(err)};if exists{continue}
  raw,err:=migrationFiles.ReadFile(filepath.Join("../../database/migrations",name));if err!=nil{log.Fatal(err)}
  if _,err=s.DB.ExecContext(ctx,string(raw));err!=nil{log.Fatalf("migration %s: %v",name,err)}
  if _,err=s.DB.ExecContext(ctx,`INSERT INTO schema_migrations(version) VALUES($1)`,name);err!=nil{log.Fatal(err)}
  fmt.Println("applied",name)
 }
}
