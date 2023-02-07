// These tests should be run against a test database, since they may truncate the frontier table.
// This is bad design, but since I did not implement the PostgresqlStorage type with transaction management
// changes are systematically committed, meaning even a rollback does not undo them.
// Going to truncate table after tests, to make them idempotent

package storageInterfaces

import (
	"fmt"

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
	storagePtr.getDatabaseConnection()

	// See comment at top of package file
	defer func() {
		storagePtr.executeStatement("TRUNCATE TABLE frontier")
	}()

	if len(storagePtr.inMemoryCache) > 0 {
		t.Error(`PostgresqlStorage must be empty initialized for this test to execute properly`)
	}

	putUrl := "http://localhost:8000/"
	_, ok := storagePtr.AddUrl(putUrl, storage.Charted)
	if !ok {
		t.Error(`PostgresqlStorage encountered an error while insertion of a new URL record`)
	}

	records := storagePtr.GetUrlsByStatus(storage.Charted, 2)
	if len(records) != 1 {
		t.Error(`PostgresqlStorage did not properly add url through AddUrl`)
	}

	err := storagePtr.UpdateUrlsStatuses([]string{putUrl, "missing1", "missing2"}, storage.Unchartable)
	if err != nil {
		t.Error(`PostgresqlStorage encountered an error while updating url status`)
	}

	if foundUrls, _ := storagePtr.UrlsExist([]string{putUrl}); len(foundUrls) == 0 {
		t.Error(`PostgresqlStorage inserted URL could not be found through UrlsExist method`)
	}
}

func TestCountUrlsFromPostgresqlStorage(t *testing.T) {
	storagePtr := getPostgresqlStorage()
	storagePtr.getDatabaseConnection()

	// See comment at top of package file
	defer func() {
		storagePtr.executeStatement("TRUNCATE TABLE frontier")
	}()

	if storagePtr.Count() != 0 {
		t.Error(`Postgresql count is different from zero when it should not`)
	}

	explorationStatuses := storage.GetPossibleExplorationStatuses()

	for index, storageType := range explorationStatuses {
		storagePtr.AddUrl(fmt.Sprintf("http://localhost:%d/", index), storageType)
	}

	if storagePtr.Count() != len(explorationStatuses) {
		t.Error(`Postgresql Count for all statuses returns an incorrect count`)
	}

	for _, storageType := range explorationStatuses {
		if storagePtr.Count(storageType) != 1 {
			t.Error(fmt.Sprintf(`Postgresql Count for status %s returns an incorrect count`, storageType))
		}
	}
}
