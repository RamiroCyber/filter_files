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
		return errors.New(strings.Join(errMessages, "; "))
	}

	return nil
}

func (r *RequestForm) ValidateExt() error {
	for _, fileHeader := range r.Files {
		filename := fileHeader.Filename
		extension := strings.ToLower(filepath.Ext(filename))

		if extension != constants.ExtensionTxt && extension != constants.ExtensionPdf {
			return errors.New("extensão de arquivo não permitida: " + extension)
		}
	}
	return nil
}
