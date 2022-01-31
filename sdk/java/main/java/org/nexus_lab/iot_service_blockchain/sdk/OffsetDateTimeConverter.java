package org.nexus_lab.iot_service_blockchain.sdk;

import com.owlike.genson.Context;
import com.owlike.genson.Converter;
import com.owlike.genson.stream.ObjectReader;
import com.owlike.genson.stream.ObjectWriter;
import java.time.OffsetDateTime;
import java.time.format.DateTimeFormatter;

/**
 * Custom {@link com.owlike.genson.Genson} serializer and deserializer for {@link
 * java.time.OffsetDateTime}.
 */
public class OffsetDateTimeConverter implements Converter<OffsetDateTime> {
  private static final DateTimeFormatter formatter =
      DateTimeFormatter.ofPattern("yyyy-MM-dd'T'HH:mm:ss.SSSXXXXX");

  @Override
  public void serialize(OffsetDateTime date, ObjectWriter writer, Context ctx) throws Exception {
    writer.writeString(date.format(formatter));
  }

  @Override
  public OffsetDateTime deserialize(ObjectReader reader, Context ctx) throws Exception {
    return OffsetDateTime.parse(reader.valueAsString());
  }
}
