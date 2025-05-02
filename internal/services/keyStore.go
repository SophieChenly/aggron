package services

import (
	"errors"
)

func StoreFileKey(fileKey FileKey) error {
	// TODO: Store file key into cache
	return errors.New("FileKey storage not implemented")
}

func GetFileKey(fileID string) (FileKey, error) {
	// TODO: Retrieve file key and stuff
	return FileKey{"", []byte{}, []byte{}, "", ""}, errors.New("FileKey retrieval not implemented")
}
