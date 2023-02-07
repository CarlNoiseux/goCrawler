// Factory pattern in charge of returning appropriate storage instance according to configuration

package storageInterfaces

import (
	"goCrawler/storage"
	"os"
	"sync"
)

type StorageType string

const (
	Map        StorageType = "map"
	Postgresql StorageType = "postgresql"
)

// GetStoragePtr currently systematically returns MapStorage instance since it's the only type currently supported.
func GetStoragePtr() *StorageInterface {
	// Retrieving factory setting through environment variable for ease of use
	storageSetting := Map
	if os.Getenv("STORAGE_TYPE") == string(Postgresql) {
		storageSetting = Postgresql
	}

	var storageI StorageInterface

	if storageSetting == Map {
		storageI = &MapStorage{
			byUrlsIndex:     map[string]*storage.UrlRecord{},
			byStatusesIndex: map[storage.ExplorationStatus]urlSet{},
			mutex:           &sync.Mutex{},
		}
	} else {
		storageI = &PostgresqlStorage{
			inMemoryCache: map[string]*storage.UrlRecord{},
		}
	}

	return &storageI
}
