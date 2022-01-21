import { parse as uuidparse } from 'uuid';

import Service from './Service';

/**
 * An IoT service request
 */
export default class ServiceRequest {
  /**
   * @param id identity of the IoT service request
   * @param time time of the IoT service request
   * @param service requested IoT service information
   * @param method IoT service request method
   * @param args IoT service request arguments
   */
  constructor(
    public id: string,
    public time: Date,
    public service: Service,
    public method: string,
    public args: string[] = [],
  ) {}

  /**
   * Get components that compose the IoT service request key
   *
   * @returns components that compose the IoT service request key
   */
  getKeyComponents() {
    return [this.id];
  }

  /**
   * Transform current service request to a plain object
   *
   * @returns a plain object containing the service request information
   */
  toObject(): { [key: string]: any } {
    return {
      id: this.id,
      time: this.time,
      service: this.service.toObject(),
      method: this.method,
      arguments: this.args,
    };
  }

  /**
   * Transform current IoT service request to JSON string
   *
   * @returns JSON representation of the service request
   */
  serialize() {
    return JSON.stringify(this.toObject());
  }

  /**
   * Check if the IoT service request properties are valid
   *
   * @throws error when service request fields are invalid
   */
  validate() {
    try {
      uuidparse(this.id);
    } catch {
      throw new Error('invalid request ID in request definition');
    }
    if (
      this.service.organizationId === '' ||
      this.service.deviceId === '' ||
      this.service.name === ''
    ) {
      throw new Error('missing requested service in request definition');
    }
    if (this.method === '') {
      throw new Error('missing request method in request definition');
    }
    if (this.time.getTime() === 0) {
      throw new Error('missing request time in request definition');
    }
  }

  /**
   * Create a new service request instance from an object
   *
   * @param obj an object that contains the service request information
   * @returns a new service request instance
   */
  static fromObject(obj: { [key: string]: any }): ServiceRequest {
    return new ServiceRequest(
      obj.id,
      new Date(obj.time),
      Service.fromObject(obj.service),
      obj.method,
      obj.arguments ? obj.arguments : [],
    );
  }

  /**
   * Create an IoT service request instance from its JSON representation
   *
   * @param data JSON string representing a service request
   * @returns a new service request instance
   */
  static deserialize(data: string) {
    return ServiceRequest.fromObject(JSON.parse(data));
  }
}
