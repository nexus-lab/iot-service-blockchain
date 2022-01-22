import { v4 as uuidv4 } from 'uuid';

import Device from '../../../sdk/javascript/Device';
import Sdk from '../../../sdk/javascript/Sdk';
import Service from '../../../sdk/javascript/Service';
import ServiceRequest from '../../../sdk/javascript/ServiceRequest';
import ServiceResponse from '../../../sdk/javascript/ServiceResponse';
import { fatal, getCredentials, log } from './util';

const ORG_ID = 'Org2MSP';
const ORG_DOMAIN = 'org2.example.com';
const USER_NAME = 'User1@org2.example.com';
const PEER_NAME = 'peer0.org2.example.com';
const PEER_ENDPOINT = 'localhost:9051';

async function getDevice(isb: Sdk) {
  const timeout = setTimeout(() => fatal('timed out getting device'), 30 * 1000);

  while (true) {
    const devices = await isb.getDeviceRegistry().getAll('Org1MSP');
    if (devices.length > 0) {
      log(`Found device ${devices[0].serialize()}`);
      clearTimeout(timeout);
      return devices[0];
    }
  }
}

async function getServices(isb: Sdk, device: Device) {
  const timeout = setTimeout(() => fatal('timed out getting services'), 30 * 1000);

  while (true) {
    const services = await isb.getServiceRegistry().getAll(device.organizationId, device.id);
    if (services.length >= 2) {
      for (const service of services) {
        log(`Found service ${service.serialize()}`);
      }
      clearTimeout(timeout);
      return services;
    }
  }
}

async function sendServiceRequests(isb: Sdk, services: Service[]) {
  const timeout = setTimeout(() => fatal('timed out waiting for responses'), 60 * 1000);

  const events = await isb.getServiceBroker().registerEvent();

  const requests: { [id: string]: ServiceRequest } = {};
  for (const service of services) {
    const request = new ServiceRequest(uuidv4(), new Date(), service, 'GET', ['1', '2', '3']);

    log(`Sending request ${request.serialize()}`);
    await isb.getServiceBroker().request(request);
    requests[request.id] = request;
  }

  log('Listening for responses');
  for await (const event of events) {
    if (event.action === 'respond') {
      const response = event.payload as ServiceResponse;
      const request = requests[response.requestId];

      log(`Recevied response ${response.serialize()}`);

      if (response.statusCode !== 0) {
        fatal(`response error, status code is ${response.statusCode}`);
      }

      const returnValue = [request.method, ...request.args].join(',');
      if (returnValue !== response.returnValue) {
        fatal(`response return value mismatch, ${returnValue} != ${response.returnValue}`);
      }

      delete requests[response.requestId];

      if (Object.keys(requests).length === 0) {
        clearTimeout(timeout);
        events.close();
        break;
      }
    }
  }
}

async function main() {
  const [certificate, privateKey, tlsCertificate] = getCredentials(
    process.env.FABRIC_ROOT ?? '.',
    ORG_DOMAIN,
    USER_NAME,
    PEER_NAME,
  );

  const isb = new Sdk({
    organizationId: ORG_ID,
    certificate,
    privateKey,
    gatewayPeerEndpoint: PEER_ENDPOINT,
    gatewayPeerServerName: PEER_NAME,
    gatewayPeerTLSCertificate: tlsCertificate,
    networkName: 'mychannel',
    chaincodeId: 'iotservice',
  });

  // start the following processes after device and services are registered to
  // avoid PHANTOM_READ_CONFLICT
  log('Waiting for device and services to be registered');
  await new Promise<void>((resolve) => setTimeout(() => resolve(), 30 * 1000));

  const device = await getDevice(isb);
  const services = await getServices(isb, device);
  await sendServiceRequests(isb, services);

  isb.close();
}

main();
