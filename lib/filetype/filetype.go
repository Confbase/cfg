package filetype

import "strings"

const (
	JSON = "JSON"
)

func Guess(fileName string) string {
	if strings.HasSuffix(strings.ToLower(fileName), ".json") {
		return JSON
	}
	return "unknown"
}
