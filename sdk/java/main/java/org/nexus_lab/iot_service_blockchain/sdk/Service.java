package org.nexus_lab.iot_service_blockchain.sdk;

import com.owlike.genson.Context;
import com.owlike.genson.annotation.JsonIgnore;
import com.owlike.genson.stream.ObjectWriter;
import java.time.OffsetDateTime;
import lombok.Data;
import lombok.NoArgsConstructor;

/** An IoT service state */
@Data
@NoArgsConstructor
public class Service {
  /** Friendly name of the IoT service */
  private String name;

  /** Identity of the device to which the IoT service belongs */
  private String deviceId;

  /** Identity of the organization to which the IoT service belongs */
  private String organizationId;

  /** Version number of the IoT service */
  private int version;

  /** A brief summary of the service's functions */
  private String description;

  /** The latest time that the service state has been updated */
  private OffsetDateTime lastUpdateTime;

  /**
   * Get components that compose the service key
   *
   * @return components that compose the service key
   */
  @JsonIgnore
  public String[] getKeyComponents() {
    return new String[] {this.organizationId, this.deviceId, this.name};
  }

  /**
   * Transform current service to JSON string
   *
   * @return JSON representation of the service
   */
  public String serialize() {
    return Json.serialize(this);
  }

  /**
   * Check if the IoT service properties are valid
   *
   * @throws IllegalArgumentException when service fields are invalid
   */
  public void validate() {
    if (this.name == null || this.name.isEmpty()) {
      throw new IllegalArgumentException("missing service name in service definition");
    }
    if (this.deviceId == null || this.deviceId.isEmpty()) {
      throw new IllegalArgumentException("missing device ID in service definition");
    }
    if (this.organizationId == null || this.organizationId.isEmpty()) {
      throw new IllegalArgumentException("missing organization ID in service definition");
    }
    if (this.version == 0) {
      throw new IllegalArgumentException("missing service version in service definition");
    }
    if (this.version < 0) {
      throw new IllegalArgumentException("service version must be a positive integer");
    }
    if (this.lastUpdateTime == null
        || (this.lastUpdateTime.toEpochSecond() == 0 && this.lastUpdateTime.getNano() == 0)) {
      throw new IllegalArgumentException("missing service last update time in service definition");
    }
  }

  /**
   * Create a new service instance from its JSON representation
   *
   * @param data JSON string representing a service
   * @return a new service instance
   */
  public static Service deserialize(String data) {
    return Json.deserialize(data, Service.class);
  }

  /** Custom JSON serializer that keeps the field order for {@link com.owlike.genson.Genson} */
  protected static final class Serializer implements com.owlike.genson.Serializer<Service> {
    @Override
    public void serialize(Service service, ObjectWriter writer, Context ctx) throws Exception {
      writer.beginObject();
      writer
          .writeString("name", service.name)
          .writeString("deviceId", service.deviceId)
          .writeString("organizationId", service.organizationId)
          .writeNumber("version", service.version)
          .writeString("description", service.description)
          .writeName("lastUpdateTime");
      ctx.genson.serialize(service.lastUpdateTime, writer, ctx);
      writer.endObject();
    }
  }
}
