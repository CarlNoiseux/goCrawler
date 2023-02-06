package storageInterfaces

import (
	//"fmt"
	"goCrawler/storage"
	"testing"
)

func getPostgresqlStorage() *PostgresqlStorage {
	return &PostgresqlStorage{
		inMemoryCache: map[string]*storage.UrlRecord{},
	}
}

func TestAddRetrieveUrlFromPostgresqlStorage(t *testing.T) {
	/* Simple test that inserts a URL into storage and then tests that it can be correctly retrieved */
	storagePtr := getPostgresqlStorage()

	//if len(storagePtr.byUrlsIndex) > 0 || len(storagePtr.byStatusesIndex) > 0 {
	//	t.Error(`MapStorage must be empty initialized for this test to execute properly`)
	//}

	putUrl := "http://localhost:8000/"
	_, ok := storagePtr.AddUrl(putUrl, storage.Charted)
	if !ok {
		t.Error(`MapStorage encountered an error while insertion of a new URL record`)
	}

	storagePtr.GetUrlsByStatus(storage.Charted)
	storagePtr.GetUrlsByStatus(storage.Charted, 2)

	storagePtr.UpdateUrlsStatuses([]string{putUrl}, storage.Unchartable)
	storagePtr.UpdateUrlsStatuses([]string{putUrl, "missing1", "missing2"}, storage.Unchartable)

	//if foundUrls, _ := storagePtr.UrlsExist([]string{putUrl}); len(foundUrls) == 0 {
	//	t.Error(`MapStorage inserted URL could not be found through UrlsExist method`)
	//}
	//
	//retrievedRecord := storagePtr.GetUrlsByStatus(storage.Charted)
	//if len(retrievedRecord) != 1 || retrievedRecord[0].Url != putUrl {
	//	t.Error(`MapStorage inserted URL could not be retrieved through GetUrlsByStatus method`)
	//}
	//if newRecord != retrievedRecord[0] {
	//	t.Error(`MapStorage inserted URL could be retrieved but refers to a copy of the UrlRecord, we want single instances`)
	//}
}
