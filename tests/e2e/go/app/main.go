package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nexus-lab/iot-service-blockchain/common"
	sdk "github.com/nexus-lab/iot-service-blockchain/sdk/go"
)

const (
	ORG_ID        = "Org2MSP"
	ORG_DOMAIN    = "org2.example.com"
	USER_NAME     = "User1@org2.example.com"
	PEER_NAME     = "peer0.org2.example.com"
	PEER_ENDPOINT = "localhost:9051"
)

func fatalOnTimeout(timeout time.Duration, message string) context.CancelFunc {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	go func() {
		<-ctx.Done()
		switch ctx.Err() {
		case context.DeadlineExceeded:
			log.Fatal(message)
		}
	}()

	return cancel
}

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

func getDevice(isb *sdk.Sdk) *common.Device {
	cancelTimeout := fatalOnTimeout(30*time.Second, "timed out getting device")
	defer cancelTimeout()

	for {
		devices, _ := isb.GetDeviceRegistry().GetAll("Org1MSP")
		if len(devices) > 0 {
			log.Printf("Found device %#v\n", devices[0])
			return devices[0]
		}
	}
}

func getServices(isb *sdk.Sdk, device *common.Device) []*common.Service {
	cancelTimeout := fatalOnTimeout(30*time.Second, "timed out getting services")
	defer cancelTimeout()

	for {
		services, _ := isb.GetServiceRegistry().GetAll(device.OrganizationId, device.Id)
		if len(services) >= 2 {
			for _, service := range services {
				log.Printf("Found service %#v\n", service)
			}
			return services
		}
	}
}

func sendServiceRequests(isb *sdk.Sdk, services []*common.Service) {
	// avoid MVCC_READ_CONFLICT of services
	time.Sleep(20 * time.Second)

	cancelTimeout := fatalOnTimeout(60*time.Second, "timed out waiting for responses")
	defer cancelTimeout()

	events, cancel, err := isb.GetServiceBroker().RegisterEvent()
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()

	requests := make(map[string]*common.ServiceRequest)
	for _, service := range services {
		request := &common.ServiceRequest{
			Id:        uuid.New().String(),
			Time:      time.Now(),
			Service:   *service,
			Method:    "GET",
			Arguments: []string{"1", "2", "3"},
		}

		log.Printf("Sending request %#v\n", request)
		err = isb.GetServiceBroker().Request(request)
		if err != nil {
			log.Fatal(err)
		}

		requests[request.Id] = request
	}

	log.Println("Listening for responses")
	for event := range events {
		if event.Action == "respond" {
			response := event.Payload.(*common.ServiceResponse)
			request := requests[response.RequestId]

			log.Printf("Recevied response %#v\n", response)

			parts := make([]string, 0)
			parts = append(parts, request.Method)
			parts = append(parts, request.Arguments...)

			if response.StatusCode != 0 {
				log.Fatalf("response error, status code is %d", response.StatusCode)
			}

			returnValue := strings.Join(parts, ",")
			if returnValue != response.ReturnValue {
				log.Fatalf("response return value mismatch, %s != %s", returnValue, response.ReturnValue)
			}

			delete(requests, response.RequestId)

			if len(requests) == 0 {
				break
			}
		}
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

	device := getDevice(isb)
	services := getServices(isb, device)
	sendServiceRequests(isb, services)

	isb.Close()
}
