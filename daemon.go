package main

import (
	"context"
	"log"
	"os"
	"fmt"
	"os/signal"
	"syscall"
	"time"

)

const defaultTick = 60 * time.Second


func start()  {
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)

	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

	go func() {
		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGHUP:
                    log.Printf("SIGHUP called")
                    //start()
					os.Exit(1)
				case os.Interrupt:
                    log.Printf("Interrupt called")
					cancel()
					os.Exit(1)
                    log.Printf("Interrupt called")
				}
			case <-ctx.Done():
				log.Printf("Done.")
				os.Exit(1)
			}
		}
	}()

    if err := run(ctx); err != nil {
        fmt.Fprintf(os.Stderr, "%s\n", err)
        os.Exit(1)
    }
}

func run(ctx context.Context) error {
    ticker := time.NewTicker(3 * time.Second)
    defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
            log.Printf("starting tick ----------------")

            log.Printf("starting task")
            time.Sleep(2 * time.Second)
            log.Printf("end task")
		}
	}
}
