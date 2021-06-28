package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	rand.Seed(98)
	sigKill := make(chan int)
	width := 50
	counter := 0

	fmt.Println("Please wait while your file is processed")

	// loading bar is a concurrent function with a non blocking for select loop
	wg.Add(1)
	go func() {
		loadingBar("=", ">", width, sigKill)
		wg.Done()
	}()

	// this for loop simulates random work on a load of known size
	for {
		z := rand.Intn(6)
		if z == 4 { // when a rand is hit
			counter++    // a counter increments
			sigKill <- 1 // and a value is sent over the chanel
		}
		if counter == width { // if this counter == the total work to be done
			sigKill <- -1 // a kill signal is sent
			break         // and th work stops
		}
		time.Sleep(40 * time.Millisecond) // sleep to mimic work
	}
	wg.Wait() // must wait for the loading bar go routine to finish
}

func loadingBar(bar, tip string, width int, killSwitch chan int) {
	totalWork := 1
	for {
		select { // if a value is recieved it is checked
		case v := <-killSwitch:
			if v == 1 { // if work is done our counter (totalWork) is incremented
				totalWork += v
			} else if v == -1 { // if work is finished the bar is cleared and the function exits
				fmt.Printf("\r")
				return
			}
		default: // if no work is done/ nor yet finished the function does not block
			// and instead displays a status bar
			fmt.Printf("\r[%s%s%s]", strings.Repeat(bar, totalWork%width), tip, strings.Repeat(" ", width-(totalWork%width)))
		}
	}
}

// loading bar needs to know how much work thre is in total
// and the workload must have access to a chanel that can communicate progres as well as completion
// -1 is a good indicator of completion
