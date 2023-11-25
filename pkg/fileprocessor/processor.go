package fileprocessor

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ProcessFile(filePath string, keywords []string, results chan<- string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	var contentBuilder strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		contentBuilder.WriteString(strings.ToUpper(scanner.Text()))
		contentBuilder.WriteString(" ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return
	}

	if containsAllKeywords(contentBuilder.String(), keywords) {
		results <- filePath
	}
}
