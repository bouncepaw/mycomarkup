package tools

import (
	"regexp"
	"strings"
)

var headingMatcher *regexp.Regexp

func init() {
	var headingM, err = regexp.Compile(`^\s*#{1,6} `)
	if err != nil {
		panic("disaster")
	}

	headingMatcher = headingM
}

// MigrateHeadings replaces all instances of the old heading syntax with the new one it can find in the given document and returns the modified version.
func MigrateHeadings(old string) string {
	var result strings.Builder
	for i, line := range strings.Split(old, "\n") {
		if i > 0 {
			result.WriteRune('\n')
		}
		if headingMatcher.MatchString(line) {
			result.WriteString(oldHeadingToNew(line))
		} else {
			result.WriteString(line)
		}
	}
	return result.String()
}

type headingReplacer int

const (
	eatingSpace headingReplacer = iota
	replacingOctothorps
	writingRest
)

// this function enforces some style
func oldHeadingToNew(heading string) string {
	var (
		result         strings.Builder
		state          = eatingSpace
		octothorpCount = 0
	)
	for _, r := range heading {
	runeCaster:
		switch state {
		case eatingSpace:
			switch r {
			case ' ', '\t':
				result.WriteRune(r)
			case '#':
				state = replacingOctothorps
				goto runeCaster
			default:
				panic("bad input")
			}
		case replacingOctothorps:
			switch r {
			case '#':
				octothorpCount++
			case ' ':
				switch octothorpCount {
				case 1, 2:
					result.WriteRune('=')
				case 3:
					result.WriteString("==")
				case 4:
					result.WriteString("===")
				case 5, 6:
					result.WriteString("====")
				}
				result.WriteRune(' ')
				state = writingRest
			default:
				panic("bad input")
			}
		case writingRest:
			result.WriteRune(r)
		}
	}
	return result.String()
}
