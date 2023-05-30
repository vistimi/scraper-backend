package util

import (
	"fmt"
	"golang.org/x/exp/slices"
	"regexp"
)

// Find index of arrayNeedle matching one arrayRegExp
func FindIndexRegExp(arrayRegExp []string, arrayNeedle []string) int {
	return slices.IndexFunc(arrayNeedle, func(needle string) bool {
		// pass through all tags of the image and its derived tags to match an unwated tag
		idx := slices.IndexFunc(arrayRegExp, func(strRegExp string) bool {
			regexpMatch := fmt.Sprintf(`[\-\_\w\d]*%s[\-\_\w\d]*`, strRegExp)
			matched, err := regexp.Match(regexpMatch, []byte(needle)) // e.g. match if strRegExp has `art` and needle has `artmodel`
			if err != nil {
				return false
			}
			return matched
		})
		if idx == -1 {
			return false
		} else {
			return true // if needle is present return true
		}
	})
}
