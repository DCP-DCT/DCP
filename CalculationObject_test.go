package DCP

import (
	"fmt"
	"math/big"
	"testing"
)

func TestCalculationObjectPaillier_KeyGen(t *testing.T) {
	nodes := make([]CtNode, 10)

	for _, node := range nodes {
		node = CtNode{
			Co:             &CalculationObjectPaillier{},
			Ids:            nil,
			ReachableNodes: nil,
		}

		e := node.Co.KeyGen()

		if e != nil {
			fmt.Println(e.Error())
			t.Fail()
		}
	}
}

func TestCalculationObjectPaillier_Encrypt(t *testing.T) {
	node := CtNode{
		Co:             &CalculationObjectPaillier{},
		Ids:            nil,
		ReachableNodes: nil,
	}

	e := node.Co.KeyGen()
	if e != nil {
		fmt.Println(e.Error())
		t.Fail()
	}

	c, e := node.Co.Encrypt(24)
	if e != nil {
		fmt.Print(e.Error())
		t.Fail()
	}

	decrypted := node.Co.Decrypt(c)

	if decrypted.Cmp(big.NewInt(24)) != 0 {
		t.Fail()
	}
}
