// A storage solution that wraps implementation using postgresql DBMS.
// Unlike MapStorage this provides a truly persistent solution without keeping everything in memory

package storageInterfaces

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"goCrawler/storage"
	"os"
)

type PostgresqlStorage struct {
	dbEngine *sql.DB

	// Sorta going to function like an ORM that keeps loaded instances in memory.
	inMemoryCache map[string]*storage.UrlRecord
}

// Define helper function for database connection
func (storagePtr *PostgresqlStorage) getConnectionString() string {
	host := os.Getenv("POSTGRESQL_HOST")
	port := os.Getenv("POSTGRESQL_PORT")
	user := os.Getenv("POSTGRESQL_USER")
	pwd := os.Getenv("POSTGRESQL_PWD")
	dbname := os.Getenv("POSTGRESQL_DBNAME")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pwd, dbname)
}

func (storagePtr *PostgresqlStorage) getDatabaseConnection() *sql.DB {
	if storagePtr.dbEngine != nil {
		return storagePtr.dbEngine
	}

	connectionString := storagePtr.getConnectionString()

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	storagePtr.dbEngine = db
	return storagePtr.dbEngine
}

// Define statement management helper functions
func (storagePtr *PostgresqlStorage) executeStatement(statement string, args ...any) error {
	storagePtr.getDatabaseConnection()

	_, err := storagePtr.dbEngine.Exec(statement, args...)

	return err
}

func (storagePtr *PostgresqlStorage) queryStatement(statement string, args ...any) []*storage.UrlRecord {
	storagePtr.getDatabaseConnection()

	rows, err := storagePtr.dbEngine.Query(statement, args...)
	defer rows.Close()

	if err != nil {
		panic(err)
	}

	records := make([]*storage.UrlRecord, 0)
	for rows.Next() {
		record := storage.UrlRecord{}
		rows.Scan(&record.Url, &record.Status)

		storagePtr.inMemoryCache[record.Url] = &record

		records = append(records, &record)
	}

	return records
}

// AddUrl see storageTypes.StorageInterface
func (storagePtr *PostgresqlStorage) AddUrl(url string, explorationStatus storage.ExplorationStatus) (*storage.UrlRecord, bool) {
	storagePtr.executeStatement(`INSERT INTO frontier (url, status) VALUES ($1, $2)`, url, string(explorationStatus))

	record := storage.UrlRecord{Url: url, Status: explorationStatus}

	storagePtr.inMemoryCache[url] = &record

	return storagePtr.inMemoryCache[url], true
}

// GetUrlsByStatus see storageTypes.StorageInterface
func (storagePtr *PostgresqlStorage) GetUrlsByStatus(explorationStatus storage.ExplorationStatus, limit ...int) []*storage.UrlRecord {

	statement := `
	SELECT *
	FROM frontier
	WHERE status = $1
	`

	//records := []*storage.UrlRecord
	var records []*storage.UrlRecord
	if len(limit) > 0 {
		statement += "\nLIMIT $2"
		records = storagePtr.queryStatement(statement, string(explorationStatus), limit[0])
	} else {
		records = storagePtr.queryStatement(statement, string(explorationStatus))
	}

	return records
}

// UpdateUrlsStatuses see storageTypes.StorageInterface
func (storagePtr *PostgresqlStorage) UpdateUrlsStatuses(urls []string, newExplorationStatus storage.ExplorationStatus) ([]*storage.UrlRecord, []string) {
	accumulator := ""
	for index := range urls[:len(urls)-1] {
		accumulator += fmt.Sprintf("$%d, ", index+1)
	}
	accumulator += fmt.Sprintf("$%d", len(urls))

	statement2 := "UPDATE frontier SET status = '%s' WHERE url IN (%s)"

	statement := fmt.Sprintf("SELECT * FROM frontier WHERE url IN (%s)", accumulator)

	urlsAny := make([]interface{}, len(urls))
	for i, v := range urls {
		urlsAny[i] = v
	}

	updated := storagePtr.queryStatement(statement, urlsAny...)
	updatedSet := map[string]bool{}
	for _, key := range updated {
		updatedSet[key.Url] = true
	}

	missing := make([]string, 0)
	for _, url := range urls {
		if _, ok := updatedSet[url]; ok {
			continue
		}

		missing = append(missing, url)
	}

	return updated, missing
}

// UpdateUrlStatus see storageTypes.StorageInterface
func (storagePtr *PostgresqlStorage) UpdateUrlStatus(url string, newExplorationStatus storage.ExplorationStatus) (*storage.UrlRecord, error) {
	records, err := storagePtr.UpdateUrlsStatuses([]string{url}, newExplorationStatus)

	if len(err) > 0 {
		return nil, errors.New("could not find requested url in storage")
	}

	return records[0], nil
}

// UrlsExist see storageTypes.StorageInterface
func (storagePtr *PostgresqlStorage) UrlsExist(urls []string) ([]*storage.UrlRecord, []string) {
	accumulator := ""
	for index := range urls[:len(urls)-1] {
		accumulator += fmt.Sprintf("$%d, ", index+1)
	}
	accumulator += fmt.Sprintf("$%d", len(urls))

	statement := fmt.Sprintf("SELECT * FROM frontier WHERE url IN (%s)", accumulator)

	urlsAny := make([]interface{}, len(urls))
	for i, v := range urls {
		urlsAny[i] = v
	}

	updated := storagePtr.queryStatement(statement, urlsAny...)
	updatedSet := map[string]bool{}
	for _, key := range updated {
		updatedSet[key.Url] = true
	}

	missing := make([]string, 0)
	for _, url := range urls {
		if _, ok := updatedSet[url]; ok {
			continue
		}

		missing = append(missing, url)
	}

	return found, missing
}

//func (storagePtr *PostgresqlStorage) Count(statuses ...storage.ExplorationStatus) int {
//	count := 0
//
//	_statuses := make([]storage.ExplorationStatus, 0)
//	if len(statuses) > 0 {
//		for _, status := range statuses {
//			_statuses = append(_statuses, status)
//		}
//	} else {
//		_statuses = storage.GetPossibleExplorationStatuses()
//	}
//
//	for _, storageType := range _statuses {
//		count += len(storagePtr.byStatusesIndex[storageType])
//	}
//
//	return count
//}

// TODO: hook to storageFactory
