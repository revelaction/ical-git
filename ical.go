package main

import (
	"fmt"
	"strings"
	"github.com/arran4/golang-ical"
	"github.com/teambition/rrule-go"
)

func main() {
	// Define the ICS file contents
	icsData := []byte(`
BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Your Company//Your Product//EN
BEGIN:VEVENT
UID:123456789
DTSTART;TZID=Europe/Berlin:20231226T150000
DTEND;TZID=Europe/Berlin:20231226T160000
RRULE:FREQ=MONTHLY;BYMONTHDAY=-6
SUMMARY:Pay Rent
DESCRIPTION:Remember to pay the rent.
END:VEVENT
END:VCALENDAR
`)

// RRULE:FREQ=MONTHLY;BYMONTHDAY=25

	// Parse the ICS file
	cal, err := ics.ParseCalendar(strings.NewReader(string(icsData)))
	if err != nil {
		panic(err)
	}

	// Get the events from the calendar
	events := cal.Events()

    // serialize
    // hasRule
    // Dstart
    // Rrule
    // TODO find if RRULE, if not, use golang ical to detemine next 
    // TODO if RRULE,

	// Print all recurring dates of the first event
	if len(events) > 0 {

        //ical can parse TZID
        startical, _ := events[0].GetStartAt()
        fmt.Println("ical DTSTART: ", startical)



        //get the summary like this
		rule := events[0].GetProperty(ics.ComponentPropertyRrule).Value
		start := events[0].GetProperty(ics.ComponentPropertyDtStart).Value

        icalSerialized := events[0].Serialize()
        fmt.Printf("lipo %v+", icalSerialized)


        icalStr := "DTSTART:" + start + "\nRRULE:" + rule
        fmt.Println("ical str: ", icalStr)

        icalStr = "DTSTART;TZID=Europe/Berlin:20231226T150000\nRRULE:FREQ=MONTHLY;BYMONTHDAY=-6"
        s, _ := rrule.StrToRRuleSet(icalStr)
        fmt.Printf("lipo %v#", events[0].GetProperty(ics.ComponentPropertyDtStart))
        fmt.Println("Date of the next event:", s.String())
        next := s.Iterator()

        for i := 0; i < 20; i++ {
            time, _ := next()
            fmt.Println("Date of the next event:", time)
        }
        //for _, t := range s.All() {
        //    fmt.Println(t)
        //}
	}
}

