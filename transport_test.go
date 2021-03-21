package DCP

import (
	"encoding/json"
	"github.com/google/uuid"
	"testing"
	"time"
)

type Packet struct {
	Msg string `json:"msg"`
}

func TestChannelTransport_Broadcast(t *testing.T) {
	chT := ChannelTransport{
		DataCh:         make(chan []byte),
		StopCh:         make(chan struct{}),
		ReachableNodes: make(map[chan []byte]chan struct{}),
	}

	node1Chan := make(chan []byte)

	chT.ReachableNodes[node1Chan] = nil

	packet := &Packet{
		Msg: "foobar",
	}

	b, _ := json.Marshal(packet)
	chT.Broadcast(uuid.New(), b, func() {
		return
	})

	received := <-node1Chan

	var packetReceived Packet
	_ = json.Unmarshal(received, &packetReceived)

	if packetReceived.Msg != "foobar" {
		t.Fail()
	}
}

func TestChannelTransport_Listener(t *testing.T) {
	chT := ChannelTransport{
		DataCh:         make(chan []byte),
		StopCh:         make(chan struct{}),
		ReachableNodes: make(map[chan []byte]chan struct{}),
	}

	go chT.Listen(uuid.New(), func(i []byte) bool {
		var packetReceived Packet
		_ = json.Unmarshal(i, &packetReceived)

		if packetReceived.Msg != "foobar" {
			t.Fail()
			return false
		} else {
			return true
		}
	})

	packet := &Packet{
		Msg: "foobar",
	}

	b, _ := json.Marshal(packet)

	chT.DataCh <- b
}

func TestChannelTransport_ListenerAndBroadcast(t *testing.T) {
	broadcaster := ChannelTransport{
		DataCh:         make(chan []byte),
		StopCh:         make(chan struct{}),
		ReachableNodes: make(map[chan []byte]chan struct{}),
	}

	listener := ChannelTransport{
		DataCh:         make(chan []byte),
		StopCh:         make(chan struct{}),
		ReachableNodes: make(map[chan []byte]chan struct{}),
	}

	broadcaster.ReachableNodes[listener.DataCh] = listener.StopCh
	listener.ReachableNodes[broadcaster.DataCh] = broadcaster.StopCh

	go listener.Listen(uuid.New(), func(i []byte) bool {
		return true
	})

	packet := &Packet{
		Msg: "foobar",
	}

	b, _ := json.Marshal(packet)

	broadcaster.Broadcast(uuid.New(), b, func() {
		return
	})

	time.Sleep(1 * time.Millisecond)
}
