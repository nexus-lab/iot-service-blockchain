package org.nexus_lab.iot_service_blockchain.sdk;

import com.owlike.genson.GenericType;
import com.owlike.genson.Genson;
import com.owlike.genson.GensonBuilder;

/** The default JSON serializer and deserializer. */
public class Json {
  private static final Genson genson =
      new GensonBuilder()
          .withSerializers(
              new Device.Serializer(),
              new Service.Serializer(),
              new ServiceRequest.Serializer(),
              new ServiceResponse.Serializer(),
              new ServiceRequestResponse.Serializer())
          .withConverters(new OffsetDateTimeConverter())
          .create();

  /**
   * Serializes the object into a json string.
   *
   * @see com.owlike.genson.Genson#serialize(Object)
   */
  public static String serialize(Object object) {
    return genson.serialize(object);
  }

  /**
   * Deserializes fromSource String into an instance of toClass.
   *
   * @see com.owlike.genson.Genson#deserialize(String, Class)
   */
  public static <T> T deserialize(String json, Class<T> type) {
    return genson.deserialize(json, type);
  }

  /**
   * Deserializes fromSource String into an instance of toClass.
   *
   * @see com.owlike.genson.Genson#deserialize(String, GenericType)
   */
  public static <T> T deserialize(String json, GenericType<T> type) {
    return genson.deserialize(json, type);
  }
}
