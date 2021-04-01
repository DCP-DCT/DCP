package DCP

import (
	"encoding/json"
	"fmt"
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
	config := NewCtNodeConfig()
	config.NodeVisitDecryptThreshold = 2

	node1 := NewCtNode([]string{uuid.New().String()}, config)
	node2 := NewCtNode([]string{uuid.New().String()}, config)

	_ = node1.Co.KeyGen()
	_ = node2.Co.KeyGen()

	node1.TransportLayer.ReachableNodes[node2.TransportLayer.DataCh] = struct{}{}
	node2.TransportLayer.ReachableNodes[node1.TransportLayer.DataCh] = struct{}{}

	_ = InitRoutine(PrepareIdLenCalculation, node1)

	node1.Listen()
	node2.Listen()
	node1.Broadcast(nil)

	time.Sleep(1 * time.Millisecond)
}

func TestCtNode_HandleCalculationObjectAbortAlreadyHandled(t *testing.T) {
	node1 := NewCtNode([]string{uuid.New().String(), uuid.New().String()}, NewCtNodeConfig())
	node2 := NewCtNode([]string{uuid.New().String(), uuid.New().String()}, NewCtNodeConfig())
	node3 := NewCtNode([]string{uuid.New().String(), uuid.New().String()}, NewCtNodeConfig())

	_ = node1.Co.KeyGen()
	_ = node2.Co.KeyGen()
	_ = node3.Co.KeyGen()

	node1.TransportLayer.ReachableNodes[node2.TransportLayer.DataCh] = struct{}{}
	node2.TransportLayer.ReachableNodes[node3.TransportLayer.DataCh] = struct{}{}
	node3.TransportLayer.ReachableNodes[node2.TransportLayer.DataCh] = struct{}{}

	_ = InitRoutine(PrepareIdLenCalculation, node1)

	node1.Listen()
	node2.Listen()
	node3.Listen()
	node1.Broadcast(nil)

	time.Sleep(1 * time.Millisecond)
}

func TestCtNode_HandleCalculationObjectUpdateSelfNodeCo(t *testing.T) {
	node1 := NewCtNode([]string{uuid.New().String(), uuid.New().String()}, NewCtNodeConfig())
	node2 := NewCtNode([]string{uuid.New().String(), uuid.New().String()}, NewCtNodeConfig())

	node1.Co.Counter = defaultNodeVisitDecryptThreshold - 1

	_ = node1.Co.KeyGen()
	_ = node2.Co.KeyGen()

	node1.TransportLayer.ReachableNodes[node2.TransportLayer.DataCh] = struct{}{}
	node2.TransportLayer.ReachableNodes[node1.TransportLayer.DataCh] = struct{}{}

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

func Test_CtNodeMarshal(t *testing.T) {
	c := NewCtNodeConfig()
	node := NewCtNode([]string{uuid.New().String()}, c)

	_, e := json.Marshal(node)
	if e != nil {
		fmt.Println(e.Error())
		t.Fail()
	}
}
