package common

// ServiceRequestResponse a wrapper of a pair of IoT service request and response
type ServiceRequestResponse struct {
	// Request IoT service request
	Request *ServiceRequest `json:"request"`

	// Request IoT service response
	Response *ServiceResponse `json:"response"`
}
