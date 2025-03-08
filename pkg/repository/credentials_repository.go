package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"

	"keyless-auth/storage"
	"log"
)

type CredentialsRepository struct {
	db *storage.Redis
}

func NewCredentialsRepository(db *storage.Redis) *CredentialsRepository {
	return &CredentialsRepository{db: db}
}

func (cred *CredentialsRepository) SaveCredential(credential string, address string) error {
	ctx := context.Background()
	// Add leaf to redis set (for fast membership check)
	if err := cred.db.Client.SAdd(ctx, fmt.Sprintf("merkle:%s:credentials:set", address), credential).Err(); err != nil {
		log.Printf("Failed to add credential to redis set: %v", err)
		return err
	}

	// Add leaf to redis list (for ordered retrieval)
	// if err := cred.db.Client.RPush(ctx, fmt.Sprintf("merkle:%s:credentials:set", address), credential).Err(); err != nil {
	// 	log.Printf("Failed to add credential to redis list: %v", err)
	// 	return err
	// }
	return nil
}

func (cred *CredentialsRepository) Exists(credential string) (bool, error) {
	ctx := context.Background()
	return cred.db.Client.SIsMember(ctx, fmt.Sprintf("merkle:credentials:set"), credential).Result()
}

func (cred *CredentialsRepository) GetNodesByRoot(ctx context.Context, root string) ([]MerkleNode, error) {
	rootKey := fmt.Sprintf("root:%s:merkleNode", root)

	nodeIDs, err := cred.db.Client.SMembers(ctx, rootKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch node IDs for root %q: %w", root, err)
	}

	var nodes []MerkleNode
	for _, nodeID := range nodeIDs {
		nodeKey := fmt.Sprintf("merkleNode:%s", nodeID)

		jsonStr, err := cred.db.Client.Get(ctx, nodeKey).Result()
		if errors.Is(err, redis.Nil) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to fetch merkle node %q: %w", nodeID, err)
		}

		var node MerkleNode
		if err := json.Unmarshal([]byte(jsonStr), &node); err != nil {
			return nil, fmt.Errorf("failed to unmarshal merkle node JSON for %q: %w", nodeID, err)
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (cred *CredentialsRepository) GetNodesByWallet(ctx context.Context, address string) ([]MerkleNode, error) {
	walletKey := fmt.Sprintf("wallet:%s:merkleNode", address)

	nodeIDs, err := cred.db.Client.SMembers(ctx, walletKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch node IDs for wallet %q: %w", address, err)
	}

	var nodes []MerkleNode
	for i, nodeID := range nodeIDs {
		nodeKey := fmt.Sprintf("merkleNode:%s", nodeID)

		jsonStr, err := cred.db.Client.LRange(ctx, nodeKey, 0, -1).Result()
		if errors.Is(err, redis.Nil) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to fetch merkle node %q: %w", nodeID, err)
		}

		var node MerkleNode
		if err := json.Unmarshal([]byte(jsonStr[i]), &node); err != nil {
			return nil, fmt.Errorf("failed to unmarshal merkle node JSON for %q: %w", nodeID, err)
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (cred *CredentialsRepository) GetRootsByWallet(ctx context.Context, address string) ([]string, error) {
	key := fmt.Sprintf("wallet:%s:root", address)
	roots, err := cred.db.Client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	return roots, nil
}

func (cred *CredentialsRepository) SetWalletForNode(ctx context.Context, wallet string, root string) error {
	k := fmt.Sprintf("wallet:%s:root:", wallet)

	// Add wallet→root
	if err := cred.db.Client.LPush(ctx, k, root).Err(); err != nil {
		return err
	}

	return nil
}

func (cred *CredentialsRepository) AddNodeToRoot(ctx context.Context, root string, node *MerkleNode) error {
	rootKey := fmt.Sprintf("root:%s:merkleNode", root)

	jsonStr, err := json.Marshal(node)
	if err != nil {
		return err
	}
	if err := cred.db.Client.LPush(ctx, rootKey, jsonStr).Err(); err != nil {
		return err
	}

	return nil
}

func (cred *CredentialsRepository) GetLatestNode(ctx context.Context, root string) (*MerkleNode, error) {
	rootKey := fmt.Sprintf("root:%s:merkleNode", root)

	node, err := cred.db.Client.LIndex(ctx, rootKey, 0).Result()
	if err != nil {
		return nil, err
	}

	if len(node) == 0 {
		return nil, errors.New("no nodes found")
	}

	var _node MerkleNode
	if err := json.Unmarshal([]byte(node), &_node); err != nil {
		return nil, fmt.Errorf("failed to unmarshal merkle node JSON for %q: %w", node, err)
	}

	return &_node, nil
}

func (cred *CredentialsRepository) GetLatestRoot(ctx context.Context, wallet string) (string, error) {
	key := fmt.Sprintf("wallet:%s:root", wallet)

	root, err := cred.db.Client.LIndex(ctx, key, 0).Result()
	if err != nil {
		return "", err
	}

	return root, nil
}

func (cred *CredentialsRepository) GetCredentialsByWallet(ctx context.Context, address string) ([]string, error) {
	key := fmt.Sprintf("wallet:%s:credentials", address)
	creds, err := cred.db.Client.SMembers(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	return creds, nil
}
