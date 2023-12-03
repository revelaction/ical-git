package ical

import (
	"github.com/arran4/golang-ical"
	"time"
)

type Parser struct{}


func (p *Parser) Parse(data []byte) error {
    reader := bytes.NewReader(data)
    cal, err := ics.ParseCalendar(reader)
    if err != nil {
        return err
    }

    _ := cal.Events()

    return err
}

func (p *Parser) Notifications() ([]Notification, error) {

	notifications := make([]Notification, 0)
	return notifications, nil
}

func NewParser() *Parser {
	return &Parser{}
}
