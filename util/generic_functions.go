package util

import "strings"

func SeparateWords(wordsString []string) (keywords []string) {
	if len(wordsString) > 0 {
		fields := strings.Fields(wordsString[0])
		for _, field := range fields {
			if field != "" {
				keywords = append(keywords, field)
			}
		}
	}
	return keywords
}
