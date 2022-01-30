package org.nexus_lab.iot_service_blockchain.sdk;

import static org.junit.Assert.assertArrayEquals;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertThrows;
import static org.junit.Assert.assertTrue;

import com.owlike.genson.JsonBindingException;
import java.time.OffsetDateTime;
import org.junit.Test;

public class DeviceTest {
  @Test
  public void testGetKeyComponents() {
    Device device = new Device();
    device.setId("device1");
    device.setOrganizationId("org1");
    device.setName("device1");
    device.setDescription("Device of Org1 User1");
    device.setLastUpdateTime(OffsetDateTime.parse("2021-12-12T17:34:00-05:00"));

    assertArrayEquals(
        device.getKeyComponents(), new String[] {device.getOrganizationId(), device.getId()});
  }

  @Test
  public void testSerialize() {
    Device device = new Device();
    device.setId("device1");
    device.setOrganizationId("org1");
    device.setName("device1");
    device.setDescription("Device of Org1 User1");
    device.setLastUpdateTime(OffsetDateTime.parse("2021-12-12T17:34:00-05:00"));

    String serialized =
        "{\"id\":\"device1\",\"organizationId\":\"org1\",\"name\":\"device1\","
            + "\"description\":\"Device of Org1 User1\","
            + "\"lastUpdateTime\":\"2021-12-12T17:34:00.000-05:00\"}";

    assertEquals(serialized, device.serialize());
  }

  @Test
  public void testValidate() {
    Device device = new Device();

    IllegalArgumentException exception =
        assertThrows(IllegalArgumentException.class, () -> device.validate());
    assertTrue(exception.getMessage().contains("device ID"));
    device.setId("device1");

    exception = assertThrows(IllegalArgumentException.class, () -> device.validate());
    assertTrue(exception.getMessage().contains("organization ID"));
    device.setOrganizationId("org1");

    exception = assertThrows(IllegalArgumentException.class, () -> device.validate());
    assertTrue(exception.getMessage().contains("device name"));
    device.setName("device1");

    exception = assertThrows(IllegalArgumentException.class, () -> device.validate());
    assertTrue(exception.getMessage().contains("last update time"));
    device.setLastUpdateTime(OffsetDateTime.parse("2021-12-12T17:34:00-05:00"));

    device.validate();
  }

  @Test
  public void testDeserialize() {
    Device expected = new Device();
    expected.setId("device1");
    expected.setOrganizationId("org1");
    expected.setName("device1");
    expected.setDescription("Device of Org1 User1");
    expected.setLastUpdateTime(OffsetDateTime.parse("2021-12-12T17:34:00-05:00"));

    String serialized =
        "{\"id\":\"device1\",\"organizationId\":\"org1\",\"name\":\"device1\","
            + "\"description\":\"Device of Org1 User1\","
            + "\"lastUpdateTime\":\"2021-12-12T17:34:00.000-05:00\"}";

    Device actual = Device.deserialize(serialized);
    assertEquals(expected, actual);

    assertThrows(JsonBindingException.class, () -> Device.deserialize("\u0000"));
  }
}
