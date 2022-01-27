package common

import (
	"encoding/json"
	"fmt"
	"time"
)

// Service an IoT service state
type Service struct {
	// Name name of the IoT service
	Name string `json:"name"`

	// DeviceId identity of the device to which the IoT service belongs
	DeviceId string `json:"deviceId"`

	// OrganizationId identity of the organization to which the IoT service belongs
	OrganizationId string `json:"organizationId"`

	// Version version number of the IoT service
	Version int32 `json:"version"`

	// Description a brief summary of the service's functions
	Description string `json:"description"`

	// LastUpdateTime the latest time that the service state has been updated
	LastUpdateTime time.Time `json:"lastUpdateTime"`
}

// GetKeyComponents return components that compose the IoT service key
func (s *Service) GetKeyComponents() []string {
	return []string{s.OrganizationId, s.DeviceId, s.Name}
}

// Serialize transform current IoT service to JSON string
func (s *Service) Serialize() ([]byte, error) {
	return json.Marshal(s)
}

// Validate check if the IoT service properties are valid
func (s *Service) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("missing service name in service definition")
	}
	if s.DeviceId == "" {
		return fmt.Errorf("missing device ID in service definition")
	}
	if s.OrganizationId == "" {
		return fmt.Errorf("missing organization ID in service definition")
	}
	if s.Version == 0 {
		return fmt.Errorf("missing service version in service definition")
	}
	if s.Version < 0 {
		return fmt.Errorf("service version must be a positive integer")
	}
	if s.LastUpdateTime.IsZero() {
		return fmt.Errorf("missing service last update time in device definition")
	}

	return nil
}

// DeserializeService create an IoT service instance from its JSON representation
func DeserializeService(data []byte) (*Service, error) {
	service := new(Service)

	if err := json.Unmarshal(data, service); err != nil {
		return nil, err
	}

	return service, nil
}
