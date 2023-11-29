package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"read_files/pkg/filemanager"
	"read_files/pkg/fileprocessor"
	"read_files/structs"
	"read_files/util"
	"read_files/util/constants"
)

func SendFiles(c *fiber.Ctx) error {
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

	matchedFiles, err := fileprocessor.ProcessorFilesAll(request)

	if err != nil {
		util.CustomLogger(constants.Error, fmt.Sprintf("fileprocessor.ProcessorFile: %v", err))
		return c.Status(fiber.StatusInternalServerError).SendString("Error ao processar arquivos")
	}

	if matchedFiles == nil {
		return c.Status(fiber.StatusNotFound).SendString("Palavras-chave não encontradas")
	}

	zipFiles, err := filemanager.CreateZipFile(matchedFiles)
	if err != nil {
		util.CustomLogger(constants.Error, fmt.Sprintf("fileprocessor.CreateZipFile: %v", err))
		return c.Status(fiber.StatusInternalServerError).SendString("Error ao criar arquivo zip")

	}
	return c.Status(fiber.StatusOK).SendStream(zipFiles)
}
