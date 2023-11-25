package handler

import (
	"github.com/gofiber/fiber/v2"
	"read_files/models"
	"read_files/pkg/fileprocessor"
	"read_files/util"
)

func SendFiles(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Erro ao processar formulário"})
	}

	request := models.RequestForm{
		Files:    form.File["documents"],
		Keywords: util.SeparateWords(form.Value["keywords"]),
	}

	if len(request.Files) == 0 {
		return c.Status(fiber.StatusBadRequest).SendString("Envie os arquivos para analise")
	}
	if len(request.Keywords) == 0 {
		return c.Status(fiber.StatusBadRequest).SendString("Envie as palavras chave")
	}

	matchedFiles, err := fileprocessor.ProcessorFile(request)
	if err != nil {
		// Trate o erro
	}

	zipFiles, err := fileprocessor.CreateAndSendZipFile(matchedFiles)
	if err != nil {
		// Trate o erro
	}

	c.Set("Content-Type", "application/zip")
	c.Set("Content-Disposition", "attachment; filename=matched_files.zip")
	return c.SendFile(zipFiles)

}
