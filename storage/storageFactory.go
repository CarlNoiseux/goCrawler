// Factory pattern in charge of returning appropriate storage instance according to configuration

package storage

import "goCrawler/storage/storageTypes"

// GetStorage currently systematically returns MapStorage instance since it's the only type currently supported.
func GetStorage() *storageTypes.StorageInterface {
	var storagePtr storageTypes.StorageInterface

	storagePtr = MapStorage{byUrlsIndex: map[string]storageTypes.UrlRecord{}, byStatusesIndex: map[storageTypes.UrlExplorationStatus]urlSet{}}

	return &storagePtr
}
