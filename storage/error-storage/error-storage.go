package errorstorage

import (
	"fmt"
	"io/ioutil"
	"os"
)

type ErrorStorage struct {
	limit  byte
	root   string
	Errors []string
}

func New(path string, limit byte) ErrorStorage {
	return ErrorStorage{
		limit: limit,
		root:  "__" + path,
	}
}

func (s ErrorStorage) Save(content []byte) error {
	filePath := s.root + ".txt"

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		if _, err := os.Create(filePath); err != nil {
			return fmt.Errorf("error_storage.Save: error on os.Create -> %v\n", err)
		}
	}

	if err := ioutil.WriteFile(filePath, content, 0644); err != nil {
		return fmt.Errorf("error_storage.Save: error on ioutil.WriteFile -> %v\n", err)
	}

	return nil
}

func (s ErrorStorage) Limit() byte {
	return s.limit
}
