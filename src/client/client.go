package client

import (
	"log"
	"net"

	"base"
)

type Client struct {
	conn *base.Connection
}

func NewClient() *Client {
	c := &Client{}
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

}

func (c *Client) Stop() {
	c.conn.Stop()
}

func (c *Client) OnMessage(msg []byte) {
	log.Println(string(msg))
	c.conn.SendMessage([]byte("ping"))
}
