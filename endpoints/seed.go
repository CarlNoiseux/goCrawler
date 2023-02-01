package endpoints

import (
	"fmt"
	"goCrawler/context"
	"goCrawler/storage/storageTypes"
	"net/http"
)

func Seed(ctx context.Context, writer http.ResponseWriter, req *http.Request) {
	url := req.URL.Query().Get("url")

	if len(url) == 0 {
		fmt.Fprintf(writer, "No url specified in \"url\" key")
		return
	}

	(*ctx.Storage).AddUrl(url, storageTypes.Uncharted)

	message := fmt.Sprintf("Added %s url to uncharted frontier.", url)

	fmt.Fprintf(writer, message)
	fmt.Printf(message)
}
