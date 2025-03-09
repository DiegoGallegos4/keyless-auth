package service

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"keyless-auth/repository"

	"github.com/wealdtech/go-merkletree"
	"github.com/wealdtech/go-merkletree/keccak256"
)

type MerkleTreeService struct {
	credRepo *repository.CredentialsRepository
}

func NewMerkleTreeService(credRepo *repository.CredentialsRepository) *MerkleTreeService {
	return &MerkleTreeService{
		credRepo: credRepo,
	}
}

func (s *MerkleTreeService) GetRoot(wallet string) (*repository.MerkleNode, error) {
	root, err := s.credRepo.GetLatestRoot(context.Background(), wallet)
	if err != nil {
		return nil, err
	}

	node, err := s.credRepo.GetLatestNode(context.Background(), root)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (s *MerkleTreeService) AddNodeToTree(address, newCredential string) (*merkletree.MerkleTree, *repository.MerkleNode, error) {
	if address == "" || newCredential == "" {
		return nil, nil, errors.New("address and credential must not be empty")
	}

	ctx := context.Background()

	existingCredentials, err := s.credRepo.GetCredentialsByWallet(ctx, address)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get existing credentials: %w", err)
	}

	hashedCred := hashCredential(newCredential)
	data := make([][]byte, 0, len(existingCredentials)+1)

	for _, cred := range existingCredentials {
		data = append(data, []byte(cred))
	}
	data = append(data, hashedCred)

	tree, err := merkletree.NewUsing(data, keccak256.New(), []byte{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build Merkle tree: %w", err)
	}

	proofIndex := uint64(len(data) - 1)
	proof, err := tree.GenerateProof(data[proofIndex])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate proof: %w", err)
	}

	exists, err := s.credRepo.Exists(hex.EncodeToString(hashedCred))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to check for duplicates: %w", err)
	}
	if exists {
		return nil, nil, errors.New("credential already exists")
	}

	newNode := &repository.MerkleNode{
		ID:          uuid.New().String(),
		NodeType:    repository.Credential,
		Hash:        hex.EncodeToString(hashedCred),
		ProofIndex:  proofIndex,
		ProofHashes: proof.Hashes,
		TreeRoot:    tree.Root(),
	}

	err = s.credRepo.AddNodeToRoot(ctx, address, newNode)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update root and node: %w", err)
	}

	return tree, newNode, nil
}

func hashCredential(cred string) []byte {
	salt := []byte{0x1c, 0x9d, 0x3c, 0x4f}
	h := keccak256.New()
	credHash := h.Hash([]byte(cred))
	saltedHash := h.Hash(append(credHash, salt...))
	return saltedHash
}
