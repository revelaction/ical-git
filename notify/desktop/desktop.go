package desktop

import (
	"fmt"
	"github.com/your/package/notify"
)

type Desktop struct{}

func (d *Desktop) Notify(notification ical.Notification) error {
	// Implement your desktop notification logic here
	return nil
}
