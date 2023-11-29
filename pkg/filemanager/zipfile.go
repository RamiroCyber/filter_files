package filemanager

import (
	"archive/zip"
	"bytes"
	"io"
	"read_files/structs"
)

func CreateZipFile(matchedFiles []structs.FileReader) (io.Reader, error) {
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
