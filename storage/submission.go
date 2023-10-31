package storage

import (
	"fmt"
	"github.com/ArcticOJ/blizzard/v0/config"
	"io"
	"mime/multipart"
	"os"
	"path"
)

type submissionStorage struct {
	root string
}

func init() {
	Submission = submissionStorage{root: config.Config.Storage[config.Submissions]}
}

func (s submissionStorage) Create(id uint32, ext string) string {
	return path.Join(s.root, fmt.Sprintf("%d.%s", id, ext))
}

func (submissionStorage) Write(path string, f multipart.File) error {
	file, e := os.Create(path)
	if e != nil {
		panic(e)
		return e
	}
	_, e = io.Copy(file, f)
	return e
}

func (s submissionStorage) GetPath(id uint32, ext string) string {
	return path.Join(s.root, fmt.Sprintf("%d.%s", id, ext))
}
