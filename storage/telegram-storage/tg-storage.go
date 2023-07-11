package telegram_storage

import (
	"fmt"
	"io/ioutil"
	"os"
)

type Storage struct {
	root string
}

func New(root string) (*Storage, error) {
	root = "__" + root
	_, err := os.Stat(root)
	if os.IsNotExist(err) {
		if err := os.Mkdir(root, 0755); err != nil {
			return nil, fmt.Errorf("telegram_storage.New: error on mkdir -> %v\n", err)
		}
	}

	return &Storage{
		root: root,
	}, nil
}

func (s Storage) Save(username string, content []byte) error {
	filePath := s.root + "/" + username + ".txt"

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		if _, err := os.Create(filePath); err != nil {
			return fmt.Errorf("telegram_storage.Save: error on os.Create -> %v\n", err)
		}
	}

	if err := ioutil.WriteFile(filePath, content, 0644); err != nil {
		return fmt.Errorf("telegram_storage.Save: error on ioutil.WriteFile -> %v\n", err)
	}

	return nil
}

func (s Storage) Load(username string) ([]byte, error) {
	filePath := s.root + "/" + username + ".txt"

	content, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		return []byte{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("telegram_storage.Load: error on os.ReadFile -> %v\n", err)
	}

	return content, nil
}

func (s Storage) Remove(username string) error {
	filePath := s.root + "/" + username + ".txt"

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("telegram_storage.Remove: error on os.Remove -> %v\n", err)
	}

	return nil
}
