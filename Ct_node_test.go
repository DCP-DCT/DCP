package DCP

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestCtPrint(t *testing.T) {
	node := CtNode{
		Id: 			uuid.New(),
		Co:             CalculationObjectPaillier{},
		Ids:            nil,
		ReachableNodes: nil,
		Channel:        nil,
	}
	node.Co.KeyGen()

	node.Print()
}

func TestCtChannel(t *testing.T) {
	node1 := &CtNode{
		Id: 			uuid.New(),
		Co:             CalculationObjectPaillier{},
		Ids:            nil,
		ReachableNodes: nil,
		Channel:     make(chan CalculationObjectPaillier),
	}

	node1.Co.Counter = 10

	node2 := &CtNode{
		Id: 			uuid.New(),
		Co:             CalculationObjectPaillier{},
		Ids:            nil,
		ReachableNodes: nil,
		Channel:     make(chan CalculationObjectPaillier),
	}

	node1.ReachableNodes = append(node1.ReachableNodes, node2.Channel)

	node2.Listen()
	node1.Broadcast()

	time.Sleep(1 * time.Second)
}
