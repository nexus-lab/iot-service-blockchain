package org.nexus_lab.iot_service_blockchain.sdk;

import java.util.List;
import org.hyperledger.fabric.client.CallOption;
import org.hyperledger.fabric.client.CloseableIterator;
import org.hyperledger.fabric.client.CommitException;
import org.hyperledger.fabric.client.CommitStatusException;
import org.hyperledger.fabric.client.EndorseException;
import org.hyperledger.fabric.client.SubmitException;

/** Interface of core utilities for managing services on the ledger. */
public interface ServiceRegistryInterface {
  /**
   * Create or update a service in the ledger.
   *
   * @param service service to be created or updated
   * @throws EndorseException if the endorse invocation fails
   * @throws SubmitException if the submit invocation fails
   * @throws CommitStatusException if the commit status invocation fails
   * @throws CommitException if the transaction commits unsuccessfully
   */
  public void register(Service service)
      throws EndorseException, SubmitException, CommitStatusException, CommitException;

  /**
   * Get a service by its organization ID, device ID, and name.
   *
   * @param organizationId service's organization ID
   * @param deviceId device's ID
   * @param serviceName name of the service
   * @return the device
   * @throws EndorseException if the endorse invocation fails
   * @throws SubmitException if the submit invocation fails
   * @throws CommitStatusException if the commit status invocation fails
   * @throws CommitException if the transaction commits unsuccessfully
   */
  public Service get(String organizationId, String deviceId, String serviceName)
      throws EndorseException, SubmitException, CommitStatusException, CommitException;

  /**
   * Get a list of services by their organization ID and device ID.
   *
   * @param organizationId device's organization ID
   * @param deviceId device's ID
   * @return all services of the device in the organization
   * @throws EndorseException if the endorse invocation fails
   * @throws SubmitException if the submit invocation fails
   * @throws CommitStatusException if the commit status invocation fails
   * @throws CommitException if the transaction commits unsuccessfully
   */
  public List<Service> getAll(String organizationId, String deviceId)
      throws EndorseException, SubmitException, CommitStatusException, CommitException;

  /**
   * Remove a service from the ledger.
   *
   * @param service the service to be removed
   * @throws EndorseException if the endorse invocation fails
   * @throws SubmitException if the submit invocation fails
   * @throws CommitStatusException if the commit status invocation fails
   * @throws CommitException if the transaction commits unsuccessfully
   */
  public void deregister(Service service)
      throws EndorseException, SubmitException, CommitStatusException, CommitException;

  /**
   * Registers for service registry events.
   *
   * @param options chaincode event options
   * @return an async iterable of the events
   */
  public CloseableIterator<ServiceEvent> registerEvent(CallOption... options);
}
