package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"read_files/pkg/filemenager"
	"read_files/pkg/fileprocessor"
	"sync"
)

func SendFiles(c *fiber.Ctx) error {

	var keywords = []string{"Teste", "golang", "typescript", "RaMiro"}

	files := []string{"file1.txt", "file2.txt", "file3.txt"}

	err := filemenager.CreateDirIfNotExist()
	if err != nil {
		fmt.Println("Erro ao criar a pasta:", err)
	}

	var wg sync.WaitGroup
	results := make(chan string, len(files))

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
	return c.SendStatus(fiber.StatusOK)
}
