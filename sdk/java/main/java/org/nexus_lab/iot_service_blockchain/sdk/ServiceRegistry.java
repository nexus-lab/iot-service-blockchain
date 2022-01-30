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
public class ServiceRegistry implements ServiceRegistryInterface {
  private ContractInterface contract;

  /** @param contract smart contract */
  public ServiceRegistry(ContractInterface contract) {
    this.contract = contract;
  }

  @Override
  public void register(Service service)
      throws EndorseException, SubmitException, CommitStatusException, CommitException {
    String serialized = service.serialize();
    this.contract.submitTransaction("Register", serialized);
  }

  @Override
  public Service get(String organizationId, String deviceId, String serviceName)
      throws EndorseException, SubmitException, CommitStatusException, CommitException {
    byte[] data = this.contract.submitTransaction("Get", organizationId, deviceId, serviceName);
    return Service.deserialize(new String(data));
  }

  @Override
  public List<Service> getAll(String organizationId, String deviceId)
      throws EndorseException, SubmitException, CommitStatusException, CommitException {
    byte[] data = this.contract.submitTransaction("GetAll", organizationId, deviceId);
    return Json.deserialize(new String(data), new GenericType<List<Service>>() {});
  }

  @Override
  public void deregister(Service service)
      throws EndorseException, SubmitException, CommitStatusException, CommitException {
    String serialized = service.serialize();
    this.contract.submitTransaction("Deregister", serialized);
  }

  @Override
  public CloseableIterator<ServiceEvent> registerEvent(CallOption... options) {
    return new ServiceEventIterator(this.contract.registerEvent(options));
  }

  private static final class ServiceEventIterator
      extends TransformCloseableIterator<ChaincodeEvent, ServiceEvent> {
    private static final Pattern EVENT_NAME_PATTERN =
        Pattern.compile("^service:\\/\\/(.+?)\\/(.+?)\\/(.+?)\\/(.+?)$");

    public ServiceEventIterator(CloseableIterator<ChaincodeEvent> iterator) {
      super(iterator);
    }

    @Override
    public boolean canTransform(ChaincodeEvent event) {
      Matcher matcher = ServiceEventIterator.EVENT_NAME_PATTERN.matcher(event.getEventName());
      return matcher.matches() && matcher.groupCount() == 4;
    }

    @Override
    public ServiceEvent transform(ChaincodeEvent event) {
      Matcher matcher = ServiceEventIterator.EVENT_NAME_PATTERN.matcher(event.getEventName());
      matcher.matches();

      ServiceEvent serviceEvent = new ServiceEvent();
      serviceEvent.setOrganizationId(matcher.group(1));
      serviceEvent.setDeviceId(matcher.group(2));
      serviceEvent.setServiceName(matcher.group(3));
      serviceEvent.setAction(matcher.group(4));

      if ("register".equals(serviceEvent.getAction())
          || "deregister".equals(serviceEvent.getAction())) {
        try {
          serviceEvent.setPayload(Service.deserialize(new String(event.getPayload())));
        } catch (Exception e) {
          System.err.println(
              String.format(
                  "bad service event payload %s, action is %s",
                  new String(event.getPayload()), serviceEvent.getAction()));
        }
      } else {
        serviceEvent.setPayload(event.getPayload());
      }

      return serviceEvent;
    }
  }

  /**
   * The default factory for creating service registries
   *
   * @param network Hyperledger Fabric network
   * @param chaincodeId ID/name of the chaincode
   * @return The service registry
   */
  public static ServiceRegistry create(Network network, String chaincodeId) {
    return new ServiceRegistry(new Contract(network, chaincodeId, "service_registry"));
  }
}
