package DCP

import "time"

var NrOfBranches int

type CtNodeConfig struct {
	NodeVisitDecryptThreshold int
	SuppressLogging           bool
	Throttle                  *time.Duration
	CoTTL                     int
	IncludeHistory            bool
	DropHandledAfter          int
}

func NewCtNodeConfig() CtNodeConfig {
	return CtNodeConfig{
		NodeVisitDecryptThreshold: defaultNodeVisitDecryptThreshold,
		SuppressLogging:           false,
		Throttle:                  nil,
		CoTTL:                     defaultCalculationObjectTTL,
		IncludeHistory:            true,
		DropHandledAfter:          -1,
	}
}

func GetNrOfActiveBranches() int {
	return NrOfBranches
}
