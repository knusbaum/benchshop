package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/knusbaum/benchshop/chatsrv/msgproto"
	"google.golang.org/grpc"

	"net/http"
	_ "net/http/pprof"
)

var port = flag.Int("p", 8800, "The port the server will listen on")

func main() {
	flag.Parse()

	go func() {
		// This is for profiling.
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	msgproto.RegisterChatServer(s, NewChatSrv())
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
