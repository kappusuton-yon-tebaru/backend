package utils

import (
	"fmt"
	"strings"
)

func ArrayWithComma(items []string, conjunction string) string {
	return fmt.Sprintf("%s %s %s", strings.Join(items[:len(items)-1], ", "), conjunction, items[len(items)-1])
}
