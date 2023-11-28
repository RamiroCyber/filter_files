package fileprocessor

import (
	"mime/multipart"
	"path/filepath"
	"read_files/models"
	"read_files/util/constants"
	"strings"
	"sync"
)

func ProcessorFilesAll(request models.RequestForm) ([]models.FileReader, error) {
	pdfChannel := make(chan *multipart.FileHeader, len(request.Files))
	textChannel := make(chan *multipart.FileHeader, len(request.Files))
	results := make(chan models.FileReader, len(request.Files))
	errChan := make(chan error, len(request.Files))

	for _, fileHeader := range request.Files {
		filename := fileHeader.Filename
		extension := strings.ToLower(filepath.Ext(filename))

		if extension == constants.ExtensionPdf {
			pdfChannel <- fileHeader
		} else {
			textChannel <- fileHeader
		}
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		processorFilesPdf(pdfChannel, request.Keywords, results, errChan)
	}()

	go func() {
		defer wg.Done()
		processorFilesText(textChannel, request.Keywords, results, errChan)
	}()

	close(pdfChannel)
	close(textChannel)

	go func() {
		wg.Wait()
		close(results)
		close(errChan)
	}()

	var matchedFiles []models.FileReader
	for {
		select {
		case file, ok := <-results:
			if !ok {
				return matchedFiles, nil
			}
			matchedFiles = append(matchedFiles, file)
		case err := <-errChan:
			return matchedFiles, err
		}
	}
}
