import { ChaincodeEvent, CloseableAsyncIterable } from '@hyperledger/fabric-gateway';
import { TextEncoder } from 'util';

import Service from './Service';
import ServiceRegistry from './ServiceRegistry';

const utf8Encoder = new TextEncoder();
const mockContract = () => ({ submitTransaction: jest.fn(), registerEvent: jest.fn() });

test('serviceRegistry.register()', async () => {
  const contract = mockContract();
  const serviceRegistry = new ServiceRegistry(contract);

  const service = new Service('service1', 'device1', 'org1');
  await serviceRegistry.register(service);
  expect(contract.submitTransaction).toHaveBeenCalledWith('Register', service.serialize());

  contract.submitTransaction = jest.fn().mockRejectedValue(new Error());

  expect(serviceRegistry.register(new Service('service2', 'device2', 'org2'))).rejects.toThrow();
});

test('serviceRegistry.get()', () => {
  const contract = mockContract();
  const serviceRegistry = new ServiceRegistry(contract);

  const service = new Service('service1', 'device1', 'org1');
  contract.submitTransaction = jest.fn().mockResolvedValue(utf8Encoder.encode(service.serialize()));

  expect(serviceRegistry.get('org1', 'device1', 'service1')).resolves.toEqual(service);
  expect(contract.submitTransaction).toHaveBeenCalledWith('Get', 'org1', 'device1', 'service1');

  contract.submitTransaction = jest.fn().mockRejectedValue(new Error());

  expect(serviceRegistry.get('org2', 'device2', 'service2')).rejects.toThrow();
});

test('serviceRegistry.getAll()', () => {
  const contract = mockContract();
  const serviceRegistry = new ServiceRegistry(contract);

  const services = [
    new Service('service1', 'device1', 'org1').toObject(),
    new Service('service2', 'device2', 'org2').toObject(),
  ];
  contract.submitTransaction = jest
    .fn()
    .mockResolvedValue(utf8Encoder.encode(JSON.stringify(services)));

  expect(serviceRegistry.getAll('org1', 'device1')).resolves.toEqual(services);
  expect(contract.submitTransaction).toHaveBeenCalledWith('GetAll', 'org1', 'device1');

  contract.submitTransaction = jest.fn().mockResolvedValue(utf8Encoder.encode('[]'));

  expect(serviceRegistry.getAll('org2', 'device2')).resolves.toEqual([]);

  contract.submitTransaction = jest.fn().mockRejectedValue(new Error());

  expect(serviceRegistry.getAll('org3', 'device3')).rejects.toThrow();
});

test('serviceRegistry.deregister()', async () => {
  const contract = mockContract();
  const serviceRegistry = new ServiceRegistry(contract);

  const service = new Service('service1', 'device1', 'org1');
  await serviceRegistry.deregister(service);
  expect(contract.submitTransaction).toHaveBeenCalledWith('Deregister', service.serialize());

  contract.submitTransaction = jest.fn().mockRejectedValue(new Error());

  expect(serviceRegistry.deregister(new Service('service2', 'device2', 'org2'))).rejects.toThrow();
});

test('serviceRegistry.registerEvent()', async () => {
  const contract = mockContract();
  const serviceRegistry = new ServiceRegistry(contract);

  contract.registerEvent = jest.fn().mockResolvedValue({
    async *[Symbol.asyncIterator]() {
      for (let i = 0; i < 5; i++) {
        const service = new Service(`service${i}`, `device${i}`, `org${i}`);

        yield {
          eventName: `service://org${i}/device${i}/service${i}/register`,
          payload: utf8Encoder.encode(service.serialize()),
        };
      }
    },
    close: jest.fn(),
  } as CloseableAsyncIterable<ChaincodeEvent>);

  const events = await serviceRegistry.registerEvent();

  try {
    let i = 0;
    for await (const event of events) {
      expect(event.action).toBe('register');
      expect(event.organizationId).toBe(`org${i}`);
      expect(event.deviceId).toBe(`device${i}`);
      expect(event.serviceName).toBe(`service${i}`);
      expect((event.payload as Service).deviceId).toBe(`device${i}`);
      expect((event.payload as Service).organizationId).toBe(`org${i}`);
      expect((event.payload as Service).name).toBe(`service${i}`);

      i++;
    }
    expect(i).toBe(5);
  } finally {
    events.close();
  }

  contract.registerEvent = jest.fn().mockRejectedValue(new Error());

  expect(serviceRegistry.registerEvent()).rejects.toThrow();
});
