package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/revelaction/ical-git/config"
	"github.com/revelaction/ical-git/fetch/filesystem"
	"github.com/revelaction/ical-git/ical"
	"github.com/revelaction/ical-git/notify/schedule"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// configFile is the config file path
const configPathDefault = "icalgit.toml"


func main() {

	//flag.Usage = func() { fmt.Fprintf(os.Stderr, "%s\n", usage) }
	var configPath string
	flag.StringVar(&configPath, "c", configPathDefault, "the config file")
	flag.StringVar(&configPath, "config", configPathDefault, "the config file")
	flag.Parse()

	// logger
	initializeLogger()

	// signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)

	defer func() {
		signal.Stop(signalChan)
		//cancel() TODO
	}()

	go func() {

		cancel, scheduler := initialize(configPath)

		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGHUP:
					slog.Info("⚙️  SIGHUP called")
					slog.Info("⚙️  canceling previous ctx")
					cancel()
					slog.Info("⚙️  stop previous timers")
					scheduler.StopTimers()
					cancel, scheduler = initialize(configPath)

				case os.Interrupt:
					slog.Info("Interrupt called")
					cancel()
					os.Exit(1)
				}
				//case <-ctx.Done():
				//	log.Printf("Done.")
				//	os.Exit(1)
				//}
			}
		}
	}()

	select {}
}

// TODO make struct Daemon
func tick(ctx context.Context, conf config.Config, sc *schedule.Scheduler, start time.Time) {

	ticker := time.NewTicker(conf.DaemonTick)
	defer ticker.Stop()

	run(conf, start, sc)

	for {

		select {
		case <-ctx.Done():
			slog.Info("⚙️  ticker goroutine: received cancel. Ending")
			return
		case <-ticker.C:
			slog.Info("⚙️  starting new tick work")
			run(conf, start, sc)
			slog.Info("⚙️  ending tick work")
		}
	}
}

func run(conf config.Config, start time.Time, sc *schedule.Scheduler) {

	slog.Info("🚀 starting run")
	f := filesystem.New(conf.FetcherFilesystem.Directory)
	ch := f.GetCh()

	p := ical.NewParser(conf, start)
	for f := range ch {
		if f.Error != nil {
			slog.Info("fetch Error", "error", f.Error)
			os.Exit(1) // TODO
		}
		err := p.Parse(f.Content)
		if err != nil {
			fmt.Printf("error: %v+", err)
		}
	}

	sc.Schedule(p.Notifications())
	slog.Info("🏁 ending run")
}

func initialize(path string) (context.CancelFunc, *schedule.Scheduler) {
	slog.Info("⚙️  initializing: loading config", "path", path)
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	conf, err := config.Load(data)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	slog.Info("⚙️  initializing: creating new context")
	ctx, cancel := context.WithCancel(context.Background())

	now := time.Now()
	slog.Info("⚙️  initializing: 🕒 start time", "start", now.Format(time.RFC3339))

	slog.Info("⚙️  initializing: creating new scheduler")
	sc := schedule.NewScheduler(conf, now)

	slog.Info("⚙️  initializing: creating new ticker goroutine")
	go tick(ctx, conf, sc, now)
	return cancel, sc
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
