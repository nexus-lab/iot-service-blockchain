package org.nexus_lab.iot_service_blockchain.sdk;

import static org.junit.Assert.assertArrayEquals;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertThrows;
import static org.junit.Assert.assertTrue;

import com.owlike.genson.JsonBindingException;
import java.time.OffsetDateTime;
import org.junit.Test;

public class ServiceTest {
  @Test
  public void testGetKeyComponents() {
    Service service = new Service();
    service.setName("service1");
    service.setDeviceId("device1");
    service.setOrganizationId("org1");
    service.setVersion(1);
    service.setDescription("Service of Device1");
    service.setLastUpdateTime(OffsetDateTime.parse("2021-12-12T17:34:00-05:00"));

    assertArrayEquals(
        service.getKeyComponents(),
        new String[] {service.getOrganizationId(), service.getDeviceId(), service.getName()});
  }

  @Test
  public void testSerialize() {
    Service service = new Service();
    service.setName("service1");
    service.setDeviceId("device1");
    service.setOrganizationId("org1");
    service.setVersion(1);
    service.setDescription("Service of Device1");
    service.setLastUpdateTime(OffsetDateTime.parse("2021-12-12T17:34:00-05:00"));

    String serialized =
        "{\"name\":\"service1\",\"deviceId\":\"device1\","
            + "\"organizationId\":\"org1\",\"version\":1,\"description\":\"Service of Device1\","
            + "\"lastUpdateTime\":\"2021-12-12T17:34:00.000-05:00\"}";

    assertEquals(serialized, service.serialize());
  }

  @Test
  public void testValidate() {
    Service service = new Service();

    IllegalArgumentException exception =
        assertThrows(IllegalArgumentException.class, () -> service.validate());
    assertTrue(exception.getMessage().contains("service name"));
    service.setName("service1");

    exception = assertThrows(IllegalArgumentException.class, () -> service.validate());
    assertTrue(exception.getMessage().contains("device ID"));
    service.setDeviceId("device1");

    exception = assertThrows(IllegalArgumentException.class, () -> service.validate());
    assertTrue(exception.getMessage().contains("organization ID"));
    service.setOrganizationId("org1");

    exception = assertThrows(IllegalArgumentException.class, () -> service.validate());
    assertTrue(exception.getMessage().contains("service version"));
    service.setVersion(-1);

    exception = assertThrows(IllegalArgumentException.class, () -> service.validate());
    assertTrue(exception.getMessage().contains("positive integer"));
    service.setVersion(1);

    exception = assertThrows(IllegalArgumentException.class, () -> service.validate());
    assertTrue(exception.getMessage().contains("last update time"));
    service.setLastUpdateTime(OffsetDateTime.parse("2021-12-12T17:34:00-05:00"));

    service.validate();
  }

  @Test
  public void testDeserialize() {
    Service expected = new Service();
    expected.setName("service1");
    expected.setDeviceId("device1");
    expected.setOrganizationId("org1");
    expected.setVersion(1);
    expected.setDescription("Service of Device1");
    expected.setLastUpdateTime(OffsetDateTime.parse("2021-12-12T17:34:00-05:00"));

    String serialized =
        "{\"name\":\"service1\",\"deviceId\":\"device1\","
            + "\"organizationId\":\"org1\",\"version\":1,\"description\":\"Service of Device1\","
            + "\"lastUpdateTime\":\"2021-12-12T17:34:00.000-05:00\"}";

    Service actual = Service.deserialize(serialized);
    assertEquals(expected, actual);

    assertThrows(JsonBindingException.class, () -> Service.deserialize("\u0000"));
  }
}
