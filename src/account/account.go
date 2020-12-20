package account

import (
	"bytes"
	"crypto/sha256"
	"github.com/arnaucube/go-snark"
	"github.com/arnaucube/go-snark/circuitcompiler"
	"math/big"
	"src/transactions"
	"src/zk"
)

var privateInputs = make(map[*Account]int)

type Account struct{
	Address string
	Pk int
}

func CreateAccount(address string, pk int,privateInput int) *Account {
	acc:=&Account{Address: address, Pk: pk}
	privateInputs[acc] = privateInput
	return acc
}

func (account *Account) SendTx(amount int,sender *Account, receiver *Account,txs *[]transactions.Transaction,circuit *circuitcompiler.Circuit){

	var inputs []transactions.TxInput
	var outputs []transactions.TxOutput

	amountCollected:=0

	witness := zk.GenerateWitness(privateInputs[sender],sender.Pk,circuit)

	setup,px:= zk.GenerateSetup(witness,circuit)

	proof, _ := snark.GenerateProofs(*circuit, setup.Pk, witness, px)

	for _, actualTx :=range *txs{
		isUnspent:=true
		receiverAsBytes := []byte(receiver.Address)
		txAsBytes := actualTx.Serialize()
		hashTx:= sha256.Sum256(append(receiverAsBytes,txAsBytes...))

		for _,tx:=range *txs{
			for _,input:=range tx.Inputs{
				if bytes.Equal(input.Nullifier, hashTx[:]){
					isUnspent = false
				}
			}
		}

		if isUnspent==true{
			for _,out:=range actualTx.Outputs{
				b35Verify := big.NewInt(int64(out.PubInfo))
				publicSignalsVerify := []*big.Int{b35Verify}
				if snark.VerifyProof(setup.Vk, proof, publicSignalsVerify, true){
					in:= transactions.TxInput{Nullifier: hashTx[:]}
					amountCollected+=out.Value
					inputs = append(inputs,in)
				}
			}
		}
		if amountCollected >= amount{
			out1:= transactions.TxOutput{Value: amount, PubInfo: receiver.Pk}
			out2:= transactions.TxOutput{Value: amountCollected-amount, PubInfo: sender.Pk}

			outputs = append(outputs,out1,out2)
		}
	}
	newTx:=&transactions.Transaction{Inputs:inputs,Outputs: outputs}

	*txs = append(*txs, *newTx)
}

