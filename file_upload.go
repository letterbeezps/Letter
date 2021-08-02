package letter

import (
	"io"
	"os"
)

type FileUpload struct {
	FileName     string
	FileContents io.ReadCloser
}

func FileUploadFromDisk(fileName string) (*FileUpload, error) {
	fd, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	return &FileUpload{
		FileContents: fd,
		FileName:     fileName,
	}, nil
}
