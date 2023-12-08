package handler

import (
	"github.com/gofiber/fiber/v2"
	"read_files/pkg/file_manager"
	"read_files/pkg/file_processor"
	"read_files/structs"
	"read_files/util"
)

func Upload(c *fiber.Ctx) error {
	form, err := c.MultipartForm()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error ao processar formulário")
	}

	request := structs.RequestForm{
		Files:    form.File["documents"],
		Keywords: util.SeparateWords(form.Value["keywords"]),
	}

	if err = request.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err = request.ValidateExt(); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	matchedFiles, err := file_processor.ProcessorFilesAll(request)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error ao processar arquivos")
	}

	if matchedFiles == nil {
		return c.Status(fiber.StatusNotFound).SendString("Palavras-chave não encontradas")
	}

	zipFiles, err := file_manager.CreateZipFile(matchedFiles)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error ao criar arquivo zip")

	}

	return c.Status(fiber.StatusOK).SendStream(zipFiles)
}
