package DCP

import "time"

type CtNodeConfig struct {
	NodeVisitDecryptThreshold int
	SuppressLogging           bool
	Throttle                  *time.Duration
	CoTTL                     int
}

func NewCtNodeConfig() CtNodeConfig {
	return CtNodeConfig{
		NodeVisitDecryptThreshold: defaultNodeVisitDecryptThreshold,
		SuppressLogging:           false,
		Throttle:                  nil,
		CoTTL:                     defaultCalculationObjectTTL,
	}
}
