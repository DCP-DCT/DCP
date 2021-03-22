package DCP

import (
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
	config := &CtNodeConfig{
		NodeVisitDecryptThreshold: 2,
	}

	node1 := NewCtNode([]string{uuid.New().String()}, config)
	node2 := NewCtNode([]string{uuid.New().String()}, config)

	_ = node1.Co.KeyGen()
	_ = node2.Co.KeyGen()

	node1.TransportLayer.ReachableNodes[node2.TransportLayer.DataCh] = node2.TransportLayer.StopCh
	node2.TransportLayer.ReachableNodes[node1.TransportLayer.DataCh] = node1.TransportLayer.StopCh

	_ = InitRoutine(PrepareIdLenCalculation, node1)

	node1.Listen()
	node2.Listen()
	node1.Broadcast(nil)

	time.Sleep(1 * time.Millisecond)
}

func TestCtNode_HandleCalculationObjectAbortAlreadyHandled(t *testing.T) {
	node1 := NewCtNode([]string{uuid.New().String(), uuid.New().String()}, &CtNodeConfig{})
	node2 := NewCtNode([]string{uuid.New().String(), uuid.New().String()}, &CtNodeConfig{})
	node3 := NewCtNode([]string{uuid.New().String(), uuid.New().String()}, &CtNodeConfig{})

	_ = node1.Co.KeyGen()
	_ = node2.Co.KeyGen()
	_ = node3.Co.KeyGen()

	node1.TransportLayer.ReachableNodes[node2.TransportLayer.DataCh] = node2.TransportLayer.StopCh
	node2.TransportLayer.ReachableNodes[node3.TransportLayer.DataCh] = node3.TransportLayer.StopCh
	node3.TransportLayer.ReachableNodes[node2.TransportLayer.DataCh] = node2.TransportLayer.StopCh

	_ = InitRoutine(PrepareIdLenCalculation, node1)

	node1.Listen()
	node2.Listen()
	node3.Listen()
	node1.Broadcast(nil)

	time.Sleep(1 * time.Millisecond)
}

func TestCtNode_HandleCalculationObjectUpdateSelfNodeCo(t *testing.T) {
	node1 := NewCtNode([]string{uuid.New().String(), uuid.New().String()}, &CtNodeConfig{})
	node2 := NewCtNode([]string{uuid.New().String(), uuid.New().String()}, &CtNodeConfig{})

	node1.Co.Counter = defaultNodeVisitDecryptThreshold - 1

	_ = node1.Co.KeyGen()
	_ = node2.Co.KeyGen()

	node1.TransportLayer.ReachableNodes[node2.TransportLayer.DataCh] = node2.TransportLayer.StopCh
	node2.TransportLayer.ReachableNodes[node1.TransportLayer.DataCh] = node1.TransportLayer.StopCh

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

func TestCtNode_HandleCalculationObjectMerge(t *testing.T) {
	node1 := NewCtNode([]string{uuid.New().String(), uuid.New().String()}, &CtNodeConfig{})
	node2 := NewCtNode([]string{uuid.New().String(), uuid.New().String()}, &CtNodeConfig{})
	node3 := NewCtNode([]string{uuid.New().String(), uuid.New().String()}, &CtNodeConfig{})
	node4 := NewCtNode([]string{uuid.New().String(), uuid.New().String()}, &CtNodeConfig{})

	_ = node1.Co.KeyGen()
	_ = node2.Co.KeyGen()
	_ = node3.Co.KeyGen()
	_ = node4.Co.KeyGen()

	node1.TransportLayer.ReachableNodes[node2.TransportLayer.DataCh] = node2.TransportLayer.StopCh
	node1.TransportLayer.ReachableNodes[node3.TransportLayer.DataCh] = node3.TransportLayer.StopCh
	node2.TransportLayer.ReachableNodes[node1.TransportLayer.DataCh] = node1.TransportLayer.StopCh
	node3.TransportLayer.ReachableNodes[node4.TransportLayer.DataCh] = node4.TransportLayer.StopCh
	node4.TransportLayer.ReachableNodes[node1.TransportLayer.DataCh] = node1.TransportLayer.StopCh

	_ = InitRoutine(PrepareIdLenCalculation, node1)

	node1.Listen()
	node2.Listen()
	node3.Listen()
	node4.Listen()
	node1.Broadcast(nil)

	if node1.Co.Counter != 4 {
		t.Fail()
	}

	msg := node1.Co.Decrypt(node1.Co.Cipher)
	if msg.String() != "8" {
		t.Fail()
	}

	time.Sleep(1 * time.Millisecond)
}
