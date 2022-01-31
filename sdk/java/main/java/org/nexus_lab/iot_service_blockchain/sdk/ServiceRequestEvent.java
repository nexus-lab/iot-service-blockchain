package org.nexus_lab.iot_service_blockchain.sdk;

import lombok.Data;
import lombok.NoArgsConstructor;

/** The event emitted by the service broker contract notifying a service request/response update. */
@Data
@NoArgsConstructor
public class ServiceRequestEvent {
  /** Name of the action performed on the service request. */
  private String action;

  /** Organization ID of the requested service. */
  private String organizationId;

  /** Device ID of the requested service. */
  private String deviceId;

  /** Name of the requested service. */
  private String serviceName;

  /** ID of the request. */
  private String requestId;

  /** Custom event payload. */
  private Object payload;
}
