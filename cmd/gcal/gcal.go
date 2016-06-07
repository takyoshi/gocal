package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"time"

	"github.com/takyoshi/gcal"
	"google.golang.org/api/calendar/v3"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("%v", err)
	}
	confDir := usr.HomeDir + "/.config/gcal"

	var (
		gcalCmd  = kingpin.New("gcal", "google calendar events api")
		confFile = gcalCmd.Flag("conf", "config file").Short('c').Default(confDir + "/calendar.toml").String()

		evCmd = gcalCmd.Command("events", "google calendar events api")
		// GetList
		evList    = evCmd.Command("list", "insert google calendar events")
		listStart = evList.Flag("start-time", "start time of event formatted by RFC3339").
				Short('s').Default(time.Now().Add(-1 * 24 * 7 * time.Hour).Format(time.RFC3339)).String()
		listEnd = evList.Flag("end-time", "start time of event formatted by RFC3339").
			Short('e').Default(time.Now().Format(time.RFC3339)).String()

		// Insert
		evInsert    = evCmd.Command("insert", "insert google calendar events")
		eventDetail = evInsert.Flag("detail", "detail of event").Default("").String()
		eventName   = evInsert.Flag("name", "event name").Required().String()
		insertStart = evInsert.Flag("start-time", "start time of event formatted by RFC3339").
				Short('s').Default(time.Now().Format(time.RFC3339)).String()
		insertEnd = evInsert.Flag("end-time", "start time of event formatted by RFC3339").
				Short('e').Default(time.Now().Add(15 * time.Minute).Format(time.RFC3339)).String()
	)

	gcalCmd.Version("v0.0.1")
	subcmd := kingpin.MustParse(gcalCmd.Parse(os.Args[1:]))

	conf, err := gcal.LoadConfig(*confFile)
	if err != nil {
		log.Fatalf("Unable to load config file. %v", err)
	}

	gc, err := gcal.NewCalendarClient(conf, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to retrieve calendar Client %v", err)
	}

	switch subcmd {
	case evList.FullCommand():
		events, err := gc.GetEventsList(*listStart, *listEnd)
		if err != nil {
			log.Fatalf("Unable to retrieve next ten of the user's events. %v", err)
		}
		if len(events.Items) > 0 {
			for _, i := range events.Items {
				var when string
				if i.Start.DateTime != "" {
					when = i.Start.DateTime
				} else {
					when = i.Start.Date
				}
				fmt.Printf("%s\t%s\n", when, i.Summary)
			}
		} else {
			fmt.Printf("No upcoming events found.\n")
		}
	case evInsert.FullCommand():
		e := gcal.Event{
			StartTime: *insertStart,
			EndTime:   *insertEnd,
			Title:     *eventName,
			Detail:    *eventDetail,
		}

		gc.InsertEvent(e)
	}
}
