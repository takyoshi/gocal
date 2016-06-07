package gocal

import (
	"io/ioutil"
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// Gocal is interface for google calendar api
type Gocal interface {
	GetEventsList(string, string) (*calendar.Events, error)
	InsertEvent(Event) error
	InsertEvents([]Event)
}

// GocalClient is a google calenar api client
type GocalClient struct {
	Srv  *calendar.Service
	Conf Config
}

// Event is google calendar event at Gocal
type Event struct {
	Title     string `json:"title"`
	Detail    string `json:"detail"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// NewCalendarClient returns  http client google calandar api
// scope is calendar.CalendarReadonlyScope or calendar.CalendarScope
func NewCalendarClient(c Config, scope string) (*GocalClient, error) {
	var gc GocalClient
	b, err := ioutil.ReadFile(c.Credential)
	if err != nil {
		return nil, err
	}

	jc, err := google.JWTConfigFromJSON(b, scope)
	if err != nil {
		return nil, err
	}

	client := jc.Client(oauth2.NoContext)

	srv, err := calendar.New(client)
	if err != nil {
		return nil, err
	}

	gc.Srv = srv
	gc.Conf = c
	return &gc, nil
}

// GetEventsList returns event list
func (gc GocalClient) GetEventsList(startTime string, endTime string) (*calendar.Events, error) {
	events, err := gc.Srv.Events.List(gc.Conf.CalendarID).TimeMax(endTime).
		TimeMin(startTime).SingleEvents(true).OrderBy("startTime").Do()
	if err != nil {
		return nil, err
	}
	return events, nil
}

// InsertEvent insert an event to the google calendar
func (gc GocalClient) InsertEvent(event Event) error {
	start := calendar.EventDateTime{
		DateTime: event.StartTime,
	}
	end := calendar.EventDateTime{
		DateTime: event.EndTime,
	}

	ge := calendar.Event{
		Summary:     event.Title,
		Start:       &start,
		Description: event.Detail,
		End:         &end,
	}

	_, err := gc.Srv.Events.Insert(gc.Conf.CalendarID, &ge).Do()
	return err
}

// InsertEvents insert multiple events
func (gc GocalClient) InsertEvents(events []Event) {
	var err error
	for _, e := range events {
		err = gc.InsertEvent(e)
		if err != nil {
			log.Printf("[Error] %s", err.Error())
		}
	}
}
