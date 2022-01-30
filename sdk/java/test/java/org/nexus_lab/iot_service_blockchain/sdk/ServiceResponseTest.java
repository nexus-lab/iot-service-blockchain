package org.nexus_lab.iot_service_blockchain.sdk;

import static org.junit.Assert.assertArrayEquals;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertThrows;
import static org.junit.Assert.assertTrue;

import com.owlike.genson.JsonBindingException;
import java.time.OffsetDateTime;
import org.junit.Test;

public class ServiceResponseTest {
  @Test
  public void testGetKeyComponents() {
    ServiceResponse response = new ServiceResponse();
    response.setRequestId("ffbc9005-c62a-4563-a8f7-b32bba27d707");
    response.setTime(OffsetDateTime.parse("2021-12-12T17:34:00-05:00"));
    response.setStatusCode(0);
    response.setReturnValue("[\"a\",\"b\",\"c\"]");

    assertArrayEquals(response.getKeyComponents(), new String[] {response.getRequestId()});
  }

  @Test
  public void testSerialize() {
    ServiceResponse response = new ServiceResponse();
    response.setRequestId("ffbc9005-c62a-4563-a8f7-b32bba27d707");
    response.setTime(OffsetDateTime.parse("2021-12-12T17:34:00-05:00"));
    response.setStatusCode(0);
    response.setReturnValue("[\"a\",\"b\",\"c\"]");

    String serialized =
        "{\"requestId\":\"ffbc9005-c62a-4563-a8f7-b32bba27d707\","
            + "\"time\":\"2021-12-12T17:34:00.000-05:00\",\"statusCode\":0,"
            + "\"returnValue\":\"[\\\"a\\\",\\\"b\\\",\\\"c\\\"]\"}";

    assertEquals(serialized, response.serialize());
  }

  @Test
  public void testValidate() {
    ServiceResponse response = new ServiceResponse();
    response.setRequestId("123456");

    IllegalArgumentException exception =
        assertThrows(IllegalArgumentException.class, () -> response.validate());
    assertTrue(exception.getMessage().contains("request ID"));
    response.setRequestId("ffbc9005-c62a-4563-a8f7-b32bba27d707");

    exception = assertThrows(IllegalArgumentException.class, () -> response.validate());
    assertTrue(exception.getMessage().contains("response time"));
    response.setTime(OffsetDateTime.parse("2021-12-12T17:34:00-05:00"));

    response.validate();
  }

  @Test
  public void testDeserialize() {
    ServiceResponse expected = new ServiceResponse();
    expected.setRequestId("ffbc9005-c62a-4563-a8f7-b32bba27d707");
    expected.setTime(OffsetDateTime.parse("2021-12-12T17:34:00-05:00"));
    expected.setStatusCode(0);
    expected.setReturnValue("[\"a\",\"b\",\"c\"]");

    String serialized =
        "{\"requestId\":\"ffbc9005-c62a-4563-a8f7-b32bba27d707\","
            + "\"time\":\"2021-12-12T17:34:00.000-05:00\",\"statusCode\":0,"
            + "\"returnValue\":\"[\\\"a\\\",\\\"b\\\",\\\"c\\\"]\"}";

    ServiceResponse actual = ServiceResponse.deserialize(serialized);
    assertEquals(expected, actual);

    assertThrows(JsonBindingException.class, () -> ServiceResponse.deserialize("\u0000"));
  }
}
