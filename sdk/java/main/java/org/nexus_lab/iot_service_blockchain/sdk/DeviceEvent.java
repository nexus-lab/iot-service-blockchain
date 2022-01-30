package org.nexus_lab.iot_service_blockchain.sdk;

import lombok.Data;
import lombok.NoArgsConstructor;

/** The event emitted by the device registry contract notifying a device update */
@Data
@NoArgsConstructor
public class DeviceEvent {
  /** Name of the action performed on the device */
  private String action;

  /** Organization ID of the device */
  private String organizationId;

  /** ID of the device */
  private String deviceId;

  /** Custom event payload */
  private Object payload;
}
