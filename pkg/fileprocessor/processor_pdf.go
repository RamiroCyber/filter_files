package fileprocessor

import (
	"fmt"
	"mime/multipart"
	"read_files/models"
	"read_files/pkg/file_analyzer"
	"sync"
)

func processorFilesPdf(pdfChannel <-chan *multipart.FileHeader, keywords []string, results chan<- models.FileReader, errChan chan<- error) {
	var wg sync.WaitGroup

	for fileHeader := range pdfChannel {
		wg.Add(1)
		go func(fh *multipart.FileHeader) {
			defer wg.Done()
			file, err := fh.Open()
			if err != nil {
				errChan <- fmt.Errorf("error ao abrir o arquivo: %v", err)
				return
			}
			defer file.Close()

			if err := file_analyzer.SearchKeywordsInPdfFiles(file, fh.Filename, keywords, results); err != nil {
				errChan <- fmt.Errorf("searchKeywordsInPdfFiles: %v", err)
			}
		}(fileHeader)
	}
	wg.Wait()
}
