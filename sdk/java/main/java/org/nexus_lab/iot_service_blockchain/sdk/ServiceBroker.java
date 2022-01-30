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

/** Core utilities for managing IoT service requests and responses on the ledger */
public class ServiceBroker implements ServiceBrokerInterface {
  private ContractInterface contract;

  /** @param contract smart contract */
  public ServiceBroker(ContractInterface contract) {
    this.contract = contract;
  }

  @Override
  public void request(ServiceRequest request)
      throws EndorseException, SubmitException, CommitStatusException, CommitException {
    String serialized = request.serialize();
    this.contract.submitTransaction("Request", serialized);
  }

  @Override
  public void respond(ServiceResponse response)
      throws EndorseException, SubmitException, CommitStatusException, CommitException {
    String serialized = response.serialize();
    this.contract.submitTransaction("Respond", serialized);
  }

  @Override
  public ServiceRequestResponse get(String requestId)
      throws EndorseException, SubmitException, CommitStatusException, CommitException {
    byte[] data = this.contract.submitTransaction("Get", requestId);
    return ServiceRequestResponse.deserialize(new String(data));
  }

  @Override
  public List<ServiceRequestResponse> getAll(
      String organizationId, String deviceId, String serviceName)
      throws EndorseException, SubmitException, CommitStatusException, CommitException {
    byte[] data = this.contract.submitTransaction("GetAll", organizationId, deviceId, serviceName);
    return Json.deserialize(new String(data), new GenericType<List<ServiceRequestResponse>>() {});
  }

  @Override
  public void remove(String requestId)
      throws EndorseException, SubmitException, CommitStatusException, CommitException {
    this.contract.submitTransaction("Remove", requestId);
  }

  @Override
  public CloseableIterator<ServiceRequestEvent> registerEvent(CallOption... options) {
    return new ServiceRequestEventIterator(this.contract.registerEvent(options));
  }

  private static final class ServiceRequestEventIterator
      extends TransformCloseableIterator<ChaincodeEvent, ServiceRequestEvent> {
    private static final Pattern EVENT_NAME_PATTERN =
        Pattern.compile("^request:\\/\\/(.+?)\\/(.+?)\\/(.+?)\\/(.+?)\\/(.+?)$");

    public ServiceRequestEventIterator(CloseableIterator<ChaincodeEvent> iterator) {
      super(iterator);
    }

    @Override
    public boolean canTransform(ChaincodeEvent event) {
      Matcher matcher =
          ServiceRequestEventIterator.EVENT_NAME_PATTERN.matcher(event.getEventName());
      return matcher.matches() && matcher.groupCount() == 5;
    }

    @Override
    public ServiceRequestEvent transform(ChaincodeEvent event) {
      Matcher matcher =
          ServiceRequestEventIterator.EVENT_NAME_PATTERN.matcher(event.getEventName());
      matcher.matches();

      ServiceRequestEvent serviceRequestEvent = new ServiceRequestEvent();
      serviceRequestEvent.setOrganizationId(matcher.group(1));
      serviceRequestEvent.setDeviceId(matcher.group(2));
      serviceRequestEvent.setServiceName(matcher.group(3));
      serviceRequestEvent.setRequestId(matcher.group(4));
      serviceRequestEvent.setAction(matcher.group(5));

      switch (serviceRequestEvent.getAction()) {
        case "request":
          try {
            serviceRequestEvent.setPayload(
                ServiceRequest.deserialize(new String(event.getPayload())));
          } catch (Exception e) {
            System.err.println(
                String.format(
                    "bad service request event payload %s, action is %s",
                    new String(event.getPayload()), serviceRequestEvent.getAction()));
          }
          break;
        case "respond":
          try {
            serviceRequestEvent.setPayload(
                ServiceResponse.deserialize(new String(event.getPayload())));
          } catch (Exception e) {
            System.err.println(
                String.format(
                    "bad service response event payload %s, action is %s",
                    new String(event.getPayload()), serviceRequestEvent.getAction()));
          }
          break;
        case "remove":
          serviceRequestEvent.setPayload(new String(event.getPayload()));
          break;
        default:
          serviceRequestEvent.setPayload(event.getPayload());
          break;
      }

      return serviceRequestEvent;
    }
  }

  /**
   * The default factory for creating service brokers
   *
   * @param network Hyperledger Fabric network
   * @param chaincodeId ID/name of the chaincode
   * @return The service broker
   */
  public static ServiceBroker create(Network network, String chaincodeId) {
    return new ServiceBroker(new Contract(network, chaincodeId, "service_broker"));
  }
}
