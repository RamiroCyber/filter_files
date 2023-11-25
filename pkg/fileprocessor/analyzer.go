package fileprocessor

import (
	"bufio"
	"fmt"
	"mime/multipart"
	"read_files/models"
	"strings"
)

func containsAllKeywords(line string, keywords []string) bool {
	keywordFound := make(map[string]bool, len(keywords))
	for _, keyword := range keywords {
		if strings.Contains(line, strings.ToUpper(keyword)) {
			keywordFound[keyword] = true
		}
	}
	return len(keywordFound) == len(keywords)
}

func searchKeywordsInFiles(file multipart.File, filename string, keywords []string, results chan<- models.FileReader) {
	scanner := bufio.NewScanner(file)
	var contentBuilder strings.Builder

	for scanner.Scan() {
		contentBuilder.WriteString(strings.ToUpper(scanner.Text()))
		contentBuilder.WriteString(" ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return
	}

	if containsAllKeywords(contentBuilder.String(), keywords) {
		results <- models.FileReader{Filename: filename, Reader: file}
	}
}
