package hobbits

import "github.com/pkg/errors"

func (h *HobbitsNode) processHobbitsMessage(message HobbitsMessage) error {
	switch message.Protocol {
	case "RPC":
		err := h.processRPC(message)
		if err != nil {
			return errors.Wrap(err, "error processing an RPC hobbits message")
		}
	case "GOSSIP":
		err := h.processGossip(message)
		if err != nil {
			return errors.Wrap(err, "error processing a GOSSIP hobbits message")
		}
	}

	return nil
}

func (h *HobbitsNode) processRPC(message HobbitsMessage) error {

}

func (h *HobbitsNode) processGossip(message HobbitsMessage) error {

}
