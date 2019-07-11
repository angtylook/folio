package client

import (
	"base"
	"errors"
	"log"
	"time"
)

type Service struct {
	exec *base.Executor
}

func NewService() *Service {
	s := &Service{
		exec: base.NewExecutor(100),
	}
	return s
}

func (s *Service) Hello() (string, error) {
	time.Sleep(time.Second * 5)
	s.Access()
	return "hello", nil
}

func (s *Service) RandomStr() (string, error) {
	time.Sleep(time.Second * 5)
	s.Access()
	return "", errors.New("not implement")
}

func (s *Service) Access() {
	log.Println("access s ok")
}

func (s *Service) Start() {
	s.exec.Start()
}

func (s *Service) Stop() {
	s.exec.Stop()
}

type FireService struct {
	Service
}

func NewFireService() *FireService {
	f := &FireService{
		Service: *NewService(),
	}
	return f
}
