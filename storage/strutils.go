package storage

import (
	"strings"
)

func SubStr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func StrLen(s string) int {
	return len([]rune(s))
}

func StringIsLessIgnoreCase(a, b string) bool {
	return strings.ToLower(a) < strings.ToLower(b)
}
