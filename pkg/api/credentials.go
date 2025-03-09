package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"keyless-auth/repository"
	"keyless-auth/service"
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
)

type CredentialRequest struct {
	HashedCredential string `json:"hashed_credential"`
}

type CredentialResponse struct {
	WalletAddress string `json:"wallet_address"`
	MerkleRoot    string `json:"merkle_root"`
}

type CredentialsHandler struct {
	credRepo   *repository.CredentialsRepository
	walletRepo *repository.WalletRepository
	merkleTree *service.MerkleTreeService
}

func NewCredentialsHandler(credRepo *repository.CredentialsRepository, walletRepo *repository.WalletRepository) *CredentialsHandler {
	return &CredentialsHandler{
		credRepo:   credRepo,
		walletRepo: walletRepo,
		merkleTree: service.NewMerkleTreeService(credRepo),
	}
}

// generate wallet address
func GenerateWalletAddress() (string, []byte, error) {
	privKey, err := crypto.GenerateKey()
	if err != nil {
		return "", nil, err
	}

	return crypto.PubkeyToAddress(privKey.PublicKey).Hex(), privKey.D.Bytes(), nil
}

func (h *CredentialsHandler) GenerateCredential(w http.ResponseWriter, r *http.Request) {
	var req CredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// check if credential already exists
	if exists, err := h.credRepo.Exists(req.HashedCredential); err != nil || exists {
		// TODO: better return 200 with message to fetch the wallet address or credential on a different endpoint
		http.Error(w, "Credential already exists", http.StatusBadRequest)
		return
	}

	root, node, err := h.merkleTree.AddNodeToTree(req.WalletAddress, req.HashedCredential)
	if err != nil {
		http.Error(w, "Failed to generate merkle tree root", http.StatusInternalServerError)
		return
	}

	if err := h.credRepo.AddNodeToRoot(context.Background(), root.String(), node); err != nil {
		http.Error(w, "Failed to save credential", http.StatusInternalServerError)
		return
	}

	walletAddress := ""

	if len(strings.TrimSpace(req.WalletAddress)) == 0 {
		walletAddress, _, err = GenerateWalletAddress()
		if err != nil {
			http.Error(w, "Failed to generate wallet address", http.StatusInternalServerError)
			return
		}
	// generate wallet address
	walletAddress, privKey, err := GenerateWalletAddress()
	if err != nil {
		http.Error(w, "Failed to generate wallet address", http.StatusInternalServerError)
		return
	}

	// store wallet
	if err := h.walletRepo.Save(walletAddress, privKey, req.HashedCredential); err != nil {
		http.Error(w, "Failed to save wallet", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(CredentialResponse{WalletAddress: walletAddress,
		MerkleRoot: root.String()})
}

func (h *CredentialsHandler) GetWalletByCredential(w http.ResponseWriter, r *http.Request) {
	credential := mux.Vars(r)["credential"]

	wallet, err := h.walletRepo.GetWalletByCredential(credential)
	if err != nil {
		http.Error(w, "Failed to get wallet", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(CredentialResponse{MerkleRoot: wallet.MerkleRoot})
}

func (h *CredentialsHandler) GetWalletIfExists(w http.ResponseWriter, r *http.Request) {
	credential := mux.Vars(r)["credential"]

	wallet, err := h.walletRepo.GetWalletByCredential(credential)
	if err != nil {
		http.Error(w, "Failed to get wallet", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(CredentialResponse{WalletAddress: wallet.Address, MerkleRoot: wallet.MerkleRoot})
}
