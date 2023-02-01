// Academic project to implement a web crawler in golang.
// Sets up a REST web server to accept seeds from clients and store these seeds into a storage

package main

import (
	"log"
	"net/http"

	"goCrawler/context"
	"goCrawler/endpoints"
	"goCrawler/frontierExplorer"
	"goCrawler/parser"
	"goCrawler/storage"
)

func main() {
	storagePtr := storage.GetStoragePtr()
	explorerStateController := make(chan frontierExplorer.State)

	// Create context to pass to other processes
	ctx := context.Context{Storage: storagePtr, FrontierStateManager: explorerStateController}

	mux := http.NewServeMux()

	mux.HandleFunc("/", endpoints.Root)

	urlsToExploreChannel := make(chan string, 10)

	// Create a goroutine in charge of "exploring" the urls that have not been charted in the storage yet
	go frontierExplorer.FrontierExplorer(ctx.Storage, urlsToExploreChannel, explorerStateController)

	// Could scale up here, by creating several goroutine that consume from the urlsToExploreChannel channel
	go parser.ParsePageUrls(ctx, urlsToExploreChannel)

	// Wrap the handle functions to pass a context manager containing global settings, storage, etc... to endpoints
	mux.HandleFunc("/seed", func(w http.ResponseWriter, r *http.Request) {
		endpoints.Seed(ctx, w, r)
	})

	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}
