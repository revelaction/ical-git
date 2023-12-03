package ical

import (
	"bytes"
	"github.com/arran4/golang-ical"
)

type Parser struct{}

func (p *Parser) Parse(data []byte) error {
	reader := bytes.NewReader(data)
	cal, err := ics.ParseCalendar(reader)
	if err != nil {
		return err
	}

	cal.Events()

	return nil
}

//func (p *Parser) Notifications() ([]Notification, error) {
//
//	notifications := make([]Notification, 0)
//	return notifications, nil
//}

func NewParser() *Parser {
	return &Parser{}
}
