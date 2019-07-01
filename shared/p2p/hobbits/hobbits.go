package hobbits

import (
	"reflect"

	"github.com/renaynay/go-hobbits/encoding"
	"github.com/renaynay/go-hobbits/tcp"
	"github.com/renaynay/prysm/shared/p2p"
)

type HobbitsNode struct {
	host        string
	port        int
	staticPeers []string
	feeds       map[reflect.Type]p2p.Feed
	server      *tcp.Server
}

type HobbitsMessage encoding.Message

type RPCMethod int

const ( // TODO: should I integrate this with messages.proto? would that make sense?
	HELLO RPCMethod = iota
	GOODBYE
	GET_STATUS
	GET_BLOCK_HEADERS = iota + 62
	BLOCK_HEADERS
	GET_BLOCK_BODIES
	BLOCK_BODIES
	GET_ATTESTATION  //TODO: define in the spec what hex this corresponds to
	ATTESTATION      // TODO: define in the spec what this means
)

var topicMapping map[reflect.Type]string
// TODO: initialize with a const?

type GossipHeader struct {
	// TODO: if i don't care about the other values, how do I remove them so
	// TODO:  I don't spend too much unnecessary time unmarshaling?
	method_id      uint16   `bson:"method_id"`
	topic          string   `bson:"topic"`
	timestamp      uint32   `bson:"timestamp"`
	message_hash   [32]byte `bson:"message_hash"`
	hash_signature [32]byte `bson:"hash_signature"`
}
