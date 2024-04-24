package main

import (
	"fmt"
	"os"
	"os/signal"
	"github.com/Phund4/l0wb/api"
	"github.com/Phund4/l0wb/cmd/config"
	"github.com/Phund4/l0wb/internal/db"
	"github.com/Phund4/l0wb/internal/streaming"
)

func main() {
	config.ConfigSetup()
	csh := db.NewCache()
	dbObject := db.NewDB(csh)
	sh := streaming.NewStreamingHandler(dbObject)

	myApi := api.NewApi(csh)

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			sh.Finish()
			myApi.Finish()

			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
