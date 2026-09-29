"""Deterministic validation for CareerOS seed Career Brain files.

This validates structural and safety invariants only. It does not decide whether
a professional claim is true.
"""

from pathlib import Path
import re
import sys

ROOT = Path(__file__).resolve().parents[1]
claims_text = (ROOT / "career-profile" / "claims.yaml").read_text()
evidence_text = (ROOT / "career-profile" / "evidence-index.yaml").read_text()

VALID_STATUS = {"needs_verification", "user_asserted", "verified", "rejected"}
VALID_CONF = {"public", "application_safe", "internal", "confidential"}


def claim_blocks(text: str):
    return re.findall(
        r"(?ms)^  - id: .*?(?=^  - id: |\Z)",
        text,
    )


def main() -> int:
    errors = []
    blocks = claim_blocks(claims_text)

    if not blocks:
        errors.append("no Career Brain claims found")

    for block in blocks:
        status_match = re.search(r"^    status:\s*([^\n]+)$", block, re.MULTILINE)
        confidentiality_match = re.search(
            r"^    confidentiality:\s*([^\n]+)$", block, re.MULTILINE
        )

        status = status_match.group(1).strip().strip('"') if status_match else ""
        confidentiality = (
            confidentiality_match.group(1).strip().strip('"')
            if confidentiality_match
            else ""
        )

        if status not in VALID_STATUS:
            errors.append(f"invalid claim status: {status or '<missing>'}")

        if confidentiality not in VALID_CONF:
            errors.append(
                f"invalid claim confidentiality: {confidentiality or '<missing>'}"
            )

        if status == "verified" and "evidence_refs:" not in block:
            errors.append("verified claim is missing evidence_refs")

        if status == "verified" and confidentiality == "confidential":
            errors.append("confidential claim cannot be verified for application use")

    if "verification: verified" in evidence_text and "ev.user.seed" in evidence_text:
        errors.append("user seed evidence must not be marked verified")

    if errors:
        for error in errors:
            print(f"ERROR: {error}")
        return 1

    print(f"Career Brain profile validation passed ({len(blocks)} claims).")
    return 0


if __name__ == "__main__":
    sys.exit(main())
