package main

import (
	"context"
	"fmt"
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

const defaultTick = 3 * time.Second

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
		go run(ctx)

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
					go run(ctx)

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

func run(ctx context.Context) {
	ticker := time.NewTicker(defaultTick)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("run: received call for Done. returning")
			return
		case <-ticker.C:
			log.Printf("starting tick ----------------")

			// get the ical fields (from git or local fylesystem), parse them, find next, build struct, create the AfterFunc, save the cancel method.
			// get thepath of the contents
			// https://gist.github.com/sethamclean/9475737
			f := filesystem.New("../ical-testdata")
			// TODO
			ch := f.GetCh()
			p := ical.NewParser()
			for content := range ch {
				err := p.Parse(content)
				if err != nil {
					fmt.Printf("error: %v+", err)
				}

			}

			notifier := desktop.New("logo.png")
			for _, n := range p.Notifications() {
				_ = notifier.Notify(n)
			}

			// //ch has channels of contents <-chan []byte
			// parser := ical.NewParser()
			// notifications := []Notification
			// for _, calendar := range ch() {
			//     nts := parser.Parse(calendarEvent) // [] do we allow many events in the calendar?. thro errors, not supported feautres for each file
			//     // create afterfunc
			//     //pass the time.time to the AfterFunc -> get Timer
			//     //put the timer in the notification
			//     for _, nt := range nts {
			//         nt.Timer:= notficationTimer() //many event in a calendar??? NO for MVP
			//         notifications := append(notifications, nt)
			//     }
			// }
			//

			// notifications := retrieverer.Get() // []notifications
			// rad each ical fields
			// retriever interface getPaths

			log.Printf("end task")
		}
	}
}
