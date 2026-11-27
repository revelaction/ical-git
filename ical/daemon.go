package ical

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/caasmo/go-daemon-runner/daemon"
	"github.com/revelaction/ical-git/config"
	"github.com/revelaction/ical-git/fetch"
	"github.com/revelaction/ical-git/fetch/filesystem"
	"github.com/revelaction/ical-git/fetch/git"
	"github.com/revelaction/ical-git/notify/schedule"
)

// Daemon fetches the ical files, parses them and schedules the notifications
// periodically. It reloads the configuration on SIGHUP.
type Daemon struct {
	daemon.Base
	daemon.ReloadBase

	configPath string
	conf       atomic.Pointer[config.Config]
	sc         *schedule.Scheduler
}

// NewDaemon creates the daemon. The configuration is loaded from configPath by Run.
func NewDaemon(configPath string) *Daemon {
	return &Daemon{
		Base:       daemon.NewBase("ical.Daemon", nil),
		ReloadBase: daemon.NewReloadBase(),
		configPath: configPath,
	}
}

// Run loads the configuration and starts the goroutine that performs the first
// run and then waits for the next tick, a reload request or the shutdown
// signal.
func (d *Daemon) Run() error {
	loadErr := d.loadConfig()
	if loadErr != nil {
		return loadErr
	}

	d.sc = schedule.New(&d.conf)

	go func() {
		defer close(d.ShutdownDone)

		ticker := time.NewTicker(d.conf.Load().DaemonTick)
		defer ticker.Stop()

		err := d.run(d.Ctx)
		if err != nil {
			slog.Error("Run error", "error", err)
		}

		for {
			select {
			case <-d.Ctx.Done():
				d.sc.StopTimers()
				return

			case reloadCtx := <-d.ReloadCh:
				reloadErr := d.reload(reloadCtx)
				ticker.Reset(d.conf.Load().DaemonTick)
				d.ReloadDone <- reloadErr

			case <-ticker.C:
				err = d.run(d.Ctx)
				if err != nil {
					slog.Error("Run error", "error", err)
				}
			}
		}
	}()

	return nil
}

// loadConfig reads and validates the configuration file and stores it in the
// atomic pointer.
func (d *Daemon) loadConfig() error {
	data, err := os.ReadFile(d.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	conf, err := config.Load(data)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	d.conf.Store(&conf)

	slog.Info("📝 Config:", "tick_time", conf.DaemonTick)
	slog.Info("📝 Config:", "Loc", conf.Location)
	slog.Info("📝 Config:", "Fetcher", conf.Fetcher())
	if conf.IsFetcherGit() {
		slog.Info("📝 Config:", "Git Repo URL", conf.FetcherGit.Url)
		slog.Info("📝 Config:", "Private SSH Key Path", conf.FetcherGit.PrivateKeyPath)
	} else {
		slog.Info("📝 Config:", "ical_directory", conf.FetcherFilesystem.Directory)
	}
	slog.Info("📝 Config:", "notifiers", strings.Join(conf.NotifierTypes, ", "))
	for _, alarm := range conf.Alarms {
		slog.Info("📝 Config: 🔔", "type", alarm.Action, "durIso", alarm.DurIso8601, "dur", alarm.Dur)
	}
	for _, image := range conf.Images {
		slog.Info("📝 Config: 📸", "name", image.Name, "type", image.Type, "value", image.TruncatedValue())
	}

	return nil
}

// run fetches the ical files, parses them and schedules the notifications for
// the tick period. Schedule replaces the previous timers, and it is called only
// after a successful fetch, so a failed fetch keeps the previous alarms.
func (d *Daemon) run(ctx context.Context) error {
	conf := d.conf.Load()

	slog.Info("🚀 starting run")

	tickStart := time.Now()
	tickStartInConfigLoc := tickStart.In(conf.Location.Location)

	layout := "2006-01-02 15:04:05 MST"
	slog.Info("⏰ Tick start time", "start_local_location", tickStart.Format(layout), "start_config_location", tickStartInConfigLoc.Format(layout))

	var f fetch.Fetcher
	if conf.IsFetcherGit() {
		slog.Info("🧲 Fetch: git")
		f = git.New(conf.FetcherGit.Url, conf.FetcherGit.PrivateKeyPath)
	} else {
		slog.Info("🧲 Fetch: filesystem")
		f = filesystem.New(conf.FetcherFilesystem.Directory)
	}

	ch := f.GetCh(ctx)

	p := NewParser(*conf, tickStart)
	for file := range ch {
		if file.Error != nil {
			slog.Error("🧲 Fetch:", "🚨 Error:", file.Error)
			return fmt.Errorf("error: %w", file.Error)
		}

		err := p.Parse(file)
		if err != nil {
			slog.Info("👀 Parse: error. Skipping", "error", err)
		}
	}

	ctxErr := ctx.Err()
	if ctxErr != nil {
		return ctxErr
	}

	schedErr := d.sc.Schedule(p.Notifications(), tickStart)
	if schedErr != nil {
		return schedErr
	}

	slog.Info("🏁 ending run")
	return nil
}

// reload loads the configuration and runs the fetch and schedule with it. A
// failed load returns before the run, so the previous configuration and timers
// stay. A successful load is live even when the run fails; the loop resets the
// ticker from the live configuration afterwards.
func (d *Daemon) reload(ctx context.Context) error {
	loadErr := d.loadConfig()
	if loadErr != nil {
		return loadErr
	}

	return d.run(ctx)
}