package main

import (
	"flag"

	"github.com/erikreppel/blogk"
	. "github.com/tendermint/go-common"
	"github.com/tendermint/tmsp/server"
)

func main() {
	addrPtr := flag.String("addr", "tcp://0.0.0.0:46658", "Listen address")
	tmspPtr := flag.String("tmsp", "socket", "socket | grpc")
	flag.Parse()

	// Start the listener
	_, err := server.NewServer(*addrPtr, *tmspPtr, blogk.NewBlogK())
	if err != nil {
		Exit(err.Error())
	}

	// Wait forever
	TrapSignal(func() {
		// Cleanup
	})
}
