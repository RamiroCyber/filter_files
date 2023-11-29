package structs

import "io"

type FileReader struct {
	Filename string
	Reader   io.Reader
}
