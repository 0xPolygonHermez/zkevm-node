package db

// KeyValuer interface that a compatible DB backend should implement
type KeyValuer interface {
	Put(key string, data []byte) error
	// Get IMPORTANT: should return nil, nil in case of no data found!!!!
	Get(key string) ([]byte, error)
}
