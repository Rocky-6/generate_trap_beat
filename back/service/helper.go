package service

import (
	"regexp"
)

func check_regexp(reg, str string) bool {
	r := regexp.MustCompile(reg)
	return r.MatchString(str)
}
