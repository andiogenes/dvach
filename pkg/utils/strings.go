package utils

import (
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"html"
	"strings"
)

var (
	tagTrimmer = bluemonday.StrictPolicy()
)

// formatHTML prepares HTML data to console output.
//
// Indents data if corresponding value passed.
func FormatHTML(data string, indent bool) string {
	separator := "\n"
	if indent {
		separator = "\n\t"
	}

	// Replace <br> tags with escape characters
	dest := strings.ReplaceAll(data, "<br>", separator)
	// Remove HTML tags
	dest = tagTrimmer.Sanitize(dest)
	// Replace HTML escape entities with corresponding symbols
	dest = html.UnescapeString(dest)

	return dest
}

// BlueColor wraps string with ANSI escape code for blue colouring.
func BlueColor(str string) string {
	return fmt.Sprintf("\033[1;36m%s\033[0m", str)
}

// GreenColor wraps string with ANSI escape code for green colouring.
func GreenColor(str string) string {
	return fmt.Sprintf("\033[1;32m%s\033[0m", str)
}

// YellowColor wraps string with ANSI escape code for yellow colouring.
func YellowColor(str string) string {
	return fmt.Sprintf("\033[1;33m%s\033[0m", str)
}
