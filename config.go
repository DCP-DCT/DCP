package DCP

type CtNodeConfig struct {
	NodeVisitDecryptThreshold int
	SuppressLogging bool
}

func (conf *CtNodeConfig) GetThreshold() int {
	if conf.NodeVisitDecryptThreshold == 0 {
		conf.NodeVisitDecryptThreshold = defaultNodeVisitDecryptThreshold
	}

	return conf.NodeVisitDecryptThreshold
}
