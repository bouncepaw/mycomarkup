package tools

import (
	"regexp"
	"strings"
)

var rocketLinkMatcher *regexp.Regexp

func init() {
	var rlm, err = regexp.Compile(`^\s*=>`)
	if err != nil {
		panic("disaster")
	}

	rocketLinkMatcher = rlm
}

// MigrateRocketLinks replaces all instances of the old rocket link syntax with the new one it can find in the given document and returns the modified version.
func MigrateRocketLinks(old string) string {
	/*
		Read line
		If no next line, finish
		If the line does not match \s*=>, write the line as is, goto 1
		Write the part until >, look at the rest of the line
		If, trimmed of whitespace on the left and on the right, there is a whitespace in the string, add | after it and write the string
		Goto 1
	*/

	var newLines []string

	for _, line := range strings.Split(old, "\n") {
		if rocketLinkMatcher.MatchString(line) {
			newLines = append(newLines, oldRocketToNew(line))
		} else {
			newLines = append(newLines, line)
		}
	}

	return strings.Join(newLines, "\n")
}

// note no \n in rocket
// this function enforces some style
func oldRocketToNew(rocket string) string {
	gtpos := strings.IndexRune(rocket, '>')
	newRocket := rocket[:gtpos+1] + " "
	rocket = strings.TrimSpace(rocket[gtpos+1:])
	if wspos := strings.IndexRune(rocket, ' '); wspos == -1 {
		newRocket += rocket
	} else {
		newRocket += rocket[:wspos] + " |" + rocket[wspos:]
	}

	return newRocket
}
