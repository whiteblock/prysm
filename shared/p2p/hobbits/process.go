package hobbits

import (
	"context"
	"fmt"
	"net"
	"reflect"

	"github.com/prysmaticlabs/prysm/shared/hashutil"
	"github.com/renaynay/go-hobbits/encoding"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"

	"github.com/pkg/errors"
)

func (h *HobbitsNode) processHobbitsMessage(message HobbitsMessage, conn net.Conn) error {
	switch message.Protocol {
	case encoding.RPC:
		log.Debug("beginning to process the RPC message...")

		err := h.processRPC(message, conn)
		if err != nil {
			log.WithFields(log.Fields{"err":err}).Trace("there was an error processing an RPC hobbits msg ") // TODO DELETE
			return errors.Wrap(err, "error processing an RPC hobbits message")
		}
	case encoding.GOSSIP:
		log.Debug("beginning to process the GOSSIP message...")

		err := h.processGossip(message)
		if err != nil {
			return errors.Wrap(err, "error processing a GOSSIP hobbits message")
		}

		return nil
	}

	log.Debug("protocol unsupported")
	return errors.New("protocol unsupported")
}

func (h *HobbitsNode) processRPC(message HobbitsMessage, conn net.Conn) error {
	method, err := h.parseMethodID(message.Header)
	if err != nil {
		log.Debug("method id could not be parsed from message header")
		return errors.Wrap(err, "could not parse method_id: ")
	}

	switch method {
	case HELLO:
		fmt.Println("HELLO received") // TODO delete
		log.Debug("HELLO received")

		response := h.rpcHello()

		fmt.Println("a HELLO response has been built...") // TODO Delete
		log.Debug("a HELLO response has been built...")

		responseBody, err := bson.Marshal(response)

		responseMessage := HobbitsMessage{
			Version:  message.Version,
			Protocol: message.Protocol,
			Header:   message.Header,
			Body:     responseBody,
		}

		err = h.Server.SendMessage(conn, encoding.Message(responseMessage))
		if err != nil {
			log.Debug("error sending a HELLO back") // TODO delete
			fmt.Println("sending HELLO...") // TODO delete
			return errors.Wrap(err, "error sending hobbits message: ")
		}

		log.Debug("sending HELLO...")
	case GOODBYE:
		log.Debug("hobbs: GOODBYTE")
		err := h.removePeer(conn)
		if err != nil {
			return errors.Wrap(err, "error handling GOODBYE: ")
		}
	case GET_STATUS:
		log.Debug("hobbs: GET_STATUS")
		// TODO: retrieve data and call h.Send
	case GET_BLOCK_HEADERS:
		log.Debug("hobbs: GET_BLOCK_HEADERS")
		// TODO: retrieve data and call h.Send
	case BLOCK_HEADERS:
		// TODO: call Broadcast?
	case GET_BLOCK_BODIES:
		// TODO: how will this work if batchedBeaconBlockRequest uses finalizedRoot and CanonicalRoot?
	case BLOCK_BODIES:
		// TODO: call Broadcast?
	case GET_ATTESTATION:
		// TODO: retrieve data and call h.Send
	case ATTESTATION:
		// TODO: retrieve data and call h.Send
	}

	return nil
}

func (h *HobbitsNode) rpcHello() Hello { // TODO: this is garbage
	log.Debug("building a HELLO response...")

	var response Hello

	response.NodeID = h.NodeId

	//response.BestRoot = h.DB.HeadStateRoot()
	response.BestRoot = [32]byte{}

	//headState, err := h.DB.HeadState(context.Background())
	//if err != nil {
	//	response.BestSlot = 0
	//}
	//if headState == nil {
	//	response.BestSlot = 0
	//}

	response.BestSlot = 0

	finalizedState, err := h.DB.FinalizedState()
	if err != nil {
		finalizedState = nil
	}
	if finalizedState == nil {
		response.LatestFinalizedEpoch = 0
	}
	if finalizedState != nil {
		response.LatestFinalizedEpoch = finalizedState.Slot / 64 // finalized epoch
	}

	hashedFinalizedState, err := hashutil.HashProto(finalizedState) // finalized root
	if err != nil {
		response.LatestFinalizedRoot = [32]byte{}
	}
	response.LatestFinalizedRoot = hashedFinalizedState

	return response
}

func (h *HobbitsNode) removePeer(peer net.Conn) error {
	index := 0

	for i, p := range h.PeerConns {
		if reflect.DeepEqual(peer, p) {
			index = i
		}
	}
	if index == 0 {
		return errors.New("error removing peer from node's static peers")
	}

	h.PeerConns = append(h.PeerConns[:index], h.PeerConns[index+1:]...) // TODO: is there a better way to delete
	// TODO: an element from an array by its value?

	return nil
}

func (h *HobbitsNode) processGossip(message HobbitsMessage) error {
	log.Debug("processing GOSSIP message...")

	_, err := h.parseTopic(message)
	if err != nil {
		return errors.Wrap(err, "error parsing topic: ")
	}

	h.Broadcast(context.Background(), nil)

	return nil
}

func (h *HobbitsNode) parseMethodID(header []byte) (RPCMethod, error) {
	log.Debug("attempting to parse method ID from header of RPC message...")

	unmarshaledHeader := &RPCHeader{}

	err := bson.Unmarshal(header, unmarshaledHeader)
	if err != nil {
		log.Debug("could not unmarshal the header of the message: ") // TODO delete
		return RPCMethod(0), errors.Wrap(err, "could not unmarshal the header of the message: ")
	}

	return RPCMethod(unmarshaledHeader.MethodID), nil
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
