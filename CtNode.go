package DCP

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type ICtNode interface {
	InitRoutine(Prepare) error
	Broadcast()
	Listen()
	HandleCalculationObject(interface{}) bool
	Print()
}

type CtNode struct {
	Id             uuid.UUID
	Co             *CalculationObjectPaillier
	Ids            []string
	HandledCoIds   map[uuid.UUID]struct{}
	TransportLayer *ChannelTransport
	Config         *CtNodeConfig
}

func InitRoutine(fn Prepare, node *CtNode) error {
	e := fn(node)
	return e
}

func (node *CtNode) Broadcast(externalCo *CalculationObjectPaillier) {
	var objToBroadcast *CalculationObjectPaillier
	if externalCo != nil {
		objToBroadcast = externalCo
	} else {
		objToBroadcast = node.Co
	}

	b, e := json.Marshal(objToBroadcast)
	if e != nil {
		return
	}

	fmt.Printf("Broadcasting object: Id: %s, Current count: %d\n", objToBroadcast.Id, objToBroadcast.Counter)
	node.TransportLayer.Broadcast(node.Id, &b)
}

func (node *CtNode) Listen() {
	go node.TransportLayer.Listen(node.Id, node.HandleCalculationObject)
}

func (node *CtNode) HandleCalculationObject(data *[]byte) bool {
	var co *CalculationObjectPaillier = &CalculationObjectPaillier{}
	e := json.Unmarshal(*data, co)
	if e != nil {
		return false
	}

	if node.Co.PublicKey.N.Cmp(co.PublicKey.N) == 0 {
		fmt.Println("Public key match")
		if co.Counter >= node.Config.GetThreshold() {
			fmt.Println("Calculation process finished, updating internal CalculationObject")
			node.Co.Counter = co.Counter
			node.Co.Cipher = co.Cipher
			return true
		}

		fmt.Printf("Too few participants (%d) to satisfy privacy. Still listening\n", co.Counter)
		node.Broadcast(co)
		return false
	}

	if _, exist := node.HandledCoIds[co.Id]; exist {
		fmt.Printf("Calculation object with ID: %s already handled\n", co.Id.String())
		node.Broadcast(co)
		return false
	}

	fmt.Printf("Running update in node %s\n", node.Id)
	idLen := len(node.Ids)
	cipher, e := co.Encrypt(idLen)
	if e != nil {
		// No-op
		fmt.Println(e.Error())
		return false
	}

	// Add to co cipher
	co.Add(cipher)
	co.Counter = co.Counter + 1

	node.HandledCoIds[co.Id] = struct{}{}
	node.Broadcast(co)

	return false
}

func (node CtNode) Print() {
	fmt.Printf("Counter %d, PK %s, SK %s\n", node.Co.Counter, node.Co.privateKey, node.Co.PublicKey)
}
