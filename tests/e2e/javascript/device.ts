import Device from '../../../sdk/javascript/Device';
import Sdk from '../../../sdk/javascript/Sdk';
import Service from '../../../sdk/javascript/Service';
import ServiceRequest from '../../../sdk/javascript/ServiceRequest';
import ServiceResponse from '../../../sdk/javascript/ServiceResponse';
import { fatal, getCredentials, log } from './util';

const ORG_ID = 'Org1MSP';
const ORG_DOMAIN = 'org1.example.com';
const USER_NAME = 'User1@org1.example.com';
const PEER_NAME = 'peer0.org1.example.com';
const PEER_ENDPOINT = 'localhost:7051';

async function registerDevice(isb: Sdk) {
  const expected = new Device(
    isb.getDeviceId(),
    isb.getOrganizationId(),
    'device1',
    'My first device',
    new Date(),
  );

  await isb.getDeviceRegistry().register(expected);
  log(`Registered device ${expected.id}`);

  let actual = await isb.getDeviceRegistry().get(isb.getOrganizationId(), expected.id);
  if (
    expected.name !== actual.name ||
    expected.lastUpdateTime.getTime() !== actual.lastUpdateTime.getTime()
  ) {
    fatal(`inconsistent device information after registration: ${actual} != ${expected}`);
  }

  try {
    await isb.getDeviceRegistry().get(isb.getOrganizationId(), 'invalid_id');
    fatal('should return error when device is not found');
  } catch {}

  const devices = await isb.getDeviceRegistry().getAll(isb.getOrganizationId());
  if (devices.length != 1) {
    fatal(`should return only 1 device from ${isb.getOrganizationId()}`);
  }
  actual = devices[0];
  if (
    expected.name !== actual.name ||
    expected.lastUpdateTime.getTime() !== actual.lastUpdateTime.getTime()
  ) {
    fatal(
      `inconsistent device information after registration: ${actual.serialize()} != ${expected.serialize()}`,
    );
  }
}

async function registerServices(isb: Sdk) {
  const services = [
    new Service(
      'service1',
      isb.getDeviceId(),
      isb.getOrganizationId(),
      1,
      'My first service',
      new Date(),
    ),
    new Service(
      'service2',
      isb.getDeviceId(),
      isb.getOrganizationId(),
      1,
      'My second service',
      new Date(),
    ),
  ];

  for (const service of services) {
    await isb.getServiceRegistry().register(service);
    log(`Registered service ${service.serialize()}`);

    const actual = await isb
      .getServiceRegistry()
      .get(isb.getOrganizationId(), service.deviceId, service.name);
    if (
      service.name !== actual.name ||
      service.lastUpdateTime.getTime() !== actual.lastUpdateTime.getTime() ||
      service.version !== actual.version
    ) {
      fatal(
        `inconsistent service information after registration: ${actual.serialize()} != ${service.serialize()}`,
      );
    }
  }

  try {
    await isb.getServiceRegistry().get(isb.getOrganizationId(), isb.getDeviceId(), 'invalid_id');
    fatal('should return error when service is not found');
  } catch {}

  const actuals = await isb.getServiceRegistry().getAll(isb.getOrganizationId(), isb.getDeviceId());
  if (services.length !== actuals.length) {
    fatal(`should return ${services.length} service from ${isb.getOrganizationId()}`);
  }

  for (let i = 0; i < services.length; i++) {
    const expected = services[i];
    const actual = actuals[i];

    if (
      expected.name !== actual.name ||
      expected.lastUpdateTime.getTime() !== actual.lastUpdateTime.getTime() ||
      expected.version !== actual.version
    ) {
      fatal(
        `inconsistent service information after registration: ${actual.serialize()} != ${expected.serialize()}`,
      );
    }
  }
}

async function handleRequests(isb: Sdk) {
  const events = await isb.getServiceBroker().registerEvent();

  const timeout = setTimeout(() => {
    events.close();
    fatal('timed out waiting for requests');
  }, 120 * 1000);

  log('Listening for requests');
  let count = 0;
  for await (const event of events) {
    if (event.action !== 'request') {
      continue;
    }

    count++;
    const request = event.payload as ServiceRequest;

    log(`Received request ${request.serialize()}`);

    const response = new ServiceResponse(
      request.id,
      new Date(),
      0,
      [request.method, ...request.args].join(','),
    );
    await isb.getServiceBroker().respond(response);

    log(`Sent response ${response.serialize()}`);

    if (count === 2) {
      break;
    }
  }

  clearTimeout(timeout);
  events.close();
}

async function checkAndRemoveRequests(isb: Sdk) {
  const before = await isb
    .getServiceBroker()
    .getAll(isb.getOrganizationId(), isb.getDeviceId(), 'service1');
  if (before.length === 0) {
    fatal(`request/response for service1 is not found`);
  }

  const pair = before[0];
  await isb.getServiceBroker().get(pair.request.id);
  if (pair.request.id !== pair.response.requestId) {
    fatal(`request and response ID mismatch, ${pair.request.id} != ${pair.response.requestId}`);
  }

  await isb.getServiceBroker().remove(pair.request.id);

  const after = await isb
    .getServiceBroker()
    .getAll(isb.getOrganizationId(), isb.getDeviceId(), 'service1');
  if (before.length - 1 !== after.length) {
    fatal(
      `incorrect request/response count after removal: ${before.length} - 1 != ${after.length}`,
    );
  }
}

async function deregisterOneService(isb: Sdk) {
  const before = await isb.getServiceRegistry().getAll(isb.getOrganizationId(), isb.getDeviceId());

  const service = new Service('service2', isb.getDeviceId(), isb.getOrganizationId());
  await isb.getServiceRegistry().deregister(service);
  log(`Deregistered service ${service.serialize()}`);

  try {
    await isb.getServiceRegistry().get(isb.getOrganizationId(), isb.getDeviceId(), service.name);
    fatal('should return error when service is already deregistered');
  } catch {}

  const after = await isb.getServiceRegistry().getAll(isb.getOrganizationId(), isb.getDeviceId());
  if (before.length - 1 !== after.length) {
    fatal(`incorrect service count after deregistration: ${before.length} - 1 != ${after.length}`);
  }

  for (const existing of after) {
    if (existing.name === service.name) {
      fatal(`service ${service.name} has not been correctly removed`);
    }
  }
}

async function deregisterDevice(isb: Sdk) {
  const before = await isb.getDeviceRegistry().getAll(isb.getOrganizationId());

  const device = new Device(isb.getDeviceId(), isb.getOrganizationId());
  await isb.getDeviceRegistry().deregister(device);
  log(`Deregistered device ${device.serialize()}`);

  const after = await isb.getDeviceRegistry().getAll(isb.getOrganizationId());
  if (before.length - 1 !== after.length) {
    fatal(`incorrect device count after deregistration: ${before.length} - 1 != ${after.length}`);
  }

  for (const existing of after) {
    if (existing.name === device.name) {
      fatal(`device ${device.name} has not been correctly removed`);
    }
  }

  const services = await isb
    .getServiceRegistry()
    .getAll(isb.getOrganizationId(), isb.getDeviceId());
  if (services.length !== 0) {
    fatal('should have removed all services');
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

  log(`Organization ID is ${isb.getOrganizationId()}`);
  log(`Device ID is ${isb.getDeviceId()}`);

  await registerDevice(isb);
  await registerServices(isb);
  await handleRequests(isb);
  await checkAndRemoveRequests(isb);
  await deregisterOneService(isb);
  await deregisterDevice(isb);

  isb.close();
}

main();
