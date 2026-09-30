package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/jeevanragula/career-os/internal/jobs"
)

type WebFetcher struct { Client *http.Client; MaxBytes int64 }
type Page struct { URL, Body string }

func (f WebFetcher) Fetch(ctx context.Context, rawURL string) (Page, error) {
	client:=f.Client; if client==nil { client=&http.Client{Timeout:15*time.Second} }
	max:=f.MaxBytes; if max<=0 { max=2<<20 }
	req,err:=http.NewRequestWithContext(ctx,http.MethodGet,rawURL,nil); if err!=nil{return Page{},err}
	req.Header.Set("User-Agent","CareerOS/1.0 (+https://github.com/jeevanragula/career-os)")
	req.Header.Set("Accept","text/html,application/xhtml+xml,application/json")
	resp,err:=client.Do(req); if err!=nil{return Page{},err}; defer resp.Body.Close()
	if resp.StatusCode<200||resp.StatusCode>=400{return Page{},fmt.Errorf("fetch %s returned HTTP %d",rawURL,resp.StatusCode)}
	body,err:=io.ReadAll(io.LimitReader(resp.Body,max)); if err!=nil{return Page{},err}
	finalURL:=rawURL; if resp.Request!=nil&&resp.Request.URL!=nil{finalURL=resp.Request.URL.String()}
	return Page{URL:finalURL,Body:string(body)},nil
}

var hrefRE=regexp.MustCompile("(?is)<a[^>]+href=[\\\"']([^\\\"']+)[\\\"'][^>]*>(.*?)</a>")
var tagRE=regexp.MustCompile("(?is)<[^>]+>")
var spaceRE=regexp.MustCompile("\\s+")

type Link struct{URL,Label string}

func ExtractLinks(baseURL,body string)[]Link{
	base,err:=url.Parse(baseURL);if err!=nil{return nil};seen:=map[string]bool{};var out []Link
	for _,m:=range hrefRE.FindAllStringSubmatch(body,-1){
		u,err:=base.Parse(strings.TrimSpace(html.UnescapeString(m[1])));if err!=nil||(u.Scheme!="http"&&u.Scheme!="https"){continue}
		u.Fragment="";canonical:=jobs.CanonicalizeURL(u.String());if canonical==""||seen[canonical]{continue};seen[canonical]=true
		label:=strings.ToLower(strings.TrimSpace(tagRE.ReplaceAllString(html.UnescapeString(m[2])," ")))
		out=append(out,Link{URL:canonical,Label:label})
	}
	return out
}

func LooksLikeCareerURL(raw string)bool{
	u,err:=url.Parse(raw);if err!=nil{return false};s:=strings.ToLower(u.Path+"?"+u.RawQuery)
	for _,token:=range []string{"career","careers","jobs","job-openings","work-with-us","join-us","opportunities","vacancies"}{if strings.Contains(s,token){return true}}
	return false
}

func trustedJobHost(raw string)bool{
	s:=strings.ToLower(raw);for _,h:=range []string{"greenhouse.io","lever.co","ashbyhq.com","myworkdayjobs.com","smartrecruiters.com"}{if strings.Contains(s,h){return true}};return false
}

type JobPosting struct{Title,Description,URL,EmploymentType,Location,RemoteMode,Company,SourcePage string}

func ExtractJobPostings(page Page)[]JobPosting{
	re:=regexp.MustCompile("(?is)<script[^>]+type=[\\\"']application/ld\\+json[\\\"'][^>]*>(.*?)</script>")
	var out []JobPosting
	for _,m:=range re.FindAllStringSubmatch(page.Body,-1){var v any;if json.Unmarshal([]byte(html.UnescapeString(strings.TrimSpace(m[1]))),&v)!=nil{continue};walkJSONLD(v,page.URL,&out)}
	return dedupePostings(out)
}

func walkJSONLD(v any,pageURL string,out *[]JobPosting){
	switch x:=v.(type){
	case []any:for _,item:=range x{walkJSONLD(item,pageURL,out)}
	case map[string]any:
		if graph,ok:=x["@graph"];ok{walkJSONLD(graph,pageURL,out)}
		typ:=strings.ToLower(fmt.Sprint(x["@type"]));if !strings.Contains(typ,"jobposting"){return}
		p:=JobPosting{Title:stringField(x,"title"),Description:cleanHTML(stringField(x,"description")),URL:firstNonEmpty(stringField(x,"url"),pageURL),EmploymentType:stringField(x,"employmentType"),Company:nestedName(x["hiringOrganization"]),SourcePage:pageURL}
		if p.Company==""{p.Company=nestedName(x["organization"])};p.Location,p.RemoteMode=jobLocation(x["jobLocation"],x["jobLocationType"]);*out=append(*out,p)
	}
}

func stringField(m map[string]any,key string)string{if v,ok:=m[key];ok{return strings.TrimSpace(fmt.Sprint(v))};return ""}
func nestedName(v any)string{switch x:=v.(type){case map[string]any:return stringField(x,"name");case []any:for _,item:=range x{if n:=nestedName(item);n!=""{return n}}};return ""}
func jobLocation(v,remoteType any)(string,string){
	remote:=strings.TrimSpace(fmt.Sprint(remoteType));if remote=="<nil>"{remote=""};var locations []string
	add:=func(v any){a,ok:=v.(map[string]any);if !ok{return};address,ok:=a["address"].(map[string]any);if !ok{return};for _,k:=range []string{"addressLocality","addressRegion","addressCountry"}{if s:=stringField(address,k);s!=""{locations=append(locations,s)}}}
	switch x:=v.(type){case map[string]any:add(x);case []any:for _,item:=range x{add(item)}}
	return strings.Join(locations,", "),remote
}
func cleanHTML(s string)string{return strings.TrimSpace(spaceRE.ReplaceAllString(html.UnescapeString(tagRE.ReplaceAllString(s," "))," "))}
func firstNonEmpty(v ...string)string{for _,s:=range v{if strings.TrimSpace(s)!=""{return strings.TrimSpace(s)}};return ""}
func dedupePostings(in []JobPosting)[]JobPosting{seen:=map[string]bool{};out:=make([]JobPosting,0,len(in));for _,p:=range in{key:=jobs.CanonicalizeURL(p.URL);if key==""{key=strings.ToLower(p.Company+"|"+p.Title+"|"+p.Location)};if p.Title==""||seen[key]{continue};seen[key]=true;out=append(out,p)};return out}

func(e Engine)ResolveCareerPages(ctx context.Context,candidates []Candidate)[]Candidate{
	fetcher:=WebFetcher{};var out []Candidate
	for _,c:=range candidates{page,err:=fetcher.Fetch(ctx,c.URL);if err!=nil{continue};if len(ExtractJobPostings(page))>0||LooksLikeCareerURL(page.URL){c.URL=page.URL;out=append(out,c);continue};for _,link:=range ExtractLinks(page.URL,page.Body){if !LooksLikeCareerURL(link.URL)&&!trustedJobHost(link.URL){continue};c.URL=link.URL;out=append(out,c);break}}
	return out
}

type HarvestResult struct{Observations []jobs.JobObservation;Pages int}

func(e Engine)HarvestJobs(ctx context.Context,candidates []Candidate,maxCandidates int)HarvestResult{
	if maxCandidates<=0||maxCandidates>100{maxCandidates=50};fetcher:=WebFetcher{};var result HarvestResult
	for i,c:=range candidates{if i>=maxCandidates{break};page,err:=fetcher.Fetch(ctx,c.URL);if err!=nil{continue};result.Pages++;postings:=ExtractJobPostings(page)
		if len(postings)==0{checked:=0;for _,link:=range ExtractLinks(page.URL,page.Body){if checked>=30{break};if !LooksLikeCareerURL(link.URL)&&!trustedJobHost(link.URL){continue};child,err:=fetcher.Fetch(ctx,link.URL);if err!=nil{continue};checked++;postings=append(postings,ExtractJobPostings(child)...)};postings=dedupePostings(postings)}
		now:=time.Now().UTC();for _,p:=range postings{company:=p.Company;if company==""{company=c.Name};result.Observations=append(result.Observations,jobs.JobObservation{Source:"autonomous-web",SourceJobID:jobs.ObservationIdentity(jobs.NormalizeInput{URL:p.URL,Title:p.Title,Company:company,Location:p.Location}),SourceURL:p.URL,ObservedAt:now,RawTitle:p.Title,RawCompany:company,RawLocation:p.Location,RawDescription:p.Description})}
	}
	return result
}
