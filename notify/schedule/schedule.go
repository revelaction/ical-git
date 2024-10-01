package schedule

import (
	"errors"
	"fmt"
	"github.com/revelaction/ical-git/config"
	"github.com/revelaction/ical-git/notify"
	"github.com/revelaction/ical-git/notify/desktop"
	"github.com/revelaction/ical-git/notify/telegram"
	"log/slog"
	"path/filepath"
	"time"
)

type Scheduler struct {
	telegram notify.Notifier
	desktop  notify.Notifier

	conf   config.Config
	timers []*time.Timer
}

func New(c config.Config) *Scheduler {
	return &Scheduler{
		conf: c,
	}
}

// Schedule creates timers for the delivery of notifications.
// The timers are closed only on SIGHUP, not on timer ticks. That means they
// will be executed even if they are triggered in the tick boundary as they run
// in its own goroutines.
func (s *Scheduler) Schedule(notifications []notify.Notification, tickStart time.Time) error {

	s.initializeNotifiers()

	slog.Info("🚦 Schedule:", "num_notifications", len(notifications))

	for _, n := range notifications {

		dur := n.Time.Sub(tickStart)
		slog.Info("🚦 Schedule: 🔔", "📁", filepath.Base(n.EventPath), "📌", n.Time.Format("2006-01-02 15:04:05 MST"), "🔖", dur.Truncate(1*time.Second), "durIso", n.DurIso8601, "type", n.Type, "source", n.Source, "📸", n.ImageName)
		dur = 3 * time.Second // Hack

		f, err := s.getNotifyFunc(n)
		if err != nil {
			slog.Error("            :", "🚨Error:", err)
			continue

		}
		timer := time.AfterFunc(dur, f)
		s.timers = append(s.timers, timer)
	}

	slog.Info("🚦 Schedule:", "num_timers", len(s.timers))

	return nil
}

func (s *Scheduler) getNotifyFunc(n notify.Notification) (func(), error) {

	var f func()
	switch n.Type {
	case "telegram":

		if nil == s.telegram {
			return nil, errors.New("no notifier. Unable to create Notification")
		}

		f = func() {
			err := s.telegram.Notify(n)
			if err != nil {
				slog.Error("🚚 Notification:", "Send?", "🛑", "📁", filepath.Base(n.EventPath), "error", err, "📌", n.Time.Format("2006-01-02 15:04:05 MST"), "type", n.Type, "source", n.Source, "📸", n.ImageName)
				fmt.Printf("Could not deliver telegram notfication: %s", err)
				return
			}
			slog.Info("🚚 Notification:", "Send?", "✅", "📁", filepath.Base(n.EventPath), "📌", n.Time.Format("2006-01-02 15:04:05 MST"), "type", n.Type, "source", n.Source, "📸", n.ImageName)
		}

	case "desktop":
		f = func() {
			err := s.desktop.Notify(n)
			if err != nil {
				slog.Error("🚚 Notification:", "Send?", "🛑", "📁", filepath.Base(n.EventPath), "error", err, "📌", n.Time.Format("2006-01-02 15:04:05 MST"), "type", n.Type, "source", n.Source, "📸", n.ImageName)
				return
			}
			slog.Info("🚚 Notification:", "Send?", "✅", "📁", filepath.Base(n.EventPath), "📌", n.Time.Format("2006-01-02 15:04:05 MST"), "type", n.Type, "source", n.Source, "📸", n.ImageName)
		}

	}

	return f, nil

}

func (s *Scheduler) StopTimers() {
	for _, tmr := range s.timers {
		tmr.Stop()
	}

	s.timers = []*time.Timer{}
}

func (s *Scheduler) initializeNotifiers() {
	for _, t := range s.conf.NotifierTypes {
		switch t {
		case "telegram":
			tg, err := telegram.New(s.conf)
			if err != nil {
				slog.Error("🚦 Schedule: 🚨 Unable to create telegram bot client:", "error", err)
				s.telegram = nil
				break
			}

			s.telegram = tg
		case "desktop":
			s.desktop = desktop.New(s.conf)
		}
	}
}
