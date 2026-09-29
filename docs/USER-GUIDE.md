# CareerOS User Guide

## Start

Run:

    ./scripts/bootstrap.sh

Open:

    http://localhost:8080

## 1. Discover a company

Select Lever or Ashby.

### Lever

The documented public endpoint uses the postings base:

    https://api.lever.co/v0/postings/

Enter the company's Lever site name in the Name field.

### Ashby

Use:

    https://api.ashbyhq.com/posting-api/job-board/{JOB_BOARD_NAME}

The final path component is the organization's public job-board name.

## 2. Filter

Enter keywords such as:

    principal,kubernetes,cloud security

CareerOS keeps only postings matching at least one supplied keyword.

## 3. Analyze

Click Analyze on a job.

The AI receives:
- the job description
- application-safe Career Brain claims

It produces:
- job requirements
- evidence matches
- unmatched requirements
- confidence
- source quotes

## 4. Generate resume

Click Resume.

The AI generates a Markdown resume tailored to the job using only application-safe Career Brain claims.

The visible output can be copied into a document or converted to PDF.

## 5. Prepare application

Click Prepare application for last resume.

This creates an application record in the preparing state.

It does not submit anything.

## What is intentionally not automatic

CareerOS does not currently:
- submit applications automatically
- send recruiter messages automatically
- bypass login/CAPTCHA
- store job-board passwords
- fabricate missing career facts

Those are deliberate boundaries, not broken features.
