package storage

import (
	"fmt"
	"github.com/ArcticOJ/blizzard/v0/config"
	"io"
	"os"
	"path"
)

type submissionStorage struct {
	root string
}

func init() {
	Submission = submissionStorage{root: config.Config.Blizzard.Storage[config.Submissions]}
}

func (s submissionStorage) Create(id uint32, ext string) string {
	return path.Join(s.root, fmt.Sprintf("%d.%s", id, ext))
}

func (submissionStorage) Write(path string, f io.Reader) error {
	file, e := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0755)
	if e != nil {
		return e
	}
	_, e = io.Copy(file, f)
	return e
}

func (s submissionStorage) GetPath(id uint32, ext string) string {
	return path.Join(s.root, fmt.Sprintf("%d.%s", id, ext))
}
