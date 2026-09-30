package discovery

import("context";"database/sql";"encoding/json";"net/url";"strings")

type Store struct{DB *sql.DB}

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
