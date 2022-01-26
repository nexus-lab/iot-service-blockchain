import moment from './moment';
import Service from './Service';
import ServiceRequest from './ServiceRequest';

test('request.getKeyComponents()', () => {
  const request = new ServiceRequest(
    'ffbc9005-c62a-4563-a8f7-b32bba27d707',
    moment('2021-12-12T17:34:00-05:00'),
    new Service('service1', 'device1', 'org1'),
    'GET',
    ['1', '2', '3'],
  );

  expect(request.getKeyComponents()).toEqual([request.id]);
});

test('request.toObject()', () => {
  const request = new ServiceRequest(
    'ffbc9005-c62a-4563-a8f7-b32bba27d707',
    moment('2021-12-12T17:34:00-05:00'),
    new Service('service1', 'device1', 'org1'),
    'GET',
    ['1', '2', '3'],
  );
  const obj = {
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
  };

  expect(request.toObject()).toEqual(obj);
});

test('request.serialize()', () => {
  const request = new ServiceRequest(
    'ffbc9005-c62a-4563-a8f7-b32bba27d707',
    moment('2021-12-12T17:34:00-05:00'),
    new Service('service1', 'device1', 'org1'),
    'GET',
    ['1', '2', '3'],
  );
  const serialized =
    '{"id":"ffbc9005-c62a-4563-a8f7-b32bba27d707",' +
    '"time":"2021-12-12T17:34:00.000-05:00",' +
    '"service":{"name":"service1","deviceId":"device1",' +
    '"organizationId":"org1","version":0,"description":"",' +
    '"lastUpdateTime":"1969-12-31T19:00:00.000-05:00"},"method":"GET",' +
    '"arguments":["1","2","3"]}';

  expect(request.serialize()).toEqual(serialized);
});

test('request.validate()', () => {
  const request = new ServiceRequest('123456', moment(0), new Service('', '', ''), '');

  expect(() => request.validate()).toThrow(/request ID/);
  request.id = 'ffbc9005-c62a-4563-a8f7-b32bba27d707';

  expect(() => request.validate()).toThrow(/requested service/);
  request.service.organizationId = 'org1';

  expect(() => request.validate()).toThrow(/requested service/);
  request.service.deviceId = 'device1';

  expect(() => request.validate()).toThrow(/requested service/);
  request.service.name = 'service1';

  expect(() => request.validate()).toThrow(/request method/);
  request.method = 'GET';

  expect(() => request.validate()).toThrow(/request time/);
  request.time = moment('2021-12-12T17:34:00-05:00');

  expect(() => request.validate()).not.toThrow();
});

test('ServiceRequest.fromObject()', () => {
  const obj = {
    id: 'ffbc9005-c62a-4563-a8f7-b32bba27d707',
    time: moment('2021-12-12T17:34:00-05:00'),
    service: {
      name: 'service1',
      deviceId: 'device1',
      organizationId: 'org1',
    },
    method: 'GET',
    arguments: ['1', '2', '3'],
  };

  const request = ServiceRequest.fromObject(obj);
  expect(request.id).toEqual(obj.id);
  expect(request.time).toEqual(obj.time);
  expect(request.method).toEqual(obj.method);
  expect(request.args).toEqual(obj.arguments);
  expect(request.service.organizationId).toEqual(obj.service.organizationId);
  expect(request.service.deviceId).toEqual(obj.service.deviceId);
  expect(request.service.name).toEqual(obj.service.name);
});

test('ServiceRequest.deserialize()', () => {
  const expected = new ServiceRequest(
    'ffbc9005-c62a-4563-a8f7-b32bba27d707',
    moment('2021-12-12T17:34:00-05:00'),
    new Service('service1', 'device1', 'org1'),
    'GET',
    ['1', '2', '3'],
  );
  const serialized =
    '{"id":"ffbc9005-c62a-4563-a8f7-b32bba27d707",' +
    '"time":"2021-12-12T17:34:00-05:00",' +
    '"service":{"name":"service1","deviceId":"device1",' +
    '"organizationId":"org1","version":0,"description":"",' +
    '"lastUpdateTime":"1970-01-01T00:00:00.000Z"},"method":"GET",' +
    '"arguments":["1","2","3"]}';

  const actual = ServiceRequest.deserialize(serialized);
  expect(actual.id).toEqual(expected.id);
  expect(actual.time).toEqual(expected.time);
  expect(actual.method).toEqual(expected.method);
  expect(actual.args).toEqual(expected.args);
  expect(actual.service.organizationId).toEqual(expected.service.organizationId);
  expect(actual.service.deviceId).toEqual(expected.service.deviceId);
  expect(actual.service.name).toEqual(expected.service.name);

  expect(() => ServiceRequest.deserialize('\x00')).toThrow();
});
