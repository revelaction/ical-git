package schedule

import (
	"fmt"
	"github.com/revelaction/ical-git/config"
	"github.com/revelaction/ical-git/notify"
	"github.com/revelaction/ical-git/notify/desktop"
	"github.com/revelaction/ical-git/notify/telegram"
	"log/slog"
	"path/filepath"
	"time"

	const offset = 3 * time.Second
)

type Scheduler struct {
	// TODO remove NOtfier from struct
	telegram notify.Notifier
	desktop  notify.Notifier

	conf   config.Config
	timers []*time.Timer
}

func NewScheduler(c config.Config) *Scheduler {
	return &Scheduler{
		conf: c,
	}
}

func (s *Scheduler) Schedule(notifications []notify.Notification, tickStart time.Time) error {
    slog.Info("🚦 Schedule", "num_notifications", len(notifications))

	for _, n := range notifications {

		f := s.getNotifyFunc(n)
		dur := n.Time.Sub(tickStart) - offset
        slog.Info("🚦 Schedule 🔔", "📁", filepath.Base(n.EventPath), "📌", n.Time, "🔖", dur.Truncate(1*time.Second), "durIso", n.DurIso8601, "type", n.Type)
		//dur = 3 * time.Second // Hack
		timer := time.AfterFunc(dur, f)
		s.timers = append(s.timers, timer)
	}

	return nil
}

func (s *Scheduler) getNotifyFunc(n notify.Notification) func() {

	var f func()

	switch n.Type {
	case "telegram":
		if s.telegram == nil {
			s.telegram = telegram.New(s.conf)
		}

		f = func() {
			err := s.telegram.Notify(n)
			if err != nil {
				fmt.Printf("Could not deliver telegram notfication: %s", err)
			}
		}

	case "desktop":
		if s.desktop == nil {
			s.desktop = desktop.New(s.conf)
		}

		f = func() {
			err := s.desktop.Notify(n)
			if err != nil {
				fmt.Printf("Could not deliver desktop notfication: %s", err)
			}
		}
	}

	return f
}

func (s *Scheduler) StopTimers() {
	for _, tmr := range s.timers {
		tmr.Stop()
	}
}
