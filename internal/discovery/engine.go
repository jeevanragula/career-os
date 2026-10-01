package discovery

import("context";"fmt";"net/url";"strings";"time")

type Candidate struct{Name string;URL string;Evidence string;Category string;Signal string;Weight int;DiscoveredAt time.Time}
type Engine struct{Search SearchProvider}

func(e Engine)Discover(ctx context.Context,roles,domains,locations []string,year int)([]Candidate,error){
 if len(roles)==0{roles=RoleTerms[:8]};if len(domains)==0{domains=DomainTerms};if year==0{year=time.Now().Year()}
 var out []Candidate;seen:=map[string]bool{}
 for si,s:=range DefaultSignals{
  for qi,tpl:=range s.Queries{
   domain:=domains[qi%len(domains)];role:=roles[(si+qi)%len(roles)]
   q:=strings.NewReplacer("{{role}}",role,"{{domain}}",domain,"{{location}}",first(locations),"{{year}}",fmt.Sprint(year)).Replace(tpl)
   results,err:=e.Search.Search(ctx,q,8);if err!=nil{continue}
   for _,r:=range results{
    u,err:=url.Parse(r.URL);if err!=nil||u.Host==""{continue}
    key:=strings.ToLower(u.Host+"|"+r.Title);if seen[key]{continue};seen[key]=true
    out=append(out,Candidate{Name:r.Title,URL:r.URL,Evidence:r.Description,Category:s.Category,Signal:s.Name,Weight:s.Weight,DiscoveredAt:time.Now().UTC()})
   }
  }
 }
 return out,nil
}
func first(v []string)string{if len(v)>0{return v[0]};return "India"}
