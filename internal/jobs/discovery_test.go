package jobs

import "testing"

func TestQueryFingerprintIgnoresTextFormatting(t *testing.T) {
	a := DiscoveryQuery{Keywords: []string{"Go   Kubernetes"}, Companies: []string{"Zscaler"}}
	b := DiscoveryQuery{Keywords: []string{" go kubernetes "}, Companies: []string{"zscaler"}}
	if QueryFingerprint(a) != QueryFingerprint(b) { t.Fatal("equivalent discovery queries should have the same fingerprint") }
}
