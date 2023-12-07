package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/arran4/golang-ical"
	"github.com/teambition/rrule-go"
	"strings"
	"time"
	//"os"
)

func main() {
	// Define the ICS file contents
	//icsDataFoldedSpace := []byte("RDATE;TZID=Europe/Berlin:20231211T120000,20231212T120000,20231214T120000,\r\n 20231218T120000,20231211T120000,20231212T120000,20231214T120000,\r\n DESCRIPTION:Remember to pay the rent\r\n")

	icsData := []byte(`
BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Your Product//EN
CALSCALE:GREGORIAN
METHOD:PUBLISH
BEGIN:VEVENT
DTSTART;TZID=Europe/Berlin:20231206T120000
DTEND;TZID=Europe/Berlin:20231206T130000
DTSTAMP:20231201T090000Z
UID:unique-id-12345@example.com
SUMMARY: lots of events
RDATE;TZID=Europe/Berlin:20231211T120000,20231212T120000,20231214T120000,20231218T120000,20231211T120000,20231212T120000,20231214T120000,20231218T120000,20231211T120000,20231212T120000,20231214T120000,20231218T120000
END:VEVENT
END:VCALENDAR
`)

	//
	//scanner := bufio.NewScanner(bytes.NewReader(icsDataFoldedSpace))
	//scanner.Split(splitFunc)

	//for scanner.Scan() {
	//    line := scanner.Text()
	//    fmt.Printf("----%s----\n", line)
	//}
	//
	//os.Exit(1)

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
		//startical, _ := events[0].GetStartAt()
		//fmt.Println("ical DTSTART: ", startical)

		//get the summary like this
		//rule := events[0].GetProperty(ics.ComponentPropertyRrule).Value
		//start := events[0].GetProperty(ics.ComponentPropertyDtStart).Value

		icalSerialized := events[0].Serialize()
		fmt.Printf("\n%v\n", icalSerialized)

		linesStr := parseRRule(events[0])

		fmt.Printf("\n%v\n", linesStr)

		//icalStr := "DTSTART:" + start + "\nRRULE:" + rule
		//fmt.Println("ical str: ", icalStr)

		//icalStr = "DTSTART;TZID=Europe/Berlin:20231226T150000\nRRULE:FREQ=MONTHLY;BYMONTHDAY=-6"
		s, _ := rrule.StrToRRuleSet(linesStr)
		//fmt.Println("Date of the set:", s.String())
		//fmt.Printf("Set DTStart: %v#\n", s.GetDTStart())
		next := s.Iterator()

		for i := 0; i < 20; i++ {
			time, _ := next()
			fmt.Println("Date of the next event:", time)
		}
		//for _, t := range s.All() {
		//    fmt.Println(t)
		//}

		// TODO utc
		fmt.Println("The first recurrence after now is: ", s.After(time.Now(), false))
	}

}

func parseRRule(event *ics.VEvent) string {
	eventCleaned := strings.Replace(event.Serialize(), "\r\n ", "", -1)

	scanner := bufio.NewScanner(strings.NewReader(eventCleaned))
	scanner.Split(splitFunc)
	var lines []string

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("----%s----\n", line)

		if strings.HasPrefix(line, "DTSTART") {
			lines = append([]string{line}, lines...)
			continue
		}

		if strings.HasPrefix(line, "RRULE") {
			lines = append(lines, line)
			return strings.Join(lines, "\n")
		}

		if strings.HasPrefix(line, "RDATE") {

			lines = append(lines, line)
			return strings.Join(lines, "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		// TODO
	}

	return strings.Join(lines, "\n")
}

// https://icalendar.org/iCalendar-RFC-5545/3-1-content-lines.html
func splitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Look for the first occurrence of "\r\n" in the data.
	if i := bytes.Index(data, []byte("\r\n")); i >= 0 {

		// Request more data.
		if len(data) == i+1 {
			return 0, nil, nil
		}

		// Check if the next character is a space.
		if len(data) > i+2 && data[i+2] == ' ' {
			// The next character is a space, so this is not the end of the line.
			return 0, nil, nil
		}

		// Check if the next character is a HTAB.
		//if len(data) > i+2 && data[i+2] == '\t' {
		//    // The next character is a HTAB, so this is not the end of the line.
		//    return 0, nil, nil
		//}
		// We have a full "\r\n"-terminated line.
		return i + 2, data[:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
