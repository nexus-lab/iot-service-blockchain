import {
  ChaincodeEventsOptions,
  CloseableAsyncIterable,
  Network,
} from '@hyperledger/fabric-gateway';
import { TextDecoder } from 'util';

import Contract, { ContractInterface } from './Contract';
import Service from './Service';

/**
 * Interface of the event emitted by the service registry contract notifying a service update
 */
export interface ServiceEvent {
  /**
   * Name of the action performed on the service
   */
  action: string;

  /**
   * Organization ID of the service
   */
  organizationId: string;

  /**
   * ID of the device to which the service belongs
   */
  deviceId: string;

  /**
   * Name of the service
   */
  serviceName: string;

  /**
   * Custom event payload
   */
  payload: any;
}

/**
 * Interface of core utilities for managing services on the ledger
 */
export interface ServiceRegistryInterface {
  /**
   * Create or update a service in the ledger
   *
   * @param service service to be created or updated
   */
  register(service: Service): Promise<void>;

  /**
   * Get a service by its organization ID, device ID, and name
   *
   * @param organizationId service's organization ID
   * @param deviceId device's ID
   * @param serviceName name of the service
   * @returns the device
   */
  get(organizationId: string, deviceId: string, serviceName: string): Promise<Service>;

  /**
   * Get a list of services by their organization ID and device ID
   *
   * @param organizationId device's organization ID
   * @param deviceId device's ID
   * @returns all services of the device in the organization
   */
  getAll(organizationId: string, deviceId: string): Promise<Service[]>;

  /**
   * Remove a service from the ledger
   *
   * @param service the service to be removed
   */
  deregister(device: Service): Promise<void>;

  /**
   * Registers for service registry events
   *
   * @param options chaincode event options
   * @returns an async iterable of the events
   */
  registerEvent(options?: ChaincodeEventsOptions): Promise<CloseableAsyncIterable<ServiceEvent>>;
}

/**
 * Core utilities for managing devices on the ledger
 */
export default class ServiceRegistry implements ServiceRegistryInterface {
  private utf8Decoder: TextDecoder;
  /**
   * @param contract smart contract
   */
  constructor(private contract: ContractInterface) {
    this.utf8Decoder = new TextDecoder();
  }

  async register(service: Service): Promise<void> {
    const serialized = service.serialize();
    await this.contract.submitTransaction('Register', serialized);
  }

  async get(organizationId: string, deviceId: string, serviceName: string): Promise<Service> {
    const data = await this.contract.submitTransaction(
      'Get',
      organizationId,
      deviceId,
      serviceName,
    );
    const serialized = this.utf8Decoder.decode(data);
    return Service.deserialize(serialized);
  }

  async getAll(organizationId: string, deviceId: string): Promise<Service[]> {
    const data = await this.contract.submitTransaction('GetAll', organizationId, deviceId);
    const serialized = this.utf8Decoder.decode(data);
    return (JSON.parse(serialized) as any[]).map(Service.fromObject);
  }

  async deregister(service: Service): Promise<void> {
    const serialized = service.serialize();
    await this.contract.submitTransaction('Deregister', serialized);
  }

  async registerEvent(
    options?: ChaincodeEventsOptions,
  ): Promise<CloseableAsyncIterable<ServiceEvent>> {
    const events = await this.contract.registerEvent(options);
    const pattern = /^service:\/\/(.+?)\/(.+?)\/(.+?)\/(.+?)$/;
    const decoder = this.utf8Decoder;

    return {
      async *[Symbol.asyncIterator]() {
        for await (const event of events) {
          // can reuse pattern here since it has no global('g') flag
          const matches = pattern.exec(event.eventName);
          if (matches === null || matches.length != 5) {
            continue;
          }

          const serviceEvent: ServiceEvent = {
            organizationId: matches[1],
            deviceId: matches[2],
            serviceName: matches[3],
            action: matches[4],
            payload: null,
          };

          if (serviceEvent.action === 'register' || serviceEvent.action === 'deregister') {
            try {
              const service = Service.deserialize(decoder.decode(event.payload));
              serviceEvent.payload = service;
            } catch {
              console.error(
                `bad service event payload ${event.payload}, action is ${serviceEvent.action}`,
              );
              continue;
            }
          } else {
            serviceEvent.payload = event.payload;
          }

          yield serviceEvent;
        }
      },
      close() {
        events.close();
      },
    };
  }
}

/**
 * The default factory for creating service registries
 * 
 * @param network Hyperledger Fabric network
 * @param chaincodeId ID/name of the chaincode
 * @returns The service registry
 */
export function createServiceRegistry(network: Network, chaincodeId: string) {
  return new ServiceRegistry(new Contract(network, chaincodeId, 'service_registry'));
}
