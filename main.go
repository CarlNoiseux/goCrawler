// Academic project to implement a web crawler in golang.
// Sets up a REST web server to accept seeds from clients and store these seeds into a storage

package main

import (
	"log"
	"net/http"

	"goCrawler/context"
	"goCrawler/endpoints"
	"goCrawler/parser"
	"goCrawler/storage"
	"goCrawler/storage/storageTypes"
)

func main() {
	storagePtr := storage.GetStorage()

	// Create context to pass to other processes
	ctx := context.Context{Storage: storagePtr}

	mux := http.NewServeMux()

	mux.HandleFunc("/", endpoints.Root)

	// Wrap the handle function to pass a (potentially) global reference to a storage solution.
	// Might generalize this to pass a context instead, in order to wrap more stuff in it if we desire
	mux.HandleFunc("/seed", func(w http.ResponseWriter, r *http.Request) {
		endpoints.Seed(ctx, w, r)
	})

	seedUrl := "https://www.lapresse.ca/"
	(*storagePtr).WriteUrl(seedUrl, storageTypes.Uncharted)
	go parser.ParsePageUrls(ctx, seedUrl)

	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}
