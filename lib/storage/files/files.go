package files

import (
	"ASKLYBOT/lib/e"
	"ASKLYBOT/lib/storage"
	"os"
)

const defaultPrem = 0774

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.WrapIferr("can't save", err) }()

	filePath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(filepath, defaultPrem); err != nil {
		return err
	}

}
