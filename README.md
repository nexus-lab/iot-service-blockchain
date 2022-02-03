# IoT Service Blockchain

## Overview

This repository contains the chaincode and SDKs of the IoT Service Blockchain project.
IoT Service Blockchain is a secure decentralized IoT service platform based-on consortium blockchain
technology.
The chaincode provides a management interface of the IoT devices and services on the Hyperledger
Fabric blockchain.
The SDKs offer application developers to join the blockchain network with their IoT devices,
define services of their devices, and request services from IoT devices.

## Installation & Usage

- Chaincode

  IoT Service Blockchain chaincode requires Hyperledger Fabric version 2.4 or above.
  The chaincode is located under [`chaincode`](chaincode) directory.
  Follow [Hyperledger Fabric's guide](https://hyperledger-fabric.readthedocs.io/en/release-2.4/deploy_chaincode.html)
  to deploy the chaincode to your Hyperledger Fabric blockchain.

- Go SDK

  To install the Go SDK of IoT Service Blockchain, run:

  ```shell
  go get github.com/nexus-lab/iot-service-blockchain
  ```

  Refer to [`tests/e2e/go`](tests/e2e/go) for usage examples of the Go SDK.

- Java SDK

  The IoT Service Blockchain Java SDK can be installed from [JitPack](https://gitpack.io).
  Visit [https://jitpack.io/#nexus-lab/iot-service-blockchain](https://jitpack.io/#nexus-lab/iot-service-blockchain)
  for more details.
  Also, refer to [`tests/e2e/java`](tests/e2e/java) for usage examples of the Java SDK.

- JavaScript SDK

  To install the JavaScript SDK of IoT Service Blockchain, you will first need to authenticate
  to GitHub Packages.
  Follow [Authenticating to GitHub Packages](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-npm-registry#authenticating-to-github-packages)
  for more information.
  Then, follow [Installing a package](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-npm-registry#installing-a-package) 
  to set up `.npmrc` file to ask your package manager to search for this package on the GitHub
  Package repository.
  Finally, run the following command to install the SDK:

  ```shell
  npm install @nexus-lab/iot-service-blockchain@0.0.1
  ```

  Also, refer to [`tests/e2e/javascript`](tests/e2e/javascript) for usage examples of the
  JavaScript SDK.

## Testing

### Requirements

- **End-to-end tests**: Docker & Docker Compose
- **Chaincode and Go SDK**: Go version 1.16 and above
- **Java SDK**: Java version 1.8 and above, Maven
- **JavaScript SDK**: Node.js version 14 and above, Yarn

### Run Unit Tests

- Chaincode and Go SDK

  Run unit tests using the following command:

  ```shell
  go test -v ./...
  ```

- Java SDK

  Run unit tests using the following command:

  ```shell
  mvn test
  ```

- JavaScript SDK

  First, install the dependencies:

  ```shell
  yarn install
  ```

  Then, run unit tests using the following command:

  ```shell
  env TZ="America/New_York" yarn test
  ```

### Run End-to-end Tests

The [`tests/scripts/fabric`](tests/scripts/fabric) script provides simple commands that quickly
creates a testing Hyperledger Fabric blockchain network.

To download and start the network, run:

```shell
./tests/scripts/fabric download
./tests/scripts/fabric network up
```

This will download necessary files and binaries to the `.fabric` and `.explorer` directories.

Use the following command to install the chaincode to the testing blockchain network:

```shell
./tests/scripts/fabric chaincode deploy
```

Then, run end-to-end tests using the following information:

- Go SDK

  To run the end-to-end tests of Go SDK, execute:

  ```shell
  export FABRIC_ROOT=$(pwd)/.fabric
  go run ./tests/e2e/go/run.go
  ```

- Java SDK

  To run the end-to-end tests of Java SDK, execute:

  ```shell
  export FABRIC_ROOT=$(pwd)/.fabric
  mvn install
  cd tests/e2e/java
  mvn compile exec:java -Dexec.mainClass="com.example.e2e.Run"
  cd -
  ```

- JavaScript SDK

  To run the end-to-end tests of JavaScript SDK, execute:

  ```shell
  export FABRIC_ROOT=$(pwd)/.fabric
  yarn build
  env TZ="America/New_York" node ./tests/e2e/javascript/run.js
  ```

Finally, when the testings are done, use this command to shutdown the network and clean up:

```shell
./tests/scripts/fabric network down
```

## Contributing

Please refer to [CONTRIBUTING.md](CONTRIBUTING.md) for the process of submitting pull requests to
us.

## License

This project is licensed under the GPL v3 License - see [LICENSE](LICENSE) for details.
