package storage

import (
	"github.com/ArcticOJ/blizzard/v0/config"
	"github.com/google/uuid"
	"path"
)

type readmesStorage struct {
	root string
}

func init() {
	READMEs = readmesStorage{root: config.Config.Storage[config.READMEs]}
}

func (s readmesStorage) GetPath(uid uuid.UUID) string {
	return path.Join(s.root, uid.String()+".md")
}
