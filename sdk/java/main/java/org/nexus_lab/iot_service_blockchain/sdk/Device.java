package org.nexus_lab.iot_service_blockchain.sdk;

import com.owlike.genson.Context;
import com.owlike.genson.annotation.JsonIgnore;
import com.owlike.genson.stream.ObjectWriter;
import java.time.OffsetDateTime;
import lombok.Data;
import lombok.NoArgsConstructor;

/** An IoT device state */
@Data
@NoArgsConstructor
public class Device {
  /** Identity of the device */
  private String id;

  /** Identity of the organization to which the device belongs */
  private String organizationId;

  /** Friendly name of the device */
  private String name;

  /** A brief summary of the device's functions */
  private String description;

  /** The latest time that the device state has been updated */
  private OffsetDateTime lastUpdateTime;

  /**
   * Get components that compose the device key
   *
   * @return components that compose the device key
   */
  @JsonIgnore
  public String[] getKeyComponents() {
    return new String[] {this.organizationId, this.id};
  }

  /**
   * Transform current device to JSON string
   *
   * @return JSON representation of the device
   */
  public String serialize() {
    return Json.serialize(this);
  }

  /**
   * Check if the device properties are valid
   *
   * @throws IllegalArgumentException when device fields are invalid
   */
  public void validate() {
    if (this.id == null || this.id.isEmpty()) {
      throw new IllegalArgumentException("missing device ID in device definition");
    }
    if (this.organizationId == null || this.organizationId.isEmpty()) {
      throw new IllegalArgumentException("missing organization ID in device definition");
    }
    if (this.name == null || this.name.isEmpty()) {
      throw new IllegalArgumentException("missing device name in device definition");
    }
    if (this.lastUpdateTime == null
        || (this.lastUpdateTime.toEpochSecond() == 0 && this.lastUpdateTime.getNano() == 0)) {
      throw new IllegalArgumentException("missing device last update time in device definition");
    }
  }

  /**
   * Create a new device instance from its JSON representation
   *
   * @param data JSON string representing a device
   * @return a new device instance
   */
  public static Device deserialize(String data) {
    return Json.deserialize(data, Device.class);
  }

  /** Custom JSON serializer that keeps the field order for {@link com.owlike.genson.Genson} */
  protected static final class Serializer implements com.owlike.genson.Serializer<Device> {
    @Override
    public void serialize(Device device, ObjectWriter writer, Context ctx) throws Exception {
      writer.beginObject();
      writer
          .writeString("id", device.id)
          .writeString("organizationId", device.organizationId)
          .writeString("name", device.name)
          .writeString("description", device.description)
          .writeName("lastUpdateTime");
      ctx.genson.serialize(device.lastUpdateTime, writer, ctx);
      writer.endObject();
    }
  }
}
