package network

import (
	"context"
	"encoding/json"
	"meatgrinder/internal/application/dtos"
	"net"
	"sync"
)

type Client struct {
	addr string
	conn net.Conn
	mu   sync.Mutex
	ch   chan interface{}
}

func NewClient(a string) *Client {
	return &Client{
		addr: a,
		ch:   make(chan interface{}, 100),
	}
}

func (c *Client) Connect(ctx context.Context) error {
	con, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}
	c.conn = con
	go c.listen(ctx)
	return nil
}

func (c *Client) listen(ctx context.Context) {
	defer c.Close()
	d := json.NewDecoder(c.conn)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		var raw interface{}
		if err := d.Decode(&raw); err != nil {
			return
		}
		c.ch <- raw
	}
}

func (c *Client) UpdatesChannel() <-chan interface{} {
	return c.ch
}

func (c *Client) SendCommand(cmd dtos.CommandDTO) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return nil
	}
	e := json.NewEncoder(c.conn)
	return e.Encode(cmd)
}

func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
		close(c.ch)
	}
}
