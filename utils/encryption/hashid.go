package encryption

import "github.com/speps/go-hashids"

// NewHasher creates a new hasher for decoding and encoding int64 slices
func NewHasher(salt string) (*hashids.HashID, error) {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = 30

	return hashids.NewWithData(hd)
}
