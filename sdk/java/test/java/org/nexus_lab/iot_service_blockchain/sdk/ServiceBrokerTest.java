package org.nexus_lab.iot_service_blockchain.sdk;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertThrows;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import java.util.List;
import java.util.concurrent.atomic.AtomicInteger;
import org.hyperledger.fabric.client.CloseableIterator;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.Mock;
import org.mockito.junit.MockitoJUnitRunner;

@RunWith(MockitoJUnitRunner.class)
public class ServiceBrokerTest {
  @Mock Contract contract;

  @Test
  public void testRequest() throws Exception {
    ServiceBroker broker = new ServiceBroker(contract);

    ServiceRequest request1 = new ServiceRequest();
    request1.setId("request1");
    broker.request(request1);
    verify(contract).submitTransaction("Request", request1.serialize());

    assertThrows(NullPointerException.class, () -> broker.request(null));

    final ServiceRequest request2 = new ServiceRequest();
    request2.setId("request2");
    when(contract.submitTransaction("Request", request2.serialize()))
        .thenThrow(new RuntimeException());
    assertThrows(RuntimeException.class, () -> broker.request(request2));
  }

  @Test
  public void testRespond() throws Exception {
    ServiceBroker broker = new ServiceBroker(contract);

    ServiceResponse response1 = new ServiceResponse();
    response1.setRequestId("request1");
    broker.respond(response1);
    verify(contract).submitTransaction("Respond", response1.serialize());

    assertThrows(NullPointerException.class, () -> broker.request(null));

    final ServiceResponse response2 = new ServiceResponse();
    response2.setRequestId("request2");
    when(contract.submitTransaction("Respond", response2.serialize()))
        .thenThrow(new RuntimeException());
    assertThrows(RuntimeException.class, () -> broker.respond(response2));
  }

  @Test
  public void testGet() throws Exception {
    ServiceBroker broker = new ServiceBroker(contract);

    ServiceRequestResponse expected = new ServiceRequestResponse();
    when(contract.submitTransaction("Get", "request1")).thenReturn(expected.serialize().getBytes());

    ServiceRequestResponse actual = broker.get("request1");
    assertEquals(expected, actual);

    when(contract.submitTransaction("Get", "request2")).thenThrow(new RuntimeException());

    assertThrows(RuntimeException.class, () -> broker.get("request2"));
  }

  @Test
  public void testGetAll() throws Exception {
    ServiceBroker broker = new ServiceBroker(contract);

    ServiceRequestResponse[] expected =
        new ServiceRequestResponse[] {new ServiceRequestResponse(), new ServiceRequestResponse()};
    when(contract.submitTransaction("GetAll", "org1", "device1", "service1"))
        .thenReturn(Json.serialize(expected).getBytes());

    List<ServiceRequestResponse> actual = broker.getAll("org1", "device1", "service1");
    for (int i = 0; i < 2; i++) {
      assertEquals(expected[i], actual.get(i));
    }

    when(contract.submitTransaction("GetAll", "org2", "device2", "service2"))
        .thenReturn("[]".getBytes());

    actual = broker.getAll("org2", "device2", "service2");
    assertEquals(0, actual.size());

    when(contract.submitTransaction("GetAll", "org3", "device3", "service3"))
        .thenThrow(new RuntimeException());

    assertThrows(RuntimeException.class, () -> broker.getAll("org3", "device3", "service3"));
  }

  @Test
  public void testRemove() throws Exception {
    ServiceBroker broker = new ServiceBroker(contract);

    broker.remove("request1");
    verify(contract).submitTransaction("Remove", "request1");

    when(contract.submitTransaction("Remove", "request2")).thenThrow(new RuntimeException());
    assertThrows(RuntimeException.class, () -> broker.remove("request2"));
  }

  @Test
  public void testRegisterEvent() {
    ServiceBroker broker = new ServiceBroker(contract);

    when(contract.registerEvent())
        .thenReturn(
            Utils.createIterator(
                6,
                (i) -> {
                  if (i < 2) {
                    ServiceRequest request = new ServiceRequest();
                    request.setId("request" + i);

                    return new ChaincodeEvent(
                        String.format(
                            "request://org%d/device%d/service%d/request%d/request", i, i, i, i),
                        request.serialize().getBytes());
                  } else if (i < 4) {
                    ServiceResponse response = new ServiceResponse();
                    response.setRequestId("request" + i);

                    return new ChaincodeEvent(
                        String.format(
                            "request://org%d/device%d/service%d/request%d/respond", i, i, i, i),
                        response.serialize().getBytes());
                  } else {
                    return new ChaincodeEvent(
                        String.format(
                            "request://org%d/device%d/service%d/request%d/remove", i, i, i, i),
                        ("request" + i).getBytes());
                  }
                }));

    final AtomicInteger counter = new AtomicInteger(0);
    try (CloseableIterator<ServiceRequestEvent> events = broker.registerEvent()) {
      events.forEachRemaining(
          event -> {
            int i = counter.getAndIncrement();

            assertEquals("org" + i, event.getOrganizationId());
            assertEquals("device" + i, event.getDeviceId());
            assertEquals("service" + i, event.getServiceName());
            assertEquals("request" + i, event.getRequestId());

            if (i < 2) {
              assertEquals("request", event.getAction());
              assertEquals("request" + i, ((ServiceRequest) event.getPayload()).getId());
            } else if (i < 4) {
              assertEquals("respond", event.getAction());
              assertEquals("request" + i, ((ServiceResponse) event.getPayload()).getRequestId());
            } else {
              assertEquals("remove", event.getAction());
              assertEquals("request" + i, (String) event.getPayload());
            }
          });
    }
    assertEquals(6, counter.get());

    when(contract.registerEvent()).thenThrow(new RuntimeException());
    assertThrows(RuntimeException.class, () -> broker.registerEvent());
  }
}
