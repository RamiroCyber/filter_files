package fileprocessor

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"read_files/models"
	"read_files/pkg/file_analyzer"
	"read_files/util/constants"
	"runtime"
	"strings"
	"sync"
)

func ProcessorFilesAll(request models.RequestForm) ([]models.FileReader, error) {
	fileChannel := make(chan *multipart.FileHeader, len(request.Files))
	results := make(chan models.FileReader, len(request.Files))
	errChan := make(chan error, 1)

	for _, fileHeader := range request.Files {
		fileChannel <- fileHeader
	}
	close(fileChannel)

	var wg sync.WaitGroup

	numWorkers := runtime.NumCPU()
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			openFilesForAnalysis(fileChannel, request.Keywords, results, errChan)
		}()
	}

	wg.Wait()
	close(results)

	var matchedFiles []models.FileReader
	for file := range results {
		matchedFiles = append(matchedFiles, file)
	}

	if len(errChan) > 0 {
		return matchedFiles, <-errChan
	}
	return matchedFiles, nil
}

func openFilesForAnalysis(fileChannel <-chan *multipart.FileHeader, keywords []string, results chan<- models.FileReader, errChan chan<- error) {
	for fileHeader := range fileChannel {
		file, err := fileHeader.Open()
		if err != nil {
			errChan <- fmt.Errorf("error opening file: %v", err)
			return
		}

		extension := strings.ToLower(filepath.Ext(fileHeader.Filename))
		var processErr error
		if extension == constants.ExtensionPdf {
			processErr = file_analyzer.SearchKeywordsInPdfFiles(file, fileHeader.Filename, keywords, results)
		} else {
			processErr = file_analyzer.SearchKeywordsInTextFiles(file, fileHeader.Filename, keywords, results)
		}

		file.Close()

		if processErr != nil {
			errChan <- fmt.Errorf("error processing file: %v", processErr)
			return
		}
	}
}
