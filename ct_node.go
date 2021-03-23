package DCP

import (
	"encoding/json"
	"fmt"
	"github.com/didiercrunch/paillier"
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
	Id                         uuid.UUID
	Co                         *CalculationObjectPaillier
	Ids                        []string
	coProcessRunning           bool
	previousConcludedProcesses map[*paillier.PublicKey]struct{}
	HandledBranchIds           map[uuid.UUID]struct{}
	TransportLayer             *ChannelTransport
	Config                     *CtNodeConfig
	Diagnosis                  *Diagnosis
}

func NewCtNode(ids []string, config *CtNodeConfig) *CtNode {
	t := &ChannelTransport{
		DataCh:          make(chan []byte),
		StopCh:          make(chan struct{}),
		ReachableNodes:  make(map[chan []byte]chan struct{}),
		SuppressLogging: config.SuppressLogging,
		Throttle:        config.Throttle,
	}

	return &CtNode{
		Id: uuid.New(),
		Co: &CalculationObjectPaillier{
			Id:       uuid.New(),
			BranchId: nil,
			Counter:  0,
			Ttl:      config.CoTTL,
		},
		Ids:                        ids,
		coProcessRunning:           false,
		previousConcludedProcesses: make(map[*paillier.PublicKey]struct{}),
		HandledBranchIds:           make(map[uuid.UUID]struct{}),
		TransportLayer:             t,
		Config:                     config,
		Diagnosis:                  NewDiagnosis(),
	}
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

	logf(node.Config.SuppressLogging, "Broadcasting Node: %s BranchId: %s, Current count: %d\n", node.Id.String(), objToBroadcast.BranchId, objToBroadcast.Counter)
	node.Diagnosis.IncrementNumberOfBroadcasts()

	node.TransportLayer.Broadcast(node.Id, b, func() {
		if externalCo == nil && !node.coProcessRunning {
			node.coProcessRunning = true
		}
	})
}

func (node *CtNode) Listen() {
	go node.TransportLayer.Listen(node.Id, node.HandleCalculationObject)
}

func (node *CtNode) HandleCalculationObject(data []byte) bool {
	var co = &CalculationObjectPaillier{}
	e := json.Unmarshal(data, co)
	if e != nil {
		return false
	}

	co.Ttl = co.Ttl - 1
	if co.Ttl <= 0 {
		logf(node.Config.SuppressLogging, "CalculationObject: %s dropped due to expired ttl\n", co.Id.String())
		node.Diagnosis.IncrementNumberOfPacketsDropped()
		return false
	}

	if co.BranchId == nil {
		// First handle, set branch
		newBranchId := uuid.New()
		co.BranchId = &newBranchId
	} else if _, exist := node.HandledBranchIds[*co.BranchId]; exist {
		logf(node.Config.SuppressLogging, "BranchId: %s already handled\n", co.BranchId.String())
		node.Diagnosis.IncrementNumberOfDuplicates()

		if node.Co.Id != co.Id {
			node.Broadcast(co)
		}

		return false
	}

	if node.Co.PublicKey.N.Cmp(co.PublicKey.N) == 0 {
		node.Diagnosis.IncrementNumberOfPkMatches()

		if co.Counter >= node.Config.NodeVisitDecryptThreshold {
			logLn(node.Config.SuppressLogging, "Calculation process finished, updating internal CalculationObject")
			node.Diagnosis.IncrementNumberOfInternalUpdates()

			node.Co.Counter = co.Counter
			node.Co.Cipher = co.Cipher
			node.coProcessRunning = false
			node.previousConcludedProcesses[node.Co.PublicKey] = struct{}{}
			node.HandledBranchIds[*co.BranchId] = struct{}{}

			return false
		}

		logf(node.Config.SuppressLogging, "Too few participants (%d) to satisfy privacy. Still listening\n", co.Counter)
		node.Diagnosis.IncrementNumberOgRejectedDueToThreshold()

		node.Broadcast(co)
		return false
	}

	logf(node.Config.SuppressLogging, "Running update in node %s\n", node.Id)
	node.Diagnosis.IncrementNumberOfUpdates()

	idLen := len(node.Ids)
	cipher, e := co.Encrypt(idLen)
	if e != nil {
		fmt.Println(e.Error())
		return false
	}

	// Add to co cipher
	co.Add(cipher)
	co.Counter = co.Counter + 1

	node.HandledBranchIds[*co.BranchId] = struct{}{}

	node.Broadcast(co)
	return false
}

func (node *CtNode) IsCalculationProcessRunning() bool {
	return node.coProcessRunning
}

func (node CtNode) Print() {
	fmt.Printf("Counter %d, PK %s, SK %s\n", node.Co.Counter, node.Co.privateKey, node.Co.PublicKey)
}
