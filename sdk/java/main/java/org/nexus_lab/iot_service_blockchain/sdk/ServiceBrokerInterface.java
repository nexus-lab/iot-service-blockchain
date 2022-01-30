package org.nexus_lab.iot_service_blockchain.sdk;

import java.util.List;
import org.hyperledger.fabric.client.CallOption;
import org.hyperledger.fabric.client.CloseableIterator;
import org.hyperledger.fabric.client.CommitException;
import org.hyperledger.fabric.client.CommitStatusException;
import org.hyperledger.fabric.client.EndorseException;
import org.hyperledger.fabric.client.SubmitException;

/** Interface of core utilities for managing service requests on ledger */
public interface ServiceBrokerInterface {
  /**
   * Make a request to an IoT service
   *
   * @param request IoT service request to be sent
   * @throws EndorseException if the endorse invocation fails
   * @throws SubmitException if the submit invocation fails
   * @throws CommitStatusException if the commit status invocation fails
   * @throws CommitException if the transaction commits unsuccessfully
   */
  public void request(ServiceRequest request)
      throws EndorseException, SubmitException, CommitStatusException, CommitException;

  /**
   * Respond to an IoT service request
   *
   * @param response IoT service response to be sent
   * @throws EndorseException if the endorse invocation fails
   * @throws SubmitException if the submit invocation fails
   * @throws CommitStatusException if the commit status invocation fails
   * @throws CommitException if the transaction commits unsuccessfully
   */
  public void respond(ServiceResponse response)
      throws EndorseException, SubmitException, CommitStatusException, CommitException;

  /**
   * Get an IoT service request and its response (if any) by the request ID
   *
   * @param requestId service request ID
   * @return the service request and its response (if any)
   * @throws EndorseException if the endorse invocation fails
   * @throws SubmitException if the submit invocation fails
   * @throws CommitStatusException if the commit status invocation fails
   * @throws CommitException if the transaction commits unsuccessfully
   */
  public ServiceRequestResponse get(String requestId)
      throws EndorseException, SubmitException, CommitStatusException, CommitException;

  /**
   * Get a list of IoT service requests and their responses (if any) by their service organization
   * ID, service device ID, and service name
   *
   * @param organizationId organization ID of the requested service
   * @param deviceId device's ID of the requested service
   * @param serviceName name of the requested service
   * @return all services of the device in the organization
   * @throws EndorseException if the endorse invocation fails
   * @throws SubmitException if the submit invocation fails
   * @throws CommitStatusException if the commit status invocation fails
   * @throws CommitException if the transaction commits unsuccessfully
   */
  public List<ServiceRequestResponse> getAll(
      String organizationId, String deviceId, String serviceName)
      throws EndorseException, SubmitException, CommitStatusException, CommitException;

  /**
   * Remove a service request and its response (if any) from the ledger
   *
   * @param requestId ID of the service request and response to be removed
   * @throws EndorseException if the endorse invocation fails
   * @throws SubmitException if the submit invocation fails
   * @throws CommitStatusException if the commit status invocation fails
   * @throws CommitException if the transaction commits unsuccessfully
   */
  public void remove(String requestId)
      throws EndorseException, SubmitException, CommitStatusException, CommitException;

  /**
   * Registers for service request events
   *
   * @param options chaincode event options
   * @return an async iterable of the events
   */
  public CloseableIterator<ServiceRequestEvent> registerEvent(CallOption... options);
}
