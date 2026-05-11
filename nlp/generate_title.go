package nlp

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var stopWords = map[string]bool{
	"the":  true,
	"is":   true,
	"are":  true,
	"a":    true,
	"an":   true,
	"how":  true,
	"what": true,
	"why":  true,
	"can":  true,
	"i":    true,
	"to":   true,
	"of":   true,
	"for":  true,
	"and":  true,
}

func GenerateTitle(text string) string {
	text = strings.TrimSpace(text)

	if text == "" {
		return "New Chat"
	}

	// Remove markdown/code noise
	text = regexp.MustCompile("`+").ReplaceAllString(text, "")
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	// Split into words
	words := strings.Fields(strings.ToLower(text))

	var keywords []string

	for _, w := range words {
		w = strings.Trim(w, ".,!?()[]{}:\"'")

		if len(w) < 3 {
			continue
		}

		if stopWords[w] {
			continue
		}

		keywords = append(keywords, w)

		if len(keywords) >= 6 {
			break
		}
	}

	if len(keywords) == 0 {
		return truncate(text, 40)
	}

	title := strings.Join(keywords, " ")

	return capitalize(truncate(title, 50))
}

func truncate(s string, max int) string {
	if utf8.RuneCountInString(s) <= max {
		return s
	}

	r := []rune(s)
	return string(r[:max]) + "..."
}

func capitalize(s string) string {
	if s == "" {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}
