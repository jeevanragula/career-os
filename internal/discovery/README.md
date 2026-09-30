# Autonomous opportunity discovery

CareerOS does not require the user to maintain a company or job-source list.

The discovery engine starts from the user's career profile and searches multiple signal classes:

- major job portals
- startup ecosystems
- analyst and vendor ecosystems such as Gartner Peer Insights
- CNCF and cloud-native community signals
- public conference sponsorships
- cloud-event sponsorships
- security conferences and competitions
- engineering communities and engineering blogs

Signals produce candidate companies and evidence. The next verification stage resolves the company's canonical domain and career page, verifies the page, extracts public job postings, deduplicates them, and passes only verified jobs to the existing job-analysis pipeline.

The first search backend is Brave Search API. The user supplies only one search credential, BRAVE_SEARCH_API_KEY, not individual job sources or companies.

Discovery signals are evidence only. Sponsorship, Gartner presence, conference participation, or community membership does not imply that a company is a good employer.