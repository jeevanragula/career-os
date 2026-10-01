package discovery

import("context";"database/sql";"encoding/json";"net/url";"strings";"time")

type Store struct{DB *sql.DB}

type Run struct {
 ID string `json:"id"`
 StartedAt time.Time `json:"started_at"`
 FinishedAt *time.Time `json:"finished_at,omitempty"`
 Status string `json:"status"`
 Query any `json:"query"`
 DiscoveredCount int `json:"discovered_count"`
 Error string `json:"error,omitempty"`
}

func(s Store)CreateRun(ctx context.Context, query any)(string,error){
 raw,e:=json.Marshal(query);if e!=nil{return "",e}
 var id string
 e=s.DB.QueryRowContext(ctx,"INSERT INTO discovery_runs(query,status,started_at) VALUES($1,'running',now()) RETURNING id",string(raw)).Scan(&id)
 return id,e
}

func(s Store)UpdateRun(ctx context.Context,id,status string,discovered int,runErr string) error {
 _,e:=s.DB.ExecContext(ctx,"UPDATE discovery_runs SET status=$2,discovered_count=$3,error=$4,finished_at=CASE WHEN $2<>'running' THEN now() ELSE NULL END WHERE id=$1",id,status,discovered,runErr)
 return e
}

func(s Store)GetRun(ctx context.Context,id string)(Run,error){
 var r Run; var raw []byte
 e:=s.DB.QueryRowContext(ctx,"SELECT id,started_at,finished_at,status,query,discovered_count,error FROM discovery_runs WHERE id=$1",id).Scan(&r.ID,&r.StartedAt,&r.FinishedAt,&r.Status,&raw,&r.DiscoveredCount,&r.Error)
 if e!=nil{return r,e}
 _=json.Unmarshal(raw,&r.Query)
 return r,nil
}

func(s Store)SaveCandidates(ctx context.Context,candidates []Candidate)error{
 for _,c:=range candidates{
  u,e:=url.Parse(c.URL);if e!=nil||u.Host==""{continue}
  domain:=strings.ToLower(u.Hostname());if strings.HasPrefix(domain,"www."){domain=strings.TrimPrefix(domain,"www.")}
  var id string
  e=s.DB.QueryRowContext(ctx,"INSERT INTO discovered_companies(name,canonical_domain,homepage_url,last_seen_at) VALUES($1,$2,$3,now()) ON CONFLICT(canonical_domain) DO UPDATE SET name=EXCLUDED.name,last_seen_at=now() RETURNING id",c.Name,domain,c.URL).Scan(&id)
  if e!=nil{return e}
  raw,_:=json.Marshal(map[string]any{"title":c.Name,"evidence":c.Evidence})
  if _,e=s.DB.ExecContext(ctx,"INSERT INTO discovery_evidence(company_id,signal,category,source_url,evidence,weight) VALUES($1,$2,$3,$4,$5,$6)",id,c.Signal,c.Category,c.URL,string(raw),c.Weight);e!=nil{return e}
 }
 return nil
}

func(s Store)ListCandidates(ctx context.Context,limit int)([]Candidate,error){if limit<=0{limit=50};rows,e:=s.DB.QueryContext(ctx,"SELECT name,homepage_url,canonical_domain FROM discovered_companies WHERE status<>'rejected' ORDER BY last_seen_at DESC LIMIT $1",limit);if e!=nil{return nil,e};defer rows.Close();var out []Candidate;for rows.Next(){var x Candidate;var domain string;if e:=rows.Scan(&x.Name,&x.URL,&domain);e!=nil{return nil,e};x.Evidence=domain;out=append(out,x)};return out,rows.Err()}
}
