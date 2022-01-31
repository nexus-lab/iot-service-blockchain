package org.nexus_lab.iot_service_blockchain.sdk;

import org.hyperledger.fabric.client.CallOption;
import org.hyperledger.fabric.client.ChaincodeEvent;
import org.hyperledger.fabric.client.CloseableIterator;
import org.hyperledger.fabric.client.CommitException;
import org.hyperledger.fabric.client.CommitStatusException;
import org.hyperledger.fabric.client.EndorseException;
import org.hyperledger.fabric.client.SubmitException;

/** The smart contract interface. */
public interface ContractInterface {
  /**
   * Submit a transaction to the ledger.
   *
   * @param name transaction name
   * @param args transaction arguments
   * @return the result returned by the transaction function
   * @throws EndorseException if the endorse invocation fails
   * @throws SubmitException if the submit invocation fails
   * @throws CommitStatusException if the commit status invocation fails
   * @throws CommitException if the transaction commits unsuccessfully
   * @throws NullPointerException if the transaction name is null
   */
  public byte[] submitTransaction(String name, String... args)
      throws EndorseException, SubmitException, CommitStatusException, CommitException;

  /**
   * Register for chaincode events.
   *
   * @param options chaincode event options
   * @return chaincode events
   * @throws NullPointerException if the chaincode name is null
   */
  public CloseableIterator<ChaincodeEvent> registerEvent(CallOption... options);
}
