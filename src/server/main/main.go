package main

import (
	"log"
	"os"
	"os/signal"

	"wheel/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	s := server.NewServer()
	if s == nil {
		log.Println("create server fail")
		return
	}

	err := s.Init()
	if err != nil {
		log.Println("init server fail %v", err)
		return
	}

	go s.Run()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	s.Stop()
}
