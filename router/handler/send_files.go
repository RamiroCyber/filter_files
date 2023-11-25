package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
	"os"
	"read_files/pkg/filemenager"
	"read_files/pkg/fileprocessor"
	"read_files/util/constants"
	"sync"
)

var keywords = []string{"Teste", "golang"}

func SendFiles(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Erro ao processar formulário"})
	}

	files := form.File["documents"]

	var wg sync.WaitGroup
	results := make(chan string, len(files))

	if err := filemenager.CreateDirIfNotExist(constants.TempDirPath); err != nil {
		fmt.Println("Erro ao criar o diretório:", err)
		return err
	}

	for _, fileHeader := range files {
		wg.Add(1)
		go func(fh *multipart.FileHeader) {
			defer wg.Done()
			tempFilePath := constants.TempDirPath + fh.Filename
			if err := c.SaveFile(fh, tempFilePath); err != nil {
				fmt.Println("Erro ao salvar o arquivo:", err)
				return
			}
			fileprocessor.ProcessFile(tempFilePath, keywords, results)

			if err := os.Remove(tempFilePath); err != nil {
				fmt.Println("Erro ao excluir o arquivo temporário:", err)
			}
		}(fileHeader)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var matchedFiles []string
	for filename := range results {
		matchedFiles = append(matchedFiles, filename)
	}

	return c.JSON(matchedFiles)
}
