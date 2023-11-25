package models

import "mime/multipart"

type RequestForm struct {
	Files    []*multipart.FileHeader
	Keywords []string
}
