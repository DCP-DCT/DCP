package DCP

import (
	"github.com/google/uuid"
	"sync"
)

type Handler func(*[]byte) bool
type OnTrigger func()

type Transport interface {
	Listen(nodeId uuid.UUID, handler Handler)
	Broadcast(nodeId uuid.UUID, obj *[]byte, onTrigger OnTrigger)
}

type ChannelTransport struct {
	DataCh          chan *[]byte
	StopCh          chan struct{}
	ReachableNodes  map[chan *[]byte]chan struct{}
	SuppressLogging bool
}

func (chT *ChannelTransport) Listen(nodeId uuid.UUID, handler Handler) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		for obj := range chT.DataCh {
			logf(chT.SuppressLogging, "Listen triggered in node %s\n", nodeId)
			if obj != nil {
				finished := handler(obj)

				if finished {
					close(chT.StopCh)
					return
				}
			}
		}
	}()

	wg.Wait()
}

func (chT *ChannelTransport) Broadcast(nodeId uuid.UUID, obj *[]byte, onTrigger OnTrigger) {
	logf(chT.SuppressLogging, "Broadcast triggered in node %s\n", nodeId)
	onTrigger()

	for rn, stop := range chT.ReachableNodes {
		go func(rn chan *[]byte, stop chan struct{}) {
			for {
				select {
				case <-stop:
					logf(chT.SuppressLogging, "Stop channel triggered, aborting broadcast early, node: %s\n", nodeId)
					return
				case rn <- obj:
				}
			}
		}(rn, stop)
	}
}
