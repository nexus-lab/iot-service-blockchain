package org.nexus_lab.iot_service_blockchain.sdk;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertThrows;

import com.owlike.genson.JsonBindingException;
import java.time.OffsetDateTime;
import org.junit.Test;

public class ServiceRequestResponseTest {
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

    ServiceResponse response = new ServiceResponse();
    response.setRequestId("ffbc9005-c62a-4563-a8f7-b32bba27d707");
    response.setTime(OffsetDateTime.parse("2021-12-12T17:34:00-05:00"));
    response.setStatusCode(0);
    response.setReturnValue("[\"a\",\"b\",\"c\"]");

    ServiceRequestResponse pair = new ServiceRequestResponse();
    pair.setRequest(request);
    pair.setResponse(response);

    String serialized =
        "{\"request\":{\"id\":\"ffbc9005-c62a-4563-a8f7-b32bba27d707\","
            + "\"time\":\"2021-12-12T17:34:00.000-05:00\","
            + "\"service\":{\"name\":\"service1\",\"deviceId\":\"device1\","
            + "\"organizationId\":\"org1\",\"version\":0,\"description\":null,"
            + "\"lastUpdateTime\":null},\"method\":\"GET\",\"arguments\":[\"1\",\"2\",\"3\"]},"
            + "\"response\":{\"requestId\":\"ffbc9005-c62a-4563-a8f7-b32bba27d707\","
            + "\"time\":\"2021-12-12T17:34:00.000-05:00\",\"statusCode\":0,"
            + "\"returnValue\":\"[\\\"a\\\",\\\"b\\\",\\\"c\\\"]\"}}";

    assertEquals(serialized, pair.serialize());
  }

  @Test
  public void testDeserialize() {
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

    ServiceResponse response = new ServiceResponse();
    response.setRequestId("ffbc9005-c62a-4563-a8f7-b32bba27d707");
    response.setTime(OffsetDateTime.parse("2021-12-12T17:34:00-05:00"));
    response.setStatusCode(0);
    response.setReturnValue("[\"a\",\"b\",\"c\"]");

    ServiceRequestResponse expected = new ServiceRequestResponse();
    expected.setRequest(request);
    expected.setResponse(response);

    String serialized =
        "{\"request\":{\"id\":\"ffbc9005-c62a-4563-a8f7-b32bba27d707\","
            + "\"time\":\"2021-12-12T17:34:00.000-05:00\","
            + "\"service\":{\"name\":\"service1\",\"deviceId\":\"device1\","
            + "\"organizationId\":\"org1\",\"version\":0,\"description\":null,"
            + "\"lastUpdateTime\":null},\"method\":\"GET\",\"arguments\":[\"1\",\"2\",\"3\"]},"
            + "\"response\":{\"requestId\":\"ffbc9005-c62a-4563-a8f7-b32bba27d707\","
            + "\"time\":\"2021-12-12T17:34:00.000-05:00\",\"statusCode\":0,"
            + "\"returnValue\":\"[\\\"a\\\",\\\"b\\\",\\\"c\\\"]\"}}";

    ServiceRequestResponse actual = ServiceRequestResponse.deserialize(serialized);
    assertEquals(expected, actual);

    assertThrows(JsonBindingException.class, () -> ServiceRequestResponse.deserialize("\u0000"));
  }
}
