package org.nexus_lab.iot_service_blockchain.sdk;

import com.owlike.genson.GenericType;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import org.hyperledger.fabric.client.CallOption;
import org.hyperledger.fabric.client.ChaincodeEvent;
import org.hyperledger.fabric.client.CloseableIterator;
import org.hyperledger.fabric.client.CommitException;
import org.hyperledger.fabric.client.CommitStatusException;
import org.hyperledger.fabric.client.EndorseException;
import org.hyperledger.fabric.client.Network;
import org.hyperledger.fabric.client.SubmitException;

/** Core utilities for managing devices on the ledger */
public class DeviceRegistry implements DeviceRegistryInterface {
  private ContractInterface contract;

  /** @param contract smart contract */
  public DeviceRegistry(ContractInterface contract) {
    this.contract = contract;
  }

  @Override
  public void register(Device device)
      throws EndorseException, SubmitException, CommitStatusException, CommitException {
    String serialized = device.serialize();
    this.contract.submitTransaction("Register", serialized);
  }

  @Override
  public Device get(String organizationId, String deviceId)
      throws EndorseException, SubmitException, CommitStatusException, CommitException {
    byte[] data = this.contract.submitTransaction("Get", organizationId, deviceId);
    return Device.deserialize(new String(data));
  }

  @Override
  public List<Device> getAll(String organizationId)
      throws EndorseException, SubmitException, CommitStatusException, CommitException {
    byte[] data = this.contract.submitTransaction("GetAll", organizationId);
    return Json.deserialize(new String(data), new GenericType<List<Device>>() {});
  }

  @Override
  public void deregister(Device device)
      throws EndorseException, SubmitException, CommitStatusException, CommitException {
    String serialized = device.serialize();
    this.contract.submitTransaction("Deregister", serialized);
  }

  @Override
  public CloseableIterator<DeviceEvent> registerEvent(CallOption... options) {
    return new DeviceEventIterator(this.contract.registerEvent(options));
  }

  private static final class DeviceEventIterator
      extends TransformCloseableIterator<ChaincodeEvent, DeviceEvent> {
    private static final Pattern EVENT_NAME_PATTERN =
        Pattern.compile("^device:\\/\\/(.+?)\\/(.+?)\\/(.+?)$");

    public DeviceEventIterator(CloseableIterator<ChaincodeEvent> iterator) {
      super(iterator);
    }

    @Override
    public boolean canTransform(ChaincodeEvent event) {
      Matcher matcher = DeviceEventIterator.EVENT_NAME_PATTERN.matcher(event.getEventName());
      return matcher.matches() && matcher.groupCount() == 3;
    }

    @Override
    public DeviceEvent transform(ChaincodeEvent event) {
      Matcher matcher = DeviceEventIterator.EVENT_NAME_PATTERN.matcher(event.getEventName());
      matcher.matches();

      DeviceEvent deviceEvent = new DeviceEvent();
      deviceEvent.setOrganizationId(matcher.group(1));
      deviceEvent.setDeviceId(matcher.group(2));
      deviceEvent.setAction(matcher.group(3));

      if ("register".equals(deviceEvent.getAction())
          || "deregister".equals(deviceEvent.getAction())) {
        try {
          deviceEvent.setPayload(Device.deserialize(new String(event.getPayload())));
        } catch (Exception e) {
          System.err.println(
              String.format(
                  "bad device event payload %s, action is %s",
                  new String(event.getPayload()), deviceEvent.getAction()));
        }
      } else {
        deviceEvent.setPayload(event.getPayload());
      }

      return deviceEvent;
    }
  }

  /**
   * The default factory for creating device registries
   *
   * @param network Hyperledger Fabric network
   * @param chaincodeId ID/name of the chaincode
   * @return The device registry
   */
  public static DeviceRegistry create(Network network, String chaincodeId) {
    return new DeviceRegistry(new Contract(network, chaincodeId, "device_registry"));
  }
}
