package storage

import "goCrawler/storage/storageTypes"
import "goCrawler/utilities"

type urlSet map[string]bool

type MapStorage struct {
	//"Primary" key
	byUrlsIndex     map[string]storageTypes.UrlRecord
	byStatusesIndex map[storageTypes.UrlExplorationStatus]urlSet
}

func (storage MapStorage) WriteUrl(url string, explorationStatus storageTypes.UrlExplorationStatus) (storageTypes.UrlRecord, bool) {
	storage.byUrlsIndex[url] = storageTypes.UrlRecord{Url: url, Status: explorationStatus}
	storage.byStatusesIndex[explorationStatus][url] = true

	return storage.byUrlsIndex[url], true
}

func (storage MapStorage) GetUrls(explorationStatus storageTypes.UrlExplorationStatus, limit ...int) []storageTypes.UrlRecord {
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

func (storage MapStorage) UpdateUrlsStatuses(urls []string, newExplorationStatus storageTypes.UrlExplorationStatus) ([]*storageTypes.UrlRecord, []string) {
	missing := make([]string, 0)
	updated := make([]*storageTypes.UrlRecord, len(urls))
	for _, url := range urls {
		urlRecord, ok := storage.byUrlsIndex[url]
		if ok {
			// Update records with new status
			urlRecord.Status = newExplorationStatus
			updated = append(updated, &urlRecord)

			// update "index", to simulate an index re-balancing
			delete(storage.byStatusesIndex[urlRecord.Status], url)
			storage.byStatusesIndex[newExplorationStatus][url] = true
		} else {
			// Accumulate missing url
			missing = append(missing, url)
		}
	}

	return updated, missing
}
