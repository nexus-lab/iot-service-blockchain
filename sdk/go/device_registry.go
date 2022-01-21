package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/nexus-lab/iot-service-blockchain/common"
)

// DeviceEvent an event emitted by the device registry contract notifying a device update
type DeviceEvent struct {
	// Action name of the action performed on the device
	Action string

	// OrganizationId organization ID of the device
	OrganizationId string

	// DeviceId ID of the device
	DeviceId string

	// Payload custom event payload
	Payload interface{}
}

// DeviceRegistryInterface core utilities for managing devices on the ledger
type DeviceRegistryInterface interface {
	// Register create or update a device in the ledger
	Register(device *common.Device) error

	// Get return a device by its organization ID and device ID
	Get(organizationId string, deviceId string) (*common.Device, error)

	// GetAll return a list of devices by their organization ID
	GetAll(organizationId string) ([]*common.Device, error)

	// Deregister remove a device from the ledger
	Deregister(device *common.Device) error

	// RegisterEvent registers for device registry events
	RegisterEvent(options ...client.ChaincodeEventsOption) (<-chan *DeviceEvent, context.CancelFunc, error)
}

// DeviceRegistry core utilities for managing devices on the ledger
type DeviceRegistry struct {
	contract ContractInterface
}

// Register create or update a device in the ledger
func (r *DeviceRegistry) Register(device *common.Device) error {
	if device == nil {
		return fmt.Errorf("cannot register an empty device")
	}

	data, err := device.Serialize()
	if err != nil {
		return err
	}

	_, err = r.contract.SubmitTransaction("Register", string(data))
	return err
}

// Get return a device by its organization ID and device ID
func (r *DeviceRegistry) Get(organizationId string, deviceId string) (*common.Device, error) {
	data, err := r.contract.SubmitTransaction("Get", organizationId, deviceId)
	if err != nil {
		return nil, err
	}

	return common.DeserializeDevice(data)
}

// GetAll return a list of devices by their organization ID
func (r *DeviceRegistry) GetAll(organizationId string) ([]*common.Device, error) {
	data, err := r.contract.SubmitTransaction("GetAll", organizationId)
	if err != nil {
		return nil, err
	}

	results := make([]*common.Device, 0)
	if err = json.Unmarshal(data, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// Deregister remove a device from the ledger
func (r *DeviceRegistry) Deregister(device *common.Device) error {
	if device == nil {
		return fmt.Errorf("cannot deregister an empty device")
	}

	data, err := device.Serialize()
	if err != nil {
		return err
	}

	_, err = r.contract.SubmitTransaction("Deregister", string(data))
	return err
}

// RegisterEvent registers for device registry events
func (r *DeviceRegistry) RegisterEvent(options ...client.ChaincodeEventsOption) (<-chan *DeviceEvent, context.CancelFunc, error) {
	dest := make(chan *DeviceEvent)
	source, cancel, err := r.contract.RegisterEvent(options...)
	pattern := regexp.MustCompile(`^device:\/\/(.+?)\/(.+?)\/(.+?)$`)

	go func() {
		defer close(dest)

		for event := range source {
			matches := pattern.FindStringSubmatch(event.EventName)
			if len(matches) != 4 {
				continue
			}

			deviceEvent := &DeviceEvent{
				OrganizationId: matches[1],
				DeviceId:       matches[2],
				Action:         matches[3],
			}

			if deviceEvent.Action == "register" || deviceEvent.Action == "deregister" {
				device, err := common.DeserializeDevice(event.Payload)
				if err != nil {
					log.Printf("bad device event payload %#v, action is %s\n", event.Payload, deviceEvent.Action)
					continue
				}
				deviceEvent.Payload = device
			} else {
				deviceEvent.Payload = event.Payload
			}

			dest <- deviceEvent
		}
	}()

	return dest, cancel, err
}

func createDeviceRegistry(network *client.Network, chaincodeId string) DeviceRegistryInterface {
	return &DeviceRegistry{
		contract: &Contract{
			network:      network,
			chaincodeId:  chaincodeId,
			contractName: "device_registry",
		},
	}
}
