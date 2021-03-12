package DCP

import (
	"github.com/didiercrunch/paillier"
	"math/big"
)

type ICalculationObject interface {
	Add()
	Mul()
	Encrypt()
	Decrypt()
	KeyGen()
}

type CalculationObjectPaillier struct {
	counter int
	privateKey *paillier.PrivateKey
	publicKey paillier.PublicKey
	Cipher paillier.Cypher
}

func (cop *CalculationObjectPaillier) KeyGen() {
	cop.privateKey = paillier.CreatePrivateKey(big.NewInt(1), big.NewInt(3))
	cop.publicKey = cop.privateKey.PublicKey
}

func (cop CalculationObjectPaillier) Encrypt() {
	
}

func (cop CalculationObjectPaillier) Decrypt() {

}

