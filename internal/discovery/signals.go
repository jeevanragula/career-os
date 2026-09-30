package discovery

type Signal struct {
 Name string
 Category string
 Queries []string
 Weight int
}

var DefaultSignals=[]Signal{
 {Name:"technology-job-portals",Category:"job_portal",Weight:5,Queries:[]string{
  "site:linkedin.com/jobs/view {{role}} {{location}}","site:indeed.com/viewjob {{role}} {{location}}","site:wellfound.com/jobs {{role}} {{location}}","site:instahyre.com/jobs {{role}} {{location}}","site:foundit.in/job {{role}} {{location}}",
 }},
 {Name:"startup-ecosystem",Category:"startup",Weight:4,Queries:[]string{"{{domain}} startup funding {{year}} hiring","site:wellfound.com/company {{domain}}","{{domain}} startup cloud security AI hiring"}},
 {Name:"gartner-vendors",Category:"analyst_vendor",Weight:5,Queries:[]string{"site:gartner.com/reviews/market {{domain}} cybersecurity","site:gartner.com/reviews {{domain}} cloud security","{{domain}} Gartner Peer Insights careers"}},
 {Name:"cloud-native-community",Category:"community",Weight:5,Queries:[]string{"site:cncf.io {{domain}} sponsor","site:cncf.io {{domain}} ambassador","site:cncf.io {{domain}} KubeCon sponsor"}},
 {Name:"conference-sponsors",Category:"conference",Weight:4,Queries:[]string{"{{domain}} conference sponsor cloud native","{{domain}} conference sponsor security","{{domain}} conference sponsor AI infrastructure"}},
 {Name:"security-events",Category:"security_event",Weight:5,Queries:[]string{"site:blackhat.com {{domain}} sponsors","site:rsaconference.com {{domain}} sponsors","{{domain}} security competition sponsor","{{domain}} CTF sponsor security"}},
 {Name:"cloud-events",Category:"cloud_event",Weight:4,Queries:[]string{"{{domain}} AWS conference sponsor","{{domain}} Google Cloud event sponsor","{{domain}} Microsoft Ignite sponsor","{{domain}} cloud native event sponsor"}},
 {Name:"engineering-community",Category:"engineering",Weight:4,Queries:[]string{"{{domain}} engineering blog Kubernetes","{{domain}} engineering blog distributed systems","{{domain}} engineering blog AI security"}},
}

var RoleTerms=[]string{"principal engineer","staff engineer","senior staff engineer","distinguished engineer","principal architect","technical architect","software architect","security architect","AI security","AI infrastructure","cloud security","platform engineering","distributed systems","data platform","engineering manager"}
var DomainTerms=[]string{"kubernetes","cloud native","cloud security","AI security","cybersecurity","CNAPP","DSPM","CSPM","data security","AI agents","MCP","platform engineering","distributed systems","developer infrastructure","observability"}
