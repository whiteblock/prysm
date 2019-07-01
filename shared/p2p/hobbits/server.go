package hobbits

import (
	"net"
	"reflect"

	"github.com/pkg/errors"
	"github.com/renaynay/go-hobbits/encoding"
	"github.com/renaynay/go-hobbits/tcp"
	"github.com/renaynay/prysm/shared/p2p"
)

func NewHobbitsNode(host string, port int, peers []string) HobbitsNode {
	return HobbitsNode{
		host: host,
		port: port,
		staticPeers: peers,
		feeds: map[reflect.Type]p2p.Feed{},
	}
}

func (h *HobbitsNode) Listen() error { // TODO: is this how you catch the conn?
	server := tcp.NewServer(h.host, h.port)
	ch := make(chan encoding.Message)
	conns := make(chan net.Conn)

	go server.Listen(func(conn net.Conn, message encoding.Message){
		ch <- message
		conns <- conn

	})

	msg := <- ch
	conn := <-conns

	go h.processHobbitsMessage(HobbitsMessage(msg), conn)

	return nil
}

func (h *HobbitsNode) Send(message HobbitsMessage, peer string, conn net.Conn) error {
	server := tcp.NewServer(peer, h.port)

	if conn == nil {
		conn, err := net.Dial("tcp", peer)
		
		err = server.SendMessage(conn, message)
		if err != nil {
			return errors.Wrap(err, "error sending message: ")
		}
		return nil
	}

	err := server.SendMessage(conn, encoding.Message(message))
	if err != nil {
		return errors.Wrap(err, "error sending hobbits message: ")
	}

	return nil
}

func (h *HobbitsNode) Broadcast(msg HobbitsMessage) error {
	for _, peer := range h.staticPeers {
		err := h.Send(msg, peer, nil)
		if err != nil {
			return errors.Wrap(err, "error broadcasting: ")
		}
	}

	return nil
}
