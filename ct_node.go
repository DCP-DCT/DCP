package DCP

import (
	"encoding/json"
	"fmt"
	"github.com/didiercrunch/paillier"
	"github.com/google/uuid"
	"github.com/ivpusic/grpool"
	"log"
)

type ICtNode interface {
	InitRoutine(Prepare) error
	Broadcast()
	Listen()
	HandleCalculationObject(interface{}) bool
	Print()
}

type CtNode struct {
	Id               uuid.UUID                  `json:"id"`
	Do               DataObject                 `json:"data_object"`
	Co               *CalculationObjectPaillier `json:"calculation_object"`
	Ids              []string                   `json:"ids"`
	ProcessRunning   bool                       `json:"process_running"`
	BridgeNode       bool                       `json:"bridge_node"`
	HandledBranchIds map[uuid.UUID]int          `json:"-"`
	TransportLayer   *ChannelTransport          `json:"-"`
	Config           CtNodeConfig               `json:"config"`
	Diagnosis        *Diagnosis                 `json:"diagnosis"`
	WorkerPool       *grpool.Pool               `json:"-"`
}

func NewCtNode(ids []string, config CtNodeConfig, poolPtr *grpool.Pool) *CtNode {
	t := &ChannelTransport{
		DataCh:          make(chan []byte, 30000),
		ReachableNodes:  make(map[chan []byte]struct{}),
		SuppressLogging: config.SuppressLogging,
		Throttle:        config.Throttle,
	}

	return &CtNode{
		Id: uuid.New(),
		Do: DataObject{
			Plaintext:          0,
			Counter:            0,
			LatestPk:           nil,
			LatestBranchId:     uuid.Nil,
			DiscardedBranchIds: nil,
			Iteration:          0,
		},
		Co: &CalculationObjectPaillier{
			Id:       uuid.New(),
			BranchId: uuid.Nil,
			Counter:  0,
			Ttl:      config.CoTTL,
		},
		Ids:              ids,
		ProcessRunning:   false,
		BridgeNode:       false,
		HandledBranchIds: make(map[uuid.UUID]int),
		TransportLayer:   t,
		Config:           config,
		Diagnosis:        NewDiagnosis(),
		WorkerPool:       poolPtr,
	}
}

func InitRoutine(fn Prepare, node *CtNode) error {
	e := fn(node)
	return e
}

func (node *CtNode) Broadcast(externalCo *CalculationObjectPaillier) {
	var objToBroadcast CalculationObjectPaillier
	if externalCo != nil {
		objToBroadcast = *externalCo
	} else {
		objToBroadcast = *node.Co
	}

	b, e := json.Marshal(objToBroadcast)
	if e != nil {
		log.Panic(e.Error())
		return
	}

	logf(node.Config.SuppressLogging, "Broadcasting Node: %s BranchId: %s, Current count: %d, Current TTL %d\n", node.Id.String(), objToBroadcast.BranchId, objToBroadcast.Counter, objToBroadcast.Ttl)
	node.Diagnosis.IncrementNumberOfBroadcasts()

	node.WorkerPool.JobQueue <- func() {
		node.TransportLayer.Broadcast(node.Id, b, func() {
			if externalCo == nil && !node.ProcessRunning {
				node.ProcessRunning = true

				node.Do.LatestPk = new(paillier.PublicKey)
				*node.Do.LatestPk = *objToBroadcast.PublicKey
				node.Do.Iteration = node.Do.Iteration + 1
			}
		})
	}
}

func (node *CtNode) Listen() {
	go node.TransportLayer.Listen(node.Id, node.HandleCalculationObject)
}

func (node *CtNode) HandleCalculationObject(data []byte) error {

	defer node.Diagnosis.Timers.Time(NewTimer("HandleCalculationObject"))

	var co = CalculationObjectPaillier{}
	e := json.Unmarshal(data, &co)
	if e != nil {
		return e
	}

	if co.BranchId == uuid.Nil {
		// First handle, set branch
		co.BranchId = uuid.New()
		NrOfBranches++
	}

	co.Ttl = co.Ttl - 1

	logf(node.Config.SuppressLogging, "Listen triggered in node %s. CoId: %s. CoBranchId: %s\n", node.Id, co.Id, co.BranchId)

	if node.Co.PublicKey.N.Cmp(co.PublicKey.N) == 0 {
		s, st := NewTimer("PublicKeyClause")

		node.handlePkMatch(co)

		node.Diagnosis.Timers.Time(s, st)
		return nil
	}

	if co.Ttl <= 0 {
		logf(node.Config.SuppressLogging, "CalculationObject branchId: %s dropped due to expired ttl by nodeId: %s\n", co.BranchId, node.Id)
		node.Diagnosis.IncrementNumberOfPacketsDropped()
		return nil
	}

	if _, exist := node.HandledBranchIds[co.BranchId]; exist {
		logf(node.Config.SuppressLogging, "BranchId: %s already handled in node: %s\n", co.BranchId, node.Id)
		node.Diagnosis.IncrementNumberOfDuplicates()

		if node.Config.DropHandledAfter == -1 {
			node.Broadcast(&co)
			return nil
		}

		if node.HandledBranchIds[co.BranchId] >= node.Config.DropHandledAfter {
			node.Diagnosis.IncrementNumberOfPacketsDropped()
			return nil
		}

		node.Broadcast(&co)
		return nil
	}

	logf(node.Config.SuppressLogging, "Running update in node: %s on branchId: %s\n", node.Id, co.BranchId)
	node.Diagnosis.IncrementNumberOfUpdates()
	defer node.Diagnosis.Timers.Time(NewTimer("UpdateCalculationObject"))

	e = node.updateCalculationObject(&co)
	if e != nil {
		return e
	}

	node.Broadcast(&co)

	return nil
}

func (node *CtNode) updateCalculationObject(co *CalculationObjectPaillier) error {
	idLen := len(node.Ids)
	cipher, e := co.Encrypt(idLen)
	if e != nil {
		return e
	}

	// Add to co cipher
	co.Add(cipher)
	co.Counter = co.Counter + 1

	if node.Config.IncludeHistory {
		node.Diagnosis.Control.RegisterContribution(co.Id, co.BranchId, len(node.Ids))
	}

	if _, exist := node.HandledBranchIds[co.BranchId]; exist {
		node.HandledBranchIds[co.BranchId] = node.HandledBranchIds[co.BranchId] + 1
	} else {
		node.HandledBranchIds[co.BranchId] = 1
	}

	return nil
}

func (node *CtNode) UpdateDo(old CalculationObjectPaillier, new CalculationObjectPaillier) {
	oldData := node.Co.Decrypt(old.Cipher)
	newData := node.Co.Decrypt(new.Cipher)

	if node.Do.LatestBranchId != uuid.Nil {
		node.Do.DiscardedBranchIds = append(node.Do.DiscardedBranchIds, node.Do.LatestBranchId)

		node.Do.Plaintext = node.Do.Plaintext - oldData.Int64()
		node.Do.Counter = node.Do.Counter - old.Counter
	}

	node.Do.Plaintext = node.Do.Plaintext + newData.Int64()
	node.Do.Counter = node.Do.Counter + new.Counter

	node.Do.LatestBranchId = new.BranchId
}

func (node *CtNode) handlePkMatch(co CalculationObjectPaillier) {
	node.Diagnosis.IncrementNumberOfPkMatches()

	if co.Counter >= node.Config.NodeVisitDecryptThreshold {
		node.Diagnosis.IncrementNumberOfInternalUpdates()

		if node.Co.Counter < co.Counter {
			logf(node.Config.SuppressLogging, "Updating accepted DO in node %s. BranchId: %s\n", node.Id, co.BranchId)
			node.UpdateDo(*node.Co, co)

			node.Co.Counter = co.Counter

			*node.Co.Cipher = *co.Cipher
			node.ProcessRunning = false
		}
	} else {
		logf(node.Config.SuppressLogging, "Too few participants (%d) to satisfy privacy. NodeId: %s\n", co.Counter, node.Id)
		node.Diagnosis.IncrementNumberOgRejectedDueToThreshold()

		node.Broadcast(&co)
	}
}

func (node CtNode) Print() {
	fmt.Printf("Counter %d, PK %s, SK %s\n", node.Co.Counter, node.Co.privateKey, node.Co.PublicKey)
}
