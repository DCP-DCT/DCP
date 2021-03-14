package DCP

import (
	"fmt"
	"github.com/didiercrunch/paillier"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestCtPrint(t *testing.T) {
	node := CtNode{
		Id: 			uuid.New(),
		Co:             &CalculationObjectPaillier{},
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
		Co:             &CalculationObjectPaillier{
			Counter:    0,
			PrivateKey: nil,
			PublicKey:  paillier.PublicKey{},
			Cipher:     nil,
		},
		Ids:            []string{uuid.New().String(), uuid.New().String()},
		ReachableNodes: nil,
		Channel:     make(chan *CalculationObjectPaillier),
	}

	node2 := &CtNode{
		Id: 			uuid.New(),
		Co:             &CalculationObjectPaillier{
			Counter:    0,
			PrivateKey: nil,
			PublicKey:  paillier.PublicKey{},
			Cipher:     nil,
		},
		Ids:            []string{uuid.New().String(), uuid.New().String()},
		ReachableNodes: nil,
		Channel:     make(chan *CalculationObjectPaillier),
	}

	e := node1.Co.KeyGen()
	if e != nil {
		fmt.Println(e.Error())
		t.Fail()
	}
	_ = node2.Co.KeyGen()

	node1.ReachableNodes = append(node1.ReachableNodes, node2.Channel)

	_ = InitRoutine(PrepareIdLenCalculation, node1)
	fmt.Println(node1.Co.PrivateKey.Decrypt(node1.Co.Cipher))

	node2.Listen()
	node1.Broadcast()

	time.Sleep(1 * time.Second)
}
