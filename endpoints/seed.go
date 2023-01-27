package endpoints

import (
	"fmt"
	"goCrawler/storage/storageTypes"
	"net/http"
)

func Seed(writer http.ResponseWriter, req *http.Request, storageInterface *storageTypes.StorageInterface) {
	url := req.URL.Query().Get("url")

	//var arr []string
	//_ = json.Unmarshal([]byte(urlsString), &arr)
	if len(url) == 0 {
		fmt.Fprintf(writer, "No url specified in \"url\" key")
		return
	}

	(*storageInterface).WriteUrl(url, storageTypes.Uncharted)

	fmt.Fprintf(writer, "Added %s url to uncharted frontier.", url)
}