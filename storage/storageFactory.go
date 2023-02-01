// Factory pattern in charge of returning appropriate storage instance according to configuration

package storage

import (
	"goCrawler/storage/storageTypes"
	"sync"
)

// GetStoragePtr currently systematically returns MapStorage instance since it's the only type currently supported.
func GetStoragePtr() *storageTypes.StorageInterface {
	var storage storageTypes.StorageInterface

	storage = &MapStorage{
		//urls:            make([]storageTypes.UrlRecord, 0),
		byUrlsIndex:     map[string]*storageTypes.UrlRecord{},
		byStatusesIndex: map[storageTypes.UrlExplorationStatus]urlSet{},
		mutex:           &sync.Mutex{},
	}

	return &storage
}
