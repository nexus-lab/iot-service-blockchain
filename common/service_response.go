package common

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ServiceRequest an IoT service response
type ServiceResponse struct {
	// RequestId identity of the IoT service request to respond to
	RequestId string `json:"requestId"`

	// Time time of the IoT service response
	Time time.Time `json:"time"`

	// StatusCode status code of the IoT service response
	StatusCode int32 `json:"statusCode"`

	// ReturnValue return value of the IoT service response
	ReturnValue string `json:"returnValue"`
}

// GetKeyComponents return components that compose the IoT service response key
func (r *ServiceResponse) GetKeyComponents() []string {
	return []string{r.RequestId}
}

// Serialize transform current IoT service response to JSON string
func (r *ServiceResponse) Serialize() ([]byte, error) {
	return json.Marshal(r)
}

// Validate check if the IoT service response properties are valid
func (r *ServiceResponse) Validate() error {
	if _, err := uuid.Parse(r.RequestId); err != nil {
		return fmt.Errorf("invalid request ID in response definition")
	}
	if r.Time.IsZero() {
		return fmt.Errorf("missing response time in response definition")
	}

	return nil
}

// DeserializeService create an IoT service response instance from its JSON representation
func DeserializeServiceResponse(data []byte) (*ServiceResponse, error) {
	response := new(ServiceResponse)

	if err := json.Unmarshal(data, response); err != nil {
		return nil, err
	}

	return response, nil
}
