import {ethers} from 'ethers';
import { getWalletByCredential } from '../api';
import { Wallet,VerificationResult } from '../interface';
import MerkleRootStore from './MerkleRootStore.json';
export const calculateHash = async (email: string): Promise<string> => {
    const encoder = new TextEncoder();
    const data = encoder.encode(email);
    const hashBuffer = await crypto.subtle.digest('SHA-256', data);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
    return `${hashHex}`;
  };
//   export const setNewRoot = async (
//     contractAddress: string, 
//     merkleRoot: string,
//   ): Promise<boolean> => {
//     try {
        
        
//         const data=await getWalletByCredential();
        
  
//       // MerkleRootStore ABI - only the function we need
//       const abi = [
//         "function setNewRoot(bytes32 newRoot, address walletAddress) external"
//       ];
  
//       // Create contract instance
//       const contract = new ethers.Contract(contractAddress, abi, signer);
  
//       // Convert merkle root to bytes32
//       const rootBytes32 = ethers.hexlify(ethers.toUtf8Bytes(merkleRoot));
  
//       // Send transaction
//       const tx = await contract.setNewRoot(rootBytes32, walletAddress);
//       await tx.wait(); // Wait for transaction to be mined
  
//       return true;
//     } catch (error) {
//       console.error('Error setting new root:', error);
//       return false;
//     }
//   };

export const verifyProof = async (
    contractAddress: string | undefined,
    proof: string[],  // Changed to array of proofs
    leaf: string,     // Changed to leaf instead of direct proof
): Promise<VerificationResult> => {
    try {
        const wallet:Wallet = await getWalletByCredential(leaf);
        if(contractAddress===undefined){
            throw new Error('Contract address not found');
        }
        const privateKeyHex = Array.from(wallet.private_key)
            .map(b => b.toString(16).padStart(2, '0'))
            .join('');
        

        const formattedPrivateKey = privateKeyHex.startsWith('0x') 
            ? privateKeyHex 
            : `0x${privateKeyHex}`;
            const provider = new ethers.JsonRpcProvider('https://eth-sepolia.g.alchemy.com/v2/8PVgppTYpuzgtUJoV-2XrKAw_fJ1hx8V');
            const signerWallet = new ethers.Wallet(formattedPrivateKey, provider);

       
        const abi =MerkleRootStore.abi;

        const contract = new ethers.Contract(contractAddress, abi, signerWallet);

        const proofBytes32 = proof.map((p: string) => ethers.hexlify(ethers.toUtf8Bytes(p)));
        
        const leafBytes32 = ethers.hexlify(ethers.toUtf8Bytes(leaf));

        // Call verification function with correct parameters
        const tx = await contract.verifyMerkleProof(proofBytes32, leafBytes32);
        const receipt = await tx.wait();

        return {
            isValid: true,
            txHash: receipt.hash
        };
    } catch (error) {
        console.error('Error verifying merkle proof:', error);
        return {
            isValid: false,
            txHash: undefined
        };
    }
};



