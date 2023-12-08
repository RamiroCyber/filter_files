package file_analyzer

import (
	"fmt"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
	"mime/multipart"
	"read_files/structs"
	"read_files/util"
	"read_files/util/constants"
	"strings"
)

func containsAllKeywordsPdf(line string, keywords []string) bool {
	line = strings.ToUpper(line)

	for _, keyword := range keywords {
		upperKeyword := strings.ToUpper(keyword)
		if !strings.Contains(line, upperKeyword) {
			return false
		}
	}

	return true
}

func SearchKeywordsInPdfFiles(file multipart.File, filename string, keywords []string, results chan<- structs.FileReader) error {
	pdfReader, err := model.NewPdfReader(file)
	if err != nil {
		util.CustomLogger(constants.Error, fmt.Sprintf("NewPdfReader: %v", err))
		return fmt.Errorf("NewPdfReader: %v", err)
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		util.CustomLogger(constants.Error, fmt.Sprintf("GetNumPages: %v", err))
		return fmt.Errorf("GetNumPages: %v", err)
	}

	var found bool
	for i := 1; i <= numPages && !found; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			util.CustomLogger(constants.Error, fmt.Sprintf("GetPage: %v", err))
			return fmt.Errorf("GetPage: %v", err)
		}

		ex, err := extractor.New(page)
		if err != nil {
			util.CustomLogger(constants.Error, fmt.Sprintf("extractor.New: %v", err))
			return fmt.Errorf("Extractor.New: %v", err)
		}

		text, err := ex.ExtractText()
		if err != nil {
			util.CustomLogger(constants.Error, fmt.Sprintf("ExtractText: %v", err))
			return fmt.Errorf("ex.ExtractText: %v", err)
		}

		if containsAllKeywordsPdf(text, keywords) {
			found = true
			results <- structs.FileReader{Filename: filename, Reader: file}
		}
	}
	return nil
}
