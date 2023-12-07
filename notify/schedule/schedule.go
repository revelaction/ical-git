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
}

func NewScheduler(c config.Config) *Scheduler {
	return &Scheduler{
		conf: c,
	}
}

func (s *Scheduler) Notify(nt notify.Notification) error {
    
    var f func()

	switch nt.Type {
	case "telegram":
        if  s.telegram == nil {
            s.telegram = telegram.New(s.conf)
        } 

        f = func() {
            err :=  s.telegram.Notify(nt)
            if err != nil {
                fmt.Printf("Could not deliver telegram notfication: %s", err)
            }
        }

    case "desktop":
        if  s.desktop == nil {
            s.desktop = desktop.New(s.conf)
        } 

        f = func() {
            err :=  s.desktop.Notify(nt)
            if err != nil {
                fmt.Printf("Could not deliver desktop notfication: %s", err)
            }
        }
    }

    // get the start from the struct New:
    // get the tick from conf
    // if alarm
    // find the duration TODO
    timer := time.AfterFunc(2*time.Second, f)

    s.timers = append(s.timers, timer)



    return nil
}

func (s *Scheduler) StopScheduled() {
    for _, tmr := range s.timers {
        tmr.Stop()
    }
}
