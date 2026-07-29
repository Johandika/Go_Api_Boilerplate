package shared

import "strings"

func NormalizeSearch(value string) string {
	return strings.TrimSpace(value)
}
