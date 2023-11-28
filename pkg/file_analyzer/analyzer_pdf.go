package file_analyzer

import (
	"fmt"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
	"mime/multipart"
	"read_files/models"
	"strings"
)

func containsAllKeywordsPdf(line string, keywords []string) bool {
	line = strings.ToUpper(line)

	keywordFound := make(map[string]bool, len(keywords))
	for _, keyword := range keywords {
		upperKeyword := strings.ToUpper(keyword)

		if strings.Contains(line, upperKeyword) {
			keywordFound[keyword] = true
		}
	}

	for _, found := range keywordFound {
		if !found {
			return false
		}
	}

	return true
}

func SearchKeywordsInPdfFiles(file multipart.File, filename string, keywords []string, results chan<- models.FileReader) error {
	pdfReader, err := model.NewPdfReader(file)
	if err != nil {
		return fmt.Errorf("NewPdfReader: %v", err)
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return fmt.Errorf("GetNumPages: %v", err)
	}

	var found bool
	for i := 1; i <= numPages && !found; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			return fmt.Errorf("error scanner: %v", err)
		}

		ex, err := extractor.New(page)
		if err != nil {
			return fmt.Errorf("extractor.New: %v", err)
		}

		text, err := ex.ExtractText()
		if err != nil {
			fmt.Printf(text)
			return fmt.Errorf("ex.ExtractText: %v", err)
		}

		if containsAllKeywordsPdf(text, keywords) {
			found = true
			results <- models.FileReader{Filename: filename, Reader: file}
		}
	}
	return nil
}
