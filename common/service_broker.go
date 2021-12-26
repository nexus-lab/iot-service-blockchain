package common

// ServiceRequestResponse a wrapper of a pair of IoT service request and response
type ServiceRequestResponse struct {
	// Request IoT service request
	Request *ServiceRequest `json:"request"`

	// Request IoT service response
	Response *ServiceResponse `json:"response"`
}

// ServiceBrokerInterface core utilities for managing service requests on ledger
type ServiceBrokerInterface interface {
	// Request make a request to an IoT service
	Request(request *ServiceRequest) error

	// Respond respond to an IoT service request
	Respond(response *ServiceResponse) error

	// Get return an IoT service request and its response by the request ID
	Get(requestId string) (*ServiceRequestResponse, error)

	// GetAll return a list of IoT service requests and their responses by their organization ID, device ID, and service name
	GetAll(organizationId string, deviceId string, serviceName string) ([]*ServiceRequestResponse, error)

	// Remove remove a (request, response) pair from the ledger
	Remove(requestId string) error
}
