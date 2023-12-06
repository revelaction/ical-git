package main

import (
	"context"
	"fmt"
	"github.com/revelaction/ical-git/config"
	"github.com/revelaction/ical-git/notify"
	"github.com/revelaction/ical-git/fetch/filesystem"
	"github.com/revelaction/ical-git/ical"
	"github.com/revelaction/ical-git/notify/desktop"
	"github.com/revelaction/ical-git/notify/telegram"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	//"github.com/revelaction/ical-git/notify/telegram"
)

// configFile is the config file path (absolute path)
const configFile = "icalgit.toml"

func loadConfig() config.Config {
	// Config file
	var conf config.Config
	if _, err := toml.DecodeFile(configFile, &conf); err != nil {
		log.Fatal(err)
	}

	if err := conf.Validate(); err != nil {
		log.Fatal(err)
	}

    return conf
}

func main() {

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)

	defer func() {
		signal.Stop(signalChan)
		//cancel() TODO
	}()

	go func() {

        conf := loadConfig()

		ctx, cancel := context.WithCancel(context.Background())
		go tick(ctx, conf)

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
					go tick(ctx, conf)

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


    //create supported notifers in conf
	for _, n := range p.Notifications() {
        if !conf.IsNotifierAllowed(n.Type) {
            continue
        }
        fmt.Println("-----------", n)

        // get type
        // is supported? then 
        // getNotifierFortype notify.New(type)
        // save it as map here.

		// get type of notifuication build notifier.
		// AFterFunc, schedule
        not := notifier(n.Type, conf)
		_ = not.Notify(n)
	}

	log.Printf("end run()")

}

var tg *telegram.Telegram
var dk *desktop.Desktop

func notifier(notifierType string, conf config.Config) notify.Notifier {
	switch notifierType {
	case "telegram":
        if tg == nil {
            tg = telegram.New(conf)
        } 

        fmt.Println("hier notifier")
        return tg
        
	case "desktop":
        if dk == nil {
            dk = desktop.New(conf)
        }

        return dk
    }

    return nil
}
