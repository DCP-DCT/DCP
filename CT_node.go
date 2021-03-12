package DCP

import "fmt"

type ICtNode interface {
	Broadcast()
	Listen()
	Print()
}

type CtNode struct {
	Co CalculationObjectPaillier
	Ids []string
	ReachableNodes []CtNode
}

func (node CtNode) Broadcast() {

}

func (node CtNode) Listen() {

}

func (node CtNode) Print() {
	fmt.Printf("Counter %d, PK %s, SK %s\n", node.Co.Counter, node.Co.PrivateKey, node.Co.PublicKey)
}