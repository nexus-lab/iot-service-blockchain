package org.nexus_lab.iot_service_blockchain.sdk;

import static org.junit.Assert.assertArrayEquals;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertThrows;
import static org.junit.Assert.assertTrue;

import com.owlike.genson.JsonBindingException;
import java.time.OffsetDateTime;
import org.junit.Test;

public class ServiceRequestTest {
  @Test
  public void testGetKeyComponents() {
    Service service = new Service();
    service.setName("service1");
    service.setDeviceId("device1");
    service.setOrganizationId("org1");

    ServiceRequest request = new ServiceRequest();
    request.setId("ffbc9005-c62a-4563-a8f7-b32bba27d707");
    request.setTime(OffsetDateTime.parse("2021-12-12T17:34:00-05:00"));
    request.setService(service);
    request.setMethod("GET");
    request.setArguments(new String[] {"1", "2", "3"});

    assertArrayEquals(request.getKeyComponents(), new String[] {request.getId()});
  }

  @Test
  public void testSerialize() {
    Service service = new Service();
    service.setName("service1");
    service.setDeviceId("device1");
    service.setOrganizationId("org1");

    ServiceRequest request = new ServiceRequest();
    request.setId("ffbc9005-c62a-4563-a8f7-b32bba27d707");
    request.setTime(OffsetDateTime.parse("2021-12-12T17:34:00-05:00"));
    request.setService(service);
    request.setMethod("GET");
    request.setArguments(new String[] {"1", "2", "3"});

    String serialized =
        "{\"id\":\"ffbc9005-c62a-4563-a8f7-b32bba27d707\","
            + "\"time\":\"2021-12-12T17:34:00.000-05:00\","
            + "\"service\":{\"name\":\"service1\",\"deviceId\":\"device1\","
            + "\"organizationId\":\"org1\",\"version\":0,\"description\":null,"
            + "\"lastUpdateTime\":null},\"method\":\"GET\",\"arguments\":[\"1\",\"2\",\"3\"]}";

    assertEquals(serialized, request.serialize());
  }

  @Test
  public void testValidate() {
    Service service = new Service();
    ServiceRequest request = new ServiceRequest();
    request.setId("123456");
    request.setService(service);

    IllegalArgumentException exception =
        assertThrows(IllegalArgumentException.class, () -> request.validate());
    assertTrue(exception.getMessage().contains("request ID"));
    request.setId("ffbc9005-c62a-4563-a8f7-b32bba27d707");

    exception = assertThrows(IllegalArgumentException.class, () -> request.validate());
    assertTrue(exception.getMessage().contains("requested service"));
    service.setOrganizationId("org1");

    exception = assertThrows(IllegalArgumentException.class, () -> request.validate());
    assertTrue(exception.getMessage().contains("requested service"));
    service.setDeviceId("device1");

    exception = assertThrows(IllegalArgumentException.class, () -> request.validate());
    assertTrue(exception.getMessage().contains("requested service"));
    service.setName("service1");

    exception = assertThrows(IllegalArgumentException.class, () -> request.validate());
    assertTrue(exception.getMessage().contains("request method"));
    request.setMethod("GET");

    exception = assertThrows(IllegalArgumentException.class, () -> request.validate());
    assertTrue(exception.getMessage().contains("request arguments"));
    request.setArguments(new String[0]);

    exception = assertThrows(IllegalArgumentException.class, () -> request.validate());
    assertTrue(exception.getMessage().contains("request time"));
    request.setTime(OffsetDateTime.parse("2021-12-12T17:34:00-05:00"));

    request.validate();
  }

  @Test
  public void testDeserialize() {
    Service service = new Service();
    service.setName("service1");
    service.setDeviceId("device1");
    service.setOrganizationId("org1");

    ServiceRequest expected = new ServiceRequest();
    expected.setId("ffbc9005-c62a-4563-a8f7-b32bba27d707");
    expected.setTime(OffsetDateTime.parse("2021-12-12T17:34:00-05:00"));
    expected.setService(service);
    expected.setMethod("GET");
    expected.setArguments(new String[] {"1", "2", "3"});

    String serialized =
        "{\"id\":\"ffbc9005-c62a-4563-a8f7-b32bba27d707\","
            + "\"time\":\"2021-12-12T17:34:00.000-05:00\","
            + "\"service\":{\"name\":\"service1\",\"deviceId\":\"device1\","
            + "\"organizationId\":\"org1\",\"version\":0,\"description\":null,"
            + "\"lastUpdateTime\":null},\"method\":\"GET\",\"arguments\":[\"1\",\"2\",\"3\"]}";

    ServiceRequest actual = ServiceRequest.deserialize(serialized);
    assertEquals(expected, actual);

    assertThrows(JsonBindingException.class, () -> ServiceRequest.deserialize("\u0000"));
  }
}
