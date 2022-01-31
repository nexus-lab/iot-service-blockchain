package org.nexus_lab.iot_service_blockchain.sdk;

import lombok.RequiredArgsConstructor;
import org.hyperledger.fabric.client.CallOption;
import org.hyperledger.fabric.client.ChaincodeEvent;
import org.hyperledger.fabric.client.CloseableIterator;
import org.hyperledger.fabric.client.CommitException;
import org.hyperledger.fabric.client.CommitStatusException;
import org.hyperledger.fabric.client.EndorseException;
import org.hyperledger.fabric.client.Network;
import org.hyperledger.fabric.client.SubmitException;

/** Default implementation of the smart contract interface. */
@RequiredArgsConstructor
public class Contract implements ContractInterface {
  /** The Hyperledger Fabric network/channel. */
  private final Network network;

  /** Name of the chaincode. */
  private final String chaincodeId;

  /** Name of the contract. */
  private final String contractName;

  @Override
  public byte[] submitTransaction(String name, String... args)
      throws EndorseException, SubmitException, CommitStatusException, CommitException {
    org.hyperledger.fabric.client.Contract contract =
        this.network.getContract(this.chaincodeId, this.contractName);
    return contract.submitTransaction(name, args);
  }

  @Override
  public CloseableIterator<ChaincodeEvent> registerEvent(CallOption... options) {
    return this.network.getChaincodeEvents(this.chaincodeId, options);
  }
}
