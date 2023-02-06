package storageInterfaces

import (
	"fmt"
	"goCrawler/storage"
	"sync"
	"testing"
)

func getMapStorage() *MapStorage {
	return &MapStorage{
		byUrlsIndex:     map[string]*storage.UrlRecord{},
		byStatusesIndex: map[storage.ExplorationStatus]urlSet{},
		mutex:           &sync.Mutex{},
	}
}

func TestAddRetrieveUrlFromMapStorage(t *testing.T) {
	/*Simple test that inserts a URL into storage and then tests that it can be correctly retrieved */
	storagePtr := getMapStorage()

	if len(storagePtr.byUrlsIndex) > 0 || len(storagePtr.byStatusesIndex) > 0 {
		t.Error(`MapStorage must be empty initialized for this test to execute properly`)
	}

	putUrl := "http://localhost:8000/"
	newRecord, ok := storagePtr.AddUrl(putUrl, storage.Charted)
	if !ok {
		t.Error(`MapStorage encountered an error while insertion of a new URL record`)
	}

	if foundUrls, _ := storagePtr.UrlsExist([]string{putUrl}); len(foundUrls) == 0 {
		t.Error(`MapStorage inserted URL could not be found through UrlsExist method`)
	}

	retrievedRecord := storagePtr.GetUrlsByStatus(storage.Charted)
	if len(retrievedRecord) != 1 || retrievedRecord[0].Url != putUrl {
		t.Error(`MapStorage inserted URL could not be retrieved through GetUrlsByStatus method`)
	}
	if newRecord != retrievedRecord[0] {
		t.Error(`MapStorage inserted URL could be retrieved but refers to a copy of the UrlRecord, we want single instances`)
	}
}

func TestCountUrlsFromMapStorage(t *testing.T) {
	storagePtr := getMapStorage()

	if storagePtr.Count() != 0 {
		t.Error(`MapStorage count is different from zero when it should not`)
	}

	explorationStatuses := storage.GetPossibleExplorationStatuses()

	for index, storageType := range explorationStatuses {
		storagePtr.AddUrl(fmt.Sprintf("http://localhost:%d/", index), storageType)
	}

	if storagePtr.Count() != len(explorationStatuses) {
		t.Error(`MapStorage Count for all statuses returns an incorrect count`)
	}

	for _, storageType := range explorationStatuses {
		if storagePtr.Count(storageType) != 1 {
			t.Error(fmt.Sprintf(`MapStorage Count for status %s returns an incorrect count`, storageType))
		}
	}
}
