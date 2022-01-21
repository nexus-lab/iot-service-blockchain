/**
 * An IoT service state
 */
export default class Service {
  /**
   * @param name friendly name of the IoT service
   * @param deviceId identity of the device to which the IoT service belongs
   * @param organizationId identity of the organization to which the IoT service belongs
   * @param version version number of the IoT service
   * @param description a brief summary of the service's functions
   * @param lastUpdateTime the latest time that the service state has been updated
   */
  constructor(
    public name: string,
    public deviceId: string,
    public organizationId: string,
    public version: number = 0,
    public description: string = '',
    public lastUpdateTime: Date = new Date(0),
  ) {}

  /**
   * Get components that compose the service key
   *
   * @returns components that compose the service key
   */
  getKeyComponents() {
    return [this.organizationId, this.deviceId, this.name];
  }

  /**
   * Transform current service to a plain object
   *
   * @returns a plain object containing the service information
   */
  toObject(): { [key: string]: any } {
    return {
      name: this.name,
      deviceId: this.deviceId,
      organizationId: this.organizationId,
      version: this.version,
      description: this.description,
      lastUpdateTime: this.lastUpdateTime,
    };
  }

  /**
   * Transform current service to JSON string
   *
   * @returns JSON representation of the service
   */
  serialize() {
    return JSON.stringify(this.toObject());
  }

  /**
   * Check if the IoT service properties are valid
   *
   * @throws error when service fields are invalid
   */
  validate() {
    if (this.name === '') {
      throw new Error('missing service name in service definition');
    }
    if (this.deviceId === '') {
      throw new Error('missing device ID in service definition');
    }
    if (this.organizationId === '') {
      throw new Error('missing organization ID in service definition');
    }
    if (this.version === 0) {
      throw new Error('missing service version in service definition');
    }
    if (this.version !== (this.version | 0) || this.version < 0) {
      throw new Error('service version must be a positive integer');
    }
    if (this.lastUpdateTime.getTime() === 0) {
      throw new Error('missing service last update time in service definition');
    }
  }

  /**
   * Create a new device instance from an object
   *
   * @param obj an object that contains the device information
   * @returns a new device instance
   */
  static fromObject(obj: { [key: string]: any }): Service {
    return new Service(
      obj.name,
      obj.deviceId,
      obj.organizationId,
      obj.version ? obj.version : 0,
      obj.description ? obj.description : '',
      obj.lastUpdateTime ? new Date(obj.lastUpdateTime) : new Date(0),
    );
  }

  /**
   * Create an IoT service instance from its JSON representation
   *
   * @param data JSON string representing a service
   * @returns a new service instance
   */
  static deserialize(data: string) {
    return Service.fromObject(JSON.parse(data));
  }
}
