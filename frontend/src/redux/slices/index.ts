import { createSlice,PayloadAction } from "@reduxjs/toolkit";
import { UserInfoState } from "../../interface";

const initialState:UserInfoState = {
    walletAddress:'',
    merkleRoot:'',
    credential:'',
    proof:['']
}

  const userInfoSlice = createSlice({
    name: 'userInfo',
    initialState,
    reducers: {
      setWalletAddress(state, action:PayloadAction<string>) {
        state.walletAddress = action.payload;
      },
      setMerkleRoot(state, action:PayloadAction<string>) {
        state.merkleRoot = action.payload;
      },
      setCredential(state, action:PayloadAction<string>) {
        state.credential = action.payload;
      },
      setProof(state, action:PayloadAction<string[]>) {
        state.proof = action.payload;
      },
      resetAcInfo(state) {
        state.walletAddress = '';
        state.merkleRoot = '';
      },
    },
  });
  
  export const { setWalletAddress,setMerkleRoot,setCredential,setProof ,resetAcInfo } = userInfoSlice.actions;
  export default userInfoSlice.reducer;