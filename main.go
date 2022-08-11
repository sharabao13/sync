package main

import (
	"github.com/sharabao13/sync/server"
	"os"
	"os/exec"
	"os/signal"
)

func main() {
	port := "27149"
	go func() {
		server.Run()
	}()
	chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	cmd := exec.Command(chromePath, "--app=http://localhost:"+port+"/static/index.html")
	cmd.Start()

	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)

	<-chSignal
	cmd.Process.Kill()
	//ui, err := lorca.New("https://www.google.com", "", 800, 600, "--disable-translate")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//select {
	////case <-ui.Done():
	//case <-chSignal:
	//}
	//ui.Close()
	select {
	case <-chSignal:
		cmd.Process.Kill()
	}
}
