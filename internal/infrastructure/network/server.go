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
	addr        string
	gameService *services.GameService

	mu          sync.Mutex
	connections map[net.Conn]struct{}
}

func NewServer(addr string, gs *services.GameService) *Server {
	return &Server{
		addr:        addr,
		gameService: gs,
		connections: make(map[net.Conn]struct{}),
	}
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer ln.Close()
	log.Printf("Server listening on %s", s.addr)

	go s.broadcastLoop(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Println("Server context canceled, stopping accept loop.")
			return nil
		default:
		}

		conn, err := ln.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}

		s.mu.Lock()
		s.connections[conn] = struct{}{}
		s.mu.Unlock()

		log.Printf("New client connected: %s", conn.RemoteAddr())

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		s.mu.Lock()
		delete(s.connections, conn)
		s.mu.Unlock()
		conn.Close()
		log.Printf("Client disconnected: %s", conn.RemoteAddr())
	}()

	decoder := json.NewDecoder(conn)
	for {
		var cmd dtos.CommandDTO
		if err := decoder.Decode(&cmd); err != nil {
			log.Printf("Decode error from %s: %v", conn.RemoteAddr(), err)
			return
		}
		if err := s.gameService.ProcessCommandDTO(cmd); err != nil {
			log.Printf("ProcessCommandDTO error: %v", err)
			s.sendError(conn, err.Error())
		}
	}
}

func (s *Server) broadcastLoop(ctx context.Context) {
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			snapshot := s.gameService.BuildWorldSnapshot()
			s.broadcastState(snapshot)
		}
	}
}

func (s *Server) broadcastState(snapshot interface{}) {
	data, err := json.Marshal(snapshot)
	if err != nil {
		log.Println("Marshal world snapshot error:", err)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for c := range s.connections {
		if _, werr := c.Write(data); werr != nil {
			log.Printf("Error writing to client %s: %v", c.RemoteAddr(), werr)
		}
	}
}

func (s *Server) sendError(conn net.Conn, msg string) {
	resp := map[string]string{"error": msg}
	data, _ := json.Marshal(resp)
	conn.Write(data)
}
