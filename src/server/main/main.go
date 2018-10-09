package main

import (
	"log"
	"os"
	"os/signal"

	"server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	s := server.NewServer()
	if s == nil {
		log.Println("create server fail")
		return
	}
	s.Init()
	go s.Run()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	s.Stop()
}
