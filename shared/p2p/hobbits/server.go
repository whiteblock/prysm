package hobbits

import (
	"net"
	"reflect"

	"github.com/pkg/errors"
	"github.com/renaynay/go-hobbits/encoding"
	"github.com/renaynay/prysm/shared/p2p"
	"github.com/renaynay/go-hobbits/tcp"
)

func NewHobbitsNode(host string, port int, peers []string) HobbitsNode {
	return HobbitsNode{
		host: host,
		port: port,
		staticPeers: peers,
		feeds: map[reflect.Type]p2p.Feed{},
	}
}

func (h *HobbitsNode) Send(msg HobbitsMessage, peer string, conn net.Conn) error {
	server := tcp.NewServer(peer, h.port)

	err := server.SendMessage(conn, encoding.Message(msg))
	if err != nil {
		return errors.Wrap(err, "error sending hobbits message: ")
	}

	return nil
}

func (h *HobbitsNode) Broadcast(msg HobbitsMessage) error {
	for _, peer := range h.staticPeers {
		err := h.Send(msg, peer, ) //How to pass in a conn?)
		if err != nil {
			return errors.Wrap(err, "error broadcasting: ")
		}
	}

	return nil
}

func (h *HobbitsNode) Listen() (net.Conn, error) {
	server := tcp.NewServer(h.host, h.port)
	ch := make(chan encoding.Message)

	go server.Listen(func(_ net.Conn, message encoding.Message){
		ch <- message
	})

}