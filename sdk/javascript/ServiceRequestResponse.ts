import ServiceRequest from './ServiceRequest';
import ServiceResponse from './ServiceResponse';

/**
 * A wrapper of a pair of IoT service request and response
 */
export default class ServiceRequestResponse {
  /**
   * @param request IoT service request
   * @param response IoT service response
   */
  constructor(public request: ServiceRequest, public response: ServiceResponse) {}

  /**
   * Transform current service request/response pair to a plain object
   *
   * @returns a plain object containing the service request/response pair information
   */
  toObject(): { [key: string]: any } {
    return {
      request: this.request.toObject(),
      response: this.response.toObject(),
    };
  }

  /**
   * Transform current IoT service request/response pair to JSON string
   *
   * @returns JSON representation of the service request/response pair
   */
  serialize() {
    return JSON.stringify(this.toObject());
  }

  /**
   * Create a new service request/response pair instance from an object
   *
   * @param obj an object that contains the service request/response pair information
   * @returns a new service request/response pair instance
   */
  static fromObject(obj: { [key: string]: any }): ServiceRequestResponse {
    return new ServiceRequestResponse(
      ServiceRequest.fromObject(obj.request),
      ServiceResponse.fromObject(obj.response),
    );
  }

  /**
   * Create an IoT service request/response pair instance from its JSON representation
   *
   * @param data JSON string representing a service request/response pair
   * @returns a new service request/response pair instance
   */
  static deserialize(data: string) {
    return ServiceRequestResponse.fromObject(JSON.parse(data));
  }
}
