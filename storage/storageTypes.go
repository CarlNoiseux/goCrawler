// Defines an interface for various storage solutions, so that the crawler can use several implementation in a uniform
// manner.

package storage

type ExplorationStatus string

const (
	Uncharted   ExplorationStatus = "uncharted"
	Charting    ExplorationStatus = "charting"
	Charted     ExplorationStatus = "charted"
	Unchartable ExplorationStatus = "Unchartable"
)

var strToExplorationStatus = map[string]ExplorationStatus{
	string(Uncharted):   Uncharted,
	string(Charting):    Charting,
	string(Charted):     Charted,
	string(Unchartable): Unchartable,
}

//func GetExplorationStatusFromString(state string) (ExplorationStatus, bool) {
//	typedState, ok := strToExplorationStatus[strings.Trim(strings.ToLower(state), " ")]
//
//	return typedState, ok
//}

func GetPossibleExplorationStatuses() []string {
	states := make([]string, len(strToExplorationStatus))
	i := 0
	for k := range strToExplorationStatus {
		states[i] = k
		i++
	}

	return states
}

// UrlRecord definition for a record from the storage in memory
type UrlRecord struct {
	Url    string
	Status ExplorationStatus
}

//// StorageInterface definition for common interface between different storage solutions
//type StorageInterface interface {
//	// AddUrl interface method to insert an url into storage
//	AddUrl(url string, explorationStatus ExplorationStatus) (*UrlRecord, bool)
//
//	// GetUrlsByStatus interface method to retrieve a given number of urls of a given status
//	GetUrlsByStatus(explorationStatus ExplorationStatus, limit ...int) []*UrlRecord
//
//	// UpdateUrlsStatuses interface method to update status of a given list of urls
//	UpdateUrlsStatuses(urls []string, newExplorationStatus ExplorationStatus) ([]*UrlRecord, []string)
//
//	// UpdateUrlStatus interface method to update status of a given list of urls
//	UpdateUrlStatus(urls string, newExplorationStatus ExplorationStatus) (*UrlRecord, error)
//
//	// UrlsExist method to check if some urls already exist within the storage
//	UrlsExist(urls []string) (found []*UrlRecord, missing []string)
//}
