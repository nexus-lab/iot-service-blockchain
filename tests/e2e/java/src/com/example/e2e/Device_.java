package com.example.e2e;

import java.time.OffsetDateTime;
import java.util.List;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.logging.Logger;
import org.hyperledger.fabric.client.CloseableIterator;
import org.nexus_lab.iot_service_blockchain.sdk.Device;
import org.nexus_lab.iot_service_blockchain.sdk.Sdk;
import org.nexus_lab.iot_service_blockchain.sdk.SdkOptions;
import org.nexus_lab.iot_service_blockchain.sdk.Service;
import org.nexus_lab.iot_service_blockchain.sdk.ServiceRequest;
import org.nexus_lab.iot_service_blockchain.sdk.ServiceRequestEvent;
import org.nexus_lab.iot_service_blockchain.sdk.ServiceRequestResponse;
import org.nexus_lab.iot_service_blockchain.sdk.ServiceResponse;

public class Device_ {
  private static final String ORG_ID = "Org1MSP";
  private static final String ORG_DOMAIN = "org1.example.com";
  private static final String USER_NAME = "User1@org1.example.com";
  private static final String PEER_NAME = "peer0.org1.example.com";
  private static final String PEER_ENDPOINT = "localhost:7051";

  private static final Logger LOGGER = Logger.getLogger("device");

  public static void main(String[] args) throws Exception {
    String fabricRoot = System.getenv("FABRIC_ROOT");
    String[] credentials =
        Utils.getCredentials(
            fabricRoot == null ? "" : fabricRoot, ORG_DOMAIN, USER_NAME, PEER_NAME);

    SdkOptions options = new SdkOptions();
    options.setOrganizationId(ORG_ID);
    options.setCertificate(credentials[0]);
    options.setPrivateKey(credentials[1]);
    options.setGatewayPeerEndpoint(PEER_ENDPOINT);
    options.setGatewayPeerServerName(PEER_NAME);
    options.setGatewayPeerTlsCertificate(credentials[2]);
    options.setNetworkName("mychannel");
    options.setChaincodeId("iotservice");

    Sdk isb = new Sdk(options);

    Utils.log(LOGGER, "Organization ID is " + isb.getOrganizationId());
    Utils.log(LOGGER, "Device ID is " + isb.getDeviceId());

    registerDevice(isb);
    registerServices(isb);
    handleRequests(isb);
    checkAndRemoveRequests(isb);
    deregisterOneService(isb);
    deregisterDevice(isb);

    isb.close();
    System.exit(0);
  }

  private static void registerDevice(Sdk isb) throws Exception {
    Device expected = new Device();
    expected.setId(isb.getDeviceId());
    expected.setOrganizationId(isb.getOrganizationId());
    expected.setName("device1");
    expected.setDescription("My first device");
    expected.setLastUpdateTime(OffsetDateTime.now());

    isb.getDeviceRegistry().register(expected);
    Utils.log(LOGGER, "Registered device " + expected.getId());

    Device actual = isb.getDeviceRegistry().get(isb.getOrganizationId(), expected.getId());
    if (!expected.serialize().equals(actual.serialize())) {
      Utils.fatal(
          LOGGER,
          "inconsistent device information after registration: "
              + actual.serialize()
              + " != "
              + expected.serialize());
    }

    try {
      isb.getDeviceRegistry().get(isb.getOrganizationId(), "invalid_id");
      Utils.fatal(LOGGER, "should return error when device is not found");
    } catch (Exception ignored) {
    }

    List<Device> devices = isb.getDeviceRegistry().getAll(isb.getOrganizationId());
    if (devices.size() != 1) {
      Utils.fatal(LOGGER, "should return only 1 device from " + isb.getOrganizationId());
    }
    actual = devices.get(0);
    if (!expected.serialize().equals(actual.serialize())) {
      Utils.fatal(
          LOGGER,
          "inconsistent device information after registration: "
              + actual.serialize()
              + " != "
              + expected.serialize());
    }
  }

  private static void registerServices(Sdk isb) throws Exception {
    Service service1 = new Service();
    service1.setName("service1");
    service1.setDeviceId(isb.getDeviceId());
    service1.setOrganizationId(isb.getOrganizationId());
    service1.setVersion(1);
    service1.setDescription("My first service");
    service1.setLastUpdateTime(OffsetDateTime.now());

    Service service2 = new Service();
    service2.setName("service2");
    service2.setDeviceId(isb.getDeviceId());
    service2.setOrganizationId(isb.getOrganizationId());
    service2.setVersion(1);
    service2.setDescription("My second service");
    service2.setLastUpdateTime(OffsetDateTime.now());

    Service[] services = new Service[] {service1, service2};

    for (Service expected : services) {
      isb.getServiceRegistry().register(expected);
      Utils.log(LOGGER, "Registered service " + expected.serialize());

      Service actual =
          isb.getServiceRegistry()
              .get(isb.getOrganizationId(), expected.getDeviceId(), expected.getName());
      if (!expected.serialize().equals(actual.serialize())) {
        Utils.fatal(
            LOGGER,
            "inconsistent service information after registration: "
                + actual.serialize()
                + " != "
                + expected.serialize());
      }
    }

    try {
      isb.getServiceRegistry().get(isb.getOrganizationId(), isb.getDeviceId(), "invalid_id");
      Utils.fatal(LOGGER, "should return error when service is not found");
    } catch (Exception ignored) {
    }

    List<Service> actuals =
        isb.getServiceRegistry().getAll(isb.getOrganizationId(), isb.getDeviceId());
    if (actuals.size() != services.length) {
      Utils.fatal(
          LOGGER,
          "should return only " + services.length + " device from " + isb.getOrganizationId());
    }

    for (int i = 0; i < services.length; i++) {
      Service expected = services[i];
      Service actual = actuals.get(i);

      if (!expected.serialize().equals(actual.serialize())) {
        Utils.fatal(
            LOGGER,
            "inconsistent service information after registration: "
                + actual.serialize()
                + " != "
                + expected.serialize());
      }
    }
  }

  private static void handleRequests(Sdk isb) throws Exception {
    CloseableIterator<ServiceRequestEvent> events = isb.getServiceBroker().registerEvent();

    ScheduledFuture<?> timeout =
        Executors.newSingleThreadScheduledExecutor()
            .schedule(
                () -> Utils.fatal(LOGGER, "timed out waiting for responses"),
                120,
                TimeUnit.SECONDS);

    Utils.log(LOGGER, "Listening for requests");
    AtomicInteger counter = new AtomicInteger(0);
    events.forEachRemaining(
        (event) -> {
          if (!"request".equals(event.getAction())) {
            return;
          }

          ServiceRequest request = (ServiceRequest) event.getPayload();

          Utils.log(LOGGER, "Received request " + request.serialize());

          ServiceResponse response = new ServiceResponse();
          response.setRequestId(request.getId());
          response.setTime(OffsetDateTime.now());
          response.setStatusCode(0);
          response.setReturnValue(
              String.join(",", request.getMethod(), String.join(",", request.getArguments())));

          try {
            isb.getServiceBroker().respond(response);
          } catch (Exception e) {
            e.printStackTrace();
            Utils.fatal(LOGGER, "failed to respond to service request");
          }

          Utils.log(LOGGER, "Sent response " + response.serialize());

          if (counter.incrementAndGet() == 2) {
            events.close();
          }
        });

    timeout.cancel(true);
  }

  private static void checkAndRemoveRequests(Sdk isb) throws Exception {
    List<ServiceRequestResponse> before =
        isb.getServiceBroker().getAll(isb.getOrganizationId(), isb.getDeviceId(), "service1");
    if (before.size() == 0) {
      Utils.fatal(LOGGER, "request/response for service1 is not found");
    }

    ServiceRequestResponse pair = before.get(0);
    isb.getServiceBroker().get(pair.getRequest().getId());
    if (!pair.getRequest().getId().equals(pair.getResponse().getRequestId())) {
      Utils.fatal(
          LOGGER,
          "request and response ID mismatch, "
              + pair.getRequest().serialize()
              + " != "
              + pair.getResponse().serialize());
    }

    isb.getServiceBroker().remove(pair.getRequest().getId());

    List<ServiceRequestResponse> after =
        isb.getServiceBroker().getAll(isb.getOrganizationId(), isb.getDeviceId(), "service1");
    if (before.size() - 1 != after.size()) {
      Utils.fatal(
          LOGGER,
          "incorrect request/response count after removal: "
              + before.size()
              + " - 1 != "
              + after.size());
    }
  }

  private static void deregisterOneService(Sdk isb) throws Exception {
    List<Service> before =
        isb.getServiceRegistry().getAll(isb.getOrganizationId(), isb.getDeviceId());

    Service service = new Service();
    service.setName("service2");
    service.setDeviceId(isb.getDeviceId());
    service.setOrganizationId(isb.getOrganizationId());
    isb.getServiceRegistry().deregister(service);
    Utils.log(LOGGER, "Deregistered service " + service.serialize());

    try {
      isb.getServiceRegistry().get(isb.getOrganizationId(), isb.getDeviceId(), service.getName());
      Utils.fatal(LOGGER, "should return error when service is already deregistered");
    } catch (Exception ignored) {
    }

    List<Service> after =
        isb.getServiceRegistry().getAll(isb.getOrganizationId(), isb.getDeviceId());
    if (before.size() - 1 != after.size()) {
      Utils.fatal(
          LOGGER,
          "incorrect service count after deregistration: "
              + before.size()
              + " - 1 != "
              + after.size());
    }

    for (Service existing : after) {
      if (service.getName().equals(existing.getName())) {
        Utils.fatal(LOGGER, "service " + service.getName() + " has not been correctly removed");
      }
    }
  }

  private static void deregisterDevice(Sdk isb) throws Exception {
    List<Device> before = isb.getDeviceRegistry().getAll(isb.getOrganizationId());

    Device device = new Device();
    device.setId(isb.getDeviceId());
    device.setOrganizationId(isb.getOrganizationId());
    isb.getDeviceRegistry().deregister(device);
    Utils.log(LOGGER, "Deregistered device " + device.serialize());

    List<Device> after = isb.getDeviceRegistry().getAll(isb.getOrganizationId());
    if (before.size() - 1 != after.size()) {
      Utils.fatal(
          LOGGER,
          "incorrect device count after deregistration: "
              + before.size()
              + " - 1 != "
              + after.size());
    }

    for (Device existing : after) {
      if (device.getName().equals(existing.getName())) {
        Utils.fatal(LOGGER, "device " + device.getName() + " has not been correctly removed");
      }
    }

    List<Service> services =
        isb.getServiceRegistry().getAll(isb.getOrganizationId(), isb.getDeviceId());
    if (services.size() != 0) {
      Utils.fatal(LOGGER, "should have removed all services");
    }
  }
}
