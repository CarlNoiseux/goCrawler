package frontierExplorer

import (
	"goCrawler/storage/storageTypes"
	"log"
	"time"
)

type State string

const (
	Stop State = "stop"
	Run  State = "run"
	Quit State = "quit"

	maxUrlsToExplore int = 10
)

func FrontierExplorer(store *storageTypes.StorageInterface, urlsToExploreChannel chan string, stateChannel chan State) {
	urlsToExplore := make([]*storageTypes.UrlRecord, 0)
	currentState := Run

	for {
		select {
		case <-stateChannel:
			currentState = <-stateChannel
			log.Printf("new explorer state: %s", currentState)
		default:
			// Do other stuff
			switch currentState {
			case Run:
				if len(urlsToExplore) == 0 {
					urlsToExplore = (*store).GetUrlsByStatus(storageTypes.Uncharted, maxUrlsToExplore)
				}

				if len(urlsToExplore) > 0 {
					urlsToExploreChannel <- urlsToExplore[0].Url
					(*store).UpdateUrlStatus(urlsToExplore[0].Url, storageTypes.Charting)

					urlsToExplore = urlsToExplore[1:]
				} else {
					log.Print("no urls to explore found, going to sleep")
					time.Sleep(5 * time.Second)
				}

			case Stop:
				log.Print("explorer is in stopped state, going to sleep")
				time.Sleep(5 * time.Second)

			case Quit:
				log.Print("explorer is in quitting state, terminating goroutine")
				return
			}

		}
	}
}
