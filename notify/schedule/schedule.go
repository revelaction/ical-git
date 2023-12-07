package schedule

import (
	"time"
	"fmt"
	"github.com/revelaction/ical-git/notify"
	"github.com/revelaction/ical-git/config"
	"github.com/revelaction/ical-git/notify/desktop"
	"github.com/revelaction/ical-git/notify/telegram"
)

type Scheduler struct {
    telegram notify.Notifier
    desktop notify.Notifier
	conf config.Config
    timers []*time.Timer
    start time.Time
}

func NewScheduler(c config.Config, start time.Time) *Scheduler {
	return &Scheduler{
		conf: c,
        start: start,
	}
}

func (s *Scheduler) Schedule(notifications []notify.Notification) error {
    
    for _, n := range notifications {
        f := s.getNotifyFunc(n)
        dur := n.Time.Sub(s.start) 
        fmt.Printf("Notification for %#v with trigger in  %s", n, dur)
        timer := time.AfterFunc(dur, f)
        s.timers = append(s.timers, timer)
    }

    return nil
}

func (s *Scheduler) getNotifyFunc(n notify.Notification) func() {

    var f func()

	switch n.Type {
	case "telegram":
        if  s.telegram == nil {
            s.telegram = telegram.New(s.conf)
        } 

        f = func() {
            err :=  s.telegram.Notify(n)
            if err != nil {
                fmt.Printf("Could not deliver telegram notfication: %s", err)
            }
        }

    case "desktop":
        if  s.desktop == nil {
            s.desktop = desktop.New(s.conf)
        } 

        f = func() {
            err :=  s.desktop.Notify(n)
            if err != nil{
                fmt.Printf("Could not deliver desktop notfication: %s", err)
            }
        }
    }

    return f
}

func (s *Scheduler) StopScheduled() {
    for _, tmr := range s.timers {
        tmr.Stop()
    }
}
