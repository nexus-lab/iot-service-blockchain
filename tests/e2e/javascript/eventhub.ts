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

async function handleDeviceEvents(isb: Sdk) {
  log('Watching for device events');

  const events = await isb.getDeviceRegistry().registerEvent();
  const expected: { [action: string]: number } = {
    register: 1,
    deregister: 1,
  };
  const actual: { [action: string]: number } = {};

  outer: for await (const event of events) {
    actual[event.action] = (actual[event.action] ?? 0) + 1;

    if (event.action === 'register' || event.action === 'deregister') {
      const device = event.payload as Device;
      if (device.id !== event.deviceId || device.organizationId !== event.organizationId) {
        fatal('device ID or organization ID mismatch');
      }
    }

    for (const [action, value] of Object.entries(expected)) {
      if (actual[action] !== value) {
        continue outer;
      }
    }

    break;
  }

  for (const [action, value] of Object.entries(expected)) {
    if (actual[action] !== value) {
      fatal(`should have received ${value} device ${action} events`);
    }
  }

  events.close();
  log('Done watching for device events');
}

async function handleServiceEvents(isb: Sdk) {
  log('Watching for service events');

  const events = await isb.getServiceRegistry().registerEvent();
  const expected: { [action: string]: number } = {
    register: 2,
    deregister: 1,
  };
  const actual: { [action: string]: number } = {};

  outer: for await (const event of events) {
    actual[event.action] = (actual[event.action] ?? 0) + 1;

    if (event.action === 'register' || event.action === 'deregister') {
      const service = event.payload as Service;
      if (
        service.deviceId !== event.deviceId ||
        service.organizationId !== event.organizationId ||
        service.name !== event.serviceName
      ) {
        fatal('event and payload device ID, organization ID, or service name mismatch');
      }
    }

    for (const [action, value] of Object.entries(expected)) {
      if (actual[action] !== value) {
        continue outer;
      }
    }

    break;
  }

  for (const [action, value] of Object.entries(expected)) {
    if (actual[action] !== value) {
      fatal(`should have received ${value} service ${action} events`);
    }
  }

  events.close();
  log('Done watching for service events');
}

async function handleRequestEvents(isb: Sdk) {
  log('Watching for service request events');

  const events = await isb.getServiceBroker().registerEvent();
  const expected: { [action: string]: number } = {
    request: 2,
    respond: 2,
    remove: 1,
  };
  const actual: { [action: string]: number } = {};

  outer: for await (const event of events) {
    actual[event.action] = (actual[event.action] ?? 0) + 1;

    if (event.action === 'request') {
      const request = event.payload as ServiceRequest;
      if (request.id !== event.requestId) {
        fatal('event and payload request ID mismatch');
      }
    } else if (event.action === 'respond') {
      const response = event.payload as ServiceResponse;
      if (response.requestId !== event.requestId) {
        fatal('event and payload request ID mismatch');
      }
    } else if (event.action === 'remove') {
      if (event.payload !== event.requestId) {
        fatal('event and payload request ID mismatch');
      }
    }

    for (const [action, value] of Object.entries(expected)) {
      if (actual[action] !== value) {
        continue outer;
      }
    }

    break;
  }

  for (const [action, value] of Object.entries(expected)) {
    if (actual[action] !== value) {
      fatal(`should have received ${value} service request ${action} events`);
    }
  }

  events.close();
  log('Done watching for service request events');
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

  const timeout = setTimeout(() => fatal('timed out waiting for events'), 120 * 1000);
  await Promise.all([handleDeviceEvents(isb), handleServiceEvents(isb), handleRequestEvents(isb)]);
  clearTimeout(timeout);

  isb.close();
}

main();
