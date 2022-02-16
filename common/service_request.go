package common

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ServiceRequest an IoT service request
type ServiceRequest struct {
	// Id identity of the IoT service request
	Id string `json:"id"`

	// Time time of the IoT service request
	Time time.Time `json:"time"`

	// Service requested IoT service information
	Service Service `json:"service"`

	// Method IoT service request method
	Method string `json:"method"`

	// Arguments IoT service request arguments
	Arguments []string `json:"arguments"`
}

// GetKeyComponents return components that compose the IoT service request key
func (r *ServiceRequest) GetKeyComponents() []string {
	return []string{r.Id}
}

// Serialize transform current IoT service request to JSON string
func (r *ServiceRequest) Serialize() ([]byte, error) {
	return json.Marshal(r)
}

// Validate check if the IoT service request properties are valid
func (r *ServiceRequest) Validate() error {
	if _, err := uuid.Parse(r.Id); err != nil {
		return fmt.Errorf("invalid request ID in request definition")
	}
	if r.Service.OrganizationId == "" || r.Service.DeviceId == "" || r.Service.Name == "" {
		return fmt.Errorf("missing requested service in request definition")
	}
	if r.Method == "" {
		return fmt.Errorf("missing request method in request definition")
	}
	if r.Arguments == nil {
		return fmt.Errorf("request arguments cannot be null in request definition")
	}
	if r.Time.IsZero() {
		return fmt.Errorf("missing request time in request definition")
	}

	return nil
}

// DeserializeService create an IoT service request instance from its JSON representation
func DeserializeServiceRequest(data []byte) (*ServiceRequest, error) {
	request := new(ServiceRequest)

	if err := json.Unmarshal(data, request); err != nil {
		return nil, err
	}

	return request, nil
}
