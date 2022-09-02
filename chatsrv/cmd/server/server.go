package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/knusbaum/benchshop/chatsrv/msgproto"
)

type ChatSrv struct {
	clients map[string][]*SrvClient // Channel name to list of clients in that channel.
	backlog map[string][]*msgproto.OutMessage
	l       sync.RWMutex
	msgproto.UnimplementedChatServer
}

func NewChatSrv() *ChatSrv {
	return &ChatSrv{
		clients: make(map[string][]*SrvClient),
		backlog: make(map[string][]*msgproto.OutMessage),
	}
}

// Ensure our ChatSrv implements msgproto.ChatServer
var _ msgproto.ChatServer = &ChatSrv{}

func (s *ChatSrv) Send(ctx context.Context, m *msgproto.InMessage) (*msgproto.SendResponse, error) {
	om := msgproto.OutMessage{
		Timestamp: time.Now().Unix(),
		Fromnick:  m.Fromnick,
		Content:   m.Content,
	}
	s.sendToClients(m.Channame, &om)
	s.addToBacklog(m.Channame, &om)
	return &msgproto.SendResponse{}, nil
}

func (s *ChatSrv) addToBacklog(chname string, m *msgproto.OutMessage) {
	s.l.Lock()
	defer s.l.Unlock()
	s.backlog[chname] = append(s.backlog[chname], m)
}

func (s *ChatSrv) sendToClients(chname string, m *msgproto.OutMessage) {
	s.l.RLock()
	defer s.l.RUnlock()
	for _, c := range s.clients[chname] {
		c.Deliver(m)
	}
}

func (s *ChatSrv) addSrvClient(chname string, c *SrvClient) {
	s.l.Lock()
	defer s.l.Unlock()
	s.clients[chname] = append(s.clients[chname], c)
}

func (s *ChatSrv) removeSrvClient(chname string, c *SrvClient) {
	s.l.Lock()
	defer s.l.Unlock()
	clients := s.clients[chname]
	var k int
	for i := range clients {
		if clients[i] == c {
			continue
		}
		clients[k] = clients[i]
		k++
	}
	clients = clients[:k]
	s.clients[chname] = clients
}

func (s *ChatSrv) getBacklog(chname string) []*msgproto.OutMessage {
	s.l.RLock()
	defer s.l.RUnlock()
	return s.backlog[chname]
}

func (s *ChatSrv) Join(c *msgproto.Channel, stream msgproto.Chat_JoinServer) error {
	client := NewSrvClient()
	s.addSrvClient(c.Name, client)
	defer s.removeSrvClient(c.Name, client)

	bl := s.getBacklog(c.Name)
	for _, m := range bl {
		err := stream.Send(m)
		if err != nil {
			log.Printf("Failed to send message: %s\n", err)
			return err
		}
	}
	for m := range client.Receiver() {
		err := stream.Send(m)
		if err != nil {
			log.Printf("Failed to send message: %s\n", err)
			return err
		}
	}
	return nil
}
