// A "storage" solution that simply stores crawled/to crawl urls in memory.
// Useful for development when no database/cache available, and simpler implementation than storing in files on disc
// at the price of having no persistence.

package storageInterfaces

import (
	"errors"
	"goCrawler/storage"
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
	byUrlsIndex map[string]*storage.UrlRecord

	// Index by status, to facilitate searches
	byStatusesIndex map[storage.ExplorationStatus]urlSet
}

// AddUrl see storageTypes.StorageInterface
func (storagePtr *MapStorage) AddUrl(url string, explorationStatus storage.ExplorationStatus) (*storage.UrlRecord, bool) {
	storagePtr.mutex.Lock()

	record := storage.UrlRecord{Url: url, Status: explorationStatus}
	//storage.urls = append(storage.urls, record)

	storagePtr.byUrlsIndex[url] = &record

	if _, ok := storagePtr.byStatusesIndex[explorationStatus]; !ok {
		storagePtr.byStatusesIndex[explorationStatus] = make(map[string]bool)
	}
	storagePtr.byStatusesIndex[explorationStatus][url] = true

	storagePtr.mutex.Unlock()

	return storagePtr.byUrlsIndex[url], true
}

// GetUrlsByStatus see storageTypes.StorageInterface
func (storagePtr *MapStorage) GetUrlsByStatus(explorationStatus storage.ExplorationStatus, limit ...int) []*storage.UrlRecord {
	storagePtr.mutex.Lock()

	numberOfUrls := len(storagePtr.byStatusesIndex[explorationStatus])
	if len(limit) > 0 {
		numberOfUrls = utilities.Min([]int{limit[0], len(storagePtr.byStatusesIndex[explorationStatus])}...)
	}

	// This will do non-deterministic access, mirroring somewhat what would happen if we queried a database, supposing
	// indices change through time

	records := make([]*storage.UrlRecord, numberOfUrls)
	index := 0
	for k := range storagePtr.byStatusesIndex[explorationStatus] {
		record := storagePtr.byUrlsIndex[k]
		records[index] = record
		index += 1
		if numberOfUrls == index {
			break
		}
	}

	storagePtr.mutex.Unlock()

	return records
}

// UpdateUrlsStatuses see storageTypes.StorageInterface
func (storagePtr *MapStorage) UpdateUrlsStatuses(urls []string, newExplorationStatus storage.ExplorationStatus) ([]*storage.UrlRecord, []string) {
	storagePtr.mutex.Lock()

	missing := make([]string, 0)
	updated := make([]*storage.UrlRecord, len(urls))
	for _, url := range urls {
		urlRecord, ok := storagePtr.byUrlsIndex[url]
		if ok {
			// update "index", to simulate an index re-balancing
			delete(storagePtr.byStatusesIndex[urlRecord.Status], url)

			if _, ok := storagePtr.byStatusesIndex[newExplorationStatus]; !ok {
				storagePtr.byStatusesIndex[newExplorationStatus] = make(map[string]bool)
			}
			storagePtr.byStatusesIndex[newExplorationStatus][url] = true

			// Update records with new status
			urlRecord.Status = newExplorationStatus
			storagePtr.byUrlsIndex[url] = urlRecord

			updated = append(updated, urlRecord)

		} else {
			// Accumulate missing url
			missing = append(missing, url)
		}
	}

	storagePtr.mutex.Unlock()

	return updated, missing
}

// UpdateUrlStatus see storageTypes.StorageInterface
func (storagePtr *MapStorage) UpdateUrlStatus(url string, newExplorationStatus storage.ExplorationStatus) (*storage.UrlRecord, error) {
	records, err := storagePtr.UpdateUrlsStatuses([]string{url}, newExplorationStatus)

	if len(err) > 0 {
		return nil, errors.New("could not find requested url in storage")
	}

	return records[0], nil
}

// UrlsExist see storageTypes.StorageInterface
func (storagePtr *MapStorage) UrlsExist(urls []string) ([]*storage.UrlRecord, []string) {
	found := make([]*storage.UrlRecord, 0)
	missing := make([]string, 0)
	for _, url := range urls {
		if record, ok := storagePtr.byUrlsIndex[url]; ok {
			found = append(found, record)
		} else {
			missing = append(missing, url)
		}
	}

	return found, missing
}

func (storagePtr *MapStorage) Count(statuses ...storage.ExplorationStatus) int {
	count := 0

	_statuses := make([]storage.ExplorationStatus, 0)
	if len(statuses) > 0 {
		for _, status := range statuses {
			_statuses = append(_statuses, status)
		}
	} else {
		_statuses = storage.GetPossibleExplorationStatuses()
	}

	for _, storageType := range _statuses {
		count += len(storagePtr.byStatusesIndex[storageType])
	}

	return count
}
