package DCP

type Prepare func(*CtNode) error

func PrepareIdLenCalculation(node *CtNode) error {
	plain := len(node.Ids)

	c, e := node.Co.Encrypt(plain)
	if e != nil {
		return e
	}

	node.Co.Cipher = c
	node.Co.Counter = node.Co.Counter + 1

	return nil
}