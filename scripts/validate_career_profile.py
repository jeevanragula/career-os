"""Deterministic validation for CareerOS seed Career Brain files.

This intentionally validates structure and safety rules only.
It does not decide whether a professional claim is true.
"""

from pathlib import Path
import re
import sys

ROOT = Path(__file__).resolve().parents[1]
claims_text = (ROOT / "career-profile" / "claims.yaml").read_text()
evidence_text = (ROOT / "career-profile" / "evidence-index.yaml").read_text()

VALID_STATUS = {"needs_verification", "user_asserted", "verified", "rejected"}
VALID_CONF = {"public", "application_safe", "internal", "confidential"}


def values_after_key(text: str, key: str):
    return re.findall(rf"^\s+{re.escape(key)}:\s*([^\n]+)$", text, re.MULTILINE)


def main() -> int:
    errors = []

    statuses = [v.strip().strip('"') for v in values_after_key(claims_text, "status")]
    confidences = [v.strip().strip('"') for v in values_after_key(claims_text, "confidentiality")]

    for value in statuses:
        if value not in VALID_STATUS:
            errors.append(f"invalid claim status: {value}")

    for value in confidences:
        if value not in VALID_CONF:
            errors.append(f"invalid claim confidentiality: {value}")

    # Seed profile must never silently contain a verified claim without an evidence mapping.
    if "status: verified" in claims_text:
        errors.append("seed claims must not contain verified claims before evidence review")

    if "verification: verified" in evidence_text and "ev.user.seed" in evidence_text:
        errors.append("user seed evidence must not be marked verified")

    if errors:
        for error in errors:
            print(f"ERROR: {error}")
        return 1

    print("Career Brain profile validation passed.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
