package main

import (
	"github.com/sharabao13/sync/server"
	"os"
	"os/exec"
	"os/signal"
)

func startBrowser() *exec.Cmd {
	port := "27149"
	chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	cmd := exec.Command(chromePath, "--app=http://localhost:"+port+"/static/index.html")
	cmd.Start()
	return cmd
}

func listenToInterrupt() chan os.Signal {
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)
	return chSignal
}
func main() {
	go server.Run()
	cmd := startBrowser()
	chSignal := listenToInterrupt()
	
	select {
	case <-chSignal:
		cmd.Process.Kill()
	}
}
