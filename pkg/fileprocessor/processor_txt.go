package fileprocessor

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"read_files/models"
	"read_files/pkg/file_analyzer"
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

func CreateZipFile(matchedFiles []models.FileReader) (io.Reader, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for _, nr := range matchedFiles {
		if seeker, ok := nr.Reader.(io.Seeker); ok {
			_, err := seeker.Seek(0, io.SeekStart)
			if err != nil {
				return nil, err
			}
		}

		zipEntry, err := zipWriter.Create(nr.Filename)
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(zipEntry, nr.Reader)
		if err != nil {
			return nil, err
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}
