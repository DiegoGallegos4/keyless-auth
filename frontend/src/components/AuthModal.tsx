import React, { useState} from 'react';
import '../css/AuthModal.css';
import { AxiosError } from 'axios';
import { useAppDispatch} from '../redux/hooks';
import {setWalletAddress,setMerkleRoot,setCredential,setProof} from '../redux/slices/index'
import { generateCredentials,getWalletByCredential } from '../api';
import { calculateHash } from '../utils/util'
import { AuthModalProps,Wallet,CredentialResponse } from '../interface';


const AuthModal: React.FC<AuthModalProps> = ({ isOpen, onSuccess, onClose }) => {

  const [error, setError] = useState<string>('');
  const dispatch = useAppDispatch();

  const [email, setemail] = useState<string>('');

  if (!isOpen) return null;
  
  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      const hashedEmail = await calculateHash(email);
      dispatch(setCredential(hashedEmail));
      try {
        // First try to get existing wallet
        const existingWallet:Wallet = await getWalletByCredential(hashedEmail);
        dispatch(setWalletAddress(existingWallet.address));
        onSuccess();
        return;
      } catch (error) {
        // If wallet doesn't exist, continue to create new one
        if (error instanceof AxiosError && error.response?.status === 409) {

          const newWallet:CredentialResponse=await generateCredentials(hashedEmail);
          dispatch(setWalletAddress(newWallet.walletAddress));
          dispatch(setMerkleRoot(newWallet.merkleRoot));
          dispatch(setCredential(hashedEmail));
          dispatch(setProof(newWallet.proof))
          onSuccess();
          return;
        }
        throw error;
      }
    } catch (error) {
      console.error('Authentication failed:', error);
      setError('Authentication failed. Please try again.');
    }
  };

  return (
    <div className="modal-overlay">
      <div className="modal-content">
        <button className="close-button" onClick={onClose}>&times;</button>
        
        <h2>Welcome</h2>
        <p>Enter your email to continue</p>
        
        {error && <div className="error-message">{error}</div>}
        
        <form className="auth-form" onSubmit={handleSubmit}>
          <input 
            type="email" 
            value={email}
            onChange={(e) => setemail(e.target.value)}
            placeholder="Email address" 
            className="auth-input"
            required
          />
          <button type="submit" className="submit-button">
            Continue
          </button>
        </form>
      </div>
    </div>
  );
};

export default AuthModal;