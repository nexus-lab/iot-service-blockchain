package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ServiceRequest struct {
	Id        string    `json:"id"`
	Time      time.Time `json:"time"`
	Service   Service   `json:"service"`
	Method    string    `json:"method"`
	Arguments []string  `json:"arguments"`
}

type ServiceResponse struct {
	RequestId   string    `json:"requestId"`
	Time        time.Time `json:"time"`
	StatusCode  int32     `json:"statusCode"`
	ReturnValue []string  `json:"returnValue"`
}

func CheckServiceRequest(request ServiceRequest) error {
	if _, err := uuid.Parse(request.Id); err != nil {
		return fmt.Errorf("Invalid request ID in request definition")
	}
	if request.Service.OrganizationId == "" || request.Service.DeviceId == "" || request.Service.Name == "" {
		return fmt.Errorf("Missing requested service in request definition")
	}
	if request.Method == "" {
		return fmt.Errorf("Missing request method in request definition")
	}
	if request.Time.IsZero() {
		return fmt.Errorf("Missing request time in request definition")
	}

	return nil
}

func CheckServiceResponse(response ServiceResponse) error {
	if _, err := uuid.Parse(response.RequestId); err != nil {
		return fmt.Errorf("Invalid request ID in response definition")
	}
	if response.Time.IsZero() {
		return fmt.Errorf("Missing response time in response definition")
	}

	return nil
}

func CreateServiceRequestKey(ctx contractapi.TransactionContextInterface, requestId string) (string, error) {
	var err error
	var key string

	if key, err = ctx.GetStub().CreateCompositeKey("request", []string{requestId}); err != nil {
		return "", fmt.Errorf("Cannot create composite key for request")
	}

	return key, nil
}

func CreateServiceResponseKey(ctx contractapi.TransactionContextInterface, requestId string) (string, error) {
	var err error
	var key string

	if key, err = ctx.GetStub().CreateCompositeKey("response", []string{requestId}); err != nil {
		return "", fmt.Errorf("Cannot create composite key for response")
	}

	return key, nil
}

func CreateServiceRequestIndexKey(ctx contractapi.TransactionContextInterface, organizationId string, deviceId string, serviceName string, requestId string) (string, error) {
	var err error
	var key string

	if key, err = ctx.GetStub().CreateCompositeKey("requestindex", []string{organizationId, deviceId, serviceName, requestId}); err != nil {
		return "", fmt.Errorf("Cannot create composite index key for request")
	}

	return key, nil
}

func GetServiceRequest(ctx contractapi.TransactionContextInterface, requestId string) (*ServiceRequest, error) {
	var err error
	var key string
	var definition []byte

	if key, err = CreateServiceRequestKey(ctx, requestId); err != nil {
		return nil, err
	}
	if definition, err = ctx.GetStub().GetState(key); err != nil {
		return nil, fmt.Errorf("Failed to fetch request definition")
	}

	if definition != nil {
		var request ServiceRequest
		if err = json.Unmarshal(definition, &request); err != nil {
			return nil, fmt.Errorf("Cannot parse request definition")
		}
		return &request, nil
	}

	return nil, nil
}

func GetServiceRequests(ctx contractapi.TransactionContextInterface, organizationId string, deviceId string, serviceName string) ([]*ServiceRequest, error) {
	var err error
	var iterator shim.StateQueryIteratorInterface

	requests := make([]*ServiceRequest, 0)
	if iterator, err = ctx.GetStub().GetStateByPartialCompositeKey("requestindex", []string{organizationId, deviceId, serviceName}); err != nil {
		return nil, fmt.Errorf("Cannot fetch request indices")
	}
	defer iterator.Close()

	for iterator.HasNext() {
		result, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("Cannot fetch next request index")
		}

		_, keys, err := ctx.GetStub().SplitCompositeKey(result.Key)
		if err != nil || len(keys) < 3 {
			return nil, fmt.Errorf("Cannot parse request index")
		}

		requestId := keys[3]
		request, err := GetServiceRequest(ctx, requestId)
		if err != nil {
			return nil, fmt.Errorf("Failed to fetch one of the request definition")
		}

		requests = append(requests, request)
	}

	return requests, nil
}

func CreateServiceRequest(ctx contractapi.TransactionContextInterface, definition string) (*ServiceRequest, error) {
	var err error
	var key string
	var indexKey string

	var request ServiceRequest
	if err = json.Unmarshal([]byte(definition), &request); err != nil {
		return nil, fmt.Errorf("Cannot parse request definition")
	}

	if err = CheckServiceRequest(request); err != nil {
		return nil, err
	}

	service := request.Service
	if key, err = CreateServiceRequestKey(ctx, request.Id); err != nil {
		return nil, err
	}
	if indexKey, err = CreateServiceRequestIndexKey(ctx, service.OrganizationId, service.DeviceId, service.Name, request.Id); err != nil {
		return nil, err
	}

	// check if request already exists
	if data, err := ctx.GetStub().GetState(key); err == nil && data != nil {
		return nil, fmt.Errorf("Request already exists")
	}

	definition_, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("Cannot serialize request definition")
	}

	// add to global state
	if err = ctx.GetStub().PutState(key, definition_); err != nil {
		return nil, fmt.Errorf("Failed to save request definition")
	}

	// add index to global state
	if err = ctx.GetStub().PutState(indexKey, []byte{0x00}); err != nil {
		return nil, fmt.Errorf("Failed to save request index")
	}

	return &request, nil
}

func GetServiceResponse(ctx contractapi.TransactionContextInterface, requestId string) (*ServiceResponse, error) {
	var err error
	var key string
	var definition []byte

	if key, err = CreateServiceResponseKey(ctx, requestId); err != nil {
		return nil, err
	}
	if definition, err = ctx.GetStub().GetState(key); err != nil {
		return nil, fmt.Errorf("Failed to fetch response definition")
	}

	if definition != nil {
		var response ServiceResponse
		if err = json.Unmarshal(definition, &response); err != nil {
			return nil, fmt.Errorf("Cannot parse response definition")
		}
		return &response, nil
	}

	return nil, nil
}

func GetServiceResponses(ctx contractapi.TransactionContextInterface, organizationId string, deviceId string, serviceName string) ([]*ServiceResponse, error) {
	var err error
	var iterator shim.StateQueryIteratorInterface

	responses := make([]*ServiceResponse, 0)
	if iterator, err = ctx.GetStub().GetStateByPartialCompositeKey("requestindex", []string{organizationId, deviceId, serviceName}); err != nil {
		return nil, fmt.Errorf("Cannot fetch request indices")
	}
	defer iterator.Close()

	for iterator.HasNext() {
		result, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("Cannot fetch next request index")
		}

		_, keys, err := ctx.GetStub().SplitCompositeKey(result.Key)
		if err != nil || len(keys) < 3 {
			return nil, fmt.Errorf("Cannot parse request index")
		}

		requestId := keys[3]
		response, err := GetServiceResponse(ctx, requestId)
		if err != nil {
			return nil, fmt.Errorf("Failed to fetch one of the response definition")
		}

		// there may be requests without a response
		if response != nil {
			responses = append(responses, response)
		}
	}

	return responses, nil
}

func CreateServiceResponse(ctx contractapi.TransactionContextInterface, definition string) (*ServiceResponse, error) {
	var err error
	var key string
	var request *ServiceRequest
	var organizationId, deviceId string

	var response ServiceResponse
	if err = json.Unmarshal([]byte(definition), &response); err != nil {
		return nil, fmt.Errorf("Cannot parse response definition")
	}

	if err = CheckServiceResponse(response); err != nil {
		return nil, err
	}

	if request, err = GetServiceRequest(ctx, response.RequestId); err != nil {
		return nil, err
	}
	if request == nil {
		return nil, fmt.Errorf("Corresponding request does not exist")
	}

	// check if the client creating the response is the client requested for service
	if organizationId, err = GetOrganizationId(ctx); err != nil {
		return nil, err
	}
	if deviceId, err = GetDeviceId(ctx); err != nil {
		return nil, err
	}
	if request.Service.OrganizationId != organizationId || request.Service.DeviceId != deviceId {
		return nil, fmt.Errorf("Cannot create response from a device other than the requested device")
	}

	if key, err = CreateServiceResponseKey(ctx, response.RequestId); err != nil {
		return nil, err
	}

	// check if response already exists
	if data, err := ctx.GetStub().GetState(key); err == nil && data != nil {
		return nil, fmt.Errorf("Response already exists")
	}

	definition_, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("Cannot serialize response definition")
	}

	// add to global state
	if err = ctx.GetStub().PutState(key, definition_); err != nil {
		return nil, fmt.Errorf("Failed to save response definition")
	}

	return &response, nil
}

// can only be called when deregistering service
func DeleteServiceRequestAndResponse(ctx contractapi.TransactionContextInterface, requestId string) error {
	var err error
	var indexKey string
	var requestKey string
	var responseKey string
	var request *ServiceRequest
	var response *ServiceResponse

	if requestKey, err = CreateServiceRequestKey(ctx, requestId); err != nil {
		return err
	}
	if request, err = GetServiceRequest(ctx, requestId); err != nil {
		return err
	}
	if request == nil {
		return fmt.Errorf("Request does not exist")
	}

	if responseKey, err = CreateServiceResponseKey(ctx, requestId); err != nil {
		return err
	}

	service := request.Service
	if indexKey, err = CreateServiceRequestIndexKey(ctx, service.OrganizationId, service.DeviceId, service.Name, request.Id); err != nil {
		return err
	}
	if response, err = GetServiceResponse(ctx, requestId); err != nil {
		return err
	}

	// remove index from global state
	if err = ctx.GetStub().DelState(indexKey); err != nil {
		return fmt.Errorf("Failed to remove request index")
	}

	// remove response from global state, if exists
	if response != nil {
		if err = ctx.GetStub().DelState(responseKey); err != nil {
			return fmt.Errorf("Failed to request definition")
		}
	}

	// remove request from global state
	if err = ctx.GetStub().DelState(requestKey); err != nil {
		return fmt.Errorf("Failed to request definition")
	}

	return nil
}

type ServiceBrokerSmartContract struct {
	contractapi.Contract
}

func (s *ServiceBrokerSmartContract) GetServiceRequest(ctx contractapi.TransactionContextInterface, requestId string) (*ServiceRequest, error) {
	request, err := GetServiceRequest(ctx, requestId)
	if request == nil && err == nil {
		return nil, fmt.Errorf("Cannot find request")
	}
	return request, err
}

func (s *ServiceBrokerSmartContract) GetServiceRequests(ctx contractapi.TransactionContextInterface, organizationId string, deviceId string, serviceName string) ([]*ServiceRequest, error) {
	return GetServiceRequests(ctx, organizationId, deviceId, serviceName)
}

func (s *ServiceBrokerSmartContract) CreateServiceRequest(ctx contractapi.TransactionContextInterface, definition string) (*ServiceRequest, error) {
	return CreateServiceRequest(ctx, definition)
}

func (s *ServiceBrokerSmartContract) GetServiceResponse(ctx contractapi.TransactionContextInterface, requestId string) (*ServiceResponse, error) {
	response, err := GetServiceResponse(ctx, requestId)
	if response == nil && err == nil {
		return nil, fmt.Errorf("Cannot find response")
	}
	return response, err
}

func (s *ServiceBrokerSmartContract) GetServiceResponses(ctx contractapi.TransactionContextInterface, organizationId string, deviceId string, serviceName string) ([]*ServiceResponse, error) {
	return GetServiceResponses(ctx, organizationId, deviceId, serviceName)
}

func (s *ServiceBrokerSmartContract) CreateServiceResponse(ctx contractapi.TransactionContextInterface, definition string) (*ServiceResponse, error) {
	return CreateServiceResponse(ctx, definition)
}
