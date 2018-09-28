package server

import (
	"fmt"
	"net"

	"base"
)

type Server struct {
	listener *base.ServerListener
}

func NewServer() *Server {
	s := &Server{
		listener: base.NewServerListener(),
	}
	s.listener.Handler = s

	return s
}

func (s *Server) Init() {
	fmt.Println("server init")
	s.listener.ListenAt("tcp", ":2018")
}

func (s *Server) Run() {
	fmt.Println("server serve")
}

func (s *Server) Stop() {
	fmt.Println("server stop")
	s.listener.StopListen()
}

func (s *Server) Serve(conn net.Conn) {
	fmt.Println("Serve conn")
	conn.Write([]byte("hello"))
	conn.Close()
}
