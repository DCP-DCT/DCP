package DCP

import (
	"crypto/rand"
	"github.com/didiercrunch/paillier"
	"math/big"
	"time"
)

type IEval interface {
	Add()
}

type ICalculationObject interface {
	Add()
	Mul()
	Encrypt(int) error
	Decrypt()
	KeyGen() error
}

type CalculationObjectPaillier struct {
	Counter    int
	PrivateKey *paillier.PrivateKey
	PublicKey  paillier.PublicKey
	Cipher     *paillier.Cypher
}

func (cop *CalculationObjectPaillier) KeyGen() error {
	p1, _, e := paillier.GenerateSafePrime(128, 1, 1 * time.Second, rand.Reader)
	if e != nil {
		return e
	}

	p2, _, e := paillier.GenerateSafePrime(128, 1, 1 * time.Second, rand.Reader)
	if e != nil {
		return e
	}

	cop.PrivateKey = paillier.CreatePrivateKey(p1, p2)
	cop.PublicKey = cop.PrivateKey.PublicKey

	return nil
}

func (cop *CalculationObjectPaillier) Encrypt(plaintext int) (*paillier.Cypher, error) {
	c, e := cop.PrivateKey.Encrypt(big.NewInt(int64(plaintext)), rand.Reader)
	if e != nil {
		return nil, e
	}

	return c, nil
}

func (cop CalculationObjectPaillier) Decrypt() {

}

func (cop *CalculationObjectPaillier) Add(cipher *paillier.Cypher) {
	cop.Cipher = cop.PrivateKey.Add(cop.Cipher, cipher)
}
