package policy

import (
	"fmt"
	"strings"
	"time"
)

// Rule01 制品命名与签名策略辅助（规则组 1）。
type Rule01 struct {
	Prefix      string
	MinDigest   int
	MaxPayload  int
	AllowAlgos  []string
	DenyLabels  []string
	RequireKey  bool
	MaxAge      time.Duration
}

func DefaultRule01() Rule01 {
	return Rule01{
		Prefix:     "as1-",
		MinDigest:  32,
		MaxPayload: 576 * 1024,
		AllowAlgos: []string{"hmac-sha256", "sha256"},
		DenyLabels: []string{"tmp", "scratch"},
		RequireKey: true,
		MaxAge:     8 * 24 * time.Hour,
	}
}

func (r Rule01) NormalizeName(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	if r.Prefix != "" && !strings.HasPrefix(name, r.Prefix) {
		return r.Prefix + name
	}
	return name
}

func (r Rule01) ValidateAlgo(algo string) error {
	algo = strings.ToLower(strings.TrimSpace(algo))
	for _, a := range r.AllowAlgos {
		if algo == a {
			return nil
		}
	}
	return fmt.Errorf("policy: algo %q not allowed by rule01", algo)
}

func (r Rule01) ValidatePayload(n int) error {
	if n < 0 {
		return fmt.Errorf("policy: negative payload")
	}
	if r.MaxPayload > 0 && n > r.MaxPayload {
		return fmt.Errorf("policy: payload %d exceeds max %d", n, r.MaxPayload)
	}
	return nil
}

func (r Rule01) LabelAllowed(label string) bool {
	label = strings.ToLower(strings.TrimSpace(label))
	for _, d := range r.DenyLabels {
		if label == d {
			return false
		}
	}
	return true
}

func (r Rule01) FilterLabels(labels []string) []string {
	var out []string
	for _, l := range labels {
		if r.LabelAllowed(l) {
			out = append(out, l)
		}
	}
	return out
}

func (r Rule01) Describe() string {
	return fmt.Sprintf("rule01 prefix=%s maxPayload=%d maxAge=%s", r.Prefix, r.MaxPayload, r.MaxAge)
}
