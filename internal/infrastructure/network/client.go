package network

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"meatgrinder/internal/application/dtos"
	"net"
	"sync"
)

type Client struct {
	addr string
	conn net.Conn

	mu      sync.Mutex
	updates chan interface{}
}

func NewClient(addr string) *Client {
	return &Client{
		addr:    addr,
		updates: make(chan interface{}, 100),
	}
}

func (c *Client) Connect(ctx context.Context) error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}
	c.conn = conn

	go c.listenServer(ctx)

	log.Printf("Connected to server %s\n", c.addr)
	return nil
}

func (c *Client) listenServer(ctx context.Context) {
	defer func() {
		c.Close()
		log.Println("Client listen loop ended.")
	}()

	decoder := json.NewDecoder(c.conn)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		var raw interface{}
		if err := decoder.Decode(&raw); err != nil {
			log.Printf("Client decode error: %v", err)
			return
		}
		c.updates <- raw
	}
}

func (c *Client) SendCommand(cmd dtos.CommandDTO) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return fmt.Errorf("no connection")
	}
	return json.NewEncoder(c.conn).Encode(cmd)
}

func (c *Client) UpdatesChannel() <-chan interface{} {
	return c.updates
}

func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
	}
	close(c.updates)
}
