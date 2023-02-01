// A "storage" solution that simply stores crawled/to crawl urls in memory.
// Useful for development when no database/cache available, and simpler implementation than storing in files on disc
// at the price of having no persistence.

package storage

import (
	"errors"
	"goCrawler/storage/storageTypes"
	"sync"
)
import "goCrawler/utilities"

// urlSet defines our own set to simplify declarations
type urlSet map[string]bool

// MapStorage defines internal representation to store urls and their statuses.
// modelled with two "indices" so make for easier retrieval through interface methods, at the cost of insertion time (
// similarly to a database indices).
type MapStorage struct {
	mutex *sync.Mutex

	// "Primary" key
	byUrlsIndex map[string]*storageTypes.UrlRecord

	// Index by status, to facilitate searches
	byStatusesIndex map[storageTypes.UrlExplorationStatus]urlSet
}

// AddUrl see storageTypes.StorageInterface
func (storage *MapStorage) AddUrl(url string, explorationStatus storageTypes.UrlExplorationStatus) (*storageTypes.UrlRecord, bool) {
	storage.mutex.Lock()

	record := storageTypes.UrlRecord{Url: url, Status: explorationStatus}
	//storage.urls = append(storage.urls, record)

	storage.byUrlsIndex[url] = &record

	if _, ok := storage.byStatusesIndex[explorationStatus]; !ok {
		storage.byStatusesIndex[explorationStatus] = make(map[string]bool)
	}
	storage.byStatusesIndex[explorationStatus][url] = true

	storage.mutex.Unlock()

	return storage.byUrlsIndex[url], true
}

// GetUrlsByStatus see storageTypes.StorageInterface
func (storage *MapStorage) GetUrlsByStatus(explorationStatus storageTypes.UrlExplorationStatus, limit ...int) []*storageTypes.UrlRecord {
	storage.mutex.Lock()

	numberOfUrls := len(storage.byStatusesIndex[explorationStatus])
	if len(limit) > 0 {
		numberOfUrls = utilities.Min([]int{limit[0], len(storage.byStatusesIndex[explorationStatus])}...)
	}

	// This will do non-deterministic access, mirroring somewhat what would happen if we queried a database, supposing
	// indices change through time
	records := make([]*storageTypes.UrlRecord, numberOfUrls)
	for k := range storage.byStatusesIndex[explorationStatus] {
		record := storage.byUrlsIndex[k]
		records = append(records, record)
		numberOfUrls -= 1
		if numberOfUrls == 0 {
			break
		}
	}

	storage.mutex.Unlock()

	return records
}

// UpdateUrlsStatuses see storageTypes.StorageInterface
func (storage *MapStorage) UpdateUrlsStatuses(urls []string, newExplorationStatus storageTypes.UrlExplorationStatus) ([]*storageTypes.UrlRecord, []string) {
	storage.mutex.Lock()

	missing := make([]string, 0)
	updated := make([]*storageTypes.UrlRecord, len(urls))
	for _, url := range urls {
		urlRecord, ok := storage.byUrlsIndex[url]
		if ok {
			// update "index", to simulate an index re-balancing
			delete(storage.byStatusesIndex[urlRecord.Status], url)

			if _, ok := storage.byStatusesIndex[newExplorationStatus]; !ok {
				storage.byStatusesIndex[newExplorationStatus] = make(map[string]bool)
			}
			storage.byStatusesIndex[newExplorationStatus][url] = true

			// Update records with new status
			urlRecord.Status = newExplorationStatus
			storage.byUrlsIndex[url] = urlRecord

			updated = append(updated, urlRecord)

		} else {
			// Accumulate missing url
			missing = append(missing, url)
		}
	}

	storage.mutex.Unlock()

	return updated, missing
}

// UpdateUrlStatus see storageTypes.StorageInterface
func (storage *MapStorage) UpdateUrlStatus(url string, newExplorationStatus storageTypes.UrlExplorationStatus) (*storageTypes.UrlRecord, error) {
	records, err := storage.UpdateUrlsStatuses([]string{url}, newExplorationStatus)

	if len(err) > 0 {
		return nil, errors.New("could not find requested url in storage")
	}

	return records[0], nil
}

// UrlsExist see storageTypes.StorageInterface
func (storage *MapStorage) UrlsExist(urls []string) ([]*storageTypes.UrlRecord, []string) {
	found := make([]*storageTypes.UrlRecord, 0)
	missing := make([]string, 0)
	for _, url := range urls {
		if record, ok := storage.byUrlsIndex[url]; ok {
			found = append(found, record)
		} else {
			missing = append(missing, url)
		}
	}

	return found, missing
}
