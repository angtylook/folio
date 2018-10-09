package server

import (
	"fmt"
	"log"
	"net"

	"base"
)

type Server struct {
	listener *base.ServerListener
	conn     *base.Connection
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
	s.conn.Stop()
}

func (s *Server) Serve(conn net.Conn) {
	fmt.Println("Serve conn")
	s.conn = base.NewConnection(conn, s)
	go s.conn.Start()
	s.conn.SendMessage([]byte("hello"))
}

func (s *Server) OnMessage(msg []byte) {
	log.Println(string(msg))
	s.conn.SendMessage([]byte("pong"))
}
