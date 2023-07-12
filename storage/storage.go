package storage

type Storage interface {
	Save(username string, content []byte) error
	Load(username string) ([]byte, error)
	Remove(username string) error
}
