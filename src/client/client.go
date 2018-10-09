package client

import (
	"log"
	"net"
	"time"

	"base"
)

type Client struct {
	conn     *base.Connection
	executor *base.Executor
}

func NewClient() *Client {
	c := &Client{}
	c.executor = base.NewExecutor(1024)
	return c
}

func (c *Client) Init() error {
	conn, err := net.Dial("tcp", "127.0.0.1:2018")
	if err != nil {
		return err
	}

	c.conn = base.NewConnection(conn, c)
	go c.conn.Start()
	return nil
}

func (c *Client) Run() {
	go c.executor.StartWithFrame(c, time.Second/20)
}

func (c *Client) Stop() {
	c.conn.Stop()
	c.executor.Stop(base.StopOptRunAll)
}

func (c *Client) OnMessage(msg []byte) {
	log.Println(string(msg))
}

func (c *Client) Frame() {
	log.Println("client frame")
	c.conn.SendMessage([]byte("ping"))
}
