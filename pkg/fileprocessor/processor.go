package fileprocessor

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"read_files/models"
	"sync"
)

func ProcessorFile(request models.RequestForm) ([]models.FileReader, error) {
	var wg sync.WaitGroup

	results := make(chan models.FileReader, len(request.Files))
	errChan := make(chan error, len(request.Files))

	var matchedFiles []models.FileReader

	for _, fileHeader := range request.Files {
		wg.Add(1)
		go func(fh *multipart.FileHeader) {
			defer wg.Done()
			file, err := fh.Open()
			if err != nil {
				errChan <- fmt.Errorf("error ao abrir o arquivo: %v", err)
				return
			}
			defer file.Close()

			if err := searchKeywordsInFiles(file, fh.Filename, request.Keywords, results); err != nil {
				errChan <- err
			}
		}(fileHeader)
	}

	go func() {
		wg.Wait()
		close(results)
		close(errChan)
	}()

	for err := range errChan {
		return nil, err
	}

	for nr := range results {
		matchedFiles = append(matchedFiles, nr)
	}
	return matchedFiles, nil
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
