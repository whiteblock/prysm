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
		host:      host,
		port:      port,
		staticPeers:     peers,
		peerConns: []net.Conn{},
		feeds:     map[reflect.Type]p2p.Feed{},
	}
}

func (h *HobbitsNode) OpenConns() error { // TODO: check to see if there is an err that you can actually handle
														// TODO: use an atomic swap
	conns := make(chan net.Conn)

	for _, p := range h.staticPeers {
		go func(){
			conn, _ := net.Dial("tcp", p)

			conns <- conn
		}
	}

	h.peerConns = append(h.peerConns, ) // TODO:


	return nil
}

func (h *HobbitsNode) Listen() error {
	h.server = tcp.NewServer(h.host, h.port)

	return h.server.Listen(func(conn net.Conn, message encoding.Message) {
		h.processHobbitsMessage(HobbitsMessage(message), conn)
	})
}

func (h *HobbitsNode) Broadcast(message HobbitsMessage) error {
	for _, peer := range h.peerConns {
		err := h.server.SendMessage(peer, encoding.Message(message))
		if err != nil {
			return errors.Wrap(err, "error broadcasting: ")
		}

		peer.Close() // TODO: do I wanna be closing the conns?
	}

	return nil
}
