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
			Counter:   0,
			PublicKey: nil,
			Cipher:    nil,
		},
		Ids:            []string{uuid.New().String(), uuid.New().String()},
		ReachableNodes: nil,
		Channel:        make(chan *CalculationObjectPaillier),
	}

	node2 := &CtNode{
		Id: uuid.New(),
		Co: &CalculationObjectPaillier{
			Counter:   0,
			PublicKey: nil,
			Cipher:    nil,
		},
		Ids:            []string{uuid.New().String(), uuid.New().String()},
		ReachableNodes: nil,
		Channel:        make(chan *CalculationObjectPaillier),
	}

	_ = node1.Co.KeyGen()
	_ = node2.Co.KeyGen()

	node1.ReachableNodes = append(node1.ReachableNodes, node2.Channel)
	node2.ReachableNodes = append(node2.ReachableNodes, node1.Channel)

	_ = InitRoutine(PrepareIdLenCalculation, node1)

	node1.Listen()
	node2.Listen()
	node1.Broadcast()

	time.Sleep(500 * time.Nanosecond)
}
