package repository

import (
	"fmt"
	"time"
)

// wallet address is mapped to a root
// each root has l_sibling and r_sibling
// r_sibling is always a credential or sub_root
// l_sibling is always a proof
// each root is always public

type NodeType int

const (
	Root NodeType = iota
	SubRoot
	Credential
	Proof
)

// MerkleNode describes a single node in the Merkle tree.
type MerkleNode struct {
	ID          string
	NodeType    NodeType
	Hash        string   // Stores leaf hash
	ProofIndex  uint64   // Original index in tree
	ProofHashes [][]byte // Maintain binary format
	TreeRoot    []byte   // Root reference
	CreatedAt   time.Time
}

var NodeTypeNames = map[NodeType]string{
	Root:       "root",
	SubRoot:    "sroot",
	Credential: "credential",
	Proof:      "proof",
}

func (t NodeType) String() string {
	if name, ok := NodeTypeNames[t]; ok {
		return name
	}
	return fmt.Sprintf("Unknown NodeType (%d)", t)
}
