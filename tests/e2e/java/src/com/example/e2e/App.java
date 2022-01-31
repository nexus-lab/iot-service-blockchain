package com.example.e2e;

import java.time.OffsetDateTime;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.UUID;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;
import java.util.logging.Logger;
import org.hyperledger.fabric.client.CloseableIterator;
import org.nexus_lab.iot_service_blockchain.sdk.Device;
import org.nexus_lab.iot_service_blockchain.sdk.Sdk;
import org.nexus_lab.iot_service_blockchain.sdk.SdkOptions;
import org.nexus_lab.iot_service_blockchain.sdk.Service;
import org.nexus_lab.iot_service_blockchain.sdk.ServiceRequest;
import org.nexus_lab.iot_service_blockchain.sdk.ServiceRequestEvent;
import org.nexus_lab.iot_service_blockchain.sdk.ServiceResponse;

public class App {
  private static final String ORG_ID = "Org2MSP";
  private static final String ORG_DOMAIN = "org2.example.com";
  private static final String USER_NAME = "User1@org2.example.com";
  private static final String PEER_NAME = "peer0.org2.example.com";
  private static final String PEER_ENDPOINT = "localhost:9051";

  private static final Logger LOGGER = Logger.getLogger("app");

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
    options.setGatewayPeerTLSCertificate(credentials[2]);
    options.setNetworkName("mychannel");
    options.setChaincodeId("iotservice");

    Sdk isb = new Sdk(options);

    Utils.log(LOGGER, "Waiting for device and services to be registered");
    TimeUnit.SECONDS.sleep(30);

    Device device = getDevice(isb);
    List<Service> services = getServices(isb, device);
    sendServiceRequests(isb, services);

    isb.close();
    System.exit(0);
  }

  private static Device getDevice(Sdk isb) throws Exception {
    ScheduledFuture<?> timeout =
        Executors.newSingleThreadScheduledExecutor()
            .schedule(() -> Utils.fatal(LOGGER, "timed out getting device"), 30, TimeUnit.SECONDS);

    while (true) {
      List<Device> devices = isb.getDeviceRegistry().getAll("Org1MSP");
      if (devices.size() > 0) {
        Utils.log(LOGGER, "Found device " + devices.get(0).serialize());
        timeout.cancel(true);
        return devices.get(0);
      }
    }
  }

  private static List<Service> getServices(Sdk isb, Device device) throws Exception {
    ScheduledFuture<?> timeout =
        Executors.newSingleThreadScheduledExecutor()
            .schedule(
                () -> Utils.fatal(LOGGER, "timed out getting services"), 30, TimeUnit.SECONDS);

    while (true) {
      List<Service> services =
          isb.getServiceRegistry().getAll(device.getOrganizationId(), device.getId());
      if (services.size() >= 2) {
        for (Service service : services) {
          Utils.log(LOGGER, "Found service " + service.serialize());
        }
        timeout.cancel(true);
        return services;
      }
    }
  }

  private static void sendServiceRequests(Sdk isb, List<Service> services) throws Exception {
    ScheduledFuture<?> timeout =
        Executors.newSingleThreadScheduledExecutor()
            .schedule(
                () -> Utils.fatal(LOGGER, "timed out waiting for responses"), 60, TimeUnit.SECONDS);

    CloseableIterator<ServiceRequestEvent> events = isb.getServiceBroker().registerEvent();

    Map<String, ServiceRequest> requests = new HashMap<>();
    for (Service service : services) {
      ServiceRequest request = new ServiceRequest();
      request.setId(UUID.randomUUID().toString());
      request.setTime(OffsetDateTime.now());
      request.setService(service);
      request.setMethod("GET");
      request.setArguments(new String[] {"1", "2", "3"});

      Utils.log(LOGGER, "Sending request " + request.serialize());
      isb.getServiceBroker().request(request);
      requests.put(request.getId(), request);
    }

    Utils.log(LOGGER, "Listening for responses");
    events.forEachRemaining(
        event -> {
          if ("respond".equals(event.getAction())) {
            ServiceResponse response = (ServiceResponse) event.getPayload();
            ServiceRequest request = requests.get(response.getRequestId());

            Utils.log(LOGGER, "Received response " + response.serialize());

            if (response.getStatusCode() != 0) {
              Utils.fatal(LOGGER, "response error, status code is " + response.getStatusCode());
            }

            String returnValue =
                String.join(",", request.getMethod(), String.join(",", request.getArguments()));
            if (!returnValue.equals(response.getReturnValue())) {
              Utils.fatal(
                  LOGGER,
                  "response return value mismatch, "
                      + returnValue
                      + " != "
                      + response.getReturnValue());
            }

            requests.remove(response.getRequestId());

            if (requests.size() == 0) {
              timeout.cancel(true);
              events.close();
            }
          }
        });
  }
}
