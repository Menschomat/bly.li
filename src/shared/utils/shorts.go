package utils

import (
	"regexp"
)

func IsValidShort(short string) bool {
	var re = regexp.MustCompile(`(?m)^[a-zA-Z]+$`)
	return len(re.FindStringIndex(short)) > 0
}
