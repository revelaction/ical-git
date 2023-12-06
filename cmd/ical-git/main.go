package main

import (
	"context"
	"fmt"
	"github.com/revelaction/ical-git/config"
	"github.com/revelaction/ical-git/fetch/filesystem"
	"github.com/revelaction/ical-git/ical"
	"github.com/revelaction/ical-git/notify/desktop"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	//"github.com/revelaction/ical-git/notify/telegram"
)

func start() {
}

func main() {

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)

	defer func() {
		signal.Stop(signalChan)
		//cancel() TODO
	}()

	go func() {

		ctx, cancel := context.WithCancel(context.Background())
		go tick(ctx)

		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGHUP:
					log.Printf("SIGHUP called")
					log.Printf("canceling previous ctx")
					cancel()
					log.Printf("Initializing: read new conf, cancel AfterFuncs...")
					//create new context
					log.Printf("creating new task context")
					ctx, cancel = context.WithCancel(context.Background())
					go tick(ctx)

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

func tick(ctx context.Context) {
	// TODO toml
	conf := config.Config{
		TZ:         "Europe/Paris",
		DaemonTick: "15m",
	}

	tick, err := time.ParseDuration(conf.DaemonTick)
	if err != nil {
		os.Exit(1)
	}

	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	for {
		run(conf)

		select {
		case <-ctx.Done():
			log.Printf("run: received call for Done. returning")
			return
		case <-ticker.C:
			log.Printf("starting tick ----------------")
			run(conf)
			log.Printf("end tick ----------------")
		}
	}
}

func run(conf config.Config) {

	log.Printf("start run()")
	f := filesystem.New("../ical-testdata")
	ch := f.GetCh()

	p := ical.NewParser(conf)
	for content := range ch {
		err := p.Parse(content)
		if err != nil {
			fmt.Printf("error: %v+", err)
		}

	}

	notifier := desktop.New("logo.png")
	for _, n := range p.Notifications() {
		// get type of notifuication build notifier.
		// AFterFunc, schedule

		_ = notifier.Notify(n)
	}

	log.Printf("end run()")

}
