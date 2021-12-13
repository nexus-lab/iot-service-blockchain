package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Service struct {
	Name           string    `json:"name"`
	DeviceId       string    `json:"deviceId"`
	OrganizationId string    `json:"organizationId"`
	Version        int64     `json:"version"`
	Description    string    `json:"description"`
	LastUpdateTime time.Time `json:"lastUpdateTime"`
}

func CheckService(service Service) error {
	if service.Name == "" {
		return fmt.Errorf("Missing service name in service definition")
	}
	if service.DeviceId == "" {
		return fmt.Errorf("Missing device ID in service definition")
	}
	if service.OrganizationId == "" {
		return fmt.Errorf("Missing organization ID in service definition")
	}
	if service.Version == 0 {
		return fmt.Errorf("Missing service version in service definition")
	}
	if service.LastUpdateTime.IsZero() {
		return fmt.Errorf("Missing service last update time in device definition")
	}

	return nil
}

func CreateServiceKey(ctx contractapi.TransactionContextInterface, organizationId string, deviceId string, name string) (string, error) {
	var err error
	var key string

	if key, err = ctx.GetStub().CreateCompositeKey("service", []string{organizationId, deviceId, name}); err != nil {
		return "", fmt.Errorf("Cannot create composite key for service")
	}

	return key, nil
}

func GetService(ctx contractapi.TransactionContextInterface, organizationId string, deviceId string, name string) (*Service, error) {
	var err error
	var key string
	var definition []byte

	if key, err = CreateServiceKey(ctx, organizationId, deviceId, name); err != nil {
		return nil, err
	}
	if definition, err = ctx.GetStub().GetState(key); err != nil {
		return nil, fmt.Errorf("Failed to fetch service definition")
	}

	if definition != nil {
		var service Service
		if err = json.Unmarshal(definition, &service); err != nil {
			return nil, fmt.Errorf("Cannot parse service definition")
		}
		return &service, nil
	}

	return nil, nil

}

func GetServices(ctx contractapi.TransactionContextInterface, organizationId string, deviceId string) ([]*Service, error) {
	var err error
	var iterator shim.StateQueryIteratorInterface

	services := make([]*Service, 0)
	if iterator, err = ctx.GetStub().GetStateByPartialCompositeKey("service", []string{organizationId, deviceId}); err != nil {
		return nil, fmt.Errorf("Cannot fetch registered services")
	}
	defer iterator.Close()

	for iterator.HasNext() {
		result, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("Cannot fetch next registered service")
		}

		var service Service
		if err = json.Unmarshal(result.Value, &service); err != nil {
			return nil, fmt.Errorf("Cannot parse original service definition")
		}

		services = append(services, &service)
	}

	return services, nil
}

func RegisterService(ctx contractapi.TransactionContextInterface, definition string) (*Service, error) {
	var err error
	var key string
	var organizationId, deviceId string

	if organizationId, err = GetOrganizationId(ctx); err != nil {
		return nil, err
	}
	if deviceId, err = GetDeviceId(ctx); err != nil {
		return nil, err
	}

	var service Service
	if err = json.Unmarshal([]byte(definition), &service); err != nil {
		return nil, fmt.Errorf("Cannot parse service definition")
	}

	service.DeviceId = deviceId
	service.OrganizationId = organizationId

	if err = CheckService(service); err != nil {
		return nil, err
	}

	if key, err = CreateServiceKey(ctx, organizationId, deviceId, service.Name); err != nil {
		return nil, err
	}

	definition_, err := json.Marshal(service)
	if err != nil {
		return nil, fmt.Errorf("Cannot serialize service definition")
	}

	// add to global state
	if err = ctx.GetStub().PutState(key, definition_); err != nil {
		return nil, fmt.Errorf("Failed to save service definition")
	}

	return &service, nil
}

func DeregisterService(ctx contractapi.TransactionContextInterface, name string) error {
	var err error
	var key string
	var service *Service
	var organizationId, deviceId string

	if organizationId, err = GetOrganizationId(ctx); err != nil {
		return err
	}
	if deviceId, err = GetDeviceId(ctx); err != nil {
		return err
	}

	if service, err = GetService(ctx, organizationId, deviceId, name); err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Service does not exist")
	}

	// remove related requests and responses
	requests, err := GetServiceRequests(ctx, organizationId, deviceId, name)
	if err != nil {
		return err
	}
	for _, request := range requests {
		if err = DeleteServiceRequestAndResponse(ctx, request.Id); err != nil {
			return err
		}
	}

	// remove from global state
	if key, err = CreateServiceKey(ctx, organizationId, deviceId, name); err != nil {
		return err
	}
	if err = ctx.GetStub().DelState(key); err != nil {
		return fmt.Errorf("Failed to remove service definition")
	}

	return nil
}

type ServiceRegistrySmartContract struct {
	contractapi.Contract
}

func (s *ServiceRegistrySmartContract) GetService(ctx contractapi.TransactionContextInterface, organizationId string, deviceId string, name string) (*Service, error) {
	service, err := GetService(ctx, organizationId, deviceId, name)
	if service == nil && err == nil {
		return nil, fmt.Errorf("Cannot find service")
	}
	return service, err
}

func (s *ServiceRegistrySmartContract) GetServices(ctx contractapi.TransactionContextInterface, organizationId string, deviceId string) ([]*Service, error) {
	return GetServices(ctx, organizationId, deviceId)
}

func (s *ServiceRegistrySmartContract) RegisterService(ctx contractapi.TransactionContextInterface, definition string) (*Service, error) {
	return RegisterService(ctx, definition)
}

func (s *ServiceRegistrySmartContract) DeregisterService(ctx contractapi.TransactionContextInterface, name string) error {
	return DeregisterService(ctx, name)
}
