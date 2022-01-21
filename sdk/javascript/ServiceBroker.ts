import {
  ChaincodeEventsOptions,
  CloseableAsyncIterable,
  Network,
} from '@hyperledger/fabric-gateway';
import { TextDecoder } from 'util';

import Contract, { ContractInterface } from './Contract';
import ServiceRequest from './ServiceRequest';
import ServiceRequestResponse from './ServiceRequestResponse';
import ServiceResponse from './ServiceResponse';

/**
 * Interface of the event emitted by the service broker contract notifying a service request/response update
 */
export interface ServiceRequestEvent {
  /**
   * Name of the action performed on the service request
   */
  action: string;

  /**
   * Organization ID of the requested service
   */
  organizationId: string;

  /**
   * Device ID of the requested service
   */
  deviceId: string;

  /**
   * Name of the requested service
   */
  serviceName: string;

  /**
   * ID of the request
   */
  requestId: string;

  /**
   * Custom event payload
   */
  payload: any;
}

/**
 * Interface of core utilities for managing service requests on ledger
 */
export interface ServiceBrokerInterface {
  /**
   * Make a request to an IoT service
   *
   * @param request IoT service request to be sent
   */
  request(request: ServiceRequest): Promise<void>;

  /**
   * Respond to an IoT service request
   *
   * @param response IoT service response to be sent
   */
  respond(response: ServiceResponse): Promise<void>;

  /**
   * Get an IoT service request and its response (if any) by the request ID
   *
   * @param requestId service request ID
   * @returns the service request and its response (if any)
   */
  get(requestId: string): Promise<ServiceRequestResponse>;

  /**
   * Get a list of IoT service requests and their responses (if any) by their service organization ID, service device ID, and service name
   *
   * @param organizationId organization ID of the requested service
   * @param deviceId device's ID of the requested service
   * @param serviceName name of the requested service
   * @returns all services of the device in the organization
   */
  getAll(
    organizationId: string,
    deviceId: string,
    serviceName: string,
  ): Promise<ServiceRequestResponse[]>;

  /**
   * Remove a service request and its response (if any) from the ledger
   *
   * @param requestId ID of the service request and response to be removed
   */
  remove(requestId: string): Promise<void>;

  /**
   * Registers for service request events
   *
   * @param options chaincode event options
   * @returns an async iterable of the events
   */
  registerEvent(
    options: ChaincodeEventsOptions,
  ): Promise<CloseableAsyncIterable<ServiceRequestEvent>>;
}

/**
 * Core utilities for managing IoT service requests and responses on the ledger
 */
export default class ServiceBroker implements ServiceBrokerInterface {
  private utf8Decoder: TextDecoder;

  /**
   * @param contract smart contract
   */
  constructor(private contract: ContractInterface) {
    this.utf8Decoder = new TextDecoder();
  }

  async request(request: ServiceRequest): Promise<void> {
    const serialized = request.serialize();
    await this.contract.submitTransaction('Request', serialized);
  }

  async respond(response: ServiceResponse): Promise<void> {
    const serialized = response.serialize();
    await this.contract.submitTransaction('Respond', serialized);
  }

  async get(requestId: string): Promise<ServiceRequestResponse> {
    const data = await this.contract.submitTransaction('Get', requestId);
    const serialized = this.utf8Decoder.decode(data);
    return ServiceRequestResponse.deserialize(serialized);
  }

  async getAll(
    organizationId: string,
    deviceId: string,
    serviceName: string,
  ): Promise<ServiceRequestResponse[]> {
    const data = await this.contract.submitTransaction(
      'GetAll',
      organizationId,
      deviceId,
      serviceName,
    );
    const serialized = this.utf8Decoder.decode(data);
    return (JSON.parse(serialized) as any[]).map(ServiceRequestResponse.fromObject);
  }

  async remove(requestId: string): Promise<void> {
    await this.contract.submitTransaction('Remove', requestId);
  }

  async registerEvent(
    options?: ChaincodeEventsOptions,
  ): Promise<CloseableAsyncIterable<ServiceRequestEvent>> {
    const events = await this.contract.registerEvent(options);
    const pattern = /^request:\/\/(.+?)\/(.+?)\/(.+?)\/(.+?)\/(.+?)$/;
    const decoder = this.utf8Decoder;

    return {
      async *[Symbol.asyncIterator]() {
        for await (const event of events) {
          // can reuse pattern here since it has no global('g') flag
          const matches = pattern.exec(event.eventName);
          if (matches === null || matches.length != 6) {
            continue;
          }

          const serviceRequestEvent: ServiceRequestEvent = {
            organizationId: matches[1],
            deviceId: matches[2],
            serviceName: matches[3],
            requestId: matches[4],
            action: matches[5],
            payload: null,
          };

          if (serviceRequestEvent.action === 'request') {
            try {
              const request = ServiceRequest.deserialize(decoder.decode(event.payload));
              serviceRequestEvent.payload = request;
            } catch {
              console.error(
                `bad service request event payload ${event.payload}, action is ${serviceRequestEvent.action}`,
              );
              continue;
            }
          } else if (serviceRequestEvent.action === 'respond') {
            try {
              const respond = ServiceResponse.deserialize(decoder.decode(event.payload));
              serviceRequestEvent.payload = respond;
            } catch {
              console.error(
                `bad service response event payload ${event.payload}, action is ${serviceRequestEvent.action}`,
              );
              continue;
            }
          } else if (serviceRequestEvent.action === 'remove') {
            serviceRequestEvent.payload = decoder.decode(event.payload);
          } else {
            serviceRequestEvent.payload = event.payload;
          }

          yield serviceRequestEvent;
        }
      },
      close() {
        events.close();
      },
    };
  }
}

export function createServiceBroker(network: Network, chaincodeId: string) {
  return new ServiceBroker(new Contract(network, chaincodeId, 'service_broker'));
}
