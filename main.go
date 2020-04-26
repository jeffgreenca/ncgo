package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jeffgreenca/ncgo/pkg/client"
	"github.com/jeffgreenca/ncgo/pkg/server"
)

type runner interface {
	Run() error
}

var (
	listen  bool
	address string
	network string
)

func init() {
	flag.BoolVar(&listen, "l", false, "listen")
	udp := flag.Bool("u", false, "udp")
	tcp := flag.Bool("t", true, "tcp")
	flag.Parse()

	// udp wins if both set
	if *tcp {
		network = "tcp"
	}
	if *udp {
		network = "udp"
	}

	if flag.NArg() != 1 {
		fmt.Println("specify address: <hostname>:<port>")
		flag.PrintDefaults()
		os.Exit(1)
	}
	address = flag.Arg(0)
}

func main() {
	if listen {
		run(server.New(network, address))
	} else {
		run(client.New(network, address))
	}
}

func run(r runner) {
	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
