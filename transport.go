package DCP

import (
	"github.com/google/uuid"
	"log"
	"sync"
	"time"
)

type Handler func([]byte) error

type OnTrigger func()

type Transport interface {
	Listen(nodeId uuid.UUID, handler Handler)
	Broadcast(nodeId uuid.UUID, obj []byte, onTrigger OnTrigger)
}

type ChannelTransport struct {
	DataCh          chan []byte
	ReachableNodes  map[chan []byte]struct{}
	SuppressLogging bool
	Throttle        *time.Duration
}

func (chT *ChannelTransport) Listen(nodeId uuid.UUID, handler Handler) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		for obj := range chT.DataCh {
			if obj != nil {
				if chT.Throttle != nil {
					time.Sleep(*chT.Throttle)
				}

				e := handler(obj)
				if e != nil {
					log.Panic(e.Error())
				}
			}
		}
	}()

	wg.Wait()
}

func (chT *ChannelTransport) Broadcast(nodeId uuid.UUID, obj []byte, onTrigger OnTrigger) {
	onTrigger()

	for rn := range chT.ReachableNodes {
		go func(rn chan []byte) {
			rn <- obj
		}(rn)
	}
}
