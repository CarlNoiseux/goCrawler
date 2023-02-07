package storageInterfaces

import (
	"goCrawler/storage"
)

// StorageInterface definition for common interface between different storage solutions
type StorageInterface interface {
	// AddUrl interface method to insert an url into storage
	AddUrl(url string, explorationStatus storage.ExplorationStatus) (*storage.UrlRecord, bool)

	// GetUrlsByStatus interface method to retrieve a given number of urls of a given status
	GetUrlsByStatus(explorationStatus storage.ExplorationStatus, limit ...int) []*storage.UrlRecord

	// UpdateUrlsStatuses interface method to update status of a given list of urls
	UpdateUrlsStatuses(urls []string, newExplorationStatus storage.ExplorationStatus) error

	// UpdateUrlStatus interface method to update status of a given list of urls
	UpdateUrlStatus(urls string, newExplorationStatus storage.ExplorationStatus) error

	// UrlsExist method to check if some urls already exist within the storage
	UrlsExist(urls []string) (found []*storage.UrlRecord, missing []string)

	// Count the number of urls for given statuses in the storage
	Count(statuses ...storage.ExplorationStatus) int
}
