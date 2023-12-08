package main_test

import (
	"bytes"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"read_files/router/handler"
	"testing"
)

func TestUpload(t *testing.T) {
	app := fiber.New()

	app.Post("/upload", handler.Upload)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	filePaths := []string{"arquivos_test/test.doc", "arquivos_test/test.docx", "arquivos_test/test.txt"}
	for _, path := range filePaths {
		file, err := os.Open(path)
		assert.NoError(t, err)
		part, err := writer.CreateFormFile("documents", filepath.Base(file.Name()))
		assert.NoError(t, err)
		_, err = io.Copy(part, file)
		assert.NoError(t, err)
		file.Close()
	}

	keywords := []string{"golang", "ramiro"}
	for _, keyword := range keywords {
		_ = writer.WriteField("keywords", keyword)
	}

	assert.NoError(t, writer.Close())

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

}
