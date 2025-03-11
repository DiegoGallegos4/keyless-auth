import React, { useState } from 'react';
import '../css/Header.css';
import {useAppSelector } from '../redux/hooks';
import {UserInfoState} from '../interface';
import { verifyProof } from '../utils/util';
const Header: React.FC = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [txHash, setTxHash] = useState<string | undefined>(undefined);
  const user:UserInfoState = useAppSelector((state) => state.userInfo);
  const contractAddress:string | undefined = process.env.CONTRACT_ADDRESS;
  const handleVerifyProof = async () => {
    // TODO: Implement proof verification logic
    if (!user.proof || !user.merkleRoot) {
      alert('No proof available to verify');
      return;
    }
    try {
      // Add your proof verification logic here
      const result=await verifyProof(contractAddress,user.proof,user.credential);
      if(result.isValid){
        setTxHash(result.txHash);
        alert(`Proof verified successfully! View transaction on Etherscan: https://sepolia.etherscan.io/tx/${result.txHash}`);
      }
    } catch (error) {
      console.error('Error verifying proof:', error);
      setTxHash(undefined);
    }
  };


  return (
    <header className="header">
      <div className="header-container">
        <div className="logo">
          {/* <img src="/logo.svg" alt="Web3 Wallet" /> */}
          <span>ZK Wallet</span>
        </div>

        <nav className={`nav-menu ${isMenuOpen ? 'active' : ''}`}>
          <ul>
            <li><a href="#features">Features</a></li>
            <li><a href="#security">Security</a></li>
            <li><a href="#how-it-works">How it Works</a></li>
            <li><a href="#faq">FAQ</a></li>
          </ul>
        </nav>

        <div className="header-actions">
        <button 
            className="verify-proof-btn" 
            onClick={handleVerifyProof}
            disabled={!user.proof}
          >
            Verify Proof
          </button>
          {txHash && (
            <a 
              href={`https://sepolia.etherscan.io/tx/${txHash}`}
              target="_blank"
              rel="noopener noreferrer"
              className="etherscan-link"
            >
              View on Etherscan
            </a>
          )}
          <button 
            className="mobile-menu-btn"
            onClick={() => setIsMenuOpen(!isMenuOpen)}
          >
            <span></span>
            <span></span>
            <span></span>
          </button>
        </div>
      </div>
    </header>
  );
};

export default Header;