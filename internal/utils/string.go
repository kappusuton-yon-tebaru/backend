package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func ArrayWithComma(items []string, conjunction string) string {
	return fmt.Sprintf("%s %s %s", strings.Join(items[:len(items)-1], ", "), conjunction, items[len(items)-1])
}

func ToKebabCase(s string) string {
	lowercase := strings.ToLower(s)

	pattern := "[ _]"
	return regexp.MustCompile(pattern).ReplaceAllString(lowercase, "-")
}
