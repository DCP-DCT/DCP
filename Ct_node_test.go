package DCP

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestCtPrint(t *testing.T) {
	node := CtNode{
		Id:             uuid.New(),
		Co:             &CalculationObjectPaillier{},
		Ids:            nil,
		ReachableNodes: nil,
		Channel:        nil,
	}
	_ = node.Co.KeyGen()

	node.Print()
}

func TestCtChannel(t *testing.T) {
	node1 := &CtNode{
		Id: uuid.New(),
		Co: &CalculationObjectPaillier{
			Id:        uuid.New(),
			Counter:   0,
			PublicKey: nil,
			Cipher:    nil,
		},
		Ids:            []string{uuid.New().String(), uuid.New().String()},
		ReachableNodes: nil,
		Channel:        make(chan *CalculationObjectPaillier),
		HandledCoIds:   make(map[uuid.UUID]struct{}),
	}

	node2 := &CtNode{
		Id: uuid.New(),
		Co: &CalculationObjectPaillier{
			Id:        uuid.New(),
			Counter:   0,
			PublicKey: nil,
			Cipher:    nil,
		},
		Ids:            []string{uuid.New().String(), uuid.New().String()},
		ReachableNodes: nil,
		Channel:        make(chan *CalculationObjectPaillier),
		HandledCoIds:   make(map[uuid.UUID]struct{}),
	}

	_ = node1.Co.KeyGen()
	_ = node2.Co.KeyGen()

	node1.ReachableNodes = append(node1.ReachableNodes, node2.Channel)
	node2.ReachableNodes = append(node2.ReachableNodes, node1.Channel)

	_ = InitRoutine(PrepareIdLenCalculation, node1)

	node1.Listen()
	node2.Listen()
	node1.Broadcast(nil)

	time.Sleep(1 * time.Millisecond)
}

func TestAbortAlreadyHandled(t *testing.T) {
	node1 := &CtNode{
		Id: uuid.New(),
		Co: &CalculationObjectPaillier{
			Id:        uuid.New(),
			Counter:   0,
			PublicKey: nil,
			Cipher:    nil,
		},
		Ids:            []string{uuid.New().String(), uuid.New().String()},
		ReachableNodes: nil,
		Channel:        make(chan *CalculationObjectPaillier),
		HandledCoIds:   make(map[uuid.UUID]struct{}),
	}

	node2 := &CtNode{
		Id: uuid.New(),
		Co: &CalculationObjectPaillier{
			Id:        uuid.New(),
			Counter:   0,
			PublicKey: nil,
			Cipher:    nil,
		},
		Ids:            []string{uuid.New().String(), uuid.New().String()},
		ReachableNodes: nil,
		Channel:        make(chan *CalculationObjectPaillier),
		HandledCoIds:   make(map[uuid.UUID]struct{}),
	}

	node3 := &CtNode{
		Id: uuid.New(),
		Co: &CalculationObjectPaillier{
			Id:        uuid.New(),
			Counter:   0,
			PublicKey: nil,
			Cipher:    nil,
		},
		Ids:            []string{uuid.New().String(), uuid.New().String()},
		ReachableNodes: nil,
		Channel:        make(chan *CalculationObjectPaillier),
		HandledCoIds:   make(map[uuid.UUID]struct{}),
	}

	_ = node1.Co.KeyGen()
	_ = node2.Co.KeyGen()
	_ = node3.Co.KeyGen()

	node1.ReachableNodes = append(node1.ReachableNodes, node2.Channel)
	node2.ReachableNodes = append(node2.ReachableNodes, node3.Channel)
	node3.ReachableNodes = append(node3.ReachableNodes, node2.Channel)

	_ = InitRoutine(PrepareIdLenCalculation, node1)

	node1.Listen()
	node2.Listen()
	node3.Listen()
	node1.Broadcast(nil)

	time.Sleep(1 * time.Millisecond)
}
