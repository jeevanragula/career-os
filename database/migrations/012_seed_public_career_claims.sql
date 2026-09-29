-- Public/application-safe Career Brain seed.
-- These are intentionally limited to facts already represented in public evidence.
INSERT INTO evidence(source_type,source_reference,confidentiality,integrity_status,verification_notes)
SELECT 'portfolio','https://jeevanragula.github.io/','application_safe','reviewed','Public portfolio reviewed for CareerOS seed'
WHERE NOT EXISTS (SELECT 1 FROM evidence WHERE source_reference='https://jeevanragula.github.io/');

INSERT INTO evidence(source_type,source_reference,confidentiality,integrity_status,verification_notes)
SELECT 'patent','US12598208B2','public','reviewed','IaC Scanner for Infrastructure Component Security; inventor Jeevan Reddy Ragula'
WHERE NOT EXISTS (SELECT 1 FROM evidence WHERE source_reference='US12598208B2');

INSERT INTO evidence(source_type,source_reference,confidentiality,integrity_status,verification_notes)
SELECT 'patent','US20250202926A1','public','reviewed','Cloud-Based Data Security Posture Management; inventor Jeevan Reddy Ragula'
WHERE NOT EXISTS (SELECT 1 FROM evidence WHERE source_reference='US20250202926A1');

INSERT INTO claims(claim_key,text,entity_type,verification_status,confidence,confidentiality)
VALUES
('role.principal.architect','Principal Software Engineer / Architect','employment','verified',0.98,'application_safe'),
('employment.zscaler.principal','Principal Software Engineer at Zscaler India Pvt Ltd from April 2022 to present','employment','verified',0.98,'application_safe'),
('employment.zscaler.senior_staff','Senior Staff Software Engineer at Zscaler India Pvt Ltd from August 2021 to March 2022','employment','verified',0.98,'application_safe'),
('employment.mcafee.senior_staff','Senior Staff Engineer at McAfee India Pvt Ltd from August 2019 to August 2021','employment','verified',0.98,'application_safe'),
('employment.ivy.lead','Lead - Software Development at IVY Software Development Services Pvt Ltd from July 2010 to July 2019','employment','verified',0.98,'application_safe'),
('project.zpc','Zscaler Posture Control — Architect, Lead and Developer','project','verified',0.95,'application_safe'),
('project.skyhigh_casb','Skyhigh CASB — Lead and Developer','project','verified',0.95,'application_safe'),
('project.payment_gateway','PCI DSS Payment Gateway Platform for Entain — Lead and Developer','project','verified',0.95,'application_safe'),
('skill.golang','Golang','skill','verified',0.95,'application_safe'),
('skill.kubernetes','Kubernetes','skill','verified',0.95,'application_safe'),
('skill.aws','AWS','skill','verified',0.95,'application_safe'),
('skill.cloud_security','Cloud Security','skill','verified',0.95,'application_safe'),
('skill.kafka','Kafka','skill','verified',0.95,'application_safe'),
('patent.iac','Inventor on US Patent 12598208 B2, Infrastructure as Code (IaC) Scanner for Infrastructure Component Security','patent','verified',0.99,'public'),
('patent.dspm','Inventor on US Patent Application 20250202926 A1, Cloud-Based Data Security Posture Management (DSPM)','patent','verified',0.99,'public')
ON CONFLICT (claim_key) DO NOTHING;

INSERT INTO claim_evidence(claim_id,evidence_id)
SELECT c.id,e.id FROM claims c,evidence e WHERE c.claim_key IN ('role.principal.architect','employment.zscaler.principal','employment.zscaler.senior_staff','employment.mcafee.senior_staff','employment.ivy.lead','project.zpc','project.skyhigh_casb','project.payment_gateway','skill.golang','skill.kubernetes','skill.aws','skill.cloud_security','skill.kafka') AND e.source_reference='https://jeevanragula.github.io/'
ON CONFLICT DO NOTHING;

INSERT INTO claim_evidence(claim_id,evidence_id)
SELECT c.id,e.id FROM claims c,evidence e WHERE c.claim_key='patent.iac' AND e.source_reference='US12598208B2'
ON CONFLICT DO NOTHING;

INSERT INTO claim_evidence(claim_id,evidence_id)
SELECT c.id,e.id FROM claims c,evidence e WHERE c.claim_key='patent.dspm' AND e.source_reference='US20250202926A1'
ON CONFLICT DO NOTHING;
