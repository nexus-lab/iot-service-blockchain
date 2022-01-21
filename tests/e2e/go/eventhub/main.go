package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

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

func handleDeviceEvents(isb *sdk.Sdk, done chan<- string) {
	log.Println("Watching for device events")

	events, cancel, err := isb.GetDeviceRegistry().RegisterEvent()
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()

	expected := map[string]int{
		"register":   1,
		"deregister": 1,
	}

	actual := make(map[string]int)
outer:
	for event := range events {
		actual[event.Action]++

		if event.Action == "register" || event.Action == "deregister" {
			device := event.Payload.(*common.Device)
			if device.Id != event.DeviceId || device.OrganizationId != event.OrganizationId {
				log.Fatal("device ID or organization ID mismatch")
			}
		}

		for action, value := range expected {
			if actual[action] != value {
				continue outer
			}
		}

		break
	}

	for action, value := range expected {
		if actual[action] != value {
			log.Fatalf("should have received %d device %s events", value, action)
		}
	}

	done <- "device"
}

func handleServiceEvents(isb *sdk.Sdk, done chan<- string) {
	log.Println("Watching for service events")

	events, cancel, err := isb.GetServiceRegistry().RegisterEvent()
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()

	expected := map[string]int{
		"register":   2,
		"deregister": 1,
	}

	actual := make(map[string]int)
outer:
	for event := range events {
		actual[event.Action]++

		if event.Action == "register" || event.Action == "deregister" {
			service := event.Payload.(*common.Service)
			if service.DeviceId != event.DeviceId || service.OrganizationId != event.OrganizationId || service.Name != event.ServiceName {
				log.Fatal("device ID, organization ID, or service name mismatch")
			}
		}

		for action, value := range expected {
			if actual[action] != value {
				continue outer
			}
		}

		break
	}

	for action, value := range expected {
		if actual[action] != value {
			log.Fatalf("should have received %d service %s events", value, action)
		}
	}

	done <- "service"
}

func handleRequestEvents(isb *sdk.Sdk, done chan<- string) {
	log.Println("Watching for service request events")

	events, cancel, err := isb.GetServiceBroker().RegisterEvent()
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()

	expected := map[string]int{
		"request": 2,
		"respond": 2,
		"remove":  1,
	}

	actual := make(map[string]int)
outer:
	for event := range events {
		actual[event.Action]++

		if event.Action == "request" {
			request := event.Payload.(*common.ServiceRequest)
			if request.Id != event.RequestId {
				log.Fatal("request ID mismatch")
			}
		} else if event.Action == "respond" {
			response := event.Payload.(*common.ServiceResponse)
			if response.RequestId != event.RequestId {
				log.Fatal("request ID mismatch")
			}
		} else if event.Action == "remove" {
			if event.Payload.(string) != event.RequestId {
				log.Fatal("request ID mismatch")
			}
		}

		for action, value := range expected {
			if actual[action] != value {
				continue outer
			}
		}

		break
	}

	for action, value := range expected {
		if actual[action] != value {
			log.Fatalf("should have received %d service request %s events", value, action)
		}
	}

	done <- "request"
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

	ctx, cancelTimeout := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelTimeout()

	go func() {
		<-ctx.Done()
		switch ctx.Err() {
		case context.DeadlineExceeded:
			log.Fatal("timed out waiting for events")
		}
	}()

	done := make(chan string)
	defer close(done)

	go handleDeviceEvents(isb, done)
	go handleServiceEvents(isb, done)
	go handleRequestEvents(isb, done)

	count := 0
	for name := range done {
		count++
		log.Printf("Done watching for %s events\n", name)

		if count == 3 {
			break
		}
	}

	isb.Close()
}
