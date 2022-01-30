package org.nexus_lab.iot_service_blockchain.sdk;

import static org.nexus_lab.iot_service_blockchain.sdk.Identity.getClientId;
import static org.nexus_lab.iot_service_blockchain.sdk.Identity.parseCertificate;

import io.grpc.ManagedChannel;
import io.grpc.netty.shaded.io.grpc.netty.GrpcSslContexts;
import io.grpc.netty.shaded.io.grpc.netty.NettyChannelBuilder;
import java.security.InvalidKeyException;
import java.security.PrivateKey;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;
import java.util.concurrent.TimeUnit;
import javax.naming.InvalidNameException;
import javax.net.ssl.SSLException;
import lombok.Getter;
import org.hyperledger.fabric.client.CallOption;
import org.hyperledger.fabric.client.Gateway;
import org.hyperledger.fabric.client.Network;
import org.hyperledger.fabric.client.identity.Identities;
import org.hyperledger.fabric.client.identity.Identity;
import org.hyperledger.fabric.client.identity.Signer;
import org.hyperledger.fabric.client.identity.Signers;
import org.hyperledger.fabric.client.identity.X509Identity;

/** The IoT service blockchain sdk */
public class Sdk {
  private ManagedChannel grpcConnection;
  private Gateway gateway;

  /** The organization ID of the current calling application */
  @Getter private String organizationId;

  /** The device/client ID of the current calling application */
  @Getter private String deviceId;

  /** The device registry */
  @Getter private DeviceRegistryInterface deviceRegistry;

  /** The service registry */
  @Getter private ServiceRegistryInterface serviceRegistry;

  /** The service broker */
  @Getter private ServiceBrokerInterface serviceBroker;

  private static Identity newIdentity(String organizationId, String certificate)
      throws CertificateException {
    X509Certificate cert = Identities.readX509Certificate(certificate);
    return new X509Identity(organizationId, cert);
  }

  private static Signer newSigner(String privateKey) throws InvalidKeyException {
    PrivateKey key = Identities.readPrivateKey(privateKey);
    return Signers.newPrivateKeySigner(key);
  }

  private static ManagedChannel newGrpcConnection(
      String endpoint, String serverName, String certificate)
      throws CertificateException, SSLException {
    X509Certificate cert = Identities.readX509Certificate(certificate);
    return NettyChannelBuilder.forTarget(endpoint)
        .sslContext(GrpcSslContexts.forClient().trustManager(cert).build())
        .overrideAuthority(serverName)
        .build();
  }

  /**
   * @param options SDK initialization options
   * @throws CertificateException
   * @throws SSLException
   * @throws InvalidKeyException
   * @throws InvalidNameException
   */
  public Sdk(SdkOptions options)
      throws CertificateException, SSLException, InvalidKeyException, InvalidNameException {
    this.grpcConnection =
        newGrpcConnection(
            options.getGatewayPeerEndpoint(),
            options.getGatewayPeerServerName(),
            options.getGatewayPeerTLSCertificate());

    Gateway.Builder builder =
        Gateway.newInstance()
            .identity(newIdentity(options.getOrganizationId(), options.getCertificate()))
            .signer(newSigner(options.getPrivateKey()))
            .connection(this.grpcConnection)
            .evaluateOptions(CallOption.deadlineAfter(5, TimeUnit.SECONDS))
            .endorseOptions(CallOption.deadlineAfter(15, TimeUnit.SECONDS))
            .submitOptions(CallOption.deadlineAfter(5, TimeUnit.SECONDS))
            .commitStatusOptions(CallOption.deadlineAfter(1, TimeUnit.MINUTES));

    this.gateway = builder.connect();

    Network network = this.gateway.getNetwork(options.getNetworkName());
    this.connectSmartContracts(network, options.getChaincodeId());
    this.setIdentity(options.getOrganizationId(), options.getCertificate());
  }

  private void setIdentity(String organizationId, String certificate)
      throws InvalidNameException, CertificateException {
    this.deviceId = getClientId(parseCertificate(certificate.getBytes()));
    this.organizationId = organizationId;
  }

  private void connectSmartContracts(Network network, String chaincodeId) {
    this.deviceRegistry = DeviceRegistry.create(network, chaincodeId);
    this.serviceRegistry = ServiceRegistry.create(network, chaincodeId);
    this.serviceBroker = ServiceBroker.create(network, chaincodeId);
  }

  /** Close connection to the Hyperledger Fabric gateway */
  public void close() {
    this.gateway.close();
    this.grpcConnection.shutdown();
  }
}
