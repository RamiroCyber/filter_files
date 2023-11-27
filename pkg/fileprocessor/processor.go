package fileprocessor

import (
	"archive/zip"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"read_files/models"
	"read_files/pkg/filemenager"
	"read_files/util/constants"
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
				errChan <- fmt.Errorf("erro ao abrir o arquivo: %v", err)
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

func CreateZipFile(matchedFiles []models.FileReader) (string, error) {
	zipFilePath := filemenager.GenerateTempFilePath(constants.ZipFileName)
	filemenager.CreateDirIfNotExist()
	zipFile, err := os.Create(zipFilePath)

	if err != nil {
		fmt.Errorf("create zipFile: %v", err)
		return "", err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, nr := range matchedFiles {
		if seeker, ok := nr.Reader.(io.Seeker); ok {
			_, err := seeker.Seek(0, io.SeekStart)
			if err != nil {
				fmt.Errorf("seeker.Seek : %v", err)
				return "", err
			}
		}

		zipEntry, err := zipWriter.Create(nr.Filename)
		if err != nil {
			fmt.Errorf("create zipWriter : %v", err)
			return "", err
		}

		_, err = io.Copy(zipEntry, nr.Reader)
		if err != nil {
			fmt.Errorf("copy : %v", err)
			return "", err
		}
	}

	if err := zipWriter.Close(); err != nil {
		fmt.Errorf("close zipWriter : %v", err)
		return "", err
	}

	return zipFilePath, nil
}
