import moment from './moment';
import Service from './Service';

test('service.getKeyComponents()', () => {
  const service = new Service(
    'service1',
    'device1',
    'org1',
    1,
    'Service of Device1',
    moment('2021-12-12T17:34:00-05:00'),
  );

  expect(service.getKeyComponents()).toEqual([
    service.organizationId,
    service.deviceId,
    service.name,
  ]);
});

test('service.toObject()', () => {
  const service = new Service(
    'service1',
    'device1',
    'org1',
    1,
    'Service of Device1',
    moment('2021-12-12T17:34:00-05:00'),
  );
  const obj = {
    name: 'service1',
    deviceId: 'device1',
    organizationId: 'org1',
    version: 1,
    description: 'Service of Device1',
    lastUpdateTime: moment('2021-12-12T17:34:00-05:00'),
  };

  expect(service.toObject()).toEqual(obj);
});

test('service.serialize()', () => {
  const service = new Service(
    'service1',
    'device1',
    'org1',
    1,
    'Service of Device1',
    moment('2021-12-12T17:34:00-05:00'),
  );
  const serialized =
    '{"name":"service1","deviceId":"device1","organizationId":"org1",' +
    '"version":1,"description":"Service of Device1",' +
    '"lastUpdateTime":"2021-12-12T17:34:00.000-05:00"}';

  expect(service.serialize()).toEqual(serialized);
});

test('service.validate()', () => {
  const service = new Service('', '', '');

  expect(() => service.validate()).toThrow(/service name/);
  service.name = 'service1';

  expect(() => service.validate()).toThrow(/device ID/);
  service.deviceId = 'device1';

  expect(() => service.validate()).toThrow(/organization ID/);
  service.organizationId = 'org1';

  expect(() => service.validate()).toThrow(/service version/);
  service.version = 1;

  expect(() => service.validate()).toThrow(/last update time/);
  service.lastUpdateTime = moment('2021-12-12T17:34:00-05:00');

  expect(() => service.validate()).not.toThrow();
});

test('Service.fromObject()', () => {
  const obj = {
    name: 'service1',
    deviceId: 'device1',
    organizationId: 'org1',
    version: 1,
    description: 'Service of Device1',
    lastUpdateTime: moment('2021-12-12T17:34:00-05:00'),
  };

  const service = Service.fromObject(obj);
  expect(service.name).toEqual(obj.name);
  expect(service.deviceId).toEqual(obj.deviceId);
  expect(service.organizationId).toEqual(obj.organizationId);
  expect(service.version).toEqual(obj.version);
  expect(service.description).toEqual(obj.description);
  expect(service.lastUpdateTime).toEqual(obj.lastUpdateTime);
});

test('Service.deserialize()', () => {
  const expected = new Service(
    'service1',
    'device1',
    'org1',
    1,
    'Service of Device1',
    moment('2021-12-12T17:34:00-05:00'),
  );
  const serialized =
    '{"name":"service1","deviceId":"device1","organizationId":"org1",' +
    '"version":1,"description":"Service of Device1",' +
    '"lastUpdateTime":"2021-12-12T17:34:00-05:00"}';

  const actual = Service.deserialize(serialized);
  expect(actual.name).toEqual(expected.name);
  expect(actual.deviceId).toEqual(expected.deviceId);
  expect(actual.organizationId).toEqual(expected.organizationId);
  expect(actual.version).toEqual(expected.version);
  expect(actual.description).toEqual(expected.description);
  expect(actual.lastUpdateTime).toEqual(expected.lastUpdateTime);

  expect(() => Service.deserialize('\x00')).toThrow();
});
