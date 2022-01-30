package org.nexus_lab.iot_service_blockchain.sdk;

import java.util.List;
import org.hyperledger.fabric.client.CallOption;
import org.hyperledger.fabric.client.CloseableIterator;
import org.hyperledger.fabric.client.CommitException;
import org.hyperledger.fabric.client.CommitStatusException;
import org.hyperledger.fabric.client.EndorseException;
import org.hyperledger.fabric.client.SubmitException;

/** Interface of core utilities for managing devices on the ledger */
public interface DeviceRegistryInterface {
  /**
   * Create or update a device in the ledger
   *
   * @param device device to be created or updated
   * @throws EndorseException if the endorse invocation fails
   * @throws SubmitException if the submit invocation fails
   * @throws CommitStatusException if the commit status invocation fails
   * @throws CommitException if the transaction commits unsuccessfully
   */
  public void register(Device device)
      throws EndorseException, SubmitException, CommitStatusException, CommitException;

  /**
   * Get a device by its organization ID and device ID
   *
   * @param organizationId device's organization ID
   * @param deviceId device's ID
   * @return the device
   * @throws EndorseException if the endorse invocation fails
   * @throws SubmitException if the submit invocation fails
   * @throws CommitStatusException if the commit status invocation fails
   * @throws CommitException if the transaction commits unsuccessfully
   */
  public Device get(String organizationId, String deviceId)
      throws EndorseException, SubmitException, CommitStatusException, CommitException;

  /**
   * Get a list of devices by their organization ID
   *
   * @param organizationId devices' organization ID
   * @return all devices of the organization
   * @throws EndorseException if the endorse invocation fails
   * @throws SubmitException if the submit invocation fails
   * @throws CommitStatusException if the commit status invocation fails
   * @throws CommitException if the transaction commits unsuccessfully
   */
  public List<Device> getAll(String organizationId)
      throws EndorseException, SubmitException, CommitStatusException, CommitException;

  /**
   * Remove a device from the ledger
   *
   * @param device the device to be removed
   * @throws EndorseException if the endorse invocation fails
   * @throws SubmitException if the submit invocation fails
   * @throws CommitStatusException if the commit status invocation fails
   * @throws CommitException if the transaction commits unsuccessfully
   */
  public void deregister(Device device)
      throws EndorseException, SubmitException, CommitStatusException, CommitException;

  /**
   * Registers for device registry events
   *
   * @param options chaincode event options
   * @return an async iterable of the events
   */
  public CloseableIterator<DeviceEvent> registerEvent(CallOption... options);
}
