package config

import (
	"strings"
)

func parseConfigLine(line string) (key, value string) {
	key, value, _ = strings.Cut(line, " ")
	return
}

func parseConfigKey(key string, n int) (parts []string, rest string) {
	rest = key
	parts = make([]string, n)
	for i := 0; i < n; i++ {
		sepIDX := strings.IndexByte(rest, '.')
		if sepIDX == -1 {
			parts[i] = rest
			rest = ""
			return
		}
		parts[i] = rest[0:sepIDX]
		rest = rest[sepIDX+1:]
	}
	return
}
