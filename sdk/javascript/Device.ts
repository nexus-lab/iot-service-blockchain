import moment from './moment';

/**
 * An IoT device state
 */
export default class Device {
  /**
   * @param id identity of the device
   * @param organizationId identity of the organization to which the device belongs
   * @param name friendly name of the device
   * @param description a brief summary of the device's functions
   * @param lastUpdateTime the latest time that the device state has been updated
   */
  constructor(
    public id: string,
    public organizationId: string,
    public name: string = '',
    public description: string = '',
    public lastUpdateTime: moment.Moment = moment(0),
  ) {}

  /**
   * Get components that compose the device key
   *
   * @returns components that compose the device key
   */
  getKeyComponents(): string[] {
    return [this.organizationId, this.id];
  }

  /**
   * Transform current device to a plain object
   *
   * @returns a plain object containing the device information
   */
  toObject(): { [key: string]: any } {
    return {
      id: this.id,
      organizationId: this.organizationId,
      name: this.name,
      description: this.description,
      lastUpdateTime: this.lastUpdateTime,
    };
  }

  /**
   * Transform current device to JSON string
   *
   * @returns JSON representation of the device
   */
  serialize(): string {
    return JSON.stringify(this.toObject());
  }

  /**
   * Check if the device properties are valid
   *
   * @throws error when device fields are invalid
   */
  validate() {
    if (this.id === '') {
      throw new Error('missing device ID in device definition');
    }
    if (this.organizationId === '') {
      throw new Error('missing organization ID in device definition');
    }
    if (this.name === '') {
      throw new Error('missing device name in device definition');
    }
    if (this.lastUpdateTime.valueOf() === 0) {
      throw new Error('missing device last update time in device definition');
    }
  }

  /**
   * Create a new device instance from an object
   *
   * @param obj an object that contains the device information
   * @returns a new device instance
   */
  static fromObject(obj: { [key: string]: any }): Device {
    return new Device(
      obj.id,
      obj.organizationId,
      obj.name,
      obj.description ? obj.description : '',
      obj.lastUpdateTime ? moment(obj.lastUpdateTime) : moment(0),
    );
  }

  /**
   * Create a new device instance from its JSON representation
   *
   * @param data JSON string representing a device
   * @returns a new device instance
   */
  static deserialize(data: string): Device {
    return Device.fromObject(JSON.parse(data));
  }
}
