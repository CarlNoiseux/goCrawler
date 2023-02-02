// Factory pattern in charge of returning appropriate storage instance according to configuration

package storageInterfaces

import (
	"goCrawler/storage"
	"sync"
)

// GetStoragePtr currently systematically returns MapStorage instance since it's the only type currently supported.
func GetStoragePtr() *StorageInterface {
	var storageI StorageInterface

	storageI = &MapStorage{
		//urls:            make([]storageTypes.UrlRecord, 0),
		byUrlsIndex:     map[string]*storage.UrlRecord{},
		byStatusesIndex: map[storage.ExplorationStatus]urlSet{},
		mutex:           &sync.Mutex{},
	}

	return &storageI
}
