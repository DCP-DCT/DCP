package DCP

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestCalculationObjectPaillier_KeyGen(t *testing.T) {
	nodes := make([]CtNode, 10)

	for _, node := range nodes {
		node = CtNode{
			Id:               uuid.New(),
			Co:               &CalculationObjectPaillier{},
			Ids:              nil,
			HandledBranchIds: nil,
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
		Id:               uuid.New(),
		Co:               &CalculationObjectPaillier{},
		Ids:              nil,
		HandledBranchIds: nil,
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

func TestCalculationObjectGrowth(t *testing.T) {
	var values [][]string

	co := &CalculationObjectPaillier{
			Id:         uuid.UUID{},
			Counter:    0,
			privateKey: nil,
			PublicKey:  nil,
			Cipher:     nil,
		}

		_ = co.KeyGen()

		b, _ := json.Marshal(co)

		values = append(values, []string{strconv.Itoa(co.Counter), strconv.Itoa(len(b))})

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 100000; i++ {
		c, _ := co.Encrypt(rand.Intn(100-1)+1)
		co.Add(c)
		co.Counter = co.Counter + 1

		b, _ := json.Marshal(co)

		values = append(values, []string{strconv.Itoa(co.Counter), strconv.Itoa(len(b))})
	}

	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	headers := []string{"Counter", "Size"}
	if err := w.Write(headers); err != nil {
		panic(err.Error())
	}
	for _, value := range values {
		err := w.Write(value)
		if err != nil {
			panic(err.Error())
		}
	}
}
