package main

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/knusbaum/benchshop/chatsrv/msgproto"
	"golang.org/x/net/nettest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func testServer() (s *ChatSrv, c msgproto.ChatClient, shutdown func(), e error) {
	l, err := nettest.NewLocalListener("tcp")
	if err != nil {
		return nil, nil, nil, err
	}
	gs := grpc.NewServer()
	cs := NewChatSrv()
	msgproto.RegisterChatServer(gs, cs)
	go gs.Serve(l)

	conn, err := grpc.Dial(l.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, nil, err
	}
	//defer conn.Close()
	client := msgproto.NewChatClient(conn)

	shutdown = func() {
		l.Close()
		conn.Close()
	}

	return cs, client, shutdown, nil
}

// BenchmarkSend benchmarks sending
func BenchmarkSend(b *testing.B) {
	_, c, shutdown, err := testServer()
	if err != nil {
		b.Fatalf("Failed to start test server: %v\n", err)
	}
	defer shutdown()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Send(context.Background(), &msgproto.InMessage{Channame: "testchan", Fromnick: "tester", Content: "Hello, World!"})
	}
}

// BenchmarkJoin benchmarks a client connecting, receiving 10 messages from the backlog, and then disconnecting.
func BenchmarkJoin(b *testing.B) {
	s, c, shutdown, err := testServer()
	if err != nil {
		b.Fatalf("Failed to start test server: %v\n", err)
	}
	defer shutdown()

	var obs []*msgproto.OutMessage
	for i := 0; i < 100; i++ {
		obs = append(obs, &msgproto.OutMessage{
			Timestamp: time.Now().Unix(),
			Fromnick:  "someperson",
			Content:   "This is a message.",
		})
	}
	s.backlog["test"] = obs

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jc, err := c.Join(context.Background(), &msgproto.Channel{Name: "test"})
		if err != nil {
			b.Fatalf("Failed to join: %v\n", err)
		}
		for j := 0; j < 10; j++ {
			_, err := jc.Recv()
			if err != nil {
				b.Fatalf("Failed to receive: %v\n", err)
			}
		}
	}
}

func BenchmarkContention(b *testing.B) {
	l, err := nettest.NewLocalListener("tcp")
	if err != nil {
		b.Fatal(err)
	}
	defer l.Close()
	gs := grpc.NewServer()
	cs := NewChatSrv()
	msgproto.RegisterChatServer(gs, cs)
	go gs.Serve(l)

	var conns []msgproto.ChatClient
	for i := 0; i < 20; i++ {
		conn, err := grpc.Dial(l.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			b.Fatal(err)
		}
		defer conn.Close()
		client := msgproto.NewChatClient(conn)
		conns = append(conns, client)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for _, c := range conns {
			wg.Add(1)
			go func() {
				defer wg.Done()
				c.Send(context.Background(), &msgproto.InMessage{Channame: "testchan", Fromnick: "tester", Content: "Hello, World!"})
			}()
		}
		wg.Wait()
	}
}

// go test -v -run XXX -bench . -benchmem -mutexprofile mux.pprof -blockprofile block.pprof -cpuprofile cpu.pprof -memprofile mem.pprof -benchtime 20s
