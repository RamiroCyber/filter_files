package models

import "io"

type FileReader struct {
	Filename string
	Reader   io.Reader
}
