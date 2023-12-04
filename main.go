package main

import (
	"fmt"
	"time"
)

func main() {
	// Define the target time
	targetTime := time.Date(2023, 11, 26, 20, 56, 0, 0, time.UTC)

	// Calculate the duration until the target time
	now := time.Now().UTC()
	durationUntilTarget := targetTime.Sub(now)

	// Schedule the goroutine to run at the target time

	time.AfterFunc(durationUntilTarget, func() {
		fmt.Println("Running goroutine at target time")
	})

	// Cancel the timer
	//if timer.Stop() {
	//	fmt.Println("Timer was stopped before it could fire")
	//}

	// Keep the main function running
	select {}
}
