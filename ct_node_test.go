package DCP

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestCtPrint(t *testing.T) {
	node := CtNode{
		Id:  uuid.New(),
		Co:  &CalculationObjectPaillier{},
		Ids: nil,
	}
	_ = node.Co.KeyGen()

	node.Print()
}

func TestCtNode_HandleCalculationObjectCtChannel(t *testing.T) {
	node1 := &CtNode{
		Id: uuid.New(),
		Co: &CalculationObjectPaillier{
			TransactionId: uuid.New(),
			Counter:       0,
			PublicKey:     nil,
			Cipher:        nil,
		},
		Ids:          []string{uuid.New().String(), uuid.New().String()},
		HandledCoIds: make(map[uuid.UUID]struct{}),
		TransportLayer: &ChannelTransport{
			DataCh:         make(chan *[]byte),
			StopCh:         make(chan struct{}),
			ReachableNodes: make(map[chan *[]byte]chan struct{}),
		},
		Config: &CtNodeConfig{
			NodeVisitDecryptThreshold: 2,
		},
		Diagnosis: NewDiagnosis(),
	}

	node2 := &CtNode{
		Id: uuid.New(),
		Co: &CalculationObjectPaillier{
			TransactionId: uuid.New(),
			Counter:       0,
			PublicKey:     nil,
			Cipher:        nil,
		},
		Ids:          []string{uuid.New().String(), uuid.New().String()},
		HandledCoIds: make(map[uuid.UUID]struct{}),
		TransportLayer: &ChannelTransport{
			DataCh:         make(chan *[]byte),
			StopCh:         make(chan struct{}),
			ReachableNodes: make(map[chan *[]byte]chan struct{}),
		},
		Config: &CtNodeConfig{
			NodeVisitDecryptThreshold: 2,
		},
		Diagnosis: NewDiagnosis(),
	}

	_ = node1.Co.KeyGen()
	_ = node2.Co.KeyGen()

	node1.TransportLayer.ReachableNodes[node2.TransportLayer.DataCh] = node2.TransportLayer.StopCh
	node2.TransportLayer.ReachableNodes[node1.TransportLayer.DataCh] = node1.TransportLayer.StopCh

	_ = InitRoutine(PrepareIdLenCalculation, node1)

	node1.Listen()
	node2.Listen()
	node1.Broadcast(nil)

	time.Sleep(1 * time.Millisecond)
}

func TestCtNode_HandleCalculationObjectAbortAlreadyHandled(t *testing.T) {
	node1 := &CtNode{
		Id: uuid.New(),
		Co: &CalculationObjectPaillier{
			TransactionId: uuid.New(),
			Counter:       0,
		},
		Ids:          []string{uuid.New().String(), uuid.New().String()},
		HandledCoIds: make(map[uuid.UUID]struct{}),
		TransportLayer: &ChannelTransport{
			DataCh:         make(chan *[]byte),
			StopCh:         make(chan struct{}),
			ReachableNodes: make(map[chan *[]byte]chan struct{}),
		},
		Config:    &CtNodeConfig{},
		Diagnosis: NewDiagnosis(),
	}

	node2 := &CtNode{
		Id: uuid.New(),
		Co: &CalculationObjectPaillier{
			TransactionId: uuid.New(),
			Counter:       0,
		},
		Ids:          []string{uuid.New().String(), uuid.New().String()},
		HandledCoIds: make(map[uuid.UUID]struct{}),
		TransportLayer: &ChannelTransport{
			DataCh:         make(chan *[]byte),
			StopCh:         make(chan struct{}),
			ReachableNodes: make(map[chan *[]byte]chan struct{}),
		},
		Config:    &CtNodeConfig{},
		Diagnosis: NewDiagnosis(),
	}

	node3 := &CtNode{
		Id: uuid.New(),
		Co: &CalculationObjectPaillier{
			TransactionId: uuid.New(),
			Counter:       0,
		},
		Ids:          []string{uuid.New().String(), uuid.New().String()},
		HandledCoIds: make(map[uuid.UUID]struct{}),
		TransportLayer: &ChannelTransport{
			DataCh:         make(chan *[]byte),
			StopCh:         make(chan struct{}),
			ReachableNodes: make(map[chan *[]byte]chan struct{}),
		},
		Config:    &CtNodeConfig{},
		Diagnosis: NewDiagnosis(),
	}

	_ = node1.Co.KeyGen()
	_ = node2.Co.KeyGen()
	_ = node3.Co.KeyGen()

	node1.TransportLayer.ReachableNodes[node2.TransportLayer.DataCh] = node2.TransportLayer.StopCh
	node2.TransportLayer.ReachableNodes[node3.TransportLayer.DataCh] = node3.TransportLayer.StopCh
	node3.TransportLayer.ReachableNodes[node2.TransportLayer.DataCh] = node2.TransportLayer.StopCh

	_ = InitRoutine(PrepareIdLenCalculation, node1)

	node1.Listen()
	node2.Listen()
	node3.Listen()
	node1.Broadcast(nil)

	time.Sleep(1 * time.Millisecond)
}

func TestCtNode_HandleCalculationObjectUpdateSelfNodeCo(t *testing.T) {
	node1 := &CtNode{
		Id: uuid.New(),
		Co: &CalculationObjectPaillier{
			TransactionId: uuid.New(),
			Counter:       defaultNodeVisitDecryptThreshold - 1,
			PublicKey:     nil,
			Cipher:        nil,
		},
		Ids:          []string{uuid.New().String(), uuid.New().String()},
		HandledCoIds: make(map[uuid.UUID]struct{}),
		TransportLayer: &ChannelTransport{
			DataCh:         make(chan *[]byte),
			StopCh:         make(chan struct{}),
			ReachableNodes: make(map[chan *[]byte]chan struct{}),
		},
		Config:    &CtNodeConfig{},
		Diagnosis: NewDiagnosis(),
	}

	node2 := &CtNode{
		Id: uuid.New(),
		Co: &CalculationObjectPaillier{
			TransactionId: uuid.New(),
			Counter:       0,
			PublicKey:     nil,
			Cipher:        nil,
		},
		Ids:          []string{uuid.New().String(), uuid.New().String()},
		HandledCoIds: make(map[uuid.UUID]struct{}),
		TransportLayer: &ChannelTransport{
			DataCh:         make(chan *[]byte),
			StopCh:         make(chan struct{}),
			ReachableNodes: make(map[chan *[]byte]chan struct{}),
		},
		Config:    &CtNodeConfig{},
		Diagnosis: NewDiagnosis(),
	}

	_ = node1.Co.KeyGen()
	_ = node2.Co.KeyGen()

	node1.TransportLayer.ReachableNodes[node2.TransportLayer.DataCh] = node2.TransportLayer.StopCh
	node2.TransportLayer.ReachableNodes[node1.TransportLayer.DataCh] = node1.TransportLayer.StopCh

	_ = InitRoutine(PrepareIdLenCalculation, node1)

	node1.Listen()
	node2.Listen()
	node1.Broadcast(nil)

	time.Sleep(1 * time.Millisecond)

	msg := node1.Co.Decrypt(node1.Co.Cipher)
	if msg.String() != "4" {
		t.Fail()
	}
}
