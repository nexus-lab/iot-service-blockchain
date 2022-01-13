package common

import "encoding/json"

// ServiceRequestResponse a wrapper of a pair of IoT service request and response
type ServiceRequestResponse struct {
	// Request IoT service request
	Request *ServiceRequest `json:"request"`

	// Request IoT service response
	Response *ServiceResponse `json:"response"`
}

// Serialize transform current IoT service response to JSON string
func (r *ServiceRequestResponse) Serialize() ([]byte, error) {
	return json.Marshal(r)
}

// DeserializeService create an IoT service response instance from its JSON representation
func DeserializeServiceRequestResponse(data []byte) (*ServiceRequestResponse, error) {
	pair := new(ServiceRequestResponse)

	if err := json.Unmarshal(data, pair); err != nil {
		return nil, err
	}

	return pair, nil
}
