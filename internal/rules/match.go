package rules

import "regexp"

// rulePatternMatch applies a regex pattern to a line and returns the matched content.
func rulePatternMatch(line, pattern string) string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return ""
	}
	match := re.FindString(line)
	if match != "" {
		return match
	}
	return ""
}

// findMatchPosition finds the column position of the first match in the line.
func findMatchPosition(line, pattern string) int {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return 0
	}
	loc := re.FindStringIndex(line)
	if loc == nil {
		return 0
	}
	return loc[0] + 1 // 1-based column
}
