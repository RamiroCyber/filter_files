package file_processor

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"read_files/pkg/file_analyzer"
	"read_files/structs"
	"read_files/util"
	"read_files/util/constants"
	"strings"
	"sync"
)

func ProcessorFilesAll(request structs.RequestForm) ([]structs.FileReader, error) {
	fileChannel := make(chan *multipart.FileHeader, len(request.Files))
	results := make(chan structs.FileReader, len(request.Files))
	errChan := make(chan error, 1)

	for _, fileHeader := range request.Files {
		fileChannel <- fileHeader
	}
	close(fileChannel)

	var wg sync.WaitGroup

	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			openFilesForAnalysis(fileChannel, request.Keywords, results, errChan)
		}()
	}

	wg.Wait()
	close(results)

	var matchedFiles []structs.FileReader
	for file := range results {
		matchedFiles = append(matchedFiles, file)
	}

	if len(errChan) > 0 {
		return matchedFiles, <-errChan
	}
	return matchedFiles, nil
}

func openFilesForAnalysis(fileChannel <-chan *multipart.FileHeader, keywords []string, results chan<- structs.FileReader, errChan chan<- error) {
	for fileHeader := range fileChannel {
		file, err := fileHeader.Open()
		if err != nil {
			util.CustomLogger(constants.Error, fmt.Sprintf("error opening file: %v", err))
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
			util.CustomLogger(constants.Error, fmt.Sprintf("error processing file: %v", err))
			errChan <- fmt.Errorf("error processing file: %v", processErr)
			return
		}
	}
}
