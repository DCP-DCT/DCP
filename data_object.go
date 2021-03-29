package DCP

import "github.com/didiercrunch/paillier"

type DataObject struct {
	Plaintext int64
	Counter int
	LatestPk *paillier.PublicKey
}
