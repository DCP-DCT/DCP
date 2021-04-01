package DCP

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

type Handler func([]byte)

type OnTrigger func()

type Transport interface {
	Listen(nodeId uuid.UUID, handler Handler)
	Broadcast(nodeId uuid.UUID, obj []byte, onTrigger OnTrigger)
}

type ChannelTransport struct {
	DataCh          chan []byte
	StopCh          chan struct{}
	ReachableNodes  map[chan []byte]chan struct{}
	SuppressLogging bool
	Throttle        *time.Duration
}

func (chT *ChannelTransport) Listen(nodeId uuid.UUID, handler Handler) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		for obj := range chT.DataCh {
			fmt.Println(len(chT.DataCh))
			logf(chT.SuppressLogging, "Listen triggered in node %s. Queue size: %d\n", nodeId, len(chT.DataCh))
			if obj != nil {
				if chT.Throttle != nil {
					time.Sleep(*chT.Throttle)
				}

				handler(obj)
			}
		}
	}()

	wg.Wait()
}

func (chT *ChannelTransport) Broadcast(nodeId uuid.UUID, obj []byte, onTrigger OnTrigger) {
	onTrigger()

	for rn := range chT.ReachableNodes {
		go func(rn chan []byte) {
			select {
			case rn <- obj:
				return
			default:
			}
		}(rn)
	}
}
