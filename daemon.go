package main

import (
	"context"
	"log"
	"os"
	//"fmt"
	"os/signal"
	"syscall"
	"time"

)

const defaultTick = 60 * time.Second


func start()  {
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

    select{}
}

func run(ctx context.Context)  {
    ticker := time.NewTicker(3 * time.Second)
    defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
            log.Printf("run: received call for Done. returning")
			return 
		case <-ticker.C:
            log.Printf("starting tick ----------------")

            log.Printf("starting task")
            time.Sleep(2 * time.Second)
            log.Printf("end task")
		}
	}
}
