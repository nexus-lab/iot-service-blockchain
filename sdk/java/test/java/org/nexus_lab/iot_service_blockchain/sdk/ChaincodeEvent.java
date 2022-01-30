package org.nexus_lab.iot_service_blockchain.sdk;

import lombok.Data;

@Data
public class ChaincodeEvent implements org.hyperledger.fabric.client.ChaincodeEvent {
  private long blockNumber = 0;
  private String transactionId = "";
  private String chaincodeName = "";
  private final String eventName;
  private final byte[] payload;
}
