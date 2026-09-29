# Job Source Providers

This directory contains compliant provider adapters implementing `jobs.SourceAdapter`.

## Initial provider targets

### Lever Postings API
Lever documents a public Postings API for published jobs. The public endpoint is based on `https://api.lever.co/v0/postings/{SITE}` and supports JSON output, pagination, and filters. CareerOS should use only published postings from this public interface unless the user explicitly connects an authorized Lever integration.

Documentation: https://github.com/lever/postings-api

### Ashby Public Job Posting API
Ashby documents a public job-posting endpoint at `https://api.ashbyhq.com/posting-api/job-board/{JOB_BOARD_NAME}` for currently published postings. Compensation can optionally be included. CareerOS should treat missing fields as missing rather than infer them.

Documentation: https://developers.ashbyhq.com/docs/public-job-posting-api

### Authenticated provider APIs
Authenticated APIs must be represented with `authorization_mode: user_authorized` or `connector`. Secrets must remain outside CareerOS source files and database records.

Provider-specific implementation should include fixture-based tests and rate-limit handling before enabling production scheduling.
