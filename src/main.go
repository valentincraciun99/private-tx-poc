package main

import (
	"encoding/hex"
	"fmt"
	"src/account"
	"src/transactions"
	"src/transactions/merkleTree"
	"src/zk"
)

const(
	TreeHeight    = 2
	AlicePk       = 35
	AlicePv       = 3
	CoinbaseValue = 1000
	TxAmount      = 20
	Address       = "Alice"

)

func main(){
	var txs []transactions.Transaction
	merkleTree:= merkleTree.CreateEmptyMerkleTree(TreeHeight)
	circuit:=zk.GenerateCircuit()

	acc:= account.CreateAccount(Address,AlicePk,AlicePv)

	coinbaseTx:= transactions.Transaction{Outputs: []transactions.TxOutput{{CoinbaseValue, acc.Pk}}}
	txs = append(txs,coinbaseTx)
	merkleTree.AddTransaction(coinbaseTx.Serialize())

	acc.SendTx(TxAmount,acc,acc,&txs,circuit)
	merkleTree.AddTransaction(txs[len(txs)-1].Serialize())


	fmt.Println("MerkleTree Root:" + hex.EncodeToString(merkleTree.Root.Data))
}
