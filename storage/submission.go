package storage

import (
	"fmt"
	"github.com/ArcticOJ/blizzard/v0/config"
	"github.com/ArcticOJ/blizzard/v0/utils"
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

func (s submissionStorage) Create(ext string) string {
	rnd := utils.Rand(8, "")
	if rnd == "" {
		return ""
	}
	return path.Join(s.root, fmt.Sprintf("%s.%s", rnd, ext))
}

func (submissionStorage) Write(path string, f io.Reader) (error, func()) {
	file, e := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0755)
	if e != nil {
		return e, nil
	}
	_, e = io.Copy(file, f)
	return e, func() {
		os.Remove(path)
	}
}

func (s submissionStorage) GetPath(name string) string {
	return path.Join(s.root, name)
}
