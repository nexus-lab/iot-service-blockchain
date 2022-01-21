import {
  ChaincodeEventsOptions,
  CloseableAsyncIterable,
  Network,
} from '@hyperledger/fabric-gateway';
import { TextDecoder } from 'util';

import Contract, { ContractInterface } from './Contract';
import Device from './Device';

/**
 * Interface of the event emitted by the device registry contract notifying a device update
 */
export interface DeviceEvent {
  /**
   * Name of the action performed on the device
   */
  action: string;

  /**
   * Organization ID of the device
   */
  organizationId: string;

  /**
   * ID of the device
   */
  deviceId: string;

  /**
   * Custom event payload
   */
  payload: any;
}

/**
 * Interface of core utilities for managing devices on the ledger
 */
export interface DeviceRegistryInterface {
  /**
   * Create or update a device in the ledger
   *
   * @param device device to be created or updated
   */
  register(device: Device): Promise<void>;

  /**
   * Get a device by its organization ID and device ID
   *
   * @param organizationId device's organization ID
   * @param deviceId device's ID
   * @returns the device
   */
  get(organizationId: string, deviceId: string): Promise<Device>;

  /**
   * Get a list of devices by their organization ID
   *
   * @param organizationId devices' organization ID
   * @returns all devices of the organization
   */
  getAll(organizationId: string): Promise<Device[]>;

  /**
   * Remove a device from the ledger
   *
   * @param device the device to be removed
   */
  deregister(device: Device): Promise<void>;

  /**
   * Registers for device registry events
   *
   * @param options chaincode event options
   * @returns an async iterable of the events
   */
  registerEvent(options: ChaincodeEventsOptions): Promise<CloseableAsyncIterable<DeviceEvent>>;
}

/**
 * Core utilities for managing devices on the ledger
 */
export default class DeviceRegistry implements DeviceRegistryInterface {
  private utf8Decoder: TextDecoder;
  /**
   * @param contract smart contract
   */
  constructor(private contract: ContractInterface) {
    this.utf8Decoder = new TextDecoder();
  }

  async register(device: Device): Promise<void> {
    const serialized = device.serialize();
    await this.contract.submitTransaction('Register', serialized);
  }

  async get(organizationId: string, deviceId: string): Promise<Device> {
    const data = await this.contract.submitTransaction('Get', organizationId, deviceId);
    const serialized = this.utf8Decoder.decode(data);
    return Device.deserialize(serialized);
  }

  async getAll(organizationId: string): Promise<Device[]> {
    const data = await this.contract.submitTransaction('GetAll', organizationId);
    const serialized = this.utf8Decoder.decode(data);
    return (JSON.parse(serialized) as any[]).map(Device.fromObject);
  }

  async deregister(device: Device): Promise<void> {
    const serialized = device.serialize();
    await this.contract.submitTransaction('Deregister', serialized);
  }

  async registerEvent(
    options?: ChaincodeEventsOptions,
  ): Promise<CloseableAsyncIterable<DeviceEvent>> {
    const events = await this.contract.registerEvent(options);
    const pattern = /^device:\/\/(.+?)\/(.+?)\/(.+?)$/;
    const decoder = this.utf8Decoder;

    return {
      async *[Symbol.asyncIterator]() {
        for await (const event of events) {
          // can reuse pattern here since it has no global('g') flag
          const matches = pattern.exec(event.eventName);
          if (matches === null || matches.length != 4) {
            continue;
          }

          const deviceEvent: DeviceEvent = {
            organizationId: matches[1],
            deviceId: matches[2],
            action: matches[3],
            payload: null,
          };

          if (deviceEvent.action === 'register' || deviceEvent.action === 'deregister') {
            try {
              const device = Device.deserialize(decoder.decode(event.payload));
              deviceEvent.payload = device;
            } catch {
              console.error(
                `bad device event payload ${event.payload}, action is ${deviceEvent.action}`,
              );
              continue;
            }
          } else {
            deviceEvent.payload = event.payload;
          }

          yield deviceEvent;
        }
      },
      close() {
        events.close();
      },
    };
  }
}

export function createDeviceRegistry(network: Network, chaincodeId: string) {
  return new DeviceRegistry(new Contract(network, chaincodeId, 'device_registry'));
}
