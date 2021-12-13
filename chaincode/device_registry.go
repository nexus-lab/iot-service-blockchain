package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Device struct {
	Id             string    `json:"id"`
	OrganizationId string    `json:"organizationId"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	LastUpdateTime time.Time `json:"lastUpdateTime"`
}

func CheckDevice(device Device) error {
	if device.Id == "" {
		return fmt.Errorf("Missing device ID in device definition")
	}
	if device.OrganizationId == "" {
		return fmt.Errorf("Missing organization ID in device definition")
	}
	if device.Name == "" {
		return fmt.Errorf("Missing device name in device definition")
	}
	if device.LastUpdateTime.IsZero() {
		return fmt.Errorf("Missing device last update time in device definition")
	}

	return nil
}

func CreateDeviceKey(ctx contractapi.TransactionContextInterface, organizationId string, deviceId string) (string, error) {
	var err error
	var key string

	if key, err = ctx.GetStub().CreateCompositeKey("device", []string{organizationId, deviceId}); err != nil {
		return "", fmt.Errorf("Cannot create composite key for device")
	}

	return key, nil
}

func GetDevice(ctx contractapi.TransactionContextInterface, organizationId string, deviceId string) (*Device, error) {
	var err error
	var key string
	var definition []byte

	if key, err = CreateDeviceKey(ctx, organizationId, deviceId); err != nil {
		return nil, err
	}
	if definition, err = ctx.GetStub().GetState(key); err != nil {
		return nil, fmt.Errorf("Failed to get device definition")
	}

	if definition != nil {
		var device Device
		if err = json.Unmarshal(definition, &device); err != nil {
			return nil, fmt.Errorf("Cannot parse original device definition")
		}
		return &device, nil
	}

	return nil, nil
}

func GetDevices(ctx contractapi.TransactionContextInterface, organizationId string) ([]*Device, error) {
	var err error
	var iterator shim.StateQueryIteratorInterface

	devices := make([]*Device, 0)
	if iterator, err = ctx.GetStub().GetStateByPartialCompositeKey("device", []string{organizationId}); err != nil {
		return nil, fmt.Errorf("Cannot fetch registered devices")
	}
	defer iterator.Close()

	for iterator.HasNext() {
		result, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("Cannot fetch next registered device")
		}

		var device Device
		if err = json.Unmarshal(result.Value, &device); err != nil {
			return nil, fmt.Errorf("Cannot parse original device definition")
		}

		devices = append(devices, &device)
	}

	return devices, nil
}

func RegisterDevice(ctx contractapi.TransactionContextInterface, definition string) (*Device, error) {
	var err error
	var key string
	var device *Device
	var organizationId, deviceId string

	if organizationId, err = GetOrganizationId(ctx); err != nil {
		return nil, err
	}
	if deviceId, err = GetDeviceId(ctx); err != nil {
		return nil, err
	}

	if device, err = GetDevice(ctx, organizationId, deviceId); err != nil {
		return nil, err
	}
	if device == nil {
		device = new(Device)
	}

	var update Device
	if err = json.Unmarshal([]byte(definition), &update); err != nil {
		return nil, fmt.Errorf("Cannot parse device definition")
	}

	// only update name, description, and last update time from client
	device.Name = update.Name
	device.Description = update.Description
	device.LastUpdateTime = update.LastUpdateTime
	device.Id = deviceId
	device.OrganizationId = organizationId

	if err = CheckDevice(*device); err != nil {
		return nil, err
	}

	if key, err = CreateDeviceKey(ctx, device.OrganizationId, device.Id); err != nil {
		return nil, err
	}

	definition_, err := json.Marshal(*device)
	if err != nil {
		return nil, fmt.Errorf("Cannot serialize device definition")
	}

	// add to global state
	if err = ctx.GetStub().PutState(key, definition_); err != nil {
		return nil, fmt.Errorf("Failed to save device definition")
	}

	return device, nil
}

func DeregisterDevice(ctx contractapi.TransactionContextInterface) error {
	var err error
	var key string
	var device *Device
	var organizationId, deviceId string

	if organizationId, err = GetOrganizationId(ctx); err != nil {
		return err
	}
	if deviceId, err = GetDeviceId(ctx); err != nil {
		return err
	}

	if device, err = GetDevice(ctx, organizationId, deviceId); err != nil {
		return err
	}
	if device == nil {
		return fmt.Errorf("Device does not exist")
	}

	// deregister services
	services, err := GetServices(ctx, organizationId, deviceId)
	if err != nil {
		return err
	}
	for _, service := range services {
		if err = DeregisterService(ctx, service.Name); err != nil {
			return err
		}
	}

	// remove from global state
	if key, err = CreateDeviceKey(ctx, organizationId, deviceId); err != nil {
		return err
	}
	if err = ctx.GetStub().DelState(key); err != nil {
		return fmt.Errorf("Failed to remove device definition")
	}

	return nil
}

type DeviceRegistrySmartContract struct {
	contractapi.Contract
}

func (s *DeviceRegistrySmartContract) GetDevice(ctx contractapi.TransactionContextInterface, organizationId string, deviceId string) (*Device, error) {
	device, err := GetDevice(ctx, organizationId, deviceId)
	if device == nil && err == nil {
		return nil, fmt.Errorf("Cannot find device")
	}
	return device, err
}

func (s *DeviceRegistrySmartContract) GetDevices(ctx contractapi.TransactionContextInterface, organizationId string) ([]*Device, error) {
	return GetDevices(ctx, organizationId)
}

func (s *DeviceRegistrySmartContract) RegisterDevice(ctx contractapi.TransactionContextInterface, definition string) (*Device, error) {
	return RegisterDevice(ctx, definition)
}

func (s *DeviceRegistrySmartContract) DeregisterDevice(ctx contractapi.TransactionContextInterface) error {
	return DeregisterDevice(ctx)
}
