import axios from 'axios';
import { CredentialRequest, UserInfoState, Wallet } from './interface';



const API_BASE_URL = process.env.API_BASE_URL;

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const generateCredentials = async (hashed_cred:string ) => {
    const requestData: CredentialRequest = {
        hashed_credential: hashed_cred,
      };

      const { data } = await axios.post< UserInfoState>(
        `${API_BASE_URL}/credentials`,
        requestData
      );
  return data;
};

export const getWalletByCredential = async (credential: string) => {
  const response = await api.get<Wallet>(`/credentials/${credential}`);
  return response.data;
};

// export const generateProof = async (data: any) => {
//   const response = await api.post('/proof', data);
//   return response.data;
// };