package notifier

import (
	"time"
	"fmt"
	"github.com/revelaction/ical-git/notify"
	"github.com/revelaction/ical-git/config"
	"github.com/revelaction/ical-git/notify/desktop"
	"github.com/revelaction/ical-git/notify/telegram"
)

// TODO better name Scheduler
type N struct {
    telegram notify.Notifier
    desktop notify.Notifier
	conf config.Config
    timers []*time.Timer
}

// notifier receives the potencial notification
// generate possible notification accoreding to configs
// if notfication.Time is in TickPeriod schedule with AfterFunc

func NewN(c config.Config) *N {
	return &N{
		conf: c,
	}
}

func (n *N) Notify(nt notify.Notification) error {
    
    var f func()

	switch nt.Type {
	case "telegram":
        if  n.telegram == nil {
            n.telegram = telegram.New(n.conf)
        } 

        f = func() {
            err :=  n.telegram.Notify(nt)
            if err != nil {
                fmt.Printf("Could not deliver telegram notfication: %s", err)
            }
        }

    case "desktop":
        if  n.desktop == nil {
            n.desktop = desktop.New(n.conf)
        } 

        f = func() {
            err :=  n.desktop.Notify(nt)
            if err != nil {
                fmt.Printf("Could not deliver desktop notfication: %s", err)
            }
        }
    }

    // get the start from the struct New:
    // get the tick from conf
    // find the duration TODO
    timer := time.AfterFunc(2*time.Second, f)

    n.timers = append(n.timers, timer)



    return nil
}

func (n *N) StopScheduled() {
    for _, tmr := range n.timers {
        tmr.Stop()
    }
}
