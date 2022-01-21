import { ChaincodeEvent, CloseableAsyncIterable } from '@hyperledger/fabric-gateway';
import { TextEncoder } from 'util';

import Service from './Service';
import ServiceBroker from './ServiceBroker';
import ServiceRequest from './ServiceRequest';
import ServiceRequestResponse from './ServiceRequestResponse';
import ServiceResponse from './ServiceResponse';

const utf8Encoder = new TextEncoder();
const mockContract = () => ({ submitTransaction: jest.fn(), registerEvent: jest.fn() });
const createRequest = (id: number) =>
  new ServiceRequest(
    `request${id}`,
    new Date(),
    new Service(`service${id}`, `device${id}`, `org${id}`),
    'GET',
    [],
  );

test('serviceBroker.request()', async () => {
  const contract = mockContract();
  const serviceBroker = new ServiceBroker(contract);

  const request = createRequest(1);
  await serviceBroker.request(request);
  expect(contract.submitTransaction).toHaveBeenCalledWith('Request', request.serialize());

  contract.submitTransaction = jest.fn().mockRejectedValue(new Error());

  expect(serviceBroker.request(createRequest(2))).rejects.toThrow();
});

test('serviceBroker.respond()', async () => {
  const contract = mockContract();
  const serviceBroker = new ServiceBroker(contract);

  const response = new ServiceResponse('request1', new Date());
  await serviceBroker.respond(response);
  expect(contract.submitTransaction).toHaveBeenCalledWith('Respond', response.serialize());

  contract.submitTransaction = jest.fn().mockRejectedValue(new Error());

  expect(serviceBroker.respond(new ServiceResponse('request2', new Date()))).rejects.toThrow();
});

test('serviceBroker.get()', () => {
  const contract = mockContract();
  const serviceBroker = new ServiceBroker(contract);

  const pair = new ServiceRequestResponse(
    createRequest(1),
    new ServiceResponse('request1', new Date()),
  );
  contract.submitTransaction = jest.fn().mockResolvedValue(utf8Encoder.encode(pair.serialize()));

  expect(serviceBroker.get('request1')).resolves.toEqual(pair);
  expect(contract.submitTransaction).toHaveBeenCalledWith('Get', 'request1');

  contract.submitTransaction = jest.fn().mockRejectedValue(new Error());

  expect(serviceBroker.get('request2')).rejects.toThrow();
});

test('serviceBroker.getAll()', () => {
  const contract = mockContract();
  const serviceBroker = new ServiceBroker(contract);

  const pairs = [
    new ServiceRequestResponse(createRequest(1), new ServiceResponse('request1', new Date())),
    new ServiceRequestResponse(createRequest(2), new ServiceResponse('request2', new Date())),
  ];
  contract.submitTransaction = jest
    .fn()
    .mockResolvedValue(utf8Encoder.encode(JSON.stringify(pairs.map((pair) => pair.toObject()))));

  expect(serviceBroker.getAll('org1', 'device1', 'service1')).resolves.toEqual(pairs);
  expect(contract.submitTransaction).toHaveBeenCalledWith('GetAll', 'org1', 'device1', 'service1');

  contract.submitTransaction = jest.fn().mockResolvedValue(utf8Encoder.encode('[]'));

  expect(serviceBroker.getAll('org2', 'device2', 'service2')).resolves.toEqual([]);

  contract.submitTransaction = jest.fn().mockRejectedValue(new Error());

  expect(serviceBroker.getAll('org3', 'device3', 'service3')).rejects.toThrow();
});

test('serviceBroker.remove()', async () => {
  const contract = mockContract();
  const serviceBroker = new ServiceBroker(contract);

  await serviceBroker.remove('request1');
  expect(contract.submitTransaction).toHaveBeenCalledWith('Remove', 'request1');

  contract.submitTransaction = jest.fn().mockRejectedValue(new Error());

  expect(serviceBroker.remove('request2')).rejects.toThrow();
});

test('serviceBroker.registerEvent()', async () => {
  const contract = mockContract();
  const serviceBroker = new ServiceBroker(contract);

  contract.registerEvent = jest.fn().mockResolvedValue({
    async *[Symbol.asyncIterator]() {
      for (let i = 0; i < 2; i++) {
        const request = createRequest(i);

        yield {
          eventName: `request://org${i}/device${i}/service${i}/request${i}/request`,
          payload: utf8Encoder.encode(request.serialize()),
        };
      }
      for (let i = 2; i < 4; i++) {
        const response = new ServiceResponse(`request${i}`, new Date(), 1, '[]');

        yield {
          eventName: `request://org${i}/device${i}/service${i}/request${i}/respond`,
          payload: utf8Encoder.encode(response.serialize()),
        };
      }
      for (let i = 4; i < 6; i++) {
        yield {
          eventName: `request://org${i}/device${i}/service${i}/request${i}/remove`,
          payload: utf8Encoder.encode(`request${i}`),
        };
      }
    },
    close: jest.fn(),
  } as CloseableAsyncIterable<ChaincodeEvent>);

  const events = await serviceBroker.registerEvent();

  try {
    let i = 0;
    for await (const event of events) {
      expect(event.organizationId).toBe(`org${i}`);
      expect(event.deviceId).toBe(`device${i}`);
      expect(event.serviceName).toBe(`service${i}`);
      expect(event.requestId).toBe(`request${i}`);

      if (i < 2) {
        expect(event.action).toBe('request');
        expect((event.payload as ServiceRequest).id).toBe(`request${i}`);
        expect((event.payload as ServiceRequest).service.name).toBe(`service${i}`);
        expect((event.payload as ServiceRequest).service.deviceId).toBe(`device${i}`);
        expect((event.payload as ServiceRequest).service.organizationId).toBe(`org${i}`);
        expect((event.payload as ServiceRequest).method).toBe(`GET`);
        expect((event.payload as ServiceRequest).args).toEqual([]);
      } else if (i < 4) {
        expect(event.action).toBe('respond');
        expect((event.payload as ServiceResponse).requestId).toBe(`request${i}`);
        expect((event.payload as ServiceResponse).statusCode).toBe(1);
        expect((event.payload as ServiceResponse).returnValue).toBe(`[]`);
      } else {
        expect(event.action).toBe('remove');
        expect(event.payload).toBe(`request${i}`);
      }

      i++;
    }
    expect(i).toBe(6);
  } finally {
    events.close();
  }

  contract.registerEvent = jest.fn().mockRejectedValue(new Error());

  expect(serviceBroker.registerEvent()).rejects.toThrow();
});
