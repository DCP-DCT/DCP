package DCP

type CtNodeConfig struct {
	nodeVisitDecryptThreshold int
}

func (conf *CtNodeConfig) GetThreshold() int {
	if conf.nodeVisitDecryptThreshold == 0 {
		conf.nodeVisitDecryptThreshold = defaultNodeVisitDecryptThreshold
	}

	return conf.nodeVisitDecryptThreshold
}
