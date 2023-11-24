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

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if containsAllKeywords(strings.ToUpper(scanner.Text()), keywords) {
			results <- filePath
			return
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
	}
}
