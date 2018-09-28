package base

import (
	"fmt"
	"net"
)

type ConnHandler interface {
	Serve(conn net.Conn)
}

type ServerListener struct {
	listener net.Listener
	stop     bool
	Handler  ConnHandler
	exit     chan bool
}

func NewServerListener() *ServerListener {
	s := &ServerListener{
		exit: make(chan bool),
	}

	return s
}

func (s *ServerListener) ListenAt(network, addr string) error {
	l, err := net.Listen(network, addr)
	if err != nil {
		fmt.Printf("tcp listen on 2018 fail %v", err)
		return fmt.Errorf("listen fail")
	}
	s.listener = l
	s.stop = false
	go s.Accept()
	return nil
}

func (s *ServerListener) Accept() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("tcp accept error %v", err)
			if s.stop {
				break
			} else {
				continue
			}
		}
		if s.Handler != nil {
			go s.Handler.Serve(conn)
		}
		fmt.Printf("new connection from %s", conn.RemoteAddr().String())
	}

	s.listener.Close()
	s.listener = nil
	s.exit <- true
}

func (s *ServerListener) StopListen() {
	s.stop = true
	s.listener.Close()
	<-s.exit
}
