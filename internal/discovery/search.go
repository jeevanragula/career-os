package discovery

import("context";"encoding/json";"fmt";"net/http";"net/url";"os")

type SearchResult struct {Title string; URL string; Description string; Source string}
type SearchProvider interface { Search(context.Context,string,int)([]SearchResult,error) }
type Brave struct {APIKey string; Client *http.Client}

func NewBrave()*Brave{return &Brave{APIKey:os.Getenv("BRAVE_SEARCH_API_KEY"),Client:http.DefaultClient}}

func(b *Brave)Search(ctx context.Context,q string,count int)([]SearchResult,error){
 if b.APIKey==""{return nil,fmt.Errorf("BRAVE_SEARCH_API_KEY is not configured")}
 if count<=0{count=10}
 u:="https://api.search.brave.com/res/v1/web/search?"+url.Values{"q":[]string{q},"count":[]string{fmt.Sprint(count)}}.Encode()
 req,e:=http.NewRequestWithContext(ctx,http.MethodGet,u,nil);if e!=nil{return nil,e}
 req.Header.Set("Accept","application/json");req.Header.Set("X-Subscription-Token",b.APIKey)
 resp,e:=b.Client.Do(req);if e!=nil{return nil,e};defer resp.Body.Close()
 if resp.StatusCode<200||resp.StatusCode>=300{return nil,fmt.Errorf("brave search returned HTTP %d",resp.StatusCode)}
 var p struct{Web struct{Results []struct{Title string;URL string;Description string}}}
 if e=json.NewDecoder(resp.Body).Decode(&p);e!=nil{return nil,e}
 out:=make([]SearchResult,0,len(p.Web.Results));for _,x:=range p.Web.Results{out=append(out,SearchResult{Title:x.Title,URL:x.URL,Description:x.Description,Source:"brave"})};return out,nil
}
