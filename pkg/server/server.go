package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type Challenger interface {
	Create() string
	Verify(challenge string, solution int) bool
}

type DB interface {
	Add(key string)
	Exists(key string) bool
	Delete(key string)
}

type Quotes interface {
	Quote() string
}

type Server struct {
	challlenger Challenger
	db          DB
	quotes      Quotes
	port        int
	shutdown    <-chan interface{}
}

func NewServer(challenger Challenger, db DB, quotes Quotes, port int, shutdown <-chan interface{}) Server {
	return Server{challlenger: challenger, db: db, quotes: quotes, port: port, shutdown: shutdown}
}

func (s *Server) Start(startup chan interface{}) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Panic(err)
	}
	if s.port == 0 {
		parts := strings.Split(listener.Addr().String(), ":")
		port, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			log.Panic(err)
		}
		s.port = port
	}
	defer listener.Close()
	startup <- struct{}{}

	for {
		select {
		case <-s.shutdown:
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			go s.handle(conn)
		}
	}
}

func (s *Server) Port() int {
	return s.port
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Print(err)
		return
	}

	cmd := string(buf[:n])
	if strings.HasPrefix(cmd, "challenge") {
		s.challenge(conn)
	} else if strings.HasPrefix(cmd, "solve") {
		s.solve(conn, cmd)
	} else {
		_, err := conn.Write([]byte("requested message is invalid"))
		if err != nil {
			log.Println(err)
		}
		return
	}
}

func (s *Server) solve(conn net.Conn, answer string) {
	parts := strings.Split(answer, " ")
	if len(parts) < 3 {
		_, err := conn.Write([]byte("error"))
		if err != nil {
			log.Println(err)
		}
		return
	}

	if !s.db.Exists(parts[1]) {
		_, err := conn.Write([]byte("token already used"))
		if err != nil {
			log.Println(err)
		}
		return
	}

	num, err := strconv.Atoi(strings.Trim(parts[2], "\n"))
	if err != nil {
		log.Println(err)
	}
	if !s.challlenger.Verify(parts[1], num) {
		_, err := conn.Write([]byte("error"))
		if err != nil {
			log.Println(err)
		}
		return
	}

	s.db.Delete(answer)
	_, err = conn.Write([]byte(s.quotes.Quote()))
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) challenge(conn net.Conn) {
	challenge := s.challlenger.Create()
	s.db.Add(challenge)
	_, err := conn.Write([]byte(challenge))
	if err != nil {
		log.Println(err)
	}
}
