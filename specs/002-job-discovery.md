# Specification 002 — Job Discovery

- Status: Draft
- Version: 0.1
- Milestone: M2 Job Discovery

## Purpose

Discover relevant engineering opportunities from legitimate sources and convert them into normalized opportunity records.

## Source categories

- company career sites
- professional job platforms
- startup/company directories
- public hiring pages
- approved search/research sources

The system must respect applicable terms, access restrictions, authentication boundaries, and rate limits.

## Normalized job record

A job should support a stable internal ID, source, source URL, company, title, location, work mode, employment type, description, responsibilities, requirements, preferred qualifications, compensation when available, posting date, discovery time, and last-observed time.

## Deduplication

Consider canonical source URL, source job ID, company, normalized title, location, and description similarity.

The same opportunity on multiple sources should become one logical opportunity with multiple observations.

## Acceptance criteria

- [ ] Source observations retain source metadata.
- [ ] Multiple observations can map to one opportunity.
- [ ] Duplicate ingestion is idempotent.
- [ ] Original job text is preserved where legally and technically appropriate.
- [ ] Discovery respects source access rules and configured rate limits.
