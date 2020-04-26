package client

import (
	"io"
	"log"
	"net"
	"os"
)

type Client struct {
	in      io.Reader
	out     io.Writer
	network string
	address string
}

// New instantiates Client
func New(network, address string) *Client {
	return &Client{
		in:      os.Stdin,
		out:     os.Stdout,
		network: network,
		address: address,
	}
}

// NewTCP instantiates Client for TCP
func NewTCP(address string) *Client {
	return New("tcp", address)
}

// Run client by connecting and then reading and writing until done or error
func (c *Client) Run() error {
	conn, err := net.Dial(c.network, c.address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// receive
	done := make(chan int)
	go func() {
		received := doCopy(c.out, conn)
		log.Printf("total bytes received: %d\n", received)
		done <- 1
	}()
	// send
	sent := doCopy(conn, c.in)
	log.Printf("total bytes sent %d \n", sent)
	// wait for other side to finish talking
	<-done
	return nil
}

// doCopy is io.Copy with logging on error
func doCopy(dst io.Writer, src io.Reader) int64 {
	count, err := io.Copy(dst, src)
	if err != nil {
		log.Print(err)
	}
	return count
}
