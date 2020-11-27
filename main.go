package main

import (
	"log"
	"net/http"

	"github.com/osuksdev/kirby-calendar-service/calendar"
)

const basePath = "/api"

func main() {
	calendar.SetupRoutes(basePath)
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
