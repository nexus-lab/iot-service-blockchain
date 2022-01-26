import { parse as uuidparse } from 'uuid';

import moment from './moment';

/**
 * An IoT service response
 */
export default class ServiceResponse {
  /**
   * @param requestId identity of the IoT service request to respond to
   * @param time time of the IoT service response
   * @param statusCode status code of the IoT service response
   * @param returnValue return value of the IoT service response
   */
  constructor(
    public requestId: string,
    public time: moment.Moment,
    public statusCode: number = 0,
    public returnValue: string = '',
  ) {}

  /**
   * Get components that compose the IoT service response key
   *
   * @returns components that compose the IoT service response key
   */
  getKeyComponents() {
    return [this.requestId];
  }

  /**
   * Transform current service response to a plain object
   *
   * @returns a plain object containing the service response information
   */
  toObject(): { [key: string]: any } {
    return {
      requestId: this.requestId,
      time: this.time,
      statusCode: this.statusCode,
      returnValue: this.returnValue,
    };
  }

  /**
   * Transform current IoT service response to JSON string
   *
   * @returns JSON representation of the service response
   */
  serialize() {
    return JSON.stringify(this.toObject());
  }

  /**
   * Check if the IoT service response properties are valid
   *
   * @throws error when service response fields are invalid
   */
  validate() {
    try {
      uuidparse(this.requestId);
    } catch {
      throw new Error('invalid request ID in response definition');
    }
    if (this.time.valueOf() === 0) {
      throw new Error('missing response time in response definition');
    }
  }

  /**
   * Create a new service response instance from an object
   *
   * @param obj an object that contains the service response information
   * @returns a new service response instance
   */
  static fromObject(obj: { [key: string]: any }): ServiceResponse {
    return new ServiceResponse(
      obj.requestId,
      moment(obj.time),
      obj.statusCode ? obj.statusCode : 0,
      obj.returnValue ? obj.returnValue : '',
    );
  }

  /**
   * Create an IoT service response instance from its JSON representation
   *
   * @param data JSON string representing a service response
   * @returns a new service response instance
   */
  static deserialize(data: string) {
    return ServiceResponse.fromObject(JSON.parse(data));
  }
}
