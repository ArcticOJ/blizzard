package storage

import (
	"blizzard/config"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
)

type SubmissionStorage struct {
}

func (*SubmissionStorage) Create(id uint32, ext string) string {
	return path.Join(config.Config.Storage.Submissions, fmt.Sprintf("%d.%s", id, ext))
}

func (*SubmissionStorage) Write(path string, f multipart.File) error {
	file, e := os.Create(path)
	if e != nil {
		return e
	}
	_, e = io.Copy(file, f)
	return e
}
