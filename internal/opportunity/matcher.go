package opportunity

import (
 "context"
 "database/sql"
 "encoding/json"
 "fmt"
 "strings"
)

type Match struct { JobID string; Score float64; Matched []string; Gaps []string; Rationale string }

func Compute(ctx context.Context, db *sql.DB, limit int) ([]Match,error) {
 if limit<=0 || limit>200 { limit=100 }
 rows,err:=db.QueryContext(ctx, `SELECT j.id,j.title,j.company,COALESCE(v.description,''),COALESCE(v.location,'')
 FROM jobs j JOIN LATERAL (SELECT description,location FROM job_versions WHERE job_id=j.id ORDER BY observed_at DESC LIMIT 1) v ON true
 WHERE j.status='open' ORDER BY j.last_seen_at DESC LIMIT $1`,limit)
 if err!=nil{return nil,err}; defer rows.Close()
 var skills []struct{Name string; Level int; Target int}
 sr,err:=db.QueryContext(ctx,`SELECT name,level,COALESCE(target_level,level) FROM skills`); if err!=nil{return nil,err}
 for sr.Next(){var s struct{Name string;Level,Target int};if err:=sr.Scan(&s.Name,&s.Level,&s.Target);err!=nil{sr.Close();return nil,err};skills=append(skills,s)}
 sr.Close()
 var out []Match
 for rows.Next(){
  var id,title,company,desc,loc string; if err:=rows.Scan(&id,&title,&company,&desc,&loc);err!=nil{return nil,err}
  hay:=strings.ToLower(title+" "+company+" "+desc+" "+loc); matched:=[]string{}; gaps:=[]string{}; total:=0.0
  for _,s:=range skills { terms:=skillTerms(s.Name); hit:=false;for _,t:=range terms{if strings.Contains(hay,t){hit=true;break}}
   if hit {matched=append(matched,s.Name); total+=1+float64(s.Level)/10} else if s.Target>s.Level {gaps=append(gaps,s.Name)}
  }
  score:=0.0;if len(skills)>0 {score=100*total/float64(len(skills))}
  rationale:=fmt.Sprintf("Matched %d tracked skills; %d tracked skill gaps.",len(matched),len(gaps))
  out=append(out,Match{JobID:id,Score:score,Matched:matched,Gaps:gaps,Rationale:rationale})
 }
 if err:=rows.Err();err!=nil{return nil,err}; return out,nil
}

func skillTerms(s string) []string {
 s=strings.ToLower(strings.TrimSpace(s)); switch s {
 case "distributed systems": return []string{"distributed systems","distributed system","microservices"}
 case "kubernetes": return []string{"kubernetes","k8s"}
 case "cloud security": return []string{"cloud security","cspm","cnapp","iam"}
 case "ai security": return []string{"ai security","aispm","llm security","model security","agent security"}
 case "ai agent architecture": return []string{"ai agent","agent architecture","mcp","a2a","agentic"}
 case "technical leadership": return []string{"technical leadership","architect","principal","staff engineer","engineering leadership"}
 case "data platforms": return []string{"data platform","databricks","spark","iceberg","clickhouse","snowflake"}
 default: return []string{s}
 }
}

func Save(ctx context.Context, db *sql.DB, m Match) error {
 a,_:=json.Marshal(m.Matched);g,_:=json.Marshal(m.Gaps)
 _,err:=db.ExecContext(ctx,`INSERT INTO opportunity_matches(job_id,score,matched_skills,gaps,rationale,computed_at)
 VALUES($1,$2,$3,$4,$5,now())
 ON CONFLICT(job_id) DO UPDATE SET score=EXCLUDED.score,matched_skills=EXCLUDED.matched_skills,gaps=EXCLUDED.gaps,rationale=EXCLUDED.rationale,computed_at=now()`,m.JobID,m.Score,a,g,m.Rationale)
 return err
}
