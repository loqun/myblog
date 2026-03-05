package helper

import (
	"html"
	"regexp"
	"strings"
)

// SanitizeHTML sanitizes HTML content by removing dangerous elements
func SanitizeHTML(input string) string {
	// Remove script tags and their content
	scriptRe := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	input = scriptRe.ReplaceAllString(input, "")
	
	// Remove event handlers (onclick, onerror, etc.)
	eventRe := regexp.MustCompile(`(?i)\s*on\w+\s*=\s*["'][^"']*["']`)
	input = eventRe.ReplaceAllString(input, "")
	
	// Remove javascript: protocol
	jsProtocolRe := regexp.MustCompile(`(?i)javascript:`)
	input = jsProtocolRe.ReplaceAllString(input, "")
	
	return input
}

// EscapeLog sanitizes log input to prevent log injection
func EscapeLog(input string) string {
	// Replace newlines and carriage returns
	input = strings.ReplaceAll(input, "\n", "\\n")
	input = strings.ReplaceAll(input, "\r", "\\r")
	return html.EscapeString(input)
}
