package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nexus-lab/iot-service-blockchain/common"
	sdk "github.com/nexus-lab/iot-service-blockchain/sdk/go"
)

const (
	ORG_ID        = "Org1MSP"
	ORG_DOMAIN    = "org1.example.com"
	USER_NAME     = "User1@org1.example.com"
	PEER_NAME     = "peer0.org1.example.com"
	PEER_ENDPOINT = "localhost:7051"
)

func getCredentials() ([]byte, []byte, []byte) {
	root := filepath.Join(os.Getenv("FABRIC_ROOT"), "test-network/organizations/peerOrganizations/", ORG_DOMAIN)
	filepaths := []string{
		"users/" + USER_NAME + "/msp/signcerts/cert.pem",
		"users/" + USER_NAME + "/msp/keystore/priv_sk",
		"peers/" + PEER_NAME + "/tls/ca.crt",
	}

	files := make([][]byte, 0)
	for _, path := range filepaths {
		data, err := ioutil.ReadFile(filepath.Join(root, path))
		if err != nil {
			log.Fatal(err)
		}
		files = append(files, data)
	}

	return files[0], files[1], files[2]
}

func registerDevice(isb *sdk.Sdk) {
	expected := &common.Device{
		Id:             isb.GetDeviceId(),
		OrganizationId: isb.GetOrganizationId(),
		Name:           "device1",
		Description:    "My first device",
		LastUpdateTime: time.Now(),
	}

	err := isb.GetDeviceRegistry().Register(expected)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Registered device %#v\n", expected)

	actual, err := isb.GetDeviceRegistry().Get(isb.GetOrganizationId(), expected.Id)
	if err != nil {
		log.Fatal(err)
	}
	if expected.Name != actual.Name || !expected.LastUpdateTime.UTC().Equal(actual.LastUpdateTime) {
		log.Fatalf("inconsistent device information after registration: %#v != %#v", actual, expected)
	}

	_, err = isb.GetDeviceRegistry().Get(isb.GetOrganizationId(), "invalid_id")
	if err == nil {
		log.Fatal("should return error when device is not found")
	}

	devices, err := isb.GetDeviceRegistry().GetAll(isb.GetOrganizationId())
	if err != nil {
		log.Fatal(err)
	}

	if len(devices) != 1 {
		log.Fatalf("should return only 1 device from %s", isb.GetOrganizationId())
	}
	actual = devices[0]
	if expected.Name != actual.Name || !expected.LastUpdateTime.UTC().Equal(actual.LastUpdateTime) {
		log.Fatalf("inconsistent device information after registration: %#v != %#v", actual, expected)
	}
}

func registerServices(isb *sdk.Sdk) {
	services := []*common.Service{
		{
			OrganizationId: isb.GetOrganizationId(),
			DeviceId:       isb.GetDeviceId(),
			Name:           "service1",
			Description:    "My first service",
			Version:        1,
			LastUpdateTime: time.Now(),
		},
		{
			OrganizationId: isb.GetOrganizationId(),
			DeviceId:       isb.GetDeviceId(),
			Name:           "service2",
			Description:    "My second service",
			Version:        1,
			LastUpdateTime: time.Now(),
		},
	}

	for _, service := range services {
		err := isb.GetServiceRegistry().Register(service)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Registered service %#v\n", service)

		actual, err := isb.GetServiceRegistry().Get(isb.GetOrganizationId(), service.DeviceId, service.Name)
		if err != nil {
			log.Fatal(err)
		}
		if service.Name != actual.Name || !service.LastUpdateTime.UTC().Equal(actual.LastUpdateTime) || service.Version != actual.Version {
			log.Fatalf("inconsistent service information after registration: %#v != %#v", actual, service)
		}
	}

	_, err := isb.GetServiceRegistry().Get(isb.GetOrganizationId(), isb.GetDeviceId(), "invalid_id")
	if err == nil {
		log.Fatal("should return error when service is not found")
	}

	actuals, err := isb.GetServiceRegistry().GetAll(isb.GetOrganizationId(), isb.GetDeviceId())
	if err != nil {
		log.Fatal(err)
	}

	if len(services) != len(actuals) {
		log.Fatalf("should return %d service from %s", len(services), isb.GetOrganizationId())
	}

	for i := range services {
		expected := services[i]
		actual := actuals[i]
		if expected.Name != actual.Name || !expected.LastUpdateTime.UTC().Equal(actual.LastUpdateTime) {
			log.Fatalf("inconsistent service information after registration: %#v != %#v", actual, expected)
		}
	}
}

func handleRequests(isb *sdk.Sdk) {
	ctx, cancelTimeout := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelTimeout()

	go func() {
		<-ctx.Done()
		switch ctx.Err() {
		case context.DeadlineExceeded:
			log.Fatal("timed out waiting for requests")
		}
	}()

	events, cancel, err := isb.GetServiceBroker().RegisterEvent()
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()

	log.Println("Listening for requests")
	count := 0
	for event := range events {
		if event.Action != "request" {
			continue
		}

		count++
		request := event.Payload.(*common.ServiceRequest)

		log.Printf("Received request %#v\n", request)

		returnValue := make([]string, 0)
		returnValue = append(returnValue, request.Method)
		returnValue = append(returnValue, request.Arguments...)

		response := &common.ServiceResponse{
			RequestId:   request.Id,
			Time:        time.Now(),
			StatusCode:  0,
			ReturnValue: strings.Join(returnValue, ","),
		}
		err = isb.GetServiceBroker().Respond(response)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Sent response %#v\n", response)

		if count == 2 {
			break
		}
	}
}

func checkAndRemoveRequests(isb *sdk.Sdk) {
	before, err := isb.GetServiceBroker().GetAll(isb.GetOrganizationId(), isb.GetDeviceId(), "service1")
	if err != nil {
		log.Fatal(err)
	}
	if len(before) == 0 {
		log.Fatalf("request/response for %s is not found", "service1")
	}

	pair := before[0]
	_, err = isb.GetServiceBroker().Get(pair.Request.Id)
	if err != nil {
		log.Fatal(err)
	}
	if pair.Request.Id != pair.Response.RequestId {
		log.Fatalf("request and response ID mismatch, %s != %s", pair.Request.Id, pair.Response.RequestId)
	}

	err = isb.GetServiceBroker().Remove(pair.Request.Id)
	if err != nil {
		log.Fatal(err)
	}

	after, err := isb.GetServiceBroker().GetAll(isb.GetOrganizationId(), isb.GetDeviceId(), "service1")
	if err != nil {
		log.Fatal(err)
	}
	if len(before)-1 != len(after) {
		log.Fatalf("incorrect request/response count after removal: %#v - 1 != %#v", len(before), len(after))
	}
}

func deregisterOneService(isb *sdk.Sdk) {
	before, err := isb.GetServiceRegistry().GetAll(isb.GetOrganizationId(), isb.GetDeviceId())
	if err != nil {
		log.Fatal(err)
	}

	service := &common.Service{
		OrganizationId: isb.GetOrganizationId(),
		DeviceId:       isb.GetDeviceId(),
		Name:           "service2",
	}
	err = isb.GetServiceRegistry().Deregister(service)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Deregistered service %#v\n", service)

	_, err = isb.GetServiceRegistry().Get(isb.GetOrganizationId(), isb.GetDeviceId(), service.Name)
	if err == nil {
		log.Fatal("should return error when service is already deregistered")
	}

	after, err := isb.GetServiceRegistry().GetAll(isb.GetOrganizationId(), isb.GetDeviceId())
	if err != nil {
		log.Fatal(err)
	}

	if len(before)-1 != len(after) {
		log.Fatalf("incorrect service count after deregistration: %#v - 1 != %#v", len(before), len(after))
	}

	for _, existing := range after {
		if existing.Name == service.Name {
			log.Fatalf("service %s has not been correctly removed", service.Name)
		}
	}
}

func deregisterDevice(isb *sdk.Sdk) {
	before, err := isb.GetDeviceRegistry().GetAll(isb.GetOrganizationId())
	if err != nil {
		log.Fatal(err)
	}

	device := &common.Device{
		Id:             isb.GetDeviceId(),
		OrganizationId: isb.GetOrganizationId(),
	}
	err = isb.GetDeviceRegistry().Deregister(device)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Deregistered device %#v\n", device)

	after, err := isb.GetDeviceRegistry().GetAll(isb.GetOrganizationId())
	if err != nil {
		log.Fatal(err)
	}

	if len(before)-1 != len(after) {
		log.Fatalf("incorrect device count after deregistration: %#v - 1 != %#v", len(before), len(after))
	}

	for _, existing := range after {
		if existing.Name == device.Name {
			log.Fatalf("device %s has not been correctly removed", device.Name)
		}
	}

	services, err := isb.GetServiceRegistry().GetAll(isb.GetOrganizationId(), isb.GetDeviceId())
	if err != nil {
		log.Fatal(err)
	}
	if len(services) != 0 {
		log.Fatal("should have removed all services")
	}
}

func main() {
	certificate, privateKey, tlsCertificate := getCredentials()

	isb, err := sdk.NewSdk(
		&sdk.SdkOptions{
			OrganizationId:            ORG_ID,
			Certificate:               certificate,
			PrivateKey:                privateKey,
			GatewayPeerEndpoint:       PEER_ENDPOINT,
			GatewayPeerServerName:     PEER_NAME,
			GatewayPeerTLSCertificate: tlsCertificate,
			NetworkName:               "mychannel",
			ChaincodeId:               "iotservice",
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Organization ID is %s\n", isb.GetOrganizationId())
	log.Printf("Device ID is %s\n", isb.GetDeviceId())

	registerDevice(isb)
	registerServices(isb)
	handleRequests(isb)
	checkAndRemoveRequests(isb)
	deregisterOneService(isb)
	deregisterDevice(isb)

	isb.Close()
}
