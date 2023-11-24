package fileprocessor

import "strings"

func containsAllKeywords(line string, keywords []string) bool {
	keywordFound := make(map[string]bool, len(keywords))
	for _, keyword := range keywords {
		if strings.Contains(line, strings.ToUpper(keyword)) {
			keywordFound[keyword] = true
		}
	}
	return len(keywordFound) == len(keywords)
}
