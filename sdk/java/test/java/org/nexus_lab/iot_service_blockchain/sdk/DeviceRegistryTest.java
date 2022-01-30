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
public class DeviceRegistryTest {
  @Mock Contract contract;

  @Test
  public void testRegister() throws Exception {
    DeviceRegistry registry = new DeviceRegistry(contract);

    Device device1 = new Device();
    device1.setName("device1");
    registry.register(device1);
    verify(contract).submitTransaction("Register", device1.serialize());

    assertThrows(NullPointerException.class, () -> registry.register(null));

    final Device device2 = new Device();
    device2.setName("device2");
    when(contract.submitTransaction("Register", device2.serialize()))
        .thenThrow(new RuntimeException());
    assertThrows(RuntimeException.class, () -> registry.register(device2));
  }

  @Test
  public void testGet() throws Exception {
    DeviceRegistry registry = new DeviceRegistry(contract);

    Device expected = new Device();
    when(contract.submitTransaction("Get", "org1", "device1"))
        .thenReturn(expected.serialize().getBytes());

    Device actual = registry.get("org1", "device1");
    assertEquals(expected, actual);

    when(contract.submitTransaction("Get", "org2", "device2")).thenThrow(new RuntimeException());

    assertThrows(RuntimeException.class, () -> registry.get("org2", "device2"));
  }

  @Test
  public void testGetAll() throws Exception {
    DeviceRegistry registry = new DeviceRegistry(contract);

    Device[] expected = new Device[] {new Device(), new Device()};
    when(contract.submitTransaction("GetAll", "org1"))
        .thenReturn(Json.serialize(expected).getBytes());

    List<Device> actual = registry.getAll("org1");
    for (int i = 0; i < 2; i++) {
      assertEquals(expected[i], actual.get(i));
    }

    when(contract.submitTransaction("GetAll", "org2")).thenReturn("[]".getBytes());

    actual = registry.getAll("org2");
    assertEquals(0, actual.size());

    when(contract.submitTransaction("GetAll", "org3")).thenThrow(new RuntimeException());

    assertThrows(RuntimeException.class, () -> registry.getAll("org3"));
  }

  @Test
  public void testDeregister() throws Exception {
    DeviceRegistry registry = new DeviceRegistry(contract);

    Device device1 = new Device();
    device1.setName("device1");
    registry.deregister(device1);
    verify(contract).submitTransaction("Deregister", device1.serialize());

    assertThrows(NullPointerException.class, () -> registry.deregister(null));

    final Device device2 = new Device();
    device2.setName("device2");
    when(contract.submitTransaction("Deregister", device2.serialize()))
        .thenThrow(new RuntimeException());
    assertThrows(RuntimeException.class, () -> registry.deregister(device2));
  }

  @Test
  public void testRegisterEvent() {
    DeviceRegistry registry = new DeviceRegistry(contract);

    when(contract.registerEvent())
        .thenReturn(
            Utils.createIterator(
                5,
                (i) -> {
                  Device device = new Device();
                  device.setId("device" + i);
                  device.setName("device" + i);
                  device.setOrganizationId("org" + i);

                  return new ChaincodeEvent(
                      String.format("device://org%d/device%d/register", i, i),
                      device.serialize().getBytes());
                }));

    final AtomicInteger counter = new AtomicInteger(0);
    try (CloseableIterator<DeviceEvent> events = registry.registerEvent()) {
      events.forEachRemaining(
          event -> {
            int i = counter.getAndIncrement();

            assertEquals("register", event.getAction());
            assertEquals("org" + i, event.getOrganizationId());
            assertEquals("device" + i, event.getDeviceId());
            assertEquals("device" + i, ((Device) event.getPayload()).getId());
            assertEquals("org" + i, ((Device) event.getPayload()).getOrganizationId());
            assertEquals("device" + i, ((Device) event.getPayload()).getName());
          });
    }
    assertEquals(5, counter.get());

    when(contract.registerEvent()).thenThrow(new RuntimeException());
    assertThrows(RuntimeException.class, () -> registry.registerEvent());
  }
}
