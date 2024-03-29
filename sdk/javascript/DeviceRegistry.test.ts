import { ChaincodeEvent, CloseableAsyncIterable } from '@hyperledger/fabric-gateway';
import { TextEncoder } from 'util';

import Device from './Device';
import DeviceRegistry from './DeviceRegistry';

const utf8Encoder = new TextEncoder();
const mockContract = () => ({ submitTransaction: jest.fn(), registerEvent: jest.fn() });

test('deviceRegistry.register()', async () => {
  const contract = mockContract();
  const deviceRegistry = new DeviceRegistry(contract);

  const device = new Device('device1', 'org1', 'device1');
  await deviceRegistry.register(device);
  expect(contract.submitTransaction).toHaveBeenCalledWith('Register', device.serialize());

  contract.submitTransaction = jest.fn().mockRejectedValue(new Error());

  await expect(deviceRegistry.register(new Device('device2', 'org2', 'device2'))).rejects.toThrow();
});

test('deviceRegistry.get()', async () => {
  const contract = mockContract();
  const deviceRegistry = new DeviceRegistry(contract);

  const expected = new Device('device1', 'org1', 'device1');
  contract.submitTransaction = jest
    .fn()
    .mockResolvedValue(utf8Encoder.encode(expected.serialize()));

  const actual = await deviceRegistry.get('org1', 'device1');
  expect(actual.serialize()).toEqual(expected.serialize());
  expect(contract.submitTransaction).toHaveBeenCalledWith('Get', 'org1', 'device1');

  contract.submitTransaction = jest.fn().mockRejectedValue(new Error());

  await expect(deviceRegistry.get('org2', 'device2')).rejects.toThrow();
});

test('deviceRegistry.getAll()', async () => {
  const contract = mockContract();
  const deviceRegistry = new DeviceRegistry(contract);

  const expected = [
    new Device('device1', 'org1', 'device1'),
    new Device('device2', 'org2', 'device2'),
  ];
  contract.submitTransaction = jest
    .fn()
    .mockResolvedValue(
      utf8Encoder.encode(JSON.stringify(expected.map((device) => device.toObject()))),
    );

  const actual = await deviceRegistry.getAll('org1');
  for (let i = 0; i < 2; i++) {
    expect(actual[i].serialize()).toEqual(expected[i].serialize());
  }
  expect(contract.submitTransaction).toHaveBeenCalledWith('GetAll', 'org1');

  contract.submitTransaction = jest.fn().mockResolvedValue(utf8Encoder.encode('[]'));

  await expect(deviceRegistry.getAll('org2')).resolves.toEqual([]);

  contract.submitTransaction = jest.fn().mockRejectedValue(new Error());

  await expect(deviceRegistry.getAll('org3')).rejects.toThrow();
});

test('deviceRegistry.deregister()', async () => {
  const contract = mockContract();
  const deviceRegistry = new DeviceRegistry(contract);

  const device = new Device('device1', 'org1', 'device1');
  await deviceRegistry.deregister(device);
  expect(contract.submitTransaction).toHaveBeenCalledWith('Deregister', device.serialize());

  contract.submitTransaction = jest.fn().mockRejectedValue(new Error());

  await expect(
    deviceRegistry.deregister(new Device('device2', 'org2', 'device2')),
  ).rejects.toThrow();
});

test('deviceRegistry.registerEvent()', async () => {
  const contract = mockContract();
  const deviceRegistry = new DeviceRegistry(contract);

  contract.registerEvent = jest.fn().mockResolvedValue({
    async *[Symbol.asyncIterator]() {
      for (let i = 0; i < 5; i++) {
        const device = new Device(`device${i}`, `org${i}`, `device${i}`);

        yield {
          eventName: `device://org${i}/device${i}/register`,
          payload: utf8Encoder.encode(device.serialize()),
        };
      }
    },
    close: jest.fn(),
  } as CloseableAsyncIterable<ChaincodeEvent>);

  const events = await deviceRegistry.registerEvent();

  try {
    let i = 0;
    for await (const event of events) {
      expect(event.action).toBe('register');
      expect(event.organizationId).toBe(`org${i}`);
      expect(event.deviceId).toBe(`device${i}`);
      expect((event.payload as Device).id).toBe(`device${i}`);
      expect((event.payload as Device).organizationId).toBe(`org${i}`);
      expect((event.payload as Device).name).toBe(`device${i}`);

      i++;
    }
    expect(i).toBe(5);
  } finally {
    events.close();
  }

  contract.registerEvent = jest.fn().mockRejectedValue(new Error());

  await expect(deviceRegistry.registerEvent()).rejects.toThrow();
});
