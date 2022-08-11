package main

import (
	"github.com/sharabao13/sync/server"
	"github.com/sharabao13/sync/server/config"
	"os"
	"os/exec"
	"os/signal"
)

func main() {
	chChromDie := make(chan struct{})
	chBackendDie := make(chan struct{})
	go server.Run()
	go startBrowser(chChromDie, chBackendDie)
	chSignal := listenToInterrupt()
	for {
		select {
		case <-chSignal:
			chBackendDie <- struct{}{}

		case <-chChromDie:
			os.Exit(0)
		}
	}

}

func startBrowser(chChromDie chan struct{}, chBackendDie chan struct{}) {
	chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	cmd := exec.Command(chromePath, "--app=http://localhost:"+config.GetPort()+"/static/index.html")
	cmd.Start()
	go func() {
		<-chBackendDie
		cmd.Process.Kill()
	}()
	go func() {
		cmd.Wait()
		chChromDie <- struct{}{}
	}()

}

func listenToInterrupt() chan os.Signal {
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)
	return chSignal
}
