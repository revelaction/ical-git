package main

import (
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/revelaction/ical-git/config"
	"github.com/revelaction/ical-git/fetch/filesystem"
	"github.com/revelaction/ical-git/ical"
	"github.com/revelaction/ical-git/notify/schedule"
	"log/slog"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// configFile is the config file path (absolute path)
const configFile = "icalgit.toml"

func loadConfig() config.Config {
	// Config file
	var conf config.Config
	if _, err := toml.DecodeFile(configFile, &conf); err != nil {
		log.Fatal(err)
	}

	return conf
}

func initialize() context.CancelFunc {
	conf := loadConfig()
	ctx, cancel := context.WithCancel(context.Background())
	go tick(ctx, conf)
	return cancel
}

func main() {

    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
            if a.Key == slog.TimeKey {
                return slog.Attr{}
            }
            return a
        },
    }))
    slog.SetDefault(logger)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)

	defer func() {
		signal.Stop(signalChan)
		//cancel() TODO
	}()

	go func() {

		cancel := initialize()

		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGHUP:
					log.Printf("SIGHUP called")
					log.Printf("canceling previous ctx")
					cancel()
					log.Printf("Initializing: read new conf")
					cancel = initialize()

				case os.Interrupt:
					log.Printf("Interrupt called")
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

func tick(ctx context.Context, conf config.Config) {

	now := time.Now()
	ticker := time.NewTicker(conf.DaemonTick)
	defer ticker.Stop()

	for {
		run(conf, now)

		select {
		case <-ctx.Done():
			log.Printf("run: received call for Done. returning")
			// here the
			// TODO we need the real one.
			//ntf := schedule.NewScheduler(conf, now)
			//ntf.StopScheduled()
			return
		case <-ticker.C:
			log.Printf("starting tick ----------------")
			run(conf, now)
			log.Printf("end tick ----------------")
		}
	}
}

func run(conf config.Config, start time.Time) {

	log.Printf("start run()")
	f := filesystem.New("../ical-testdata")
	ch := f.GetCh()

	p := ical.NewParser(conf, start)
	for content := range ch {
		err := p.Parse(content)
		if err != nil {
			fmt.Printf("error: %v+", err)
		}

	}

	ntf := schedule.NewScheduler(conf, start)
	ntf.Schedule(p.Notifications())
	log.Printf("end run()")
}
