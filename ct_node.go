package DCP

import (
	"encoding/json"
	"fmt"
	"github.com/didiercrunch/paillier"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"time"
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
	HandledBranchIds map[uuid.UUID]struct{}     `json:"handled_branch_ids"`
	TransportLayer   *ChannelTransport          `json:"-"`
	Config           *CtNodeConfig              `json:"config"`
	Diagnosis        *Diagnosis                 `json:"diagnosis"`
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
		Do: DataObject{
			Plaintext:          0,
			Counter:            0,
			LatestPk:           nil,
			LatestBranchId:     nil,
			DiscardedBranchIds: nil,
			Iteration:          0,
		},
		Co: &CalculationObjectPaillier{
			Id:       uuid.New(),
			BranchId: nil,
			Counter:  0,
			Ttl:      config.CoTTL,
		},
		Ids:              ids,
		ProcessRunning:   false,
		HandledBranchIds: make(map[uuid.UUID]struct{}),
		TransportLayer:   t,
		Config:           config,
		Diagnosis:        NewDiagnosis(),
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

	node.TransportLayer.Broadcast(node.Id, b, func() {
		if externalCo == nil && !node.ProcessRunning {
			node.ProcessRunning = true

			node.Do.LatestPk = new(paillier.PublicKey)
			*node.Do.LatestPk = *objToBroadcast.PublicKey
			node.Do.Iteration = node.Do.Iteration + 1
		}
	})
}

func (node *CtNode) Listen() {
	go node.TransportLayer.Listen(node.Id, node.HandleCalculationObject)
}

func (node *CtNode) RunRandomTrigger(stop chan struct{}) {
	rand.Seed(time.Now().UnixNano())

	for {
		select {
		case <-stop:
			return
		default:
			if node.ProcessRunning {
				break
			}

			nr := rand.Intn(10-1) + 1

			if nr == 5 {
				e := InitRoutine(PrepareIdLenCalculation, node)
				if e != nil {
					log.Panic(e.Error())
					return
				}

				fmt.Printf("Starting process for node %s\n", node.Id)
				node.Broadcast(nil)
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func (node *CtNode) HandleCalculationObject(data []byte) error {

	defer node.Diagnosis.Timers.Time(NewTimer("HandleCalculationObject"))

	var co = CalculationObjectPaillier{}
	e := json.Unmarshal(data, &co)
	if e != nil {
		return e
	}

	if co.BranchId == nil {
		// First handle, set branch
		newBranchId := uuid.New()
		co.BranchId = &newBranchId
	}

	co.Ttl = co.Ttl - 1

	logf(node.Config.SuppressLogging, "Listen triggered in node %s. CoId: %s. CoBranchId: %s\n", node.Id, co.Id, co.BranchId)

	if node.Co.PublicKey.N.Cmp(co.PublicKey.N) == 0 {
		s, st := NewTimer("PublicKeyClause")

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

		node.Diagnosis.Timers.Time(s, st)
		return nil
	}


	if co.Ttl <= 0 {
		logf(node.Config.SuppressLogging, "CalculationObject branchId: %s dropped due to expired ttl by nodeId: %s\n", co.BranchId, node.Id)
		node.Diagnosis.IncrementNumberOfPacketsDropped()
		return nil
	}

	if _, exist := node.HandledBranchIds[*co.BranchId]; exist {
		logf(node.Config.SuppressLogging, "BranchId: %s already handled in node: %s\n", co.BranchId, node.Id)
		node.Diagnosis.IncrementNumberOfDuplicates()

		node.Broadcast(&co)

		return nil
	}

	logf(node.Config.SuppressLogging, "Running update in node: %s on branchId: %s\n", node.Id, co.BranchId)
	node.Diagnosis.IncrementNumberOfUpdates()

	defer node.Diagnosis.Timers.Time(NewTimer("UpdateCalculationObject"))

	idLen := len(node.Ids)
	cipher, e := co.Encrypt(idLen)
	if e != nil {
		return e
	}

	// Add to co cipher
	co.Add(cipher)
	co.Counter = co.Counter + 1

	node.Diagnosis.Control.RegisterContribution(co.Id, *co.BranchId, len(node.Ids))
	node.HandledBranchIds[*co.BranchId] = struct{}{}

	node.Broadcast(&co)

	return nil
}

func (node *CtNode) UpdateDo(old CalculationObjectPaillier, new CalculationObjectPaillier) {
	oldData := node.Co.Decrypt(old.Cipher)
	newData := node.Co.Decrypt(new.Cipher)

	if node.Do.LatestBranchId != nil {
		node.Do.DiscardedBranchIds = append(node.Do.DiscardedBranchIds, *node.Do.LatestBranchId)

		node.Do.Plaintext = node.Do.Plaintext - oldData.Int64()
		node.Do.Counter = node.Do.Counter - old.Counter
	}

	node.Do.Plaintext = node.Do.Plaintext + newData.Int64()
	node.Do.Counter = node.Do.Counter + new.Counter

	node.Do.LatestBranchId = new.BranchId
}

func (node CtNode) Print() {
	fmt.Printf("Counter %d, PK %s, SK %s\n", node.Co.Counter, node.Co.privateKey, node.Co.PublicKey)
}
