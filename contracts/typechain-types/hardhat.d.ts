/* Autogenerated file. Do not edit manually. */
/* tslint:disable */
/* eslint-disable */

import { ethers } from "ethers";
import {
  DeployContractOptions,
  FactoryOptions,
  HardhatEthersHelpers as HardhatEthersHelpersBase,
} from "@nomicfoundation/hardhat-ethers/types";

import * as Contracts from ".";

declare module "hardhat/types/runtime" {
  interface HardhatEthersHelpers extends HardhatEthersHelpersBase {
    getContractFactory(
      name: "MerkleRootStore",
      signerOrOptions?: ethers.Signer | FactoryOptions
    ): Promise<Contracts.MerkleRootStore__factory>;
    getContractFactory(
      name: "MerkleRootStore",
      signerOrOptions?: ethers.Signer | FactoryOptions
    ): Promise<Contracts.MerkleRootStore__factory>;
    getContractFactory(
      name: "Verifier",
      signerOrOptions?: ethers.Signer | FactoryOptions
    ): Promise<Contracts.Verifier__factory>;

    getContractAt(
      name: "MerkleRootStore",
      address: string | ethers.Addressable,
      signer?: ethers.Signer
    ): Promise<Contracts.MerkleRootStore>;
    getContractAt(
      name: "MerkleRootStore",
      address: string | ethers.Addressable,
      signer?: ethers.Signer
    ): Promise<Contracts.MerkleRootStore>;
    getContractAt(
      name: "Verifier",
      address: string | ethers.Addressable,
      signer?: ethers.Signer
    ): Promise<Contracts.Verifier>;

    deployContract(
      name: "MerkleRootStore",
      signerOrOptions?: ethers.Signer | DeployContractOptions
    ): Promise<Contracts.MerkleRootStore>;
    deployContract(
      name: "MerkleRootStore",
      signerOrOptions?: ethers.Signer | DeployContractOptions
    ): Promise<Contracts.MerkleRootStore>;
    deployContract(
      name: "Verifier",
      signerOrOptions?: ethers.Signer | DeployContractOptions
    ): Promise<Contracts.Verifier>;

    deployContract(
      name: "MerkleRootStore",
      args: any[],
      signerOrOptions?: ethers.Signer | DeployContractOptions
    ): Promise<Contracts.MerkleRootStore>;
    deployContract(
      name: "MerkleRootStore",
      args: any[],
      signerOrOptions?: ethers.Signer | DeployContractOptions
    ): Promise<Contracts.MerkleRootStore>;
    deployContract(
      name: "Verifier",
      args: any[],
      signerOrOptions?: ethers.Signer | DeployContractOptions
    ): Promise<Contracts.Verifier>;

    // default types
    getContractFactory(
      name: string,
      signerOrOptions?: ethers.Signer | FactoryOptions
    ): Promise<ethers.ContractFactory>;
    getContractFactory(
      abi: any[],
      bytecode: ethers.BytesLike,
      signer?: ethers.Signer
    ): Promise<ethers.ContractFactory>;
    getContractAt(
      nameOrAbi: string | any[],
      address: string | ethers.Addressable,
      signer?: ethers.Signer
    ): Promise<ethers.Contract>;
    deployContract(
      name: string,
      signerOrOptions?: ethers.Signer | DeployContractOptions
    ): Promise<ethers.Contract>;
    deployContract(
      name: string,
      args: any[],
      signerOrOptions?: ethers.Signer | DeployContractOptions
    ): Promise<ethers.Contract>;
  }
}
