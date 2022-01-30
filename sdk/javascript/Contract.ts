import {
  ChaincodeEvent,
  ChaincodeEventsOptions,
  CloseableAsyncIterable,
  Network,
} from '@hyperledger/fabric-gateway';

/**
 * The smart contract interface
 */
export interface ContractInterface {
  /**
   * Submit a transaction to the ledger
   *
   * @param name transaction name
   * @param args transaction arguments
   * @returns the result returned by the transaction function
   */
  submitTransaction(name: string, ...args: Array<string | Uint8Array>): Promise<Uint8Array>;

  /**
   * Register for chaincode events
   *
   * @param options chaincode event options
   * @returns chaincode events
   */
  registerEvent(options?: ChaincodeEventsOptions): Promise<CloseableAsyncIterable<ChaincodeEvent>>;
}

/**
 * Default implementation of the smart contract interface
 */
export default class Contract implements ContractInterface {
  /**
   * @param network the Hyperledger Fabric network/channel
   * @param chaincodeId name of the chaincode
   * @param contractName name of the contract
   */
  constructor(
    private network: Network,
    private chaincodeId: string,
    private contractName: string,
  ) {}

  submitTransaction(name: string, ...args: (string | Uint8Array)[]): Promise<Uint8Array> {
    const contract = this.network.getContract(this.chaincodeId, this.contractName);
    return contract.submitTransaction(name, ...args);
  }

  registerEvent(options?: ChaincodeEventsOptions): Promise<CloseableAsyncIterable<ChaincodeEvent>> {
    return this.network.getChaincodeEvents(this.chaincodeId, options);
  }
}
