package DCP

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/ivpusic/grpool"
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

	pool := grpool.NewPool(100, 10)
	defer pool.Release()

	node1 := NewCtNode([]string{uuid.New().String()}, config, pool)
	node2 := NewCtNode([]string{uuid.New().String()}, config, pool)

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
	pool := grpool.NewPool(100, 10)
	defer pool.Release()
	node1 := NewCtNode([]string{uuid.New().String(), uuid.New().String()}, NewCtNodeConfig(), pool)
	node2 := NewCtNode([]string{uuid.New().String(), uuid.New().String()}, NewCtNodeConfig(), pool)
	node3 := NewCtNode([]string{uuid.New().String(), uuid.New().String()}, NewCtNodeConfig(), pool)

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
	fmt.Println(NrOfBranches)
}

func TestCtNode_HandleCalculationObjectUpdateSelfNodeCo(t *testing.T) {
	pool := grpool.NewPool(100, 10)
	defer pool.Release()
	node1 := NewCtNode([]string{uuid.New().String(), uuid.New().String()}, NewCtNodeConfig(), pool)
	node2 := NewCtNode([]string{uuid.New().String(), uuid.New().String()}, NewCtNodeConfig(), pool)

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
	pool := grpool.NewPool(100, 10)
	defer pool.Release()

	c := NewCtNodeConfig()
	node := NewCtNode([]string{uuid.New().String()}, c, pool)

	_, e := json.Marshal(node)
	if e != nil {
		fmt.Println(e.Error())
		t.Fail()
	}
}
