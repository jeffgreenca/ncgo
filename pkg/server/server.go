package server

import (
	"io"
	"log"
	"net"
	"os"
)

type Server struct {
	in      io.Reader
	out     io.Writer
	network string
	address string
}

// New instantiates Server
func New(network, address string) *Server {
	return &Server{
		in:      os.Stdin,
		out:     os.Stdout,
		network: network,
		address: address,
	}
}

// NewTCP instantiates Server for TCP
func NewTCP(address string) *Server {
	return New("tcp", address)
}

func (s *Server) Run() error {
	listener, err := net.Listen(s.network, s.address)
	if err != nil {
		return err
	}
	// early bird gets the worm
	// TODO handle multiple clients, and copy Server.out to all in parallel
	conn, err := listener.Accept()
	if err != nil {
		return err
	}
	s.handle(conn)
	return nil
}

func (s *Server) handle(conn net.Conn) {
	log.Printf("new connection from %s\n", conn.RemoteAddr())
	defer conn.Close()

	// send (but if client disconnects, interrupt sending)
	go func() {
		sent := doCopy(conn, s.in)
		log.Printf("total bytes sent for conn %v: %d\n", conn, sent)
	}()

	// receive
	received := doCopy(s.out, conn)
	log.Printf("total bytes received for conn %v: %d\n", conn, received)
}

// doCopy is io.Copy with logging on error
func doCopy(dst io.Writer, src io.Reader) int64 {
	count, err := io.Copy(dst, src)
	if err != nil {
		log.Print(err)
	}
	return count
}
