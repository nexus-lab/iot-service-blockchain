import Service from './Service';
import ServiceRequest from './ServiceRequest';
import ServiceRequestResponse from './ServiceRequestResponse';
import ServiceResponse from './ServiceResponse';
import moment from './moment';

test('pair.toObject()', () => {
  const pair = new ServiceRequestResponse(
    new ServiceRequest(
      'ffbc9005-c62a-4563-a8f7-b32bba27d707',
      moment('2021-12-12T17:34:00-05:00'),
      new Service('service1', 'device1', 'org1'),
      'GET',
      ['1', '2', '3'],
    ),
    new ServiceResponse(
      'ffbc9005-c62a-4563-a8f7-b32bba27d707',
      moment('2021-12-12T17:34:00-05:00'),
      0,
      '["a","b","c"]',
    ),
  );
  const obj = {
    request: {
      id: 'ffbc9005-c62a-4563-a8f7-b32bba27d707',
      time: moment('2021-12-12T17:34:00-05:00'),
      service: {
        name: 'service1',
        deviceId: 'device1',
        organizationId: 'org1',
        version: 0,
        description: '',
        lastUpdateTime: moment(0),
      },
      method: 'GET',
      arguments: ['1', '2', '3'],
    },
    response: {
      requestId: 'ffbc9005-c62a-4563-a8f7-b32bba27d707',
      time: moment('2021-12-12T17:34:00-05:00'),
      statusCode: 0,
      returnValue: '["a","b","c"]',
    },
  };

  expect(pair.toObject()).toEqual(obj);
});

test('pair.serialize()', () => {
  const pair = new ServiceRequestResponse(
    new ServiceRequest(
      'ffbc9005-c62a-4563-a8f7-b32bba27d707',
      moment('2021-12-12T17:34:00-05:00'),
      new Service('service1', 'device1', 'org1'),
      'GET',
      ['1', '2', '3'],
    ),
    new ServiceResponse(
      'ffbc9005-c62a-4563-a8f7-b32bba27d707',
      moment('2021-12-12T17:34:00-05:00'),
      0,
      '["a","b","c"]',
    ),
  );

  const serialized =
    '{"request":{"id":"ffbc9005-c62a-4563-a8f7-b32bba27d707","time":"2021-12-12T17:34:00.000-05:00",' +
    '"service":{"name":"service1","deviceId":"device1","organizationId":"org1","version":0,' +
    '"description":"","lastUpdateTime":"1969-12-31T19:00:00.000-05:00"},"method":"GET",' +
    '"arguments":["1","2","3"]},"response":{"requestId":"ffbc9005-c62a-4563-a8f7-b32bba27d707",' +
    '"time":"2021-12-12T17:34:00.000-05:00","statusCode":0,"returnValue":"[\\"a\\",\\"b\\",\\"c\\"]"}}';

  expect(pair.serialize()).toEqual(serialized);
});

test('ServiceRequestResponse.fromObject()', () => {
  const obj = {
    request: {
      id: 'ffbc9005-c62a-4563-a8f7-b32bba27d707',
      time: moment('2021-12-12T17:34:00-05:00'),
      service: {
        name: 'service1',
        deviceId: 'device1',
        organizationId: 'org1',
      },
      method: 'GET',
      arguments: ['1', '2', '3'],
    },
    response: {
      requestId: 'ffbc9005-c62a-4563-a8f7-b32bba27d707',
      time: moment('2021-12-12T17:34:00-05:00'),
      statusCode: 0,
      returnValue: '["a","b","c"]',
    },
  };

  const pair = ServiceRequestResponse.fromObject(obj);
  expect(pair.request.id).toEqual(obj.request.id);
  expect(pair.request.time).toEqual(obj.request.time);
  expect(pair.request.method).toEqual(obj.request.method);
  expect(pair.request.args).toEqual(obj.request.arguments);
  expect(pair.request.service.organizationId).toEqual(obj.request.service.organizationId);
  expect(pair.request.service.deviceId).toEqual(obj.request.service.deviceId);
  expect(pair.request.service.name).toEqual(obj.request.service.name);
  expect(pair.response.requestId).toEqual(obj.response.requestId);
  expect(pair.response.time).toEqual(obj.response.time);
  expect(pair.response.returnValue).toEqual(obj.response.returnValue);
  expect(pair.response.statusCode).toEqual(obj.response.statusCode);
});

test('ServiceRequestResponse.deserialize()', () => {
  const expected = new ServiceRequestResponse(
    new ServiceRequest(
      'ffbc9005-c62a-4563-a8f7-b32bba27d707',
      moment('2021-12-12T17:34:00-05:00'),
      new Service('service1', 'device1', 'org1'),
      'GET',
      ['1', '2', '3'],
    ),
    new ServiceResponse(
      'ffbc9005-c62a-4563-a8f7-b32bba27d707',
      moment('2021-12-12T17:34:00-05:00'),
      0,
      '["a","b","c"]',
    ),
  );

  const serialized =
    '{"request":{"id":"ffbc9005-c62a-4563-a8f7-b32bba27d707","time":"2021-12-12T17:34:00-05:00",' +
    '"service":{"name":"service1","deviceId":"device1","organizationId":"org1","version":0,' +
    '"description":"","lastUpdateTime":"1969-12-31T19:00:00-05:00"},"method":"GET",' +
    '"arguments":["1","2","3"]},"response":{"requestId":"ffbc9005-c62a-4563-a8f7-b32bba27d707",' +
    '"time":"2021-12-12T17:34:00-05:00","statusCode":0,"returnValue":"[\\"a\\",\\"b\\",\\"c\\"]"}}';

  const actual = ServiceRequestResponse.deserialize(serialized);
  expect(actual.request.id).toEqual(expected.request.id);
  expect(actual.request.time).toEqual(expected.request.time);
  expect(actual.request.method).toEqual(expected.request.method);
  expect(actual.request.args).toEqual(expected.request.args);
  expect(actual.request.service.organizationId).toEqual(expected.request.service.organizationId);
  expect(actual.request.service.deviceId).toEqual(expected.request.service.deviceId);
  expect(actual.request.service.name).toEqual(expected.request.service.name);
  expect(actual.response.requestId).toEqual(expected.response.requestId);
  expect(actual.response.time).toEqual(expected.response.time);
  expect(actual.response.returnValue).toEqual(expected.response.returnValue);
  expect(actual.response.statusCode).toEqual(expected.response.statusCode);

  expect(() => ServiceResponse.deserialize('\x00')).toThrow();
});
