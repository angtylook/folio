package main

import (
	"fmt"
	"os"
	"os/signal"

	"server"
)

func main() {
	s := server.NewServer()
	if s == nil {
		fmt.Println("create server fail")
		return
	}
	s.Init()
	go s.Run()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	s.Stop()
}
