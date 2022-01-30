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
public class ServiceRegistryTest {
  @Mock Contract contract;

  @Test
  public void testRegister() throws Exception {
    ServiceRegistry registry = new ServiceRegistry(contract);

    Service service1 = new Service();
    service1.setName("service1");
    registry.register(service1);
    verify(contract).submitTransaction("Register", service1.serialize());

    assertThrows(NullPointerException.class, () -> registry.register(null));

    final Service service2 = new Service();
    service2.setName("service2");
    when(contract.submitTransaction("Register", service2.serialize()))
        .thenThrow(new RuntimeException());
    assertThrows(RuntimeException.class, () -> registry.register(service2));
  }

  @Test
  public void testGet() throws Exception {
    ServiceRegistry registry = new ServiceRegistry(contract);

    Service expected = new Service();
    when(contract.submitTransaction("Get", "org1", "device1", "service1"))
        .thenReturn(expected.serialize().getBytes());

    Service actual = registry.get("org1", "device1", "service1");
    assertEquals(expected, actual);

    when(contract.submitTransaction("Get", "org2", "device2", "service2"))
        .thenThrow(new RuntimeException());

    assertThrows(RuntimeException.class, () -> registry.get("org2", "device2", "service2"));
  }

  @Test
  public void testGetAll() throws Exception {
    ServiceRegistry registry = new ServiceRegistry(contract);

    Service[] expected = new Service[] {new Service(), new Service()};
    when(contract.submitTransaction("GetAll", "org1", "device1"))
        .thenReturn(Json.serialize(expected).getBytes());

    List<Service> actual = registry.getAll("org1", "device1");
    for (int i = 0; i < 2; i++) {
      assertEquals(expected[i], actual.get(i));
    }

    when(contract.submitTransaction("GetAll", "org2", "device2")).thenReturn("[]".getBytes());

    actual = registry.getAll("org2", "device2");
    assertEquals(0, actual.size());

    when(contract.submitTransaction("GetAll", "org3", "device3")).thenThrow(new RuntimeException());

    assertThrows(RuntimeException.class, () -> registry.getAll("org3", "device3"));
  }

  @Test
  public void testDeregister() throws Exception {
    ServiceRegistry registry = new ServiceRegistry(contract);

    Service service1 = new Service();
    service1.setName("device1");
    registry.deregister(service1);
    verify(contract).submitTransaction("Deregister", service1.serialize());

    assertThrows(NullPointerException.class, () -> registry.deregister(null));

    final Service service2 = new Service();
    service2.setName("service2");
    when(contract.submitTransaction("Deregister", service2.serialize()))
        .thenThrow(new RuntimeException());
    assertThrows(RuntimeException.class, () -> registry.deregister(service2));
  }

  @Test
  public void testRegisterEvent() {
    ServiceRegistry registry = new ServiceRegistry(contract);

    when(contract.registerEvent())
        .thenReturn(
            Utils.createIterator(
                5,
                (i) -> {
                  Service service = new Service();
                  service.setName("service" + i);
                  service.setOrganizationId("org" + i);
                  service.setDeviceId("device" + i);

                  return new ChaincodeEvent(
                      String.format("service://org%d/device%d/service%d/register", i, i, i),
                      service.serialize().getBytes());
                }));

    final AtomicInteger counter = new AtomicInteger(0);
    try (CloseableIterator<ServiceEvent> events = registry.registerEvent()) {
      events.forEachRemaining(
          event -> {
            int i = counter.getAndIncrement();

            assertEquals("register", event.getAction());
            assertEquals("org" + i, event.getOrganizationId());
            assertEquals("device" + i, event.getDeviceId());
            assertEquals("service" + i, event.getServiceName());
            assertEquals("org" + i, ((Service) event.getPayload()).getOrganizationId());
            assertEquals("device" + i, ((Service) event.getPayload()).getDeviceId());
            assertEquals("service" + i, ((Service) event.getPayload()).getName());
          });
    }
    assertEquals(5, counter.get());

    when(contract.registerEvent()).thenThrow(new RuntimeException());
    assertThrows(RuntimeException.class, () -> registry.registerEvent());
  }
}
