package DCP

import "time"

var NrOfBranches int

type CtNodeConfig struct {
	NodeVisitDecryptThreshold int
	SuppressLogging           bool
	Throttle                  *time.Duration
	CoTTL                     int
	IncludeHistory            bool
}

func NewCtNodeConfig() CtNodeConfig {
	return CtNodeConfig{
		NodeVisitDecryptThreshold: defaultNodeVisitDecryptThreshold,
		SuppressLogging:           false,
		Throttle:                  nil,
		CoTTL:                     defaultCalculationObjectTTL,
		IncludeHistory:            true,
	}
}

func GetNrOfActiveBranches() int {
	return NrOfBranches
}
