// Defines an interface

package storageTypes

type UrlExplorationStatus string

const (
	Uncharted UrlExplorationStatus = "uncharted"
	Charting  UrlExplorationStatus = "charting"
	Charted   UrlExplorationStatus = "charted"
)

type UrlRecord struct {
	Url    string
	Status UrlExplorationStatus
}

type StorageInterface interface {
	WriteUrl(url, explorationStatus UrlExplorationStatus) (UrlRecord, bool)
	GetUrlsByStatus(explorationStatus UrlExplorationStatus, limit ...int) []UrlRecord
	UpdateUrlsStatuses(urls []string, newExplorationStatus UrlExplorationStatus) ([]*UrlRecord, []string)
}
