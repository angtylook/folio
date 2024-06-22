package base

import (
	"encoding/binary"
	"log"
	"net"
	"sync"
)

const lengthByteSize = 2
const writeQueueSize = 2048

type MessageHandler interface {
	OnMessage(msg []byte)
}

type writeMessage struct {
	payload []byte
}

type Connection struct {
	conn    net.Conn
	handler MessageHandler

	readLenBuffer  []byte
	writeLenBuffer []byte

	writeQueue chan *writeMessage
	closed     bool
	done       chan bool
	wg         sync.WaitGroup
}

func NewConnection(c net.Conn, h MessageHandler) *Connection {
	con := &Connection{
		conn:           c,
		handler:        h,
		readLenBuffer:  make([]byte, lengthByteSize, lengthByteSize),
		writeLenBuffer: make([]byte, lengthByteSize, lengthByteSize),
		writeQueue:     make(chan *writeMessage, writeQueueSize),
		closed:         false,
		done:           make(chan bool, 3),
	}

	return con
}

func (c *Connection) Start() {
	go c.ReadPump()
	go c.WritePump()
	select {
	case <-c.done:
		if !c.closed {
			c.conn.Close()
			c.closed = true
		}
	}
	c.wg.Wait()
}

func (c *Connection) Stop() {
	c.done <- true
}

func (c *Connection) Closed() bool {
	return c.closed
}

func (c *Connection) SendMessage(msg []byte) {
	m := &writeMessage{
		payload: msg,
	}
	c.writeQueue <- m
}

func (c *Connection) ReadPump() {
	c.wg.Add(1)
	for {
		err := c.readBuffer(c.readLenBuffer)
		if err != nil {
			log.Println(err)
			break
		}
		length := int(binary.BigEndian.Uint16(c.readLenBuffer))
		buffer := make([]byte, length, length)
		err = c.readBuffer(buffer)
		if err != nil {
			log.Println(err)
			break
		}
		go c.handler.OnMessage(buffer)
	}
	c.done <- true
	c.wg.Done()
}

func (c *Connection) readBuffer(buffer []byte) error {
	length := len(buffer)
	readSize := 0
	for readSize < length {
		n, err := c.conn.Read(buffer[readSize:])

		if err != nil {
			return err
		}

		readSize += n
	}
	return nil
}

func (c *Connection) WritePump() {
	c.wg.Add(1)
	for {
		msg, ok := <-c.writeQueue
		if !ok {
			break
		}
		binary.BigEndian.PutUint16(c.writeLenBuffer, uint16(len(msg.payload)))
		err := c.writeBuffer(c.writeLenBuffer)
		if err != nil {
			log.Println(err)
			break
		}

		err = c.writeBuffer(msg.payload)
		if err != nil {
			log.Println(err)
			break
		}
	}
	c.done <- true
	c.wg.Done()
}

func (c *Connection) writeBuffer(buffer []byte) error {
	length := len(buffer)
	writeSize := 0
	for writeSize < length {
		n, err := c.conn.Write(buffer[writeSize:])
		if err != nil {
			return err
		}
		writeSize += n
	}
	return nil
}
