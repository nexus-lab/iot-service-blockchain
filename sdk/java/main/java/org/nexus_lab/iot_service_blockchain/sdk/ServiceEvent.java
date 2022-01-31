package org.nexus_lab.iot_service_blockchain.sdk;

import lombok.Data;
import lombok.NoArgsConstructor;

/** The event emitted by the service registry contract notifying a service update. */
@Data
@NoArgsConstructor
public class ServiceEvent {
  /** Name of the action performed on the service. */
  private String action;

  /** Organization ID of the service. */
  private String organizationId;

  /** ID of the device to which the service belongs. */
  private String deviceId;

  /** ID of the service. */
  private String serviceName;

  /** Custom event payload. */
  private Object payload;
}
