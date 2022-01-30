package org.nexus_lab.iot_service_blockchain.sdk;

import static org.junit.Assert.assertEquals;

import com.owlike.genson.Genson;
import com.owlike.genson.GensonBuilder;
import java.time.OffsetDateTime;
import org.junit.Test;

public class OffsetDateTimeConverterTest {
  private static final Genson genson =
      new GensonBuilder().withConverters(new OffsetDateTimeConverter()).create();

  @Test
  public void testSerialize() {
    assertEquals(
        "\"2021-12-12T17:34:00.000-05:00\"",
        genson.serialize(OffsetDateTime.parse("2021-12-12T17:34:00-05:00")));
    assertEquals(
        "\"2021-12-12T17:34:00.000Z\"",
        genson.serialize(OffsetDateTime.parse("2021-12-12T17:34:00Z")));
    assertEquals("null", genson.serialize(null));
  }

  @Test
  public void testDeserialize() {
    assertEquals(
        OffsetDateTime.parse("2021-12-12T17:34:00-05:00"),
        genson.deserialize("\"2021-12-12T17:34:00-05:00\"", OffsetDateTime.class));
    assertEquals(
        OffsetDateTime.parse("2021-12-12T17:34:00Z"),
        genson.deserialize("\"2021-12-12T17:34:00.000Z\"", OffsetDateTime.class));
    assertEquals(null, genson.deserialize("null", OffsetDateTime.class));
  }
}
