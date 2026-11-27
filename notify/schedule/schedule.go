package schedule

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/revelaction/ical-git/config"
	"github.com/revelaction/ical-git/notify"
	"github.com/revelaction/ical-git/notify/desktop"
	"github.com/revelaction/ical-git/notify/telegram"
)

type Scheduler struct {
	telegram notify.Notifier
	desktop  notify.Notifier

	conf   *atomic.Pointer[config.Config]
	timers []*time.Timer
}

func New(conf *atomic.Pointer[config.Config]) *Scheduler {
	return &Scheduler{conf: conf}
}

// Schedule creates timers for the delivery of notifications.
// The timers are closed only on SIGHUP, not on timer ticks. That means they
// will be executed even if they are triggered in the tick boundary as they run
// in its own goroutines.
func (s *Scheduler) Schedule(notifications []notify.Notification, tickStart time.Time) error {
	s.StopTimers()

	conf := s.conf.Load()

	s.initializeNotifiers(conf)

	slog.Info("🚦 Schedule:", "num_notifications", len(notifications))

	for _, n := range notifications {

		dur := n.Time.Sub(tickStart)
		slog.Info("🚦 Schedule: 🔔", "📁", filepath.Base(n.EventPath), "📌", n.Time.Format("2006-01-02 15:04:05 MST"), "🔖", dur.Truncate(1*time.Second), "durIso", n.DurIso8601, "type", n.Type, "source", n.Source, "📸", n.ImageName)
		//dur = 3 * time.Second // Hack

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

		telegramNotifier := s.telegram
		f = func() {
			err := telegramNotifier.Notify(n)
			if err != nil {
				slog.Error("🚚 Notification:", "Send?", "🛑", "📁", filepath.Base(n.EventPath), "error", err, "📌", n.Time.Format("2006-01-02 15:04:05 MST"), "type", n.Type, "source", n.Source, "📸", n.ImageName)
				fmt.Printf("Could not deliver telegram notfication: %s", err)
				return
			}
			slog.Info("🚚 Notification:", "Send?", "✅", "📁", filepath.Base(n.EventPath), "📌", n.Time.Format("2006-01-02 15:04:05 MST"), "type", n.Type, "source", n.Source, "📸", n.ImageName)
		}

	case "desktop":
		if nil == s.desktop {
			return nil, errors.New("no notifier. Unable to create Notification")
		}

		desktopNotifier := s.desktop
		f = func() {
			err := desktopNotifier.Notify(n)
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

func (s *Scheduler) initializeNotifiers(conf *config.Config) {
	s.telegram = nil
	s.desktop = nil

	for _, t := range conf.NotifierTypes {
		switch t {
		case "telegram":
			tg, err := telegram.New(*conf)
			if err != nil {
				slog.Error("🚦 Schedule: 🚨 Unable to create telegram bot client:", "error", err)
				break
			}

			s.telegram = tg
		case "desktop":
			s.desktop = desktop.New(*conf)
		}
	}
}
