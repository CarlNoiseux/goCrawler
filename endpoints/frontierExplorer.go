package endpoints

import (
	"fmt"
	"goCrawler/context"
	"goCrawler/frontierExplorer"
	"net/http"
)

func ExplorerRoot(writer http.ResponseWriter, _ *http.Request) {

	fmt.Fprintf(
		writer,
		"Endpoint to query and manage the explorer state. \n "+
			"Possible states for the explorer are: %q", frontierExplorer.GetPossibleStates())
}

func ExplorerState(ctx context.Context, writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		*ctx.FrontierStateManager <- frontierExplorer.Ping
		state := <-*ctx.FrontierStateManager

		fmt.Fprintf(writer, "Explorer is currently in \"%s\" state.", state)

	case "POST", "PUT":
		newState := req.URL.Query().Get("state")

		state, ok := frontierExplorer.GetStateFromString(newState)
		if !ok {
			fmt.Fprintf(writer, "Invalid state passed to the endpoint: \"%s\". Accepted values are: %q", newState, frontierExplorer.GetPossibleStates())
			return
		}

		*ctx.FrontierStateManager <- state
		fmt.Fprintf(writer, "Explorer is now set to state \"%s\"", newState)
	}
}
