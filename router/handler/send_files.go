package handler

import (
	"github.com/gofiber/fiber/v2"
	"path/filepath"
	"read_files/models"
	"read_files/pkg/fileprocessor"
	"read_files/util"
	"read_files/util/constants"
)

func SendFiles(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Erro ao processar formulário")
	}

	request := models.RequestForm{
		Files:    form.File["documents"],
		Keywords: util.SeparateWords(form.Value["keywords"]),
	}

	for _, fileHeader := range request.Files {
		filename := fileHeader.Filename
		extension := filepath.Ext(filename)

		if extension != constants.ExtensionTxt && extension != constants.ExtensionPdf {
			return c.Status(fiber.StatusBadRequest).SendString("Extensão de arquivo não permitida")
		}
	}

	if len(request.Files) == 0 {
		return c.Status(fiber.StatusBadRequest).SendString("Envie os arquivos para analise")
	}
	if len(request.Keywords) == 0 {
		return c.Status(fiber.StatusBadRequest).SendString("Envie as palavras chave")
	}

	matchedFiles, err := fileprocessor.ProcessorFile(request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Erro ao processar arquivos")

	}

	zipFiles, err := fileprocessor.CreateZipFile(matchedFiles)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Erro ao criar arquivo zip")

	}

	return c.SendFile(zipFiles)

}
