package DCP

import "testing"

func TestCtPrint(t *testing.T) {
	node := new(CtNode)
	node.Co.KeyGen()

	node.Print()
}
