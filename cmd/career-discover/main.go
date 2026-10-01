package main

import (
 "context"
 "log"
 "os"
 "strings"
 "time"
 "github.com/jeevanragula/career-os/internal/discovery"
 "github.com/jeevanragula/career-os/internal/jobs"
 "github.com/jeevanragula/career-os/internal/opportunity"
 "github.com/jeevanragula/career-os/internal/store"
)

func main(){
 dsn:=os.Getenv("DATABASE_URL");if dsn==""{log.Fatal("DATABASE_URL is required")}
 if os.Getenv("BRAVE_SEARCH_API_KEY")==""{log.Fatal("BRAVE_SEARCH_API_KEY is required")}
 ctx,cancel:=context.WithTimeout(context.Background(),20*time.Minute);defer cancel()
 s,err:=store.Open(ctx,dsn);if err!=nil{log.Fatal(err)};defer s.Close()
 engine:=discovery.Engine{Search:discovery.NewBrave()}
 started:=time.Now().UTC()
 candidates,err:=engine.Discover(ctx,split(os.Getenv("CAREEROS_DISCOVERY_ROLES")),split(os.Getenv("CAREEROS_DISCOVERY_DOMAINS")),split(os.Getenv("CAREEROS_DISCOVERY_LOCATIONS")),started.Year());if err!=nil{log.Fatal(err)}
 run,err:=discovery.StartHarvestRun(ctx,s.DB,started.Format("20060102"));if err!=nil{log.Fatal(err)}
 resolved:=engine.ResolveCareerPages(ctx,candidates);result:=engine.HarvestJobs(ctx,resolved,75)
 ingested:=0
 for _,o:=range result.Observations{
  in:=jobs.NormalizeInput{Source:o.Source,SourceJobID:o.SourceJobID,URL:o.SourceURL,Title:o.RawTitle,Company:o.RawCompany,Location:o.RawLocation,Description:o.RawDescription,ObservedAt:o.ObservedAt}
  id,err:=s.UpsertJob(ctx,in);if err!=nil{log.Printf("upsert %q: %v",o.RawTitle,err);continue}
  if err=s.InsertObservation(ctx,id,o);err!=nil{log.Printf("observation %q: %v",o.RawTitle,err);continue}
  if err=s.UpsertVersion(ctx,id,in);err!=nil{log.Printf("version %q: %v",o.RawTitle,err);continue}
  _=s.SaveJobEvidence(ctx,id,o.SourceURL,"autonomous-web",map[string]any{"source":o.Source,"observed_at":o.ObservedAt})
  ingested++
 }
 matches,err:=opportunity.Compute(ctx,s.DB,200);if err!=nil{log.Printf("matching: %v",err)} else {for _,m:=range matches{if err:=opportunity.Save(ctx,s.DB,m);err!=nil{log.Printf("save match %s: %v",m.JobID,err)}}}
 status:="completed";_ = discovery.FinishHarvestRun(ctx,s.DB,run.ID,int64(result.Pages),int64(len(result.Observations)),int64(ingested),status,"")
 log.Printf("CareerOS harvest: candidates=%d resolved=%d pages=%d jobs=%d ingested=%d matched=%d expired=deferred",len(candidates),len(resolved),result.Pages,len(result.Observations),ingested,len(matches))
}
func split(v string)[]string{var out []string;for _,x:=range strings.Split(v,","){if x=strings.TrimSpace(x);x!=""{out=append(out,x)}};return out}
