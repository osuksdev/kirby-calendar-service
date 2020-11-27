package calendar

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/osuksdev/kirby-calendar-service/cors"
)

const calendarPath = "calendar"

func handleCalendar(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		eventList := getEventList()
		j, err := json.Marshal(eventList)
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodPost:
		var event Event
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = addOrUpdateEvent(event)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleEvent(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", calendarPath))
	if len(urlPathSegments[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	eventID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		event := getEvent(eventID)
		if event == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		j, err := json.Marshal(event)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}

	case http.MethodPost:
		var event Event
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if event.EventID != eventID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = addOrUpdateEvent(event)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodDelete:
		removeEvent(eventID)

	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// SetupRoutes :
func SetupRoutes(apiBasePath string) {
	calendarHandler := http.HandlerFunc(handleCalendar)
	eventHandler := http.HandlerFunc(handleEvent)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, calendarPath), cors.Middleware(calendarHandler))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, calendarPath), cors.Middleware(eventHandler))
}
