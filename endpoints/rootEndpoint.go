package endpoints

import (
	"fmt"
	"net/http"
)

func Root(responseWriter http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(responseWriter, "goCrawler API is up and running")
}
