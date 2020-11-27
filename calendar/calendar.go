package calendar

// Event type
type Event struct {
	EventID      int    `json:"eventId"`
	EventName    string `json:"eventName"`
	EventDetails string `json:"eventDetails"`
}
