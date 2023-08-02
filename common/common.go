package common

import "regexp"

func ShouldBeExcluded(name string, exclusions []string) bool {
	for _, exclusion := range exclusions {
		if match, _ := regexp.MatchString(exclusion, name); match {
			return true
		}
	}

	return false
}
