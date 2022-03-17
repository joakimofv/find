package find

import (
	"fmt"
	"os"
	"strings"
)

// This file contains utility functions for wildcard pattern matching.

// LongestFixedPart gives the longest substring not containing any wildcards.
func LongestFixedPart(pattern string) string {
	parts := strings.Split(pattern, "*")
	if len(parts) == 0 {
		return ""
	}
	longestPart := parts[0]
	for _, part := range parts {
		if len(part) > len(longestPart) {
			longestPart = part
		}
	}
	return longestPart
}

// Replace searches for oldPattern in line and replaces it with newPattern.
// '*' counts as wildcards in the patterns.
//
// If the line was modified then it will be returned along with true,
// otherwise false will be returned and the string will be empty.
func Replace(line, oldPattern, newPattern string) (string, bool) {
	// Expose strings to the user, but use rune slices internally
	return replace(line, []rune(oldPattern), []rune(newPattern))
}
func replace(line string, oldPattern, newPattern []rune) (string, bool) {
	// Try to match line into oldPattern, if successful then insert the wildcard contents into newPattern.
	wildcardContents, matched := match(oldPattern, []rune(line))
	if !matched {
		return "", false
	}
	// The patterns might match inside the wildcards too
	for i, wildcardContent := range wildcardContents {
		if wildcardContent == "" || wildcardContent == line {
			// Skip to avoid getting stuck in infinite recursion
			continue
		}
		mod, ok := replace(wildcardContent, oldPattern, newPattern)
		if !ok {
			continue
		}
		wildcardContents[i] = mod
	}
	lineMod, err := insert(newPattern, wildcardContents)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return "", false
	}
	if line == lineMod {
		return "", false
	}
	return lineMod, true
}

func match(pattern, runes []rune) ([]string, bool) {
	var wildcardContents []string
	idx := 0
	for i, r := range pattern {
		switch r {
		case '*':
			// Consume runes until encountering the next character of the pattern
			var next rune
			if i+1 < len(pattern) {
				next = pattern[i+1]
			}
			var wildcardContent string
			for ; idx < len(runes); idx++ {
				if runes[idx] == next {
					// The match might either be the end of the wildcard sequence, or somewhere in
					// the middle. Try match the remainder of the pattern and string in order to tell.
					trailingWildcardContents, ok := match(pattern[i+1:], runes[idx:])
					if ok {
						// We are done, all has been processed and it was a match
						wildcardContents = append(wildcardContents, wildcardContent)
						wildcardContents = append(wildcardContents, trailingWildcardContents...)
						return wildcardContents, true
					}
					// Add the character to the wildcard and continue on
				}
				wildcardContent += string(runes[idx])
			}
			// Ran out of runes of the target.
			wildcardContents = append(wildcardContents, wildcardContent)

		default:
			// Must have an exact match
			if idx >= len(runes) || runes[idx] != r {
				return nil, false
			}
			idx++
		}
	}
	// If there are unmatched runes left in the target then the match is invalid
	if idx != len(runes) {
		return nil, false
	}
	return wildcardContents, true
}

func insert(pattern []rune, parts []string) (string, error) {
	var ret string
	idx := 0
	for i := 0; i < len(pattern); i++ {
		if pattern[i] == '*' {
			if idx >= len(parts) {
				// Bad pattern
				return "", fmt.Errorf("Too many of '*' in target pattern '%s'\n", string(pattern))
			}
			ret += parts[idx]
			idx++
		} else {
			ret += string(pattern[i])
		}
	}
	return ret, nil
}
