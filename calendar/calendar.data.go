package calendar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"
)

// hold our event list in memory
var eventMap = struct {
	sync.RWMutex
	m map[int]Event
}{m: make(map[int]Event)}

func init() {
	fmt.Println("loading events...")
	calMap, err := loadEventMap()
	eventMap.m = calMap
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d events loaded...\n", len(eventMap.m))
}

func loadEventMap() (map[int]Event, error) {
	fileName := "events.json"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] does not exist", fileName)
	}

	file, _ := ioutil.ReadFile(fileName)
	eventList := make([]Event, 0)
	err = json.Unmarshal([]byte(file), &eventList)
	if err != nil {
		log.Fatal(err)
	}
	calMap := make(map[int]Event)
	for i := 0; i < len(eventList); i++ {
		calMap[eventList[i].EventID] = eventList[i]
	}
	return calMap, nil
}

func getEvent(eventID int) *Event {
	eventMap.RLock()
	defer eventMap.RUnlock()
	if event, ok := eventMap.m[eventID]; ok {
		return &event
	}
	return nil
}

func removeEvent(eventID int) {
	eventMap.Lock()
	defer eventMap.Unlock()
	delete(eventMap.m, eventID)
}

func getEventList() []Event {
	eventMap.RLock()
	events := make([]Event, 0, len(eventMap.m))
	for _, value := range eventMap.m {
		events = append(events, value)
	}
	eventMap.RUnlock()
	return events
}

func getEventIds() []int {
	eventMap.RLock()
	eventIds := []int{}
	for key := range eventMap.m {
		eventIds = append(eventIds, key)
	}
	eventMap.RUnlock()
	sort.Ints(eventIds)
	return eventIds
}

func getNextEventID() int {
	eventIds := getEventIds()
	return eventIds[len(eventIds)-1] + 1
}

func addOrUpdateEvent(event Event) (int, error) {
	addOrUpdateID := -1
	if event.EventID > 0 {
		oldEvent := getEvent(event.EventID)
		// if it exists, replace it, otherwise return
		if oldEvent == nil {
			return 0, fmt.Errorf("event id [%d] doesn't exist", event.EventID)
		}
		addOrUpdateID = event.EventID
	} else {
		addOrUpdateID = getNextEventID()
		event.EventID = addOrUpdateID
	}
	eventMap.Lock()
	eventMap.m[addOrUpdateID] = event
	eventMap.Unlock()
	return addOrUpdateID, nil
}
