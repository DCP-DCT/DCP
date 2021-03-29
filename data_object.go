package DCP

import (
	"github.com/didiercrunch/paillier"
	"github.com/google/uuid"
)

type DataObject struct {
	Plaintext          int64
	Counter            int
	LatestPk           *paillier.PublicKey
	LatestBranchId     *uuid.UUID
	DiscardedBranchIds []uuid.UUID
	Iteration          int
}
