// A "storage" solution that simply stores crawled/to crawl urls in memory.
// Useful for development when no database/cache available, and simpler implementation than storing in files on disc
// at the price of having no persistence.

package storage

import (
	"errors"
	"goCrawler/storage/storageTypes"
)
import "goCrawler/utilities"

// urlSet defines our own set to simplify declarations
type urlSet map[string]bool

// MapStorage defines internal representation to store urls and their statuses.
// modelled with two "indices" so make for easier retrieval through interface methods, at the cost of insertion time (
// similarly to a database indices).
type MapStorage struct {
	//"Primary" key
	byUrlsIndex     map[string]storageTypes.UrlRecord
	byStatusesIndex map[storageTypes.UrlExplorationStatus]urlSet
}

// TODO: add mutex to control flow here, since this is a single global resource instead of an external one like a true database
// TODO: performance will be impacted, but this is more of a mock anyway

func (storage MapStorage) init() {
	if storage.byUrlsIndex == nil {
		storage.byUrlsIndex = make(map[string]storageTypes.UrlRecord)
	}

	if storage.byStatusesIndex == nil {
		storage.byStatusesIndex = make(map[storageTypes.UrlExplorationStatus]urlSet)
	}
}

// WriteUrl see storageTypes.StorageInterface
func (storage MapStorage) WriteUrl(url string, explorationStatus storageTypes.UrlExplorationStatus) (storageTypes.UrlRecord, bool) {
	storage.init()

	record := storageTypes.UrlRecord{Url: url, Status: explorationStatus}

	storage.byUrlsIndex[url] = record

	if _, ok := storage.byStatusesIndex[explorationStatus]; !ok {
		storage.byStatusesIndex[explorationStatus] = make(map[string]bool)
	}
	storage.byStatusesIndex[explorationStatus][url] = true

	return storage.byUrlsIndex[url], true
}

// GetUrlsByStatus see storageTypes.StorageInterface
func (storage MapStorage) GetUrlsByStatus(explorationStatus storageTypes.UrlExplorationStatus, limit ...int) []storageTypes.UrlRecord {
	storage.init()

	numberOfUrls := len(storage.byStatusesIndex[explorationStatus])
	if len(limit) > 0 {
		numberOfUrls = utilities.Min([]int{limit[0], len(storage.byStatusesIndex[explorationStatus])}...)
	}

	// This will do non-deterministic access, mirroring somewhat what would happen if we queried a database, supposing
	// indices change through time
	records := make([]storageTypes.UrlRecord, numberOfUrls)
	for k := range storage.byStatusesIndex[explorationStatus] {
		records = append(records, storage.byUrlsIndex[k])
		numberOfUrls -= 1
		if numberOfUrls == 0 {
			break
		}
	}

	return records
}

// UpdateUrlsStatuses see storageTypes.StorageInterface
func (storage MapStorage) UpdateUrlsStatuses(urls []string, newExplorationStatus storageTypes.UrlExplorationStatus) ([]*storageTypes.UrlRecord, []string) {
	storage.init()

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
			updated = append(updated, &urlRecord)

		} else {
			// Accumulate missing url
			missing = append(missing, url)
		}
	}

	return updated, missing
}

// UpdateUrlStatus see storageTypes.StorageInterface
func (storage MapStorage) UpdateUrlStatus(url string, newExplorationStatus storageTypes.UrlExplorationStatus) (*storageTypes.UrlRecord, error) {
	storage.init()

	records, err := storage.UpdateUrlsStatuses([]string{url}, newExplorationStatus)

	if len(err) > 0 {
		return nil, errors.New("could not find requested url in storage")
	}

	return records[0], nil
}
