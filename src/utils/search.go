package utils

import (
	"fmt"
	"regexp"
)


// Tells whether the RegExp of the array elements match the needle.
// Returns idxPtr arrRegExps, idxPtr needles, error
func ContainsRegExp(arrRegExps []string, needles []string) (*int, *int, error) {
	for i, str := range arrRegExps {
		regexpMatch := fmt.Sprintf(`[\-\_\w\d]*%s[\-\_\w\d]*`, str)
		for j, needle := range needles {
			matched, err := regexp.Match(regexpMatch, []byte(needle))
			if err != nil {
				return nil, nil, err
			}
			if matched {
				return &i, &j, nil
			}
		}
	}
	return nil, nil, nil
}
