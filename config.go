package DCP

import "time"

type CtNodeConfig struct {
	NodeVisitDecryptThreshold int
	SuppressLogging           bool
	Throttle                  *time.Duration
}

func (conf *CtNodeConfig) GetThreshold() int {
	if conf.NodeVisitDecryptThreshold == 0 {
		conf.NodeVisitDecryptThreshold = defaultNodeVisitDecryptThreshold
	}

	return conf.NodeVisitDecryptThreshold
}
