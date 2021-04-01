package DCP

import (
	"github.com/didiercrunch/paillier"
	"github.com/google/uuid"
)

type DataObject struct {
	Plaintext          int64               `json:"plaintext"`
	Counter            int                 `json:"counter"`
	LatestPk           *paillier.PublicKey `json:"latest_pk"`
	LatestBranchId     uuid.UUID           `json:"latest_branch_id"`
	DiscardedBranchIds []uuid.UUID         `json:"discarded_branch_ids"`
	Iteration          int                 `json:"iteration"`
}
