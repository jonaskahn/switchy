package rules

import (
	"net/url"
	"regexp"
	"strings"

	"switchy/internal/config"
)

// ForcePickerBrowser sentinel value returned when a rule match wants
// browser selection UI to stay visible instead of resolving to a concrete browser.
const ForcePickerBrowser = "_Switchy"

// Match resolves the first matching browser for the given URL across all rulesets.
// Returns ForcePickerBrowser if an auto-rule match wants to show the picker UI,
// or an empty string if no rule matches.
func Match(rawURL string, rulesets []config.Ruleset) string {
	parsed, _ := url.Parse(rawURL)

	for _, rs := range rulesets {
		if !rs.Enabled {
			continue
		}
		for _, rule := range rs.Rules {
			if matchRule(rule.Pattern, rawURL, parsed) {
				return rule.Browser
			}
		}
	}
	return ""
}

func matchRule(pattern, rawURL string, parsed *url.URL) bool {
	switch {
	case strings.HasPrefix(pattern, "d$"):
		return matchDomain(pattern[2:], parsed)
	case strings.HasPrefix(pattern, "r$"):
		return matchRegex(pattern[2:], rawURL)
	case strings.HasPrefix(pattern, "s$"):
		return strings.Contains(rawURL, pattern[2:])
	default:
		return strings.Contains(rawURL, pattern)
	}
}

func matchDomain(domain string, parsed *url.URL) bool {
	if parsed == nil {
		return false
	}
	host := parsed.Hostname()
	return host == domain || strings.HasSuffix(host, "."+domain)
}

func matchRegex(pattern, rawURL string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(rawURL)
}
