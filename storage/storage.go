package storage

type Storage interface {
	Save(fileName string, content []byte) error
	Load(fileName string) ([]byte, error)
	Remove(fileName string) error
}
