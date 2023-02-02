// Defines an interface for various storage solutions, so that the crawler can use several implementation in a uniform
// manner.

package storage

import "strings"

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

func GetExplorationStatusFromString(status string) (ExplorationStatus, bool) {
	typedStatus, ok := strToExplorationStatus[strings.Trim(strings.ToLower(status), " ")]

	return typedStatus, ok
}

func GetPossibleExplorationStatusesStrings() []string {
	states := make([]string, len(strToExplorationStatus))
	i := 0
	for k := range strToExplorationStatus {
		states[i] = k
		i++
	}

	return states
}

func GetPossibleExplorationStatuses() []ExplorationStatus {
	states := make([]ExplorationStatus, len(strToExplorationStatus))
	i := 0
	for _, v := range strToExplorationStatus {
		states[i] = v
		i++
	}

	return states
}

// UrlRecord definition for a record from the storage in memory
type UrlRecord struct {
	Url    string
	Status ExplorationStatus
}
