package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)
	
func main() {
	fmt.Println("Process ID: " , os.Getpid())
	// create a channel to receive system signals
	// channel will be the asynchronous endpoint for signal notifications
	sigs := make(chan os.Signal , 1)

	//resigter the signals we want to catch
	//we want to catch SIGINT (interrupt from terminal)
	signal.Notify(sigs , syscall.SIGINT , syscall.SIGTERM)

	done := make(chan bool , 1)

	go func() {
		sig := <- sigs
		fmt.Printf("\n\n [HANDLER] Received signal: %v\n" , sig)
		done <- true
	}()

	fmt.Println("[MAIN] Waiting for work or signal...")
	<- done
	fmt.Println("[MAIN] Application shut down successfully")
}