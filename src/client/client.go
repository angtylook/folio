package client

import (
	"log"
	//"net"
	"time"

	"base"
)

type Client struct {
	conn     *base.Connection
	executor *base.Executor
	service  *Service
	fire     *FireService
}

func NewClient() *Client {
	c := &Client{}
	c.executor = base.NewExecutor(1024)
	c.service = NewService()
	c.fire = NewFireService()
	return c
}

func (c *Client) Init() error {
	/*
		conn, err := net.Dial("tcp", "127.0.0.1:2018")
		if err != nil {
			return err
		}

		c.conn = base.NewConnection(conn, c)
		go c.conn.Start()
	*/
	go c.service.Start()
	go c.fire.Start()
	return nil
}

func (c *Client) Run() {
	go c.executor.StartWithFrame(c, time.Second/20)
	// test
	{
		f := base.NewFuture()
		f.ReplyAt(c.executor, func(arg interface{}, err error) {
			msg := arg.(string)
			log.Printf("reply at: %s, err: %v", msg, err)
		}).ExecuteAt(c.service.exec, func() (interface{}, error) {
			r, e := c.service.Hello()
			return r, e
		})
		f.Wait()
	}

	{
		f := base.NewFuture()
		fire := c.fire
		service := c.service
		f.ReplyAt(fire.exec, func(arg interface{}, err error) {
			msg := arg.(string)
			log.Printf("reply at: %s, err: %v", msg, err)
		}).ExecuteAt(service.exec, func() (interface{}, error) {
			r, e := service.RandomStr()
			return r, e
		})
		c.service = nil
		c.fire = nil
	}
}

func (c *Client) Stop() {
	// c.conn.Stop()
	c.executor.Stop()
}

func (c *Client) OnMessage(msg []byte) {
	log.Println(string(msg))
}

func (c *Client) Frame() {
	// log.Println("client frame")
	// c.conn.SendMessage([]byte("ping"))
}
