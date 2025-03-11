export interface AuthModalProps {
    isOpen: boolean;
    onSuccess: () => void;
    onClose: () => void;
}


export interface Wallet {
    address: string;
    private_key: Uint8Array;
    credential: string;
}

export interface CredentialRequest {
    hashed_credential: string
}
export interface CredentialResponse{
    walletAddress: string,
    merkleRoot: string,
    proof:string[]
}
export interface UserInfoState {
    walletAddress: string,
    merkleRoot: string,
    credential: string,
    proof:string[]
}
export interface VerificationResult {
    isValid: boolean;
    txHash?: string;
}