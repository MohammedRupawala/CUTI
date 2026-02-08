package core

import (
	"io"
	"syscall"
)




type Client struct {
	io.ReadWriter
	fd int
	queue RedisCmds
	isTxn bool
}

func (f *Client) Write (b []byte) (int , error){
	return syscall.Write(f.fd,b)
}

func (f *Client) Read (b []byte) (int , error){
	return syscall.Read(f.fd,b)
}


func TxExec(c *Client){
	c.queue = make(RedisCmds, 0)
	c.isTxn = false
}

func AddToQueue(c *Client,cmd *RedisCmd){
	c.queue = append(c.queue, cmd)
}
func TxDiscard(c *Client){
	c.queue = make(RedisCmds, 0)
	c.isTxn = false
}

func StartTx(c *Client){
	c.isTxn = true
}


func NewClient(fd int) *Client {
	return &Client{
		fd:     fd,
		queue: make(RedisCmds, 0),
	}

}