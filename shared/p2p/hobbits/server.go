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

func (h *HobbitsNode) Listen() error {
	h.server = tcp.NewServer(h.host, h.port)

	return h.server.Listen(func(conn net.Conn, message encoding.Message) {
		h.processHobbitsMessage(HobbitsMessage(message), conn)
	})
}

func (h *HobbitsNode) Broadcast(message HobbitsMessage) error { // TODO: can i pre-open connections and just loop over open conns instead?
	for _, peer := range h.staticPeers {
		conn, err := net.Dial("tcp", peer)

		err = h.server.SendMessage(conn, encoding.Message(message))
		if err != nil {
			return errors.Wrap(err, "error broadcasting: ")
		}

		conn.Close()
	}

	return nil
}

//if conn == nil {
//	conn, err := net.Dial("tcp", peer)
//
//	err = server.SendMessage(conn, encoding.Message(message))
//	if err != nil {
//		return errors.Wrap(err, "error sending message: ")
//	}
//	return nil
//}