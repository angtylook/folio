package base

import (
	"encoding/binary"
	"log"
	"net"
)

const lengthByteSize = 2
const writeQueueSize = 2048

type MessageHandler interface {
	OnMessage(msg []byte)
	OnConnClose()
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

	closed bool
}

func NewConnection(c net.Conn, h MessageHandler) *Connection {
	con := &Connection{
		conn:           c,
		handler:        h,
		readLenBuffer:  make([]byte, lengthByteSize, lengthByteSize),
		writeLenBuffer: make([]byte, lengthByteSize, lengthByteSize),
		writeQueue:     make(chan *writeMessage, writeQueueSize),
		closed:         false,
	}

	return con
}

func (c *Connection) Start() {
	go c.ReadPump()
	go c.WritePump()
}

func (c *Connection) SendMessage(msg []byte) {
	m := &writeMessage{
		payload: msg,
	}
	c.writeQueue <- m
}

func (c *Connection) ReadPump() {
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
