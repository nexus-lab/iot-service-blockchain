package org.nexus_lab.iot_service_blockchain.sdk;

import com.owlike.genson.Context;
import com.owlike.genson.stream.ObjectWriter;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
public class ServiceRequestResponse {
  /** IoT service request */
  private ServiceRequest request;

  /** IoT service response */
  private ServiceResponse response;

  /**
   * Transform current service request/response pair to JSON string
   *
   * @return JSON representation of the service request/response pair
   */
  public String serialize() {
    return Json.serialize(this);
  }

  /**
   * Create a new service request/response pair instance from its JSON representation
   *
   * @param data JSON string representing a service request/response pair
   * @return a new service request/response pair instance
   */
  public static ServiceRequestResponse deserialize(String data) {
    return Json.deserialize(data, ServiceRequestResponse.class);
  }

  /** Custom JSON serializer that keeps the field order for {@link com.owlike.genson.Genson} */
  protected static final class Serializer
      implements com.owlike.genson.Serializer<ServiceRequestResponse> {
    @Override
    public void serialize(ServiceRequestResponse pair, ObjectWriter writer, Context ctx)
        throws Exception {
      writer.beginObject();
      writer.writeName("request");
      ctx.genson.serialize(pair.request, writer, ctx);
      writer.writeName("response");
      ctx.genson.serialize(pair.response, writer, ctx);
      writer.endObject();
    }
  }
}
