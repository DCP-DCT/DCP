package DCP

import (
	"fmt"
	"github.com/google/uuid"
)

type ICtNode interface {
	InitRoutine(Prepare) error
	Broadcast()
	Listen()
	HandleCalculationObject(interface{})
	Print()
}

type CtNode struct {
	Id             uuid.UUID
	Co             *CalculationObjectPaillier
	Ids            []string
	ReachableNodes []chan *CalculationObjectPaillier
	Channel        chan *CalculationObjectPaillier
}

func InitRoutine(fn Prepare, node *CtNode) error {
	e := fn(node)
	return e
}

func (node *CtNode) Broadcast() {
	fmt.Printf("Broadcasting triggered in node %s\n", node.Id)
	for _, rn := range node.ReachableNodes {
		go func(rn chan *CalculationObjectPaillier) {
			rn <- node.Co
		}(rn)
	}
}

func (node *CtNode) Listen() {
	go func() {
		for {
			co := <-node.Channel
			fmt.Printf("Listen triggered in node %s\n", node.Id)
			node.HandleCalculationObject(co)
		}
	}()
}

func (node *CtNode) HandleCalculationObject(co *CalculationObjectPaillier) {
	// Run Eval
	// Broadcast

	// Check PK

	// Check counter

	if node.Co.PublicKey.N.Cmp(co.PublicKey.N) == 0 {
		fmt.Println("Public key match")
	}

	idLen := len(node.Ids)

	cipher, e := co.Encrypt(idLen)
	if e != nil {
		// No-op
		fmt.Println(e.Error())
		return
	}

	// Add to co cipher
	co.Add(cipher)
	co.Counter = co.Counter + 1

	node.Broadcast()
}

func (node CtNode) Print() {
	fmt.Printf("Counter %d, PK %s, SK %s\n", node.Co.Counter, node.Co.privateKey, node.Co.PublicKey)
}
