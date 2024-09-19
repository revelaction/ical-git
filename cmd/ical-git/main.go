package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/revelaction/ical-git/config"
	"github.com/revelaction/ical-git/fetch"
	"github.com/revelaction/ical-git/fetch/filesystem"
	"github.com/revelaction/ical-git/fetch/git"
	"github.com/revelaction/ical-git/ical"
	"github.com/revelaction/ical-git/notify/schedule"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const usage = `Usage:
    ical-git [-c CONF_FILE] 

Options:
    -c, --config                load the configuration file at CONF_FILE instead of default
    -v, --version               Print the version 
    -h, --help                  Show this

CONF_FILE is the toml configuration file 

ical-git will react to a SIGHUP signal reloading the configuration file.

Examples:
    $ ical-git --config /path/to/config/file.toml # start the daemon with the configuration file
    $ ical-git -v  # print version`

// Version can be set at link time
var BuildTag string

// configFile is the config file path
const configPathDefault = "icalgit.toml"

func main() {

	var configPath string
	var versionFlag, helpFlag bool

	flag.Usage = func() { fmt.Fprintf(os.Stderr, "%s\n", usage) }
	flag.BoolVar(&versionFlag, "v", false, "print the version")
	flag.BoolVar(&versionFlag, "version", false, "print the version")
	flag.StringVar(&configPath, "c", configPathDefault, "the config file")
	flag.StringVar(&configPath, "config", configPathDefault, "the config file")
	flag.BoolVar(&helpFlag, "h", false, "print the version")
	flag.BoolVar(&helpFlag, "help", false, "Show this")
	flag.Parse()

	if versionFlag {
		if BuildTag != "" {
			fmt.Println(BuildTag)
			return
		}
		fmt.Println("(unknown)")
		return
	}

	if helpFlag {
		flag.Usage()
		return
	}

	// logger
	initializeLogger()

	slog.Info("🏁 app:", "Version", BuildTag)

	// signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)

	defer func() {
		signal.Stop(signalChan)
		//cancel() TODO
	}()

	ctx, cancel, scheduler := initialize(configPath)

	go func() {

		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGHUP:
					slog.Info("🔧 SIGHUP called")
					slog.Info("🔧 canceling previous ctx")
					cancel()
					slog.Info("🔧 stop previous timers")
					scheduler.StopTimers()
					slog.Info("🔧 Reinizializing")
					ctx, cancel, scheduler = initialize(configPath)

				case os.Interrupt:
					slog.Info("🔧 Interrupt called. Cancelling goroutines. Exiting")
					cancel()
					os.Exit(1)
				}
			case <-ctx.Done():
				slog.Error("🚨 Tick Error. Context was cancelled. Reinizializing")
				scheduler.StopTimers()
				ctx, cancel, scheduler = initialize(configPath)
			}
		}
	}()

	select {}
}

// initialize reads the config and creates a goroutine (tick method) to
// retrieve periodically the ical files and set alarms delivery timers
// (goroutines).
// initialize is run at the start or after a SIGHUB signal
func initialize(path string) (context.Context, context.CancelFunc, *schedule.Scheduler) {
	slog.Info("🔧 Init: loading config", "path", path)
	data, err := os.ReadFile(path)
	if err != nil {
		slog.Error("Failed to read config file", "error", err)
		os.Exit(1)
	}
	conf, err := config.Load(data)
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

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

	for key, value := range conf.Images {
		slog.Info("📝 Config: 📸", "key", key, "value", value)
	}

	// Create context to cance the tick goroutine on SIGHUP
	ctx, cancel := context.WithCancel(context.Background())

	slog.Info("🔧 Init: creating new scheduler")
	sc := schedule.New(conf)

	slog.Info("🔧 Init: creating goroutine for time ticks")
	go tick(ctx, cancel, conf, sc)
	return ctx, cancel, sc
}

// tick is a goroutine to periodically retrieve the ical files an set alarms.
// tick does not stop the alarm timers at the start.
// At the start of the tick, all alarms for the tick period are scheduled.
// At the start of the next tick the is no alarms timers, so there is no need to close them.
func tick(ctx context.Context, cancel context.CancelFunc, conf config.Config, sc *schedule.Scheduler) {

	ticker := time.NewTicker(conf.DaemonTick)
	defer ticker.Stop()

	err := run(conf, sc)
	if err != nil {
        // sleep before canceling here, to allow interrupt signals to still work
        time.Sleep(time.Minute)
		cancel()
		slog.Error("Tick Error, canceling", "error", err)
		return
	}

	for {

		select {
		case <-ctx.Done():
			slog.Info("🔧 ticker goroutine: received cancel. Ending")
			return
		case <-ticker.C:
			slog.Info("🔧 starting new tick work")
			err = run(conf, sc)
			slog.Info("🔧 ending tick work")
			if err != nil {
                // sleep before canceling here, to allow interrupt signals to still work
                time.Sleep(time.Minute)
				cancel()
				slog.Error("Tick Error, canceling", "error", err)
				return
			}
		}
	}
}

func run(conf config.Config, sc *schedule.Scheduler) error {

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

	ch := f.GetCh()

	p := ical.NewParser(conf, tickStart)
	for f := range ch {
		if f.Error != nil {
            slog.Error("🧲 Fetch:", "🚨 Error:", f.Error)
			return fmt.Errorf("error: %w", f.Error)
		}
		err := p.Parse(f)
		if err != nil {
            slog.Info("👀 Parse: error. Skipping", "error", err)
		}
	}

	sc.Schedule(p.Notifications(), tickStart)
	slog.Info("🏁 ending run")
	return nil
}

func initializeLogger() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))

	slog.SetDefault(logger)
}
