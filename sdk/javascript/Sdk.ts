import * as grpc from '@grpc/grpc-js';
import { Gateway, Identity, Network, connect, signers } from '@hyperledger/fabric-gateway';
import * as crypto from 'crypto';
import { TextEncoder } from 'util';

import { DeviceRegistryInterface, createDeviceRegistry } from './DeviceRegistry';
import { ServiceBrokerInterface, createServiceBroker } from './ServiceBroker';
import { ServiceRegistryInterface, createServiceRegistry } from './ServiceRegistry';
import { getClientId, parseCertificate } from './identity';

/**
 * SDK initialization options
 */
export interface SdkOptions {
  /**
   * Organization/MSP ID
   */
  organizationId: string;

  /**
   * PEM-formated X509 client certificate
   */
  certificate: string;

  /**
   * PEM-formated client private key
   */
  privateKey: string;

  /**
   * Network address of the gateway peer
   */
  gatewayPeerEndpoint: string;

  /**
   * Server name of the gateway peer
   */
  gatewayPeerServerName: string;

  /**
   * PEM-formated X509 TLS certificate of the gateway peer
   */
  gatewayPeerTLSCertificate: string;

  /**
   * Blockchain network channel name
   */
  networkName: string;

  /**
   * Name of the chaincode
   */
  chaincodeId: string;
}

function newIdentity(organizationId: string, certificate: string): Identity {
  const encoder = new TextEncoder();
  return {
    mspId: organizationId,
    credentials: encoder.encode(certificate),
  };
}

function newSigner(privateKey: string) {
  return signers.newPrivateKeySigner(crypto.createPrivateKey(Buffer.from(privateKey)));
}

function newGrpcConnection(endpoint: string, serverName: string, certificate: string) {
  const credentials = grpc.credentials.createSsl(Buffer.from(certificate));
  return new grpc.Client(endpoint, credentials, { 'grpc.ssl_target_name_override': serverName });
}

/**
 * The IoT service blockchain sdk
 */
export default class Sdk {
  private grpcConnection: grpc.Client;
  private gateway: Gateway;
  private organizationId!: string;
  private deviceId!: string;
  private deviceRegistry!: DeviceRegistryInterface;
  private serviceRegistry!: ServiceRegistryInterface;
  private serviceBroker!: ServiceBrokerInterface;

  /**
   * @param options SDK initialization options
   */
  constructor(options: SdkOptions) {
    this.grpcConnection = newGrpcConnection(
      options.gatewayPeerEndpoint,
      options.gatewayPeerServerName,
      options.gatewayPeerTLSCertificate,
    );
    const identity = newIdentity(options.organizationId, options.certificate);
    const signer = newSigner(options.privateKey);

    this.gateway = connect({
      client: this.grpcConnection,
      identity,
      signer,
      evaluateOptions: () => ({ deadline: Date.now() + 5000 }), // 5 seconds
      endorseOptions: () => ({ deadline: Date.now() + 15000 }), // 15 seconds
      submitOptions: () => ({ deadline: Date.now() + 5000 }), // 5 seconds
      commitStatusOptions: () => ({ deadline: Date.now() + 60000 }), // 1 minute
    });

    const network = this.gateway.getNetwork(options.networkName);
    this.connectSmartContracts(network, options.chaincodeId);
    this.setIdentity(options.organizationId, options.certificate);
  }

  private setIdentity(organizationId: string, certificate: string) {
    this.deviceId = getClientId(parseCertificate(certificate));
    this.organizationId = organizationId;
  }

  private connectSmartContracts(network: Network, chaincodeId: string) {
    this.deviceRegistry = createDeviceRegistry(network, chaincodeId);
    this.serviceRegistry = createServiceRegistry(network, chaincodeId);
    this.serviceBroker = createServiceBroker(network, chaincodeId);
  }

  /**
   * Get the device/client ID of the current calling application
   *
   * @returns the device/client ID of the current calling application
   */
  getDeviceId() {
    return this.deviceId;
  }

  /**
   * Get the organization ID of the current calling application
   *
   * @returns the organization ID of the current calling application
   */
  getOrganizationId() {
    return this.organizationId;
  }

  /**
   * Get the device registry
   *
   * @returns the device registry
   */
  getDeviceRegistry() {
    return this.deviceRegistry;
  }

  /**
   * Get the server registry
   *
   * @returns the server registry
   */
  getServiceRegistry() {
    return this.serviceRegistry;
  }

  /**
   * Get the service broker
   *
   * @returns the service broker
   */
  getServiceBroker() {
    return this.serviceBroker;
  }

  /**
   * Close connection to the Hyperledger Fabric gateway
   */
  close() {
    this.gateway.close();
    this.grpcConnection.close();
  }
}
