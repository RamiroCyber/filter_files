package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var keywords = []string{"teste", "golang", "typescript"}

func containsKeyword(line string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(line, keyword) {
			return true
		}
	}
	return false
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.Mkdir(dir, 0755)
	}
	return nil
}

func moveFile(filePath, destination string) error {
	cmd := exec.Command("mv", filePath, destination)
	return cmd.Run()
}

func processFile(filePath string, results chan<- string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if containsKeyword(scanner.Text(), keywords) {
			results <- filePath
			return
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
	}
}

func main() {
	files := []string{"file1.txt", "file2.txt", "file3.txt"}
	targetDir := "pasta"

	err := createDirIfNotExist(targetDir)
	if err != nil {
		fmt.Println("Erro ao criar a pasta:", err)
		return
	}

	var wg sync.WaitGroup
	results := make(chan string, 100)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, filePath := range files {
			processFile(filePath, results)
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println("Arquivo", result, "contém palavras-chave:")
		err := moveFile(result, targetDir)
		if err != nil {
			fmt.Printf("Erro ao mover o arquivo %s: %v\n", result, err)
		}
	}
}
