package jobs

import (
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

var whitespace = regexp.MustCompile("\\s+")

var trackingParams = map[string]bool{
	"utm_source": true,
	"utm_medium": true,
	"utm_campaign": true,
	"utm_content": true,
	"utm_term": true,
	"gh_src": true,
	"source": true,
}

func NormalizeText(value string) string {
	return strings.ToLower(strings.TrimSpace(whitespace.ReplaceAllString(value, " ")))
}

func CanonicalizeURL(raw string) string {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Scheme == "" || u.Host == "" {
		return strings.TrimSpace(raw)
	}

	q := u.Query()
	for key := range q {
		if trackingParams[strings.ToLower(key)] {
			q.Del(key)
		}
	}
	u.RawQuery = q.Encode()
	u.Fragment = ""
	u.Host = strings.ToLower(u.Host)
	u.Scheme = strings.ToLower(u.Scheme)
	return strings.TrimRight(u.String(), "/")
}

func ObservationIdentity(in NormalizeInput) string {
	canonicalURL := CanonicalizeURL(in.URL)
	if canonicalURL != "" {
		return "url:" + canonicalURL
	}

	if in.SourceJobID != "" {
		return "source-id:" + NormalizeText(in.Source) + ":" + NormalizeText(in.SourceJobID)
	}

	parts := []string{
		NormalizeText(in.Company),
		NormalizeText(in.Title),
		NormalizeText(in.Location),
		NormalizeText(in.RemoteMode),
	}
	return "derived:" + strings.Join(parts, "|")
}

func ContentHash(in NormalizeInput) string {
	parts := []string{
		NormalizeText(in.Title),
		NormalizeText(in.Company),
		NormalizeText(in.Location),
		NormalizeText(in.RemoteMode),
		NormalizeText(in.Employment),
		NormalizeText(in.Description),
	}
	sort.Strings(parts)
	h := sha256.Sum256([]byte(strings.Join(parts, "\n")))
	return hex.EncodeToString(h[:])
}
