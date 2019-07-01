package hobbits

import (
	"net"
	"github.com/renaynay/go-hobbits/encoding"
	"gopkg.in/mgo.v2/bson"

	"github.com/pkg/errors"
)

func (h *HobbitsNode) processHobbitsMessage(message HobbitsMessage, conn net.Conn) error {
	switch message.Protocol {
	case "RPC":
		err := h.processRPC(message, conn)
		if err != nil {
			return errors.Wrap(err, "error processing an RPC hobbits message")
		}

		return nil
	case "GOSSIP":
		err := h.processGossip(message)
		if err != nil {
			return errors.Wrap(err, "error processing a GOSSIP hobbits message")
		}

		return nil
	}

	return errors.New("protocol unsupported")
}

func (h *HobbitsNode) processRPC(message HobbitsMessage, conn net.Conn) error {
	method, err := h.parseMethodID(message.Header)
	if err != nil {
		return errors.Wrap(err, "could not parse method_id: ")
	}

	switch method {
	case HELLO:
		// TODO: retrieve data and send it

		h.staticPeers = append(h.staticPeers, conn.RemoteAddr().String())

		err := h.server.SendMessage(conn, encoding.Message(message))
		if err != nil {
			return errors.Wrap(err, "error sending hobbits message: ")
		}
	case GOODBYE:
		rem := conn.RemoteAddr().String()
		err := h.removePeer(rem)
		if err != nil {
			return errors.Wrap(err, "error handling GOODBYE: ")
		}
	case GET_STATUS:
		// TODO: retrieve data and call h.Send
	case GET_BLOCK_HEADERS:
		// TODO: retrieve data and call h.Send
	case BLOCK_HEADERS:
		// TODO: call Broadcast?
	case GET_BLOCK_BODIES:
		// TODO: retrieve data and call h.Send
	case BLOCK_BODIES:
		// TODO: call Broadcast?
	case GET_ATTESTATION:
		// TODO: retrieve data and call h.Send
	case ATTESTATION:
		// TODO: retrieve data and call h.Send
	}

	return nil
}

func (h *HobbitsNode) removePeer(peer string) error {
	index := 0

	for i, p := range h.staticPeers {
		if peer == p {
			index = i
		}
	}
	if index == 0 {
		return errors.New("error removing peer from node's static peers")
	}

	h.staticPeers = append(h.staticPeers[:index], h.staticPeers[index+1:]...) // TODO: is there a better way to delete
																			// TODO: an element from an array by its value?

	return nil
}

func (h *HobbitsNode) processGossip(message HobbitsMessage) error {
	_, err := h.parseTopic(message)
	if err != nil {
		return errors.Wrap(err, "error parsing topic: ")
	}

	err := h.Broadcast(message)
	if err != nil {
		return errors.Wrap(err, "error broadcasting: ")
	}

	return nil
}

func (h *HobbitsNode) parseMethodID(header []byte) (RPCMethod, error) {

}

// parseTopic takes care of parsing the topic and updating the node's feeds
func (h *HobbitsNode) parseTopic(message HobbitsMessage) (string, error) {
	header := GossipHeader{}

	err := bson.Unmarshal(message.Header, header)
	if err != nil {
		return "", errors.Wrap(err, "error unmarshaling gossip message header: ")
	}

	// TODO: checks against topicMapping?
	// TODO: somehow updates h.Feeds?
	return header.topic, nil
}
