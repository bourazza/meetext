package pdf

import (
	"regexp"
	"strings"
)

var (
	// Common header/footer patterns
	pageNumberPattern = regexp.MustCompile(`(?m)^\s*\d+\s*$`)
	multiSpacePattern = regexp.MustCompile(`[ \t]+`)
	multiNewlinePattern = regexp.MustCompile(`\n{3,}`)
	
	// Common meeting document artifacts
	headerFooterPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)page \d+ of \d+`),
		regexp.MustCompile(`(?i)confidential`),
		regexp.MustCompile(`(?i)proprietary`),
		regexp.MustCompile(`(?i)draft`),
	}
)

// CleanText removes noise, duplicates, headers/footers, and normalizes whitespace.
func CleanText(raw string) string {
	if raw == "" {
		return ""
	}

	// Step 1: Normalize line endings
	text := strings.ReplaceAll(raw, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	// Step 2: Remove page numbers (standalone lines with just numbers)
	text = pageNumberPattern.ReplaceAllString(text, "")

	// Step 3: Remove common header/footer patterns
	for _, pattern := range headerFooterPatterns {
		text = pattern.ReplaceAllString(text, "")
	}

	// Step 4: Collapse multiple spaces into single space
	text = multiSpacePattern.ReplaceAllString(text, " ")

	// Step 5: Remove leading/trailing whitespace from each line
	lines := strings.Split(text, "\n")
	var cleaned []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleaned = append(cleaned, line)
		}
	}
	text = strings.Join(cleaned, "\n")

	// Step 6: Collapse multiple newlines into double newline (paragraph separation)
	text = multiNewlinePattern.ReplaceAllString(text, "\n\n")

	// Step 7: Final trim
	return strings.TrimSpace(text)
}
