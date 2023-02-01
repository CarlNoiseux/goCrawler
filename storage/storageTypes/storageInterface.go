// Defines an interface for various storage solutions, so that the crawler can use several implementation in a uniform
// manner.

package storageTypes

type UrlExplorationStatus string

const (
	Uncharted UrlExplorationStatus = "uncharted"
	Charting  UrlExplorationStatus = "charting"
	Charted   UrlExplorationStatus = "charted"
)

// UrlRecord definition for a record from the storage in memory
type UrlRecord struct {
	Url    string
	Status UrlExplorationStatus
}

// StorageInterface definition for common interface between different storage solutions
type StorageInterface interface {
	// WriteUrl interface method to insert an url into storage
	WriteUrl(url string, explorationStatus UrlExplorationStatus) (UrlRecord, bool)

	// GetUrlsByStatus interface method to retrieve a given number of urls of a given status
	GetUrlsByStatus(explorationStatus UrlExplorationStatus, limit ...int) []UrlRecord

	// UpdateUrlsStatuses interface method to update status of a given list of urls
	UpdateUrlsStatuses(urls []string, newExplorationStatus UrlExplorationStatus) ([]*UrlRecord, []string)

	// UpdateUrlStatus interface method to update status of a given list of urls
	UpdateUrlStatus(urls string, newExplorationStatus UrlExplorationStatus) (*UrlRecord, error)
}
