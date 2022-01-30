package org.nexus_lab.iot_service_blockchain.sdk;

import com.owlike.genson.Context;
import com.owlike.genson.annotation.JsonIgnore;
import com.owlike.genson.stream.ObjectWriter;
import java.time.OffsetDateTime;
import java.util.UUID;
import java.util.regex.Pattern;
import lombok.Data;
import lombok.NoArgsConstructor;

/** An IoT service request */
@Data
@NoArgsConstructor
public class ServiceRequest {
  /** Identity of the IoT service request */
  private String id;

  /** Time of the IoT service request */
  private OffsetDateTime time;

  /** Requested IoT service information */
  private Service service;

  /** IoT service request method */
  private String method;

  /** IoT service request arguments */
  private String[] arguments;

  /**
   * Get components that compose the IoT service request key
   *
   * @return components that compose the IoT service request key
   */
  @JsonIgnore
  public String[] getKeyComponents() {
    return new String[] {this.id};
  }

  /**
   * Transform current service request to JSON string
   *
   * @return JSON representation of the service request
   */
  public String serialize() {
    return Json.serialize(this);
  }

  /**
   * Check if the IoT service request properties are valid
   *
   * @throws IllegalArgumentException when service request fields are invalid
   */
  public void validate() {
    Pattern uuidPattern =
        Pattern.compile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$");
    if (this.id == null || !uuidPattern.matcher(this.id).matches()) {
      throw new IllegalArgumentException("invalid request ID in request definition");
    }
    try {
      UUID.fromString(this.id);
    } catch (IllegalArgumentException exception) {
      throw new IllegalArgumentException("invalid request ID in request definition");
    }
    if (this.service.getDeviceId() == null
        || this.service.getDeviceId().isEmpty()
        || this.service.getOrganizationId() == null
        || this.service.getOrganizationId().isEmpty()
        || this.service.getName() == null
        || this.service.getName().isEmpty()) {
      throw new IllegalArgumentException("missing requested service in request definition");
    }
    if (this.method == null || this.method.isEmpty()) {
      throw new IllegalArgumentException("missing request method in request definition");
    }
    if (this.time == null || (this.time.toEpochSecond() == 0 && this.time.getNano() == 0)) {
      throw new IllegalArgumentException("missing request time in request definition");
    }
  }

  /**
   * Create a new service request instance from its JSON representation
   *
   * @param data JSON string representing a service request
   * @return a new service request instance
   */
  public static ServiceRequest deserialize(String data) {
    return Json.deserialize(data, ServiceRequest.class);
  }

  /** Custom JSON serializer that keeps the field order for {@link com.owlike.genson.Genson} */
  protected static final class Serializer implements com.owlike.genson.Serializer<ServiceRequest> {
    @Override
    public void serialize(ServiceRequest request, ObjectWriter writer, Context ctx)
        throws Exception {
      writer.beginObject();
      writer.writeString("id", request.id);
      writer.writeName("time");
      ctx.genson.serialize(request.time, writer, ctx);
      writer.writeName("service");
      ctx.genson.serialize(request.service, writer, ctx);
      writer.writeString("method", request.method);
      writer.writeName("arguments");
      ctx.genson.serialize(request.arguments, writer, ctx);
      writer.endObject();
    }
  }
}
