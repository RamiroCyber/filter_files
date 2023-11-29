package file_analyzer

import (
	"bufio"
	"fmt"
	"mime/multipart"
	"read_files/structs"
	"strings"
)

func containsAllKeywordsTxt(line string, keywords []string) bool {
	keywordFound := make(map[string]bool, len(keywords))
	for _, keyword := range keywords {

		if strings.Contains(line, strings.ToUpper(keyword)) {
			keywordFound[keyword] = true
		}
	}
	return len(keywordFound) == len(keywords)
}

func SearchKeywordsInTextFiles(file multipart.File, filename string, keywords []string, results chan<- structs.FileReader) error {
	scanner := bufio.NewScanner(file)
	var contentBuilder strings.Builder

	for scanner.Scan() {
		contentBuilder.WriteString(strings.ToUpper(scanner.Text()))
		contentBuilder.WriteString(" ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Errorf("error scanner: %v", err)
		return err
	}

	if containsAllKeywordsTxt(contentBuilder.String(), keywords) {
		results <- structs.FileReader{Filename: filename, Reader: file}
	}
	return nil
}
