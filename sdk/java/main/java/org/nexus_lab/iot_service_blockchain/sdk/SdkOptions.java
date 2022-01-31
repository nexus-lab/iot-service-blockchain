package org.nexus_lab.iot_service_blockchain.sdk;

import lombok.Data;
import lombok.NoArgsConstructor;

/** SDK initialization options. */
@Data
@NoArgsConstructor
public class SdkOptions {
  /** Organization/MSP ID. */
  private String organizationId;

  /** PEM-formated X509 client certificate. */
  private String certificate;

  /** PEM-formated client private key. */
  private String privateKey;

  /** Network address of the gateway peer. */
  private String gatewayPeerEndpoint;

  /** Server name of the gateway peer. */
  private String gatewayPeerServerName;

  /** PEM-formated X509 TLS certificate of the gateway peer. */
  private String gatewayPeerTlsCertificate;

  /** Blockchain network channel name. */
  private String networkName;

  /** Name of the chaincode. */
  private String chaincodeId;
}
