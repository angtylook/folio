package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"client"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	c := client.NewClient()
	if c == nil {
		fmt.Println("create client fail")
		return
	}
	c.Init()
	go c.Run()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	c.Stop()
}
