package common

import (
	"encoding/json"
	"fmt"
	"time"
)

// Device an IoT device state
type Device struct {
	// Id identity of the device
	Id string `json:"id"`

	// Id identity of the organization to which the device belongs
	OrganizationId string `json:"organizationId"`

	// Name friendly name of the device
	Name string `json:"name"`

	// Description a brief summary of the device's functions
	Description string `json:"description"`

	// LastUpdateTime the latest time that the device state has been updated
	LastUpdateTime time.Time `json:"lastUpdateTime"`
}

// GetKeyComponents return components that compose the device key
func (d *Device) GetKeyComponents() []string {
	return []string{d.OrganizationId, d.Id}
}

// Serialize transform current device to JSON string
func (d *Device) Serialize() ([]byte, error) {
	return json.Marshal(d)
}

// Validate check if the device properties are valid
func (d *Device) Validate() error {
	if d.Id == "" {
		return fmt.Errorf("missing device ID in device definition")
	}
	if d.OrganizationId == "" {
		return fmt.Errorf("missing organization ID in device definition")
	}
	if d.Name == "" {
		return fmt.Errorf("missing device name in device definition")
	}
	if d.LastUpdateTime.IsZero() {
		return fmt.Errorf("missing device last update time in device definition")
	}

	return nil
}

// DeserializeDevice create a new device instance from its JSON representation
func DeserializeDevice(data []byte) (*Device, error) {
	device := new(Device)

	if err := json.Unmarshal(data, device); err != nil {
		return nil, err
	}

	return device, nil
}
