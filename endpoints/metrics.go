package endpoints

import (
	"fmt"
	"goCrawler/context"
	"goCrawler/storage"
	"net/http"
)

func MetricsRoot(writer http.ResponseWriter, _ *http.Request) {

	fmt.Fprintf(
		writer,
		"Endpoint to query and manage the explorer state. \n "+
			"Possible states for the explorer are: %q", storage.GetPossibleExplorationStatusesStrings())
}

func MetricsCounts(ctx context.Context, writer http.ResponseWriter, req *http.Request) {
	status := req.URL.Query().Get("status")

	count := 0
	if len(status) > 0 {
		statusTyped, ok := storage.GetExplorationStatusFromString(status)

		if ok {
			count = (*ctx.Storage).Count(statusTyped)
		} else {
			fmt.Fprintf(writer, "Invalid status passed to the endpoint: \"%s\". Accepted values are: %q", status, storage.GetPossibleExplorationStatusesStrings())
			return
		}
	} else {
		count = (*ctx.Storage).Count()
	}

	fmt.Fprintf(writer, "Found %d urls satisfying search term", count)
}
