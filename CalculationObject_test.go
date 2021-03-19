package DCP

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"math/big"
	"testing"
)

func TestCalculationObjectPaillier_KeyGen(t *testing.T) {
	nodes := make([]CtNode, 10)

	for _, node := range nodes {
		node = CtNode{
			Id:           uuid.New(),
			Co:           &CalculationObjectPaillier{},
			Ids:          nil,
			HandledCoIds: nil,
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
		Id:           uuid.New(),
		Co:           &CalculationObjectPaillier{},
		Ids:          nil,
		HandledCoIds: nil,
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

func TestCalculationObjectPaillier_Json(t *testing.T) {
	co := &CalculationObjectPaillier{
		Id:         uuid.UUID{},
		Counter:    10,
		privateKey: nil,
		PublicKey:  nil,
		Cipher:     nil,
	}

	_ = co.KeyGen()

	co.Cipher, _ = co.Encrypt(16)

	serialized, _ := json.Marshal(co)

	var deserialized CalculationObjectPaillier
	e := json.Unmarshal(serialized, &deserialized)
	if e != nil {
		t.Fail()
	}

	if deserialized.Counter != 10 {
		t.Fail()
	}

	decryptedCipher := co.Decrypt(deserialized.Cipher)

	if decryptedCipher.String() != "16" {
		t.Fail()
	}
}
