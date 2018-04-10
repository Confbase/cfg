package filetype

import "strings"

const (
	Json    = "json"
	Toml    = "toml"
	Yaml    = "yaml"
	Unknown = "unknown"
)

func Guess(fileName string) string {
	lowered := strings.ToLower(fileName)
	if strings.HasSuffix(lowered, ".json") {
		return Json
	} else if strings.HasSuffix(lowered, ".toml") {
		return Toml
	} else if strings.HasSuffix(lowered, ".yaml") || strings.HasSuffix(lowered, ".yml") {
		return Yaml
	}
	return "unknown"
}
