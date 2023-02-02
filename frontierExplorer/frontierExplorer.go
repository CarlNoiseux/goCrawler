package frontierExplorer

import (
	"goCrawler/storage/storageTypes"
	"log"
	"strings"
	"time"
)

type State string

const (
	Stop State = "stop"
	Run  State = "run"
	Quit State = "quit"
	Ping State = "ping"

	maxUrlsToExplore int = 10
)

var strToState = map[string]State{
	string(Stop): Stop,
	string(Run):  Run,
	string(Quit): Quit,
	string(Ping): Ping,
}

func GetStateFromString(state string) (State, bool) {
	typedState, ok := strToState[strings.Trim(strings.ToLower(state), " ")]

	return typedState, ok
}

func GetPossibleStates() []string {
	states := make([]string, len(strToState))
	i := 0
	for k := range strToState {
		states[i] = k
		i++
	}

	return states
}

func FrontierExplorer(store *storageTypes.StorageInterface, urlsToExploreChannel *chan string, stateChannel *chan State) {
	urlsToExplore := make([]*storageTypes.UrlRecord, 0)
	previousState, currentState := Run, Run

	for {
		select {
		case newState := <-*stateChannel:
			previousState = currentState
			currentState = newState
			log.Printf("new explorer state: %s", currentState)
		default:
			// Do other stuff
			switch currentState {
			case Ping:
				*stateChannel <- previousState
				currentState = previousState

			case Run:
				if len(urlsToExplore) == 0 {
					urlsToExplore = (*store).GetUrlsByStatus(storageTypes.Uncharted, maxUrlsToExplore)
				}

				if len(urlsToExplore) > 0 {
					*urlsToExploreChannel <- urlsToExplore[0].Url
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
