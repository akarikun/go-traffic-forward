package common

import "strings"

func IsNullOrEmpty(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}
