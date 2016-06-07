package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/takyoshi/gocal"
	"google.golang.org/api/calendar/v3"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	gocalCmd = kingpin.New("gocal", "google calendar events api")
	confFile = gocalCmd.Flag("conf", "config file").Short('c').Default(os.Getenv("HOME") + "/.config/gocal/calendar.toml").String()

	// Events command
	evCmd = gocalCmd.Command("events", "google calendar events api")
	// GetList
	evList    = evCmd.Command("list", "insert google calendar events")
	listStart = evList.Flag("start-time", "start time of event formatted by RFC3339").
			Short('s').Default(time.Now().Add(-1 * 24 * 7 * time.Hour).Format(time.RFC3339)).String()
	listEnd = evList.Flag("end-time", "start time of event formatted by RFC3339").
		Short('e').Default(time.Now().Format(time.RFC3339)).String()
	// Insert
	evInsert    = evCmd.Command("insert", "insert google calendar events")
	evJSON      = evInsert.Flag("json-file", "insert google events from json file").Default("").String()
	eventDetail = evInsert.Flag("detail", "detail of event").Default("").String()
	eventName   = evInsert.Flag("name", "event name").Default("").String()
	insertStart = evInsert.Flag("start-time", "start time of event formatted by RFC3339").
			Short('s').Default(time.Now().Format(time.RFC3339)).String()
	insertEnd = evInsert.Flag("end-time", "start time of event formatted by RFC3339").
			Short('e').Default(time.Now().Add(15 * time.Minute).Format(time.RFC3339)).String()
)

func main() {
	gocalCmd.Version("v0.0.1")
	subcmd := kingpin.MustParse(gocalCmd.Parse(os.Args[1:]))

	conf, err := gocal.LoadConfig(*confFile)
	if err != nil {
		log.Fatalf("Unable to load config file. %v", err)
	}

	gc, err := gocal.NewCalendarClient(conf, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to retrieve calendar Client %v", err)
	}

	switch subcmd {
	case evList.FullCommand():
		events, err := gc.GetEventsList(*listStart, *listEnd)
		if err != nil {
			log.Fatalf("Unable to retrieve next ten of the user's events. %v", err)
		}
		res, err := events.MarshalJSON()
		if err != nil {
			log.Fatalf("%v", err)
		}
		fmt.Printf("%s", string(res))
	case evInsert.FullCommand():
		if *eventName == "" && *evJSON == "" {
			log.Fatalln("Empty both --name and --json-file parameters.")
		}

		if *evJSON != "" {
			var evts []gocal.Event
			f, err := ioutil.ReadFile(*evJSON)
			if err != nil {
				log.Fatalf("%s", err)
			}

			if err = json.Unmarshal(f, &evts); err != nil {
				log.Fatalf("%s", err)
			}

			for index, e := range evts {
				e.StartTime = *insertStart
				e.EndTime = *insertEnd
				evts[index] = e
			}
			gc.InsertEvents(evts)
		} else {
			e := gocal.Event{
				StartTime: *insertStart,
				EndTime:   *insertEnd,
				Title:     *eventName,
				Detail:    *eventDetail,
			}

			gc.InsertEvent(e)
		}
	}
}
