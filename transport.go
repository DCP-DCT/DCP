package DCP

import (
	"github.com/google/uuid"
	"log"
	"time"
)

type Handler func([]byte) error

type OnTrigger func()

type Transport interface {
	Listen(nodeId uuid.UUID, handler Handler)
	Broadcast(nodeId uuid.UUID, obj []byte, onTrigger OnTrigger)
}

type ChannelTransport struct {
	DataCh          chan []byte `json:"-"`
	ReachableNodes  map[chan []byte]struct{} `json:"-"`
	SuppressLogging bool `json:"-"`
	Throttle        *time.Duration `json:"-"`
}

func (chT *ChannelTransport) Listen(nodeId uuid.UUID, handler Handler) {
	go func() {
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
}

func (chT *ChannelTransport) Broadcast(nodeId uuid.UUID, obj []byte, onTrigger OnTrigger) {
	onTrigger()

	for rn := range chT.ReachableNodes {
		go func(rn chan []byte) {
			select {
			case rn <- obj:
			default:
			}
		}(rn)
	}
}
