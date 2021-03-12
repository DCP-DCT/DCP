package DCP

type ICtNode interface {
	Broadcast()
	Listen()
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