package repository

import (
	"context"
<<<<<<< HEAD
	"log"
=======
	"encoding/json"
	"fmt"
>>>>>>> 483d9215152da2ad6883daaa0789698081fed34d
	"time"

	"keyless-auth/domain"
	"keyless-auth/storage"
)

type WalletRepository struct {
	db *storage.Redis
}

func NewWalletRepository(db *storage.Redis) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Save(address string, privKey []byte, credential string, merkleRoot string) error {
	wallet := &domain.Wallet{
		Address:    address,
		PrivateKey: privKey,
		Credential: credential,
		MerkleRoot: merkleRoot,
	}
<<<<<<< HEAD
	log.Printf("Saving wallet: %v", wallet)

	serializedWallet, err := storage.Serialize(wallet)
	if err != nil {
		log.Printf("Failed to serialize wallet: %v", err)
		return err
	}

	err = r.db.Save(context.Background(), storage.GenerateCacheKey("wallet", credential), serializedWallet, time.Hour*24)
	if err != nil {
		log.Printf("Failed to save wallet: %v", err)
		return err
	}
	return nil
=======

	_wallet, err := json.Marshal(wallet)
	if err != nil {
		return fmt.Errorf("unable to marshal data for redis: %w", err)
	}

	return r.db.Save(context.Background(), storage.GenerateCacheKey("wallet", credential), _wallet, time.Hour*24)
>>>>>>> 483d9215152da2ad6883daaa0789698081fed34d
}

func (r *WalletRepository) GetWalletByCredential(hashedCredential string) (*domain.Wallet, error) {
	value, err := r.db.Get(context.Background(), storage.GenerateCacheKey("wallet", hashedCredential))
	if err != nil {
<<<<<<< HEAD
		log.Printf("Failed to get wallet by credential: %v", err)
=======
>>>>>>> 483d9215152da2ad6883daaa0789698081fed34d
		return nil, err
	}
	var wallet domain.Wallet
	err = storage.Deserialize(string(value), &wallet)
	return &wallet, err
}
