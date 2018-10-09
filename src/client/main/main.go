package main

import (
	"log"
	"os"
	"os/signal"

	"client"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	c := client.NewClient()
	if c == nil {
		log.Println("create client fail")
		return
	}

	err := c.Init()
	if err != nil {
		log.Println("init client fail %v", err)
		return
	}

	go c.Run()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	c.Stop()
}
