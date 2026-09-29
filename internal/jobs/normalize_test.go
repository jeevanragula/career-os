package jobs

import "testing"

func TestCanonicalizeURLRemovesTrackingParameters(t *testing.T) {
	got := CanonicalizeURL("https://example.com/jobs/123/?utm_source=linkedin&utm_medium=social&foo=bar#section")
	want := "https://example.com/jobs/123?foo=bar"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestObservationIdentityUsesCanonicalURL(t *testing.T) {
	a := NormalizeInput{URL: "https://jobs.example.com/1?utm_source=a"}
	b := NormalizeInput{URL: "https://jobs.example.com/1?utm_source=b"}
	if ObservationIdentity(a) != ObservationIdentity(b) {
		t.Fatal("tracking parameters should not change identity")
	}
}

func TestObservationIdentitySeparatesCompanies(t *testing.T) {
	a := NormalizeInput{Company: "Acme", Title: "Principal Engineer", Location: "Hyderabad"}
	b := NormalizeInput{Company: "Beta", Title: "Principal Engineer", Location: "Hyderabad"}
	if ObservationIdentity(a) == ObservationIdentity(b) {
		t.Fatal("different companies must remain separate")
	}
}

func TestContentHashIsStableForEquivalentText(t *testing.T) {
	a := NormalizeInput{Title: "Principal Engineer", Company: "Acme", Description: "Build systems."}
	b := NormalizeInput{Title: " principal   engineer ", Company: "ACME", Description: "Build systems."}
	if ContentHash(a) != ContentHash(b) {
		t.Fatal("equivalent normalized content should hash identically")
	}
}
