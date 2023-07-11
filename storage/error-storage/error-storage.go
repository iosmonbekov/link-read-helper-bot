package errorstorage

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type ErrorStorage struct {
	limit  byte
	root   string
	Errors []error
}

func New(path string, limit byte) ErrorStorage {
	return ErrorStorage{
		limit:  limit,
		root:   "__" + path,
		Errors: []error{},
	}
}

func (s ErrorStorage) Save() error {
	filePath := s.root + ".txt"

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		if _, err := os.Create(filePath); err != nil {
			return fmt.Errorf("error_storage.Save: error on os.Create -> %v\n", err)
		}
	}

	content := []string{}

	for _, e := range s.Errors {
		content = append(content, e.Error())
	}

	if err := ioutil.WriteFile(filePath, []byte(strings.Join(content, "\n")), 0644); err != nil {
		return fmt.Errorf("error_storage.Save: error on ioutil.WriteFile -> %v\n", err)
	}

	return nil
}

func (s ErrorStorage) Append(err error) {
	if len(s.Errors) == 0 {
		s.Errors = append([]error{}, err)
		return
	}

	if s.Errors[0].Error() == err.Error() {
		s.Errors = append(s.Errors, err)
		return
	} else {
		s.Errors = append([]error{}, err)
	}
}

func (s ErrorStorage) Limit() byte {
	return s.limit
}
