package fileprocessor

import (
	"fmt"
	"mime/multipart"
	"read_files/models"
	"read_files/pkg/file_analyzer"
	"sync"
)

func processorFilesText(textChannel <-chan *multipart.FileHeader, keywords []string, results chan<- models.FileReader, errChan chan<- error) {
	var wg sync.WaitGroup

	for fileHeader := range textChannel {
		wg.Add(1)
		go func(fh *multipart.FileHeader) {
			defer wg.Done()
			file, err := fh.Open()
			if err != nil {
				errChan <- fmt.Errorf("error ao abrir o arquivo: %v", err)
				return
			}
			defer file.Close()

			if err := file_analyzer.SearchKeywordsInTextFiles(file, fh.Filename, keywords, results); err != nil {
				errChan <- fmt.Errorf("searchKeywordsInTextFiles: %v", err)
			}
		}(fileHeader)
	}
	wg.Wait()
}
