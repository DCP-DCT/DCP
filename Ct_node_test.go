package DCP

import (
	"testing"
	"time"
)

func TestCtPrint(t *testing.T) {
	node := CtNode{
		Co:             CalculationObjectPaillier{},
		Ids:            nil,
		ReachableNodes: nil,
		channel:        nil,
	}
	node.Co.KeyGen()

	node.Print()
}

func TestCtChannel(t *testing.T) {
	node1 := &CtNode{
		Co:             CalculationObjectPaillier{},
		Ids:            nil,
		ReachableNodes: nil,
		channel:     make(chan CalculationObjectPaillier),
	}

	node1.Co.Counter = 10

	node2 := &CtNode{
		Co:             CalculationObjectPaillier{},
		Ids:            nil,
		ReachableNodes: nil,
		channel:     make(chan CalculationObjectPaillier),
	}

	node1.ReachableNodes = append(node1.ReachableNodes, node2.channel)

	node2.Listen()
	node1.Broadcast()

	time.Sleep(1 * time.Second)
}
