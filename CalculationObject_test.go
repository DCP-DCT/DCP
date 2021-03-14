package DCP

import (
	"fmt"
	"testing"
)

func TestCalculationObjectPaillier_KeyGen(t *testing.T) {
	nodes := make([]CtNode, 10)

	for _, node := range nodes {
		e := node.Co.KeyGen()

		if e != nil {
			fmt.Println(e.Error())
			t.Fail()
		}
	}
}

func TestCalculationObjectPaillier_Encrypt(t *testing.T) {
	node := CtNode{
		Co:             CalculationObjectPaillier{},
		Ids:            nil,
		ReachableNodes: nil,
	}

	_ = node.Co.KeyGen()
	_, e := node.Co.Encrypt(24)
	if e != nil {
		fmt.Print(e.Error())
		t.Fail()
	}
}
