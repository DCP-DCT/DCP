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
