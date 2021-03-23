package DCP

import (
	"crypto/rand"
	"github.com/didiercrunch/paillier"
	"math/big"
)

func EncryptIdsPaillier(pk *paillier.PublicKey, idsLen int) (*paillier.Cypher, error) {
	cipher, e := pk.Encrypt(big.NewInt(int64(idsLen)), rand.Reader)
	if e != nil {
		return nil, e
	}

	return cipher, nil
}
