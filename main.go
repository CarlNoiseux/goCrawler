// Academic project to implement a web crawler in golang.
// Sets up a REST web server to accept seeds from clients and store these seeds into a storage

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", root)
	mux.HandleFunc("/seed", seed)

	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

func root(responseWriter http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(responseWriter, "goCrawler API is up and running")
}

func seed(_ http.ResponseWriter, _ *http.Request) {

}
