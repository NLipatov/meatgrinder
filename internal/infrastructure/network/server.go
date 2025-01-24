package network

import (
	"context"
	"encoding/json"
	"log"
	"meatgrinder/internal/application/dtos"
	"meatgrinder/internal/application/services"
	"net"
	"sync"
	"time"
)

type Server struct {
	addr  string
	game  *services.GameService
	mu    sync.Mutex
	conns map[net.Conn]struct{}
}

func NewServer(a string, g *services.GameService) *Server {
	return &Server{
		addr:  a,
		game:  g,
		conns: make(map[net.Conn]struct{}),
	}
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer ln.Close()
	log.Printf("Server on %s", s.addr)
	go s.broadcast(ctx)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		s.mu.Lock()
		s.conns[conn] = struct{}{}
		s.mu.Unlock()
		go s.handle(conn)
	}
}

func (s *Server) handle(c net.Conn) {
	defer func() {
		s.mu.Lock()
		delete(s.conns, c)
		s.mu.Unlock()
		c.Close()
	}()
	d := json.NewDecoder(c)
	for {
		var cmd dtos.CommandDTO
		if err := d.Decode(&cmd); err != nil {
			return
		}
		s.game.ProcessCommandDTO(cmd)
	}
}

func (s *Server) broadcast(ctx context.Context) {
	t := time.NewTicker(200 * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			ss := s.game.BuildWorldSnapshot()
			b, _ := json.Marshal(ss)
			s.mu.Lock()
			for c := range s.conns {
				c.Write(b)
			}
			s.mu.Unlock()
		}
	}
}
