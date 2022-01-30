package org.nexus_lab.iot_service_blockchain.sdk;

import com.owlike.genson.Context;
import com.owlike.genson.annotation.JsonIgnore;
import com.owlike.genson.stream.ObjectWriter;
import java.time.OffsetDateTime;
import java.util.UUID;
import java.util.regex.Pattern;
import lombok.Data;
import lombok.NoArgsConstructor;

/** An IoT service response */
@Data
@NoArgsConstructor
public class ServiceResponse {
  /** Identity of the IoT service request to respond to */
  private String requestId;

  /** Time of the IoT service response */
  private OffsetDateTime time;

  /** Status code of the IoT service response */
  private int statusCode;

  /** Return value of the IoT service response */
  private String returnValue;

  /**
   * Get components that compose the IoT service response key
   *
   * @return components that compose the IoT service response key
   */
  @JsonIgnore
  public String[] getKeyComponents() {
    return new String[] {this.requestId};
  }

  /**
   * Transform current service response to JSON string
   *
   * @return JSON representation of the service response
   */
  public String serialize() {
    return Json.serialize(this);
  }

  /**
   * Check if the IoT service response properties are valid
   *
   * @throws IllegalArgumentException when service response fields are invalid
   */
  public void validate() {
    Pattern uuidPattern =
        Pattern.compile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$");
    if (this.requestId == null || !uuidPattern.matcher(this.requestId).matches()) {
      throw new IllegalArgumentException("invalid request ID in response definition");
    }
    try {
      UUID.fromString(this.requestId);
    } catch (IllegalArgumentException exception) {
      throw new IllegalArgumentException("invalid request ID in request definition");
    }
    if (this.time == null || (this.time.toEpochSecond() == 0 && this.time.getNano() == 0)) {
      throw new IllegalArgumentException("missing response time in response definition");
    }
  }

  /**
   * Create a new service response instance from its JSON representation
   *
   * @param data JSON string representing a service response
   * @return a new service response instance
   */
  public static ServiceResponse deserialize(String data) {
    return Json.deserialize(data, ServiceResponse.class);
  }

  /** Custom JSON serializer that keeps the field order for {@link com.owlike.genson.Genson} */
  protected static final class Serializer implements com.owlike.genson.Serializer<ServiceResponse> {
    @Override
    public void serialize(ServiceResponse response, ObjectWriter writer, Context ctx)
        throws Exception {
      writer.beginObject();
      writer.writeString("requestId", response.requestId);
      writer.writeName("time");
      ctx.genson.serialize(response.time, writer, ctx);
      writer.writeNumber("statusCode", response.statusCode);
      writer.writeString("returnValue", response.returnValue);
      writer.endObject();
    }
  }
}
