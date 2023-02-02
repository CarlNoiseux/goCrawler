package endpoints

import (
	"fmt"
	"goCrawler/storage"
	"net/http"
)

func MetricsRoot(writer http.ResponseWriter, _ *http.Request) {

	fmt.Fprintf(
		writer,
		"Endpoint to query and manage the explorer state. \n "+
			"Possible states for the explorer are: %q", storage.GetPossibleExplorationStatuses())
}
