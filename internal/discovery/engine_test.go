package discovery

import("context";"testing")

type fakeSearch struct{queries []string}
func(f *fakeSearch)Search(ctx context.Context,q string,n int)([]SearchResult,error){f.queries=append(f.queries,q);return []SearchResult{{Title:q,URL:"https://example.com/careers",Description:"evidence"}},nil}

func TestDiscoverRotatesTargets(t *testing.T){
 f:=&fakeSearch{};e:=Engine{Search:f}
 _,err:=e.Discover(context.Background(),[]string{"principal engineer","AI security"},[]string{"kubernetes","DSPM"},[]string{"India"},2026)
 if err!=nil{t.Fatal(err)}
 if len(f.queries)<4{t.Fatalf("expected multiple queries, got %d",len(f.queries))}
 seenPrincipal,seenAI:=false,false
 for _,q:=range f.queries{if contains(q,"principal engineer"){seenPrincipal=true};if contains(q,"AI security"){seenAI=true}}
 if !seenPrincipal||!seenAI{t.Fatalf("expected both roles to be represented")}
}
func contains(s,sub string)bool{for i:=0;i+len(sub)<=len(s);i++{if s[i:i+len(sub)]==sub{return true}};return false}
