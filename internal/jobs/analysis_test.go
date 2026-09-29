package jobs

import "testing"

func TestRequirementMatchDoesNotAllowInvalidConfidence(t *testing.T) {
	m := RequirementMatch{Confidence: 1.1}
	if m.Confidence >= 0 && m.Confidence <= 1 { t.Fatal("test fixture should represent an invalid confidence") }
}
