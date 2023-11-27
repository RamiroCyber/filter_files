package models

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"read_files/util/constants"
	"strings"
)

type RequestForm struct {
	Files    []*multipart.FileHeader
	Keywords []string
}

func (r *RequestForm) Validate() error {
	var errMessages []string

	if len(r.Files) == 0 {
		errMessages = append(errMessages, "envie os arquivos para análise")
	}
	if len(r.Keywords) == 0 {
		errMessages = append(errMessages, "envie as palavras-chave")
	}

	if len(errMessages) > 0 {
		return errors.New(strings.Join(errMessages, "; \n"))
	}

	return nil
}

func (r *RequestForm) ValidateExt() error {

	allowedExtensions := map[string]bool{
		constants.ExtensionTxt:  true,
		constants.ExtensionPdf:  true,
		constants.ExtensionDoc:  true,
		constants.ExtensionDocx: true,
	}

	for _, fileHeader := range r.Files {
		filename := fileHeader.Filename
		extension := strings.ToLower(filepath.Ext(filename))

		if !allowedExtensions[extension] {
			return errors.New("extensão de arquivo não permitida: " + extension + " no arquivo " + filename)
		}
	}
	return nil
}
