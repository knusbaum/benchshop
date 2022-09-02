package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/knusbaum/benchshop/chatsrv/msgproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr = flag.String("addr", "localhost:8800", "the address to connect to")
var nick = flag.String("name", "anonymous", "chat nickname")
var channame = flag.String("chan", "lobby", "channel to join")

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := msgproto.NewChatClient(conn)

	channel := &msgproto.Channel{Name: *channame}
	go func() {
		// Read messages in the background.
		stream, err := c.Join(context.Background(), channel)
		if err != nil {
			fmt.Printf("Failed to join channel: %v\n", err)
			os.Exit(1)
		}
		for {
			m, err := stream.Recv()
			if err != nil {
				fmt.Printf("Receive failed: %v\n", err)
				os.Exit(1)
			}
			t := time.Unix(m.Timestamp, 0)

			fmt.Printf("%s: (%s) %s\n", t.Format(time.Kitchen), m.Fromnick, m.Content)
		}
	}()

	in := bufio.NewReader(os.Stdin)
	for {
		line, _, err := in.ReadLine()
		if err != nil {
			fmt.Printf("Failed to read input: %v\n", err)
			os.Exit(1)
		}
		_, err = c.Send(context.Background(), &msgproto.InMessage{Channame: *channame, Fromnick: *nick, Content: string(line)})
		if err != nil {
			fmt.Printf("Failed to send message: %v\n", err)
			os.Exit(1)
		}
	}
}
