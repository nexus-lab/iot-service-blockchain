package com.example.e2e;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.Callable;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.ThreadPoolExecutor;
import java.util.concurrent.TimeUnit;
import java.util.logging.Logger;
import org.hyperledger.fabric.client.CloseableIterator;
import org.nexus_lab.iot_service_blockchain.sdk.Device;
import org.nexus_lab.iot_service_blockchain.sdk.DeviceEvent;
import org.nexus_lab.iot_service_blockchain.sdk.Sdk;
import org.nexus_lab.iot_service_blockchain.sdk.SdkOptions;
import org.nexus_lab.iot_service_blockchain.sdk.Service;
import org.nexus_lab.iot_service_blockchain.sdk.ServiceEvent;
import org.nexus_lab.iot_service_blockchain.sdk.ServiceRequest;
import org.nexus_lab.iot_service_blockchain.sdk.ServiceRequestEvent;
import org.nexus_lab.iot_service_blockchain.sdk.ServiceResponse;

public class EventHub {
  private static final String ORG_ID = "Org2MSP";
  private static final String ORG_DOMAIN = "org2.example.com";
  private static final String USER_NAME = "User1@org2.example.com";
  private static final String PEER_NAME = "peer0.org2.example.com";
  private static final String PEER_ENDPOINT = "localhost:9051";

  private static final Logger LOGGER = Logger.getLogger("eventhub");

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

    ScheduledFuture<?> timeout =
        Executors.newSingleThreadScheduledExecutor()
            .schedule(
                () -> Utils.fatal(LOGGER, "timed out waiting for responses"),
                120,
                TimeUnit.SECONDS);

    ThreadPoolExecutor executor = (ThreadPoolExecutor) Executors.newFixedThreadPool(3);

    List<Callable<Object>> tasks = new ArrayList<>();
    tasks.add(Executors.callable(new DeviceEventHandler(isb)));
    tasks.add(Executors.callable(new ServiceEventHandler(isb)));
    tasks.add(Executors.callable(new RequestEventHandler(isb)));
    executor.invokeAll(tasks);

    executor.shutdown();
    timeout.cancel(true);

    isb.close();
    System.exit(0);
  }

  private static class DeviceEventHandler implements Runnable {
    private static final Logger LOGGER = Logger.getLogger("device_event_handler");

    private Sdk isb;

    DeviceEventHandler(Sdk isb) {
      this.isb = isb;
    }

    @Override
    public void run() {
      Utils.log(LOGGER, "Watching for device events");

      CloseableIterator<DeviceEvent> events = isb.getDeviceRegistry().registerEvent();
      Map<String, Integer> expected = new HashMap<>();
      expected.put("register", 1);
      expected.put("deregister", 1);
      Map<String, Integer> actual = new HashMap<>();

      events.forEachRemaining(
          (event) -> {
            actual.put(event.getAction(), actual.getOrDefault(event.getAction(), 0) + 1);

            if ("register".equals(event.getAction()) || "deregister".equals(event.getAction())) {
              Device device = (Device) event.getPayload();
              if (!event.getDeviceId().equals(device.getId())
                  || !event.getOrganizationId().equals(device.getOrganizationId())) {
                Utils.fatal(LOGGER, "device ID or organization ID mismatch");
              }
            }

            for (Map.Entry<String, Integer> entry : expected.entrySet()) {
              if (actual.get(entry.getKey()) != entry.getValue()) {
                return;
              }
            }

            events.close();
          });

      for (Map.Entry<String, Integer> entry : expected.entrySet()) {
        if (actual.get(entry.getKey()) != entry.getValue()) {
          Utils.fatal(
              LOGGER,
              "should have received " + entry.getValue() + " device " + entry.getKey() + " events");
        }
      }

      Utils.log(LOGGER, "Done watching for device events");
    }
  }

  private static class ServiceEventHandler implements Runnable {
    private static final Logger LOGGER = Logger.getLogger("service_event_handler");

    private Sdk isb;

    ServiceEventHandler(Sdk isb) {
      this.isb = isb;
    }

    @Override
    public void run() {
      Utils.log(LOGGER, "Watching for service events");

      CloseableIterator<ServiceEvent> events = isb.getServiceRegistry().registerEvent();
      Map<String, Integer> expected = new HashMap<>();
      expected.put("register", 2);
      expected.put("deregister", 1);
      Map<String, Integer> actual = new HashMap<>();

      events.forEachRemaining(
          (event) -> {
            actual.put(event.getAction(), actual.getOrDefault(event.getAction(), 0) + 1);

            if ("register".equals(event.getAction()) || "deregister".equals(event.getAction())) {
              Service service = (Service) event.getPayload();
              if (!event.getServiceName().equals(service.getName())
                  || !event.getDeviceId().equals(service.getDeviceId())
                  || !event.getOrganizationId().equals(service.getOrganizationId())) {
                Utils.fatal(
                    LOGGER,
                    "event and payload device ID, organization ID, or service name mismatch");
              }
            }

            for (Map.Entry<String, Integer> entry : expected.entrySet()) {
              if (actual.get(entry.getKey()) != entry.getValue()) {
                return;
              }
            }

            events.close();
          });

      for (Map.Entry<String, Integer> entry : expected.entrySet()) {
        if (actual.get(entry.getKey()) != entry.getValue()) {
          Utils.fatal(
              LOGGER,
              "should have received "
                  + entry.getValue()
                  + " service "
                  + entry.getKey()
                  + " events");
        }
      }

      Utils.log(LOGGER, "Done watching for service events");
    }
  }

  private static class RequestEventHandler implements Runnable {
    private static final Logger LOGGER = Logger.getLogger("request_event_handler");

    private Sdk isb;

    RequestEventHandler(Sdk isb) {
      this.isb = isb;
    }

    @Override
    public void run() {
      Utils.log(LOGGER, "Watching for service request events");

      CloseableIterator<ServiceRequestEvent> events = isb.getServiceBroker().registerEvent();
      Map<String, Integer> expected = new HashMap<>();
      expected.put("request", 2);
      expected.put("respond", 2);
      expected.put("remove", 1);
      Map<String, Integer> actual = new HashMap<>();

      events.forEachRemaining(
          (event) -> {
            actual.put(event.getAction(), actual.getOrDefault(event.getAction(), 0) + 1);

            switch (event.getAction()) {
              case "request":
                ServiceRequest request = (ServiceRequest) event.getPayload();
                if (!event.getRequestId().equals(request.getId())) {
                  Utils.fatal(LOGGER, "event and payload request ID mismatch");
                }
                break;
              case "respond":
                ServiceResponse response = (ServiceResponse) event.getPayload();
                if (!event.getRequestId().equals(response.getRequestId())) {
                  Utils.fatal(LOGGER, "event and payload request ID mismatch");
                }
                break;
              case "remove":
                if (!event.getRequestId().equals(event.getPayload())) {
                  Utils.fatal(LOGGER, "event and payload request ID mismatch");
                }
                break;
            }

            for (Map.Entry<String, Integer> entry : expected.entrySet()) {
              if (actual.get(entry.getKey()) != entry.getValue()) {
                return;
              }
            }

            events.close();
          });

      for (Map.Entry<String, Integer> entry : expected.entrySet()) {
        if (actual.get(entry.getKey()) != entry.getValue()) {
          Utils.fatal(
              LOGGER,
              "should have received "
                  + entry.getValue()
                  + " service request "
                  + entry.getKey()
                  + " events");
        }
      }

      Utils.log(LOGGER, "Done watching for service request events");
    }
  }
}
