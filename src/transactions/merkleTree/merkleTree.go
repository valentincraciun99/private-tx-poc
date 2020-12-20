package merkleTree

import "crypto/sha256"

type MerkleTree struct {
	Root *MerkleNode
	height int
	leaves []*MerkleNode
}

type MerkleNode struct {
	Data  []byte
	Left  *MerkleNode
	Right *MerkleNode
}

func CreateMerkleNode( data []byte) *MerkleNode {
	node := MerkleNode{data,nil,nil}
	return &node
}

func CreateEmptyMerkleTree(height int) *MerkleTree {

	mk:=&MerkleTree{CreateMerkleNode(nil),height,nil}
	var nodes []*MerkleNode
	nodes = append(nodes,mk.Root)


	for level :=1; level <height; level++{
		var nextLevelNodes []*MerkleNode

		for nodeIndex :=0;nodeIndex<len(nodes);nodeIndex++{
			left:=CreateMerkleNode(nil)
			right:=CreateMerkleNode(nil)

			nextLevelNodes = append(nextLevelNodes,left,right)

			nodes[nodeIndex].Left = left
			nodes[nodeIndex].Right = right

		}
		nodes = nextLevelNodes
	}

	mk.leaves = append(mk.leaves,nodes...)

	return mk
}

func HashTree(root *MerkleNode){

	if root==nil {
		return
	}

	HashTree(root.Left)
	HashTree(root.Right)

	if root.Left!=nil && root.Right!=nil{
		childrenHashes := append(root.Left.Data, root.Right.Data...)
		hash := sha256.Sum256(childrenHashes)
		root.Data = hash[:]
	}
}

func (tree *MerkleTree) AddTransaction(data []byte){

	var currentNode *MerkleNode
	var leafNumber int

	for leafIndex :=range tree.leaves{
		if tree.leaves[leafIndex].Data==nil{
			currentNode = tree.leaves[leafIndex]
			leafNumber = leafIndex+1
			break
		}
	}
	hash := sha256.Sum256(data)
	currentNode.Data = hash[:]

	if leafNumber == (1<<(tree.height-1)){
		HashTree(tree.Root)
	}
}
