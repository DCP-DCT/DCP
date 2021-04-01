package DCP

import (
	"crypto/rand"
	"fmt"
	"github.com/didiercrunch/paillier"
	"github.com/google/uuid"
	"math/big"
	"time"
)

type IEval interface {
	Add(cipher interface{})
}

type ICalculationObject interface {
	Add(cipher interface{})
	Mul(cipher interface{})
	Encrypt(int) error
	Decrypt(cipher interface{}) *big.Int
	KeyGen() error
	Serialize()
}

type CalculationObjectPaillier struct {
	Id         uuid.UUID  `json:"id"`
	BranchId   uuid.UUID `json:"branch_id"`
	Counter    int        `json:"counter"`
	privateKey *paillier.PrivateKey
	PublicKey  *paillier.PublicKey `json:"public_key"`
	Cipher     *paillier.Cypher    `json:"cipher"`
	Ttl        int                 `json:"ttl"`
}

func (cop *CalculationObjectPaillier) KeyGen() error {
	p1, _, e := paillier.GenerateSafePrime(128, 1, 1*time.Second, rand.Reader)
	if e != nil {
		return e
	}

	p2, _, e := paillier.GenerateSafePrime(128, 1, 1*time.Second, rand.Reader)
	if e != nil {
		return e
	}

	cop.privateKey = paillier.CreatePrivateKey(p1, p2)
	cop.PublicKey = &cop.privateKey.PublicKey

	return nil
}

func (cop *CalculationObjectPaillier) Encrypt(plaintext int) (*paillier.Cypher, error) {
	c, e := cop.PublicKey.Encrypt(big.NewInt(int64(plaintext)), rand.Reader)

	if e != nil {
		return nil, e
	}

	return c, nil
}

func (cop *CalculationObjectPaillier) Decrypt(cipher *paillier.Cypher) *big.Int {
	return cop.privateKey.Decrypt(cipher)
}

func (cop *CalculationObjectPaillier) Add(cipher *paillier.Cypher) {
	if cop.Cipher == nil {
		fmt.Println("Own Cipher nil")
		cop.Cipher = cipher
		return
	}

	if cipher == nil {
		fmt.Println("Supplied Cipher nil")
		return
	}

	cop.Cipher = cop.PublicKey.Add(cop.Cipher, cipher)
}

func (cop *CalculationObjectPaillier) Serialize() string {
	return ""
}
