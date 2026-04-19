package jdparser

import (
	"regexp"
	"strings"
	"unicode"
)

// normalizeText normalizes text for parsing by lowercasing and
// normalizing whitespace.
func normalizeText(text string) string {
	// Replace various dash types with standard dash
	text = strings.ReplaceAll(text, "–", "-")
	text = strings.ReplaceAll(text, "—", "-")

	// Normalize bullet points
	text = strings.ReplaceAll(text, "•", "-")
	text = strings.ReplaceAll(text, "●", "-")
	text = strings.ReplaceAll(text, "○", "-")
	text = strings.ReplaceAll(text, "◦", "-")
	text = strings.ReplaceAll(text, "▪", "-")
	text = strings.ReplaceAll(text, "■", "-")

	// Normalize whitespace within lines
	var lines []string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		// Collapse multiple spaces
		spaceRe := regexp.MustCompile(`\s+`)
		line = spaceRe.ReplaceAllString(line, " ")
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// containsWord checks if text contains a word (case-insensitive, word boundary aware).
func containsWord(text, word string) bool {
	lower := strings.ToLower(text)
	wordLower := strings.ToLower(word)

	// Handle special cases like C++, C# where we can't use word boundaries
	if strings.ContainsAny(word, "+#") {
		return strings.Contains(lower, wordLower)
	}

	// For short words or acronyms, use word boundary matching
	if len(word) <= 3 || isAcronym(word) {
		// Build pattern with word boundaries
		pattern := `(?i)\b` + regexp.QuoteMeta(wordLower) + `\b`
		re, err := regexp.Compile(pattern)
		if err != nil {
			return strings.Contains(lower, wordLower)
		}
		return re.MatchString(lower)
	}

	// For longer words, simple contains is usually fine
	return strings.Contains(lower, wordLower)
}

// isAcronym checks if a word is likely an acronym (all caps or mostly caps).
func isAcronym(word string) bool {
	if len(word) == 0 {
		return false
	}

	upper := 0
	lower := 0
	for _, r := range word {
		if unicode.IsUpper(r) {
			upper++
		} else if unicode.IsLower(r) {
			lower++
		}
	}

	// All uppercase, or mostly uppercase with some numbers
	return upper > 0 && lower == 0
}

// containsAny checks if text contains any of the given substrings.
func containsAny(text string, substrs []string) bool {
	for _, s := range substrs {
		if strings.Contains(text, s) {
			return true
		}
	}
	return false
}

// mapKeys returns the keys of a map as a slice.
func mapKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// isSectionHeader checks if a line looks like a section header.
func isSectionHeader(line string) bool {
	// Empty lines are not headers
	if strings.TrimSpace(line) == "" {
		return false
	}

	lower := strings.ToLower(line)

	// Common section header patterns
	headerPatterns := []string{
		"about", "responsibilities", "qualifications", "requirements",
		"what you", "who you", "benefits", "perks", "compensation",
		"about the", "about us", "the role", "your role",
		"how to apply", "equal opportunity", "diversity",
	}

	for _, pattern := range headerPatterns {
		if strings.HasPrefix(lower, pattern) || strings.Contains(lower, ":") {
			return true
		}
	}

	// Lines that are all caps and short are likely headers
	if len(line) < 50 && line == strings.ToUpper(line) {
		return true
	}

	return false
}

// isBulletPoint checks if a line starts with a bullet point indicator.
func isBulletPoint(line string) bool {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) == 0 {
		return false
	}

	// Check for common bullet prefixes
	bulletPrefixes := []string{
		"-", "*", "•", "●", "○", "◦", "▪", "■",
		"1.", "2.", "3.", "4.", "5.", "6.", "7.", "8.", "9.",
		"1)", "2)", "3)", "4)", "5)", "6)", "7)", "8)", "9)",
	}

	for _, prefix := range bulletPrefixes {
		if strings.HasPrefix(trimmed, prefix) {
			return true
		}
	}

	return false
}

// cleanBulletPoint removes bullet point prefixes and cleans up the text.
func cleanBulletPoint(line string) string {
	trimmed := strings.TrimSpace(line)

	// Remove bullet prefixes
	bulletPrefixes := []string{
		"- ", "* ", "• ", "● ", "○ ", "◦ ", "▪ ", "■ ",
	}

	for _, prefix := range bulletPrefixes {
		if strings.HasPrefix(trimmed, prefix) {
			trimmed = strings.TrimPrefix(trimmed, prefix)
			break
		}
	}

	// Remove numbered prefixes like "1." or "1)"
	numPattern := regexp.MustCompile(`^\d+[.)]\s*`)
	trimmed = numPattern.ReplaceAllString(trimmed, "")

	return strings.TrimSpace(trimmed)
}
