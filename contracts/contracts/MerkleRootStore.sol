// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract MerkleRootStore {
    bytes32 public merkleRoot;

    event MerkleRootUpdated(bytes32 newMerkleRoot);

    function updateMerkleRoot(bytes32 newRoot) public {
        merkleRoot = newRoot;
        emit MerkleRootUpdated(newRoot);
    }

    function getMerkleRoot() public view returns (bytes32) {
        return merkleRoot;
    }

    function verifyMerkleProof(bytes32[] calldata proof, bytes32 leaf) external view returns (bool) {
        bytes32 computedHash = leaf;
        for (uint256 i = 0; i < proof.length; i++) {
            if (computedHash < proof[i]) {
                computedHash = keccak256(abi.encodePacked(computedHash, proof[i]));
            } else {
                computedHash = keccak256(abi.encodePacked(proof[i], computedHash));
            }
        }
        return computedHash == merkleRoot;
    }
}