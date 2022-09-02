package main

import "github.com/knusbaum/benchshop/chatsrv/msgproto"

type SrvClient struct {
	out chan *msgproto.OutMessage
}

func NewSrvClient() *SrvClient {
	return &SrvClient{
		out: make(chan *msgproto.OutMessage, 10),
	}
}

func (s *SrvClient) Deliver(m *msgproto.OutMessage) {
	select {
	case s.out <- m:
	default:
	}
}

func (s *SrvClient) Receiver() <-chan *msgproto.OutMessage {
	return s.out
}
