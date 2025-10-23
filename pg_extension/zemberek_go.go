// package pgextension
package main

import "C"
import (
	"fmt"
	"sync"
	"time"
)

// GoHello is a demo function that uses goroutines to process the input
//
//export GoHello
func GoHello(name *C.char) *C.char {
	goName := C.GoString(name)

	// Create a channel to receive results from goroutines
	resultChan := make(chan string, 3)
	var wg sync.WaitGroup

	// Spawn multiple goroutines to demonstrate concurrent processing
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// Simulate some work
			time.Sleep(time.Millisecond * 10)
			resultChan <- fmt.Sprintf("Goroutine-%d processed '%s'", id, goName)
		}(i)
	}

	// Close channel when all goroutines complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results from all goroutines
	var results string
	for msg := range resultChan {
		if results != "" {
			results += "; "
		}
		results += msg
	}

	finalMessage := fmt.Sprintf("Hello from Zemberek-Go! Results: [%s]", results)

	return C.CString(finalMessage)
}

func main() {
	// Required for CGO build
}
