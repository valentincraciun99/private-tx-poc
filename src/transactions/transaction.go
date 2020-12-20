package transactions

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Transaction struct {
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxOutput struct {
	Value  int
	PubInfo int
}

type TxInput struct {
	Nullifier []byte
}


func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}
