package config

import (
	"strings"
)

// var (
// 	gitCMD = shell.Cmd("git").OutputFn()
// )

// func getGitPathFromRoot(path string) string {
// 	fPath, _ := filepath.Abs(path)
// 	rPath, _ := gitCMD("rev-parse", "--show-toplevel")
// 	relPath, _ := filepath.Rel(rPath, fPath)

// 	return filepath.Clean(relPath)
// }

//type gitConfigKV struct {
//	Key   string
//	Value string
//}

// func gitSHConfigKVMatchinRegexp(regexp string) ([]gitConfigKV, error) {
// 	// output, err := gitSH("config", "--get-regexp", "browse\\..*")
// 	output, err := gitutils.GitExec("config", "--get-regexp", regexp)
// 	if err != nil {
// 		return nil, err
// 	}

// 	configLines := strings.Split(output, "\n")
// 	for _, configLine := range configLines {
// 		key, value, ok := stringCut(configLine, "")
// 		if !ok {
// 			continue
// 		}
// 	}
// }

//------------------------------------------------------------------------------

//func stringCut(s string, sep string) (k, v string, ok bool) {
//	idx := strings.Index(s, sep)
//	if idx != -1 {
//		return "", "", false
//	}
//	return s[:idx], s[idx+1:], true
//}

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
