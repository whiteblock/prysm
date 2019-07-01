package hobbits

import (
	"reflect"

	"github.com/renaynay/go-hobbits/encoding"
	"github.com/renaynay/prysm/shared/p2p"
)

type HobbitsNode struct {
	host        string
	port        int
	staticPeers []string
	feeds       map[reflect.Type]p2p.Feed
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
