package main

import (
	"fmt"
	"read_files/pkg/filemenager"
	"read_files/pkg/fileprocessor"
	"sync"
)

var keywords = []string{"Teste", "golang", "typescript", "RaMiro"}

func main() {
	files := []string{"file1.txt", "file2.txt", "file3.txt"}

	err := filemenager.CreateDirIfNotExist()
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
			fileprocessor.ProcessFile(filePath, keywords, results)
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println("Arquivo", result, "contém palavras-chave:")
		err := filemenager.MoveFile(result)
		if err != nil {
			fmt.Printf("Erro ao mover o arquivo %s: %v\n", result, err)
		}
	}
}
